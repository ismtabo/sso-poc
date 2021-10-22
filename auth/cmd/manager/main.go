package main

import (
	"context"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/ismtabo/sso-poc/auth/pkg/cfg"
	"github.com/ismtabo/sso-poc/auth/pkg/repository"
	"github.com/ismtabo/sso-poc/auth/pkg/service"
	"github.com/rs/zerolog"
	"github.com/sonyarouje/simdb"
)

const (
	USAGE = `Fourth Platform HaaC Subscription Controller

 usage:
   manager [options] <command> -- [<cmdArgs>...]
   manager --version
   manager --help

 Options:
  -h, --help                  show this help message and exit
  --version                   show version and exit
  -c CONFIG, --config CONFIG  config file [default: ./config.yml]
  -q, --quiet                 silent (without output messages)

commands:
  create-user      Create user in application
  help             Show more information on a specific command

See 'manager help -- <command>' for more information on a specific command.
`
)

var (
	version = "0.0.0" // meant to be set from outside
)

func main() {
	log := getLogger()
	args := os.Args[1:]
	globalArgs, _ := docopt.ParseArgs(USAGE, args, version)
	cmdName := globalArgs["<command>"].(string)
	cmdArgs := globalArgs["<cmdArgs>"].([]string)
	commandArgs := []string{cmdName}
	if cmdName != "help" {
		commandArgs = append(commandArgs, cmdArgs...)
	} else {
		commandArgs = append(cmdArgs, "--help")
		if len(cmdArgs) < 1 {
			log.Fatal().Msg("missing command arg")
		}
		cmdName = cmdArgs[0]
	}
	var config Config
	if err := cfg.Load(globalArgs["--config"].(string), &config); err != nil {
		log.Fatal().Msgf("Error loading configuration. %s", err)
	}
	log.Debug().Msgf("Configuration loaded: %+v", config)
	if err := configLogger(&config); err != nil {
		log.Fatal().Msgf("Error configuring the logger. %s", err)
	}
	driver := mustCreateSimdbDriver(&config, log)
	userRepo := repository.NewSimdbUserRepository(driver)
	authSvc := service.NewAuthService(userRepo)
	ctx := withAppCtx(context.Background(), log, authSvc)
	switch cmdName {
	case "create-user":
		if err := createUserCmd(ctx, globalArgs, commandArgs); err != nil {
			log.Fatal().Err(err).Msg("")
		}
	default:
		log.Fatal().Msgf("invalid command: %s", cmdName)
	}
}

func getLogger() *zerolog.Logger {
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000Z07:00"
	zerolog.TimestampFieldName = "time"
	zerolog.LevelFieldName = "lvl"
	zerolog.MessageFieldName = "msg"
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &log
}

func configLogger(config *Config) error {
	lvl, err := zerolog.ParseLevel(strings.ToLower(config.Log.Level))
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(lvl)
	return nil
}

func mustCreateSimdbDriver(config *Config, log *zerolog.Logger) *simdb.Driver {
	driver, err := simdb.New(config.Database.File)
	if err != nil {
		log.Fatal().Msgf("Error creating simdb driver. %s", err)
	}
	return driver
}

func createUserCmd(ctx context.Context, globalArgs map[string]interface{}, commandArgs []string) error {
	const usage = `Create user in application

Create user in application.

 usage:
   manager create-user [options] <username> <password>

 Options:
  -h, --help               show this help message and exit
`
	args, err := docopt.ParseArgs(usage, commandArgs, "")
	if err != nil {
		return err
	}
	appCtx := getAppCtx(ctx)
	uid := args["<username>"].(string)
	password := args["<password>"].(string)
	if err := appCtx.authSvc.Register(uid, password); err != nil {
		return err
	}
	appCtx.log.Info().Msgf("successfully create user '%s'", uid)
	return nil
}
