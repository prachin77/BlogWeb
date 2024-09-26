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
    BlogTitle        string `bson:"blog_title" json:"blogTitle"`
    BlogContent      string `bson:"blog_content" json:"blogContent"`
    AuthorName       string `bson:"author_name" json:"authorname"`
    BlogCreationDate string `bson:"blog_creation_date" json:"blog_creation_date"`
    Tag              string `bson:"blog_tag" json:"tag"`
    Views            int    `bson:"views" json:"views"`
    Likes            int    `bson:"likes" json:"likes"`
    Comments         string `bson:"comments" json:"comments"`
}

