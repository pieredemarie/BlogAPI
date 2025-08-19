package auth

import (
	"BlogAPI/pkg/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) RegisterHandler(c *gin.Context) { //auth/register
	var newUser RegisterRequest

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest,ErrorResponce{
			Message: err.Error(),
		})
		return
	}

	user := storage.User{
		Username: newUser.Username,
        Email:    newUser.Email,
        Password: newUser.Password,
	}
	if err := h.Storage.Register(user); err != nil {
		c.JSON(http.StatusBadRequest,ErrorResponce{
			Message: err.Error(),
		})
		return
	}
	
	c.Status(http.StatusCreated)
}

func (h *Handler) LoginHandler(c *gin.Context) {
	var log LoginRequest

	if err := c.BindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest,ErrorResponce{
			Message: err.Error(),
		})
		return
	}
	token, err := h.Storage.Login(log.Email,log.Password); 
	if err != nil {
		c.JSON(http.StatusBadRequest,ErrorResponce{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK,LoginResponce{
		Token: token,
	})
}

func (h *Handler) GetProfile(c *gin.Context) { // GET  users/me
	userID, _ := c.Get("userID")

	user, err := h.Storage.GetUserById(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound,gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": user.ID,
		"username": user.Username,
		"email": user.Email,
	})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	UserID, _ := c.Get("userID")
	var req UpdateProfileRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error": "Invalid request"})
		return
	}

	var hashedPass string 
	if req.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password),bcrypt.DefaultCost)
		hashedPass = string(hash)
	}

	err := h.Storage.UpdateProfile(UserID.(int), req.Username,req.Email,hashedPass)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"error": "update failed"})
		return
	}

	c.Status(http.StatusOK)
}
