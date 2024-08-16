package models

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	UserId   string `json:"userid"`
}

// for client side 
type AuthPageStatus struct {
    IsLogin bool
}
