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
	UserId string `json:"userid"`
	BlogId string `json:"blogid"`
	BlogCreationDate string `json:"blog_creation_date"`
	BlogTitle string `json:"blogtitle"`
	BlogContent string `json:"blogcontent"`
	BlogImage string `json:"blogimage"`
	Tags []string `json:"tags"`
	Views int `json:"views"`
	Likes int `json:"likes"`
	Comments string `json:"comments"`
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


