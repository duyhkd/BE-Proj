package handlers

import (
	"Server/db"
	"Server/httpServer"
	"Server/model"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func MakePost(c *gin.Context) {
	username := c.Query("username")

	bytedata, err := io.ReadAll(c.Request.Body)
	if err != nil {
		httpServer.BadRequest(c, err.Error())
	}
	text := string(bytedata)

	post := model.Post{
		UserName: username,
		Text:     text,
	}
	result := db.DB.Create(&post)

	if result.Error != nil {
		httpServer.BadRequest(c, result.Error.Error())
	} else {
		httpServer.Ok(c, "Posted!")
	}
}

func RemovePost(c *gin.Context) {
	username := c.Value("username").(string)
	postId, _ := uuid.Parse(c.Query("post"))

	var post model.Post
	result := db.DB.First(&post, "id = ?", postId)
	if result.Error != nil {
		httpServer.NotFound(c, "Post doesn't exist")
	}

	if username != post.UserName {
		httpServer.Unauthorized(c, "User is not allow to delete this post")
		return
	}

	if err := db.DB.Where("post_id = ?", postId).Delete(&model.Comment{}).Error; err != nil {
		httpServer.BadRequest(c, "Failed to delete comments associated with the post")
		return
	}

	if err := db.DB.Where("post_id = ?", postId).Delete(&model.Like{}).Error; err != nil {
		httpServer.BadRequest(c, "Failed to delete likes associated with the post")
		return
	}

	if err := db.DB.Delete(post).Error; err != nil {
		httpServer.BadRequest(c, "Failed to delete post")
		return
	}

	httpServer.Ok(c, "Post deleted!")
}

func EditPost(c *gin.Context) {
	username := c.Value("username").(string)
	postId, err := uuid.Parse(c.Query("post"))
	if err != nil {
		httpServer.BadRequest(c, err.Error())
	}

	bytedata, err := io.ReadAll(c.Request.Body)
	if err != nil {
		httpServer.BadRequest(c, err.Error())
	}
	text := string(bytedata)

	var post model.Post

	result := db.DB.First(&post, "id = ?", postId)
	if result.Error != nil {
		httpServer.NotFound(c, "Post doesn't exist")
	}

	if username != post.UserName {
		httpServer.Unauthorized(c, "User is not allow to updated this post")
		return
	}

	post.Text = text
	if err = db.DB.Save(&post).Error; err != nil {
		httpServer.BadRequest(c, "Failed to update post")
		return
	}

	httpServer.Ok(c, "Post updated!")
}

func MakeComment(c *gin.Context) {
	username := c.Value("username").(string)
	postId, err := uuid.Parse(c.Query("post"))
	if err != nil {
		httpServer.BadRequest(c, err.Error())
	}

	bytedata, err := io.ReadAll(c.Request.Body)
	if err != nil {
		httpServer.BadRequest(c, err.Error())
	}
	text := string(bytedata)

	comment := model.Comment{
		UserName: username,
		Text:     text,
		PostId:   postId,
	}
	result := db.DB.Create(&comment)

	if result.Error != nil {
		httpServer.BadRequest(c, result.Error.Error())
	} else {
		httpServer.Ok(c, "Commented!")
	}
}

func LikePost(c *gin.Context) {
	username := c.Value("username").(string)
	postId, err := uuid.Parse(c.Query("post"))
	if err != nil {
		httpServer.BadRequest(c, err.Error())
		return
	}

	var post model.Post
	result := db.DB.First(&post, "id = ?", postId)
	if result.Error != nil {
		httpServer.NotFound(c, "Post doesn't exist")
		return
	}

	var like model.Like
	result = db.DB.Where("post_id = ?", postId).Where("user_name = ?", username).First(&like)

	if result.RowsAffected > 0 {
		if err = db.DB.Where("post_id = ?", postId).Where("user_name = ?", username).Delete(&like).Error; err != nil {
			httpServer.BadRequest(c, "Failed to unlike post")
			return
		}
		httpServer.Ok(c, "Unliked Post!")
		return
	} else {
		if err = db.DB.Create(&model.Like{
			PostId:   postId,
			UserName: username,
		}).Error; err != nil {
			httpServer.BadRequest(c, "Failed to like post")
			return
		}
		httpServer.Ok(c, "Liked Post!")
		return
	}
}
