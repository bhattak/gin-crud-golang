package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"project/model"
	"project/util"

	"github.com/gin-gonic/gin"
)

// UserAuthService represents the interface for the userauth service.
type UserAuthService interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type userAuthServiceImpl struct {
	userAuthRepo model.UserAuthRepository
}

// NewUserAuthService creates a new user service with the given userAuth repository.
func NewUserAuthService(userAuthRepo model.UserAuthRepository) UserAuthService {
	return &userAuthServiceImpl{userAuthRepo}
}

// Login functanality
func (service *userAuthServiceImpl) Login(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		util.WriteError(c.Writer, http.StatusBadRequest, "failed to read request body")
		return
	}

	user := &model.LoginRequest{}
	if err := json.Unmarshal(body, user); err != nil {
		util.WriteError(c.Writer, http.StatusBadRequest, "failed to parse request body")
		return
	}
	token, err := service.userAuthRepo.Login(user)
	if err != nil {
		util.WriteError(c.Writer, http.StatusInternalServerError, fmt.Sprintf("failed to login user: %v", err))
		return
	}

	util.WriteJSON(c.Writer, http.StatusOK, gin.H{
		"message":       "Logged in succesfully",
		"Authorization": token,
	})
}

func (service *userAuthServiceImpl) Register(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		util.WriteError(c.Writer, http.StatusBadRequest, "failed to read request body")
		return
	}

	user := &model.LoginRequest{}
	if err := json.Unmarshal(body, user); err != nil {
		util.WriteError(c.Writer, http.StatusBadRequest, "failed to parse request body")
		return
	}
	if err := service.userAuthRepo.Register(user); err != nil {
		util.WriteError(c.Writer, http.StatusInternalServerError, fmt.Sprintf("failed to create user: %v", err))
		return
	}

	util.WriteJSON(c.Writer, http.StatusCreated, gin.H{"message": "user added successfully. please login to get full access"})

}
