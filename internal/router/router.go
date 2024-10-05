package routes

import (
	"github.com/Sokcheatsrorng/go-clean-architecture/internal/delivery" 
	"github.com/gin-gonic/gin"
)

func InitRoutes(userHandler *delivery.UserHandler) *gin.Engine {

	r := gin.Default()

	userRouter := r.Group("/users") 
	{
		userRouter.POST("", userHandler.CreateUser)
		userRouter.GET("", userHandler.GetAllUsers)
		userRouter.GET("/:_id", userHandler.GetUserById)
		userRouter.PUT("/:_id", userHandler.UpdateUserById)
		userRouter.DELETE("/:_id",userHandler.DeleteUserById)
	}

	return r

}