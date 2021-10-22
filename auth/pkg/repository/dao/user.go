package dao

type User struct {
	UID      string `json:"username"`
	Password string `json:"password"`
}

func (u User) ID() (jsonField string, value interface{}) {
	return "username", u.UID
}
