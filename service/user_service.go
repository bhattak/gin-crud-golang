// service/user_service.go
package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"project/model"
	"project/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserService represents the interface for the user service.
type UserService interface {
	AddUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	FetchUserByID(c *gin.Context)
	FetchAllUsers(c *gin.Context)
	// Login(c *gin.Context)
}

type userServiceImpl struct {
	userRepo model.UserRepository
}

// NewUserService creates a new user service with the given user repository.
func NewUserService(userRepo model.UserRepository) UserService {
	return &userServiceImpl{userRepo}
}

// AddUser creates a new user.
func (service *userServiceImpl) AddUser(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		util.WriteError(c.Writer, http.StatusBadRequest, "failed to read request body")
		return
	}

	user := &model.User{}
	if err := json.Unmarshal(body, user); err != nil {
		util.WriteError(c.Writer, http.StatusBadRequest, "failed to parse request body")
		return
	}

	if err := service.userRepo.AddUser(user); err != nil {
		util.WriteError(c.Writer, http.StatusInternalServerError, fmt.Sprintf("failed to create user: %v", err))
		return
	}

	util.WriteJSON(c.Writer, http.StatusCreated, gin.H{"message": "user created successfully"})
}

// UpdateUser updates an existing user.
func (service *userServiceImpl) UpdateUser(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		util.WriteError(c.Writer, http.StatusBadRequest, "failed to read request body")
		return
	}

	user := &model.User{}
	if err := json.Unmarshal(body, user); err != nil {
		util.WriteError(c.Writer, http.StatusBadRequest, "failed to parse request body")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.WriteError(c.Writer, http.StatusBadRequest, "invalid user ID")
		return
	}

	existingUser, err := service.userRepo.GetUserByID(id)
	if err != nil {
		util.WriteError(c.Writer, http.StatusInternalServerError, fmt.Sprintf("failed to update user: %v", err))
		return
	}

	if existingUser == nil {
		util.WriteError(c.Writer, http.StatusNotFound, "user not found")
		return
	}

	existingUser.Name = user.Name
	existingUser.Email = user.Email

	if err := service.userRepo.UpdateUser(existingUser); err != nil {
		util.WriteError(c.Writer, http.StatusInternalServerError, fmt.Sprintf("failed to update user: %v", err))
		return
	}

	util.WriteJSON(c.Writer, http.StatusOK, gin.H{"message": "user updated successfully"})
}

// FetchAllUsers fetches all users.
func (service *userServiceImpl) FetchAllUsers(c *gin.Context) {
	users, err := service.userRepo.GetAllUsers()
	if err != nil {
		util.WriteError(c.Writer, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch users: %v", err))
		return
	}

	util.WriteJSON(c.Writer, http.StatusOK, users)
}

// FetchUserByID fetches a user by ID.
func (service *userServiceImpl) FetchUserByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.WriteError(c.Writer, http.StatusBadRequest, "invalid user ID")
		return
	}

	user, err := service.userRepo.GetUserByID(id)
	if err != nil {
		util.WriteError(c.Writer, http.StatusInternalServerError, fmt.Sprintf("failed to fetch user: %v", err))
		return
	}

	if user == nil {
		util.WriteError(c.Writer, http.StatusNotFound, "user not found")
		return
	}

	util.WriteJSON(c.Writer, http.StatusOK, user)
}

// DeleteUser deletes an existing user.
func (service *userServiceImpl) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.WriteError(c.Writer, http.StatusBadRequest, "invalid user ID")
		return
	}

	existingUser, err := service.userRepo.GetUserByID(id)
	if err != nil {
		util.WriteError(c.Writer, http.StatusInternalServerError, fmt.Sprintf("Ffiled to get user: %v", err))
		return
	}

	if existingUser == nil {
		util.WriteError(c.Writer, http.StatusNotFound, "User not found")
		return
	}

	err = service.userRepo.DeleteUser(id)
	if err != nil {
		util.WriteError(c.Writer, http.StatusInternalServerError, fmt.Sprintf("Failed to delete user: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// Login functanality
// func (service *userServiceImpl) Login(c *gin.Context) {
// 	body, err := ioutil.ReadAll(c.Request.Body)
// 	if err != nil {
// 		util.WriteError(c.Writer, http.StatusBadRequest, "Failed to read request body")
// 		return
// 	}

// 	user := &model.User{}
// 	if err := json.Unmarshal(body, user); err != nil {
// 		util.WriteError(c.Writer, http.StatusBadRequest, "Failed to parse request body")
// 		return
// 	}
// 	token, err := service.userRepo.Login(user)
// 	if err != nil {
// 		util.WriteError(c.Writer, http.StatusInternalServerError, fmt.Sprintf("Failed to update user: %v", err))
// 		return
// 	}

// 	util.WriteJSON(c.Writer, http.StatusOK, gin.H{
// 		"message":       "User updated successfully",
// 		"Authorization": token,
// 	})
// }
