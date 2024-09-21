package models

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	UserId   string `json:"userid"`
}

// for client side authentication frontend
type AuthPageStatus struct {
	IsLogin bool
}

// for blogs
type Blog struct {
	UserId           string `bson:"userid" json:"userid"`
	BlogId           string `bson:"blogid" json:"blogid"`
	AuthorName       string `bson:"authorname" json:"authorname"`
	BlogCreationDate string `bson:"blog_creation_date" json:"blog_creation_date"`
	BlogTitle        string `bson:"blogtitle" json:"blogtitle"`
	BlogContent      string `bson:"blogcontent" json:"blogcontent"`
	BlogImage        []byte `bson:"blogimage" json:"blogimage"`
	Tags             string `bson:"tags" json:"tags"`
	Views            int    `bson:"views" json:"views"`
	Likes            int    `bson:"likes" json:"likes"`
	Comments         string `bson:"comments" json:"comments"`
}

// Blog storing format in mongo db
// userid : {
// 	// user id & than all the blogs that user has created in list format
// 	[
// 		{
// 			blogid : "",
// 			blog_creation_date:"",
// 			blogtitle:"",
// 			blogcontent:"",
// 			blogimage:"",
// 			tags:"",
// 			views:"",
// 			likes:"",
// 			comments:""
// 		},
// 		{
// 			blogid : "",
// 			blog_creation_date:"",
// 			blogtitle:"",
// 			blogcontent:"",
// 			blogimage:"",
// 			tags:"",
// 			views:"",
// 			likes:"",
// 			comments:""
// 		},
// 		........ other blogs
// 	]
// }

// NORMAL FORMAT -> insert as it is
