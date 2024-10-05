package delivery

import (
	"log"
	"net/http"

	"github.com/Sokcheatsrorng/go-clean-architecture/internal/model"
	"github.com/Sokcheatsrorng/go-clean-architecture/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateUsers([]model.User{user}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *UserHandler) GetAllUsers(ctx *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserById(ctx *gin.Context) {
	id := ctx.Param("_id")

	user, err := h.service.GetUserById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUserById(ctx *gin.Context) {
    id := ctx.Param("_id")  

    var user model.User
    if err := ctx.ShouldBindJSON(&user); err != nil {
        log.Println("Binding error:", err) 
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
        return
    }

    if user.ID != id {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID in URL does not match ID in payload"})
        return
    }

    err := h.service.UpdateUsers(id, []model.User{user})
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *UserHandler) DeleteUserById(ctx *gin.Context) {
	id := ctx.Param("_id")

	if err := h.service.DeleteUsers(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
