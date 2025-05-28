package controllers

import (
	"go-movie-api/movies/model"
	"go-movie-api/movies/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userService service.UserService
}

type UserController interface {
	CreateUser(c *gin.Context)
	GetUsers(c *gin.Context)
}

func NewUserController(userService service.UserService) UserController {
	return userController{userService: userService}
}

func (mc userController) CreateUser(ctx *gin.Context) {
	var createUserReq model.CreateUserRequest
	if err := ctx.ShouldBindJSON(&createUserReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("request is valid")

	err := mc.userService.CreateUser(createUserReq)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, model.CreateUserResponse{Status: "Success"})
}

func (mc userController) GetUsers(ctx *gin.Context) {
	resp, err := mc.userService.GetUsers()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}
