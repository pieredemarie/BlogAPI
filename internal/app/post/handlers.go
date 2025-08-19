package post

import (
	"BlogAPI/pkg/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (p *PostHandler) CreatePost(c *gin.Context) {
	var newPost storage.Post

	if err := c.BindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error": "invalid request"})
		return
	}	

	err := p.Storage.CreatePost(newPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"error": "server error"})
		return
	}

	c.Status(http.StatusCreated)
}

func (p *PostHandler) GetPosts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	posts, err := p.Storage.GetPosts(100) 
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})

}

func (p *PostHandler) GetPostById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	post, err := p.Storage.GetPostById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

func (p *PostHandler) EditPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error": "invalid id"})
		return
	}

	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authorId := int(userId.(float64))

	var req EditPostRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error": "invalid request"})
		return
	}

	err = p.Storage.EditPost(id,authorId,req.Title,req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post updated"})
}
func (p *PostHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error": "invalid id"})
		return
	}

	err = p.Storage.DeletePost(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK,gin.H{"message": "post deleted"})
}

