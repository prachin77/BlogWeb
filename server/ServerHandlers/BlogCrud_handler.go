package serverhandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prachin77/db"
	"github.com/prachin77/server/models"
	"github.com/prachin77/server/utils"
)

func PostBlog(ctx *gin.Context) {
	ctx.Header("content-Type", "application/json")

	var blog models.Blog

	if err := json.NewDecoder(ctx.Request.Body).Decode(&blog); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "data not in correct format"})
		return
	}

	if blog.AuthorName == "" || blog.BlogContent == "" || blog.BlogTitle == "" || blog.Tag == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "fields can't be empty"})
		return
	}

	if len(blog.BlogContent) < 1 || len(blog.BlogContent) > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "length of blog should be between 1 to 100"})
		return
	}

	if len(blog.BlogTitle) < 4 || len(blog.BlogTitle) > 20 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Title length should be between 4 to 20"})
		return
	}

	blog.BlogCreationDate = utils.GetCurrentDate(ctx)

	if insertedBlog, err := db.AddBlog(&blog); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to add blog"})
		return
	}  else {
		fmt.Println("inserted blog : ",insertedBlog)
	}

	respInterface := map[string]interface{}{
		"blog title":         blog.BlogTitle,
		"blog tag":           blog.Tag,
		"blog content":       blog.BlogContent,
		"Username":           blog.AuthorName,
		"Blog creation date": blog.BlogCreationDate,
	}

	fmt.Println("blog details : ", respInterface)
	ctx.JSON(http.StatusOK, respInterface)
}
