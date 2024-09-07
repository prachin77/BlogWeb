package clienthandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prachin77/db"
	"github.com/prachin77/server/models"
	"github.com/prachin77/server/utils"
)

func PostBlog(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")

	// Check for Content-Type
	contentType := ctx.Request.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		fmt.Println("Content-Type must be multipart/form-data")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be multipart/form-data"})
		return
	}

	var blog models.Blog

	blog.BlogTitle = ctx.Request.PostFormValue("blogtitle")
	blog.BlogContent = ctx.Request.PostFormValue("blogcontent")
	blog.Tags = ctx.Request.PostFormValue("tags")
	fmt.Println("blog content = ", blog.BlogContent)

	file, header, err := ctx.Request.FormFile("blogimage")
	if err != nil {
		fmt.Println("Error fetching image from user:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch image"})
		return
	}
	defer file.Close()

	imageData, err := utils.FileHeaderToBytes(header)
	if err != nil {
		fmt.Println("Error reading file:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image"})
		return
	}

	// BYTES FURTHER ENCODED TO STRING USING BASE64 TO STORE IN MONGO DB
	// Convert imageData to base64-encoded string
	// imageDataStr := base64.StdEncoding.EncodeToString(imageData)

	// // Decode the base64 string back to bytes
	// decodedImageData, err := base64.StdEncoding.DecodeString(imageDataStr)
	// if err != nil {
	// 	fmt.Println("Error decoding base64 string:", err)
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode image data"})
	// 	return
	// }

	// // Verify that the decoded data matches the original image data
	// if bytes.Equal(imageData, decodedImageData) {
	// 	fmt.Println("Decoded image data matches the original data.")
	// } else {
	// 	fmt.Println("Decoded image data does NOT match the original data.")
	// }

	// create blog id
	blog.BlogId = utils.TokenGenerator()

	// extract userid from cookie
	userid, tokenString := utils.GetCookie(ctx)
	if tokenString == "" {
		fmt.Println("cookie is null")
		RenderLoginPage(ctx)
		return
	}
	blog.UserId = userid
	userDetails, err := db.SearchUserWithId(blog.UserId)
	if err != nil {
		fmt.Println("error in finding user = ", err)
		return
	}
	fmt.Println("user details = ", userDetails)
	blog.AuthorName = userDetails.UserName

	blog.BlogCreationDate = utils.GetCurrentDate(ctx)

	respInterface = map[string]interface{}{
		"BlogTitle":           blog.BlogTitle,
		"BlogContent":         blog.BlogContent,
		"BlogTags":            blog.Tags,
		"AuthorName":          blog.AuthorName,
		"ImageDataInBytes":    imageData,
		"BlogImageDataLength": len(imageData),
		"UserId":              blog.UserId,
		"BlogId":              blog.BlogId,
		"BlogCreationDate":    blog.BlogCreationDate,
	}
	// Marshal the map to a JSON-formatted string
	jsonData, err := json.MarshalIndent(respInterface, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	db.AddBlog(respInterface, &blog)

	fmt.Println("Blog details:")
	fmt.Println(string(jsonData))

	RenderHomePage(ctx, userid)
}
