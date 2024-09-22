package handlers

import (
	"Server/db"
	"Server/httpServer"
	"Server/model"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func MakePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpServer.MethodNotAllowed(w)
	}

	username := r.URL.Query().Get("username")

	bytedata, err := io.ReadAll(r.Body)
	if err != nil {
		httpServer.BadRequest(w, err.Error())
	}
	text := string(bytedata)

	post := model.Post{
		UserName: username,
		Text:     text,
	}
	result := db.DB.Create(&post)

	if result.Error != nil {
		httpServer.BadRequest(w, result.Error.Error())
	} else {
		httpServer.Ok(w, "Posted!")
	}
}

func RemovePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpServer.MethodNotAllowed(w)
	}

	username := r.Context().Value("username").(string)
	postId, _ := uuid.Parse(r.URL.Query().Get("post"))

	var post model.Post
	result := db.DB.First(&post, "id = ?", postId)
	if result.Error != nil {
		httpServer.NotFound(w, "Post doesn't exist")
	}

	if username != post.UserName {
		httpServer.Unauthorized(w, "User is not allow to delete this post")
		return
	}

	if err := db.DB.Where("post_id = ?", postId).Delete(&model.Comment{}).Error; err != nil {
		httpServer.BadRequest(w, "Failed to delete comments associated with the post")
		return
	}

	if err := db.DB.Delete(post).Error; err != nil {
		httpServer.BadRequest(w, "Failed to delete post")
		return
	}

	httpServer.Ok(w, "Post deleted!")
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpServer.MethodNotAllowed(w)
	}

	username := r.Context().Value("username").(string)
	postId, _ := uuid.Parse(r.URL.Query().Get("post"))

	bytedata, err := io.ReadAll(r.Body)
	if err != nil {
		httpServer.BadRequest(w, err.Error())
	}
	text := string(bytedata)

	var post model.Post

	result := db.DB.First(&post, "id = ?", postId)
	if result.Error != nil {
		httpServer.NotFound(w, "Post doesn't exist")
	}

	if username != post.UserName {
		httpServer.Unauthorized(w, "User is not allow to updated this post")
		return
	}

	post.Text = text
	if err = db.DB.Save(&post).Error; err != nil {
		httpServer.BadRequest(w, "Failed to update post")
		return
	}

	httpServer.Ok(w, "Post updated!")
}

func MakeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpServer.MethodNotAllowed(w)
	}

	username := r.URL.Query().Get("username")
	postId, err := uuid.Parse(r.URL.Query().Get("post"))
	if err != nil {
		httpServer.BadRequest(w, err.Error())
	}

	bytedata, err := io.ReadAll(r.Body)
	if err != nil {
		httpServer.BadRequest(w, err.Error())
	}
	text := string(bytedata)

	comment := model.Comment{
		UserName: username,
		Text:     text,
		PostId:   postId,
	}
	result := db.DB.Create(&comment)

	if result.Error != nil {
		httpServer.BadRequest(w, result.Error.Error())
	} else {
		httpServer.Ok(w, "Commented!")
	}
}
