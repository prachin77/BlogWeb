package models

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	UserId   string `json:"userid"`
}

type AuthPageStatus struct {
    IsLogin bool
}
