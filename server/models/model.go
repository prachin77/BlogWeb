package models

import "database/sql"

type User struct {
	UserName     string         `json:"username"`
	Password     string         `json:"password"`
	SessionToken sql.NullString `json:"session_token"`
}

// for client side authentication frontend
type AuthPageStatus struct {
	IsLogin bool
}

// for blogs
type Blog struct {
	AuthorName       string `bson:"authorname" json:"authorname"`
	BlogCreationDate string `bson:"blog_creation_date" json:"blog_creation_date"`
	BlogTitle        string `bson:"blogtitle" json:"blogtitle"`
	BlogContent      string `bson:"blogcontent" json:"blogcontent"`
	Tag              string `bson:"tag" json:"tag"`
	Views            int    `bson:"views" json:"views"`
	Likes            int    `bson:"likes" json:"likes"`
	Comments         string `bson:"comments" json:"comments"`
}
