package clienthandlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/prachin77/server/models"
	"github.com/prachin77/server/utils"
)

func PostBlog(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")

	userid, tokenString := utils.GetCookie(ctx)
	if tokenString == "" {
		fmt.Println("cookie is null, session timed out")
		RenderLoginPage(ctx)
		return
	}
	fmt.Println("user id: ", userid)

	var blog models.Blog

	blog.BlogTitle = ctx.Request.PostFormValue("blogtitle")
	blog.Tag = ctx.Request.PostFormValue("tag")
	blog.BlogContent = ctx.Request.PostFormValue("blogcontent")

	blogImage, err := ctx.FormFile("blogimage")
	if err != nil {
		fmt.Println("error fetching image:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error fetching image"})
		return
	}

	file, err := blogImage.Open()
	if err != nil {
		fmt.Println("error opening image file:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening image file"})
		return
	}
	defer file.Close()

	// Define the directory for storing user images
	userImageDir := filepath.Join("P:/BlogWeb/db/user_uploaded_images", userid)
	
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(userImageDir, os.ModePerm); err != nil {
		fmt.Println("error creating directory:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating directory"})
		return
	}

	// Define the path to save the uploaded image
	imagePath := filepath.Join(userImageDir, blogImage.Filename)
	
	// Save the uploaded image	
	out, err := os.Create(imagePath)
	if err != nil {
		fmt.Println("error creating image file:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving image file"})
		return
	}
	defer out.Close()

	// Copy the image content to the file
	if _, err := io.Copy(out, file); err != nil {
		fmt.Println("error copying image content:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving image content"})
		return
	}

	// Construct the path to return in the response
	userImagePath := filepath.Join("P:/BlogWeb/db/user_uploaded_images", userid, blogImage.Filename)


	respInterface := map[string]interface{}{
		"blog title":   blog.BlogTitle,
		"blog tag":     blog.Tag,
		"blog content": blog.BlogContent,
		"blog image":   blogImage.Filename,
		"user image path" : userImagePath,
	}

	fmt.Println("blog details:", respInterface)
	ctx.JSON(http.StatusOK, respInterface)
}
