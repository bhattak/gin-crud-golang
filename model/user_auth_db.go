// model/user_db.go
package model

import (
	"database/sql"
	"fmt"

	"project/util"
)

// UserRepository provides an interface for interacting with the user database.
type UserAuthRepository interface {
	Login(*LoginRequest) (string, error)
	Register(*LoginRequest) error
}

type userAuthRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository with the given database connection.
func NewUserAuthRepository(db *sql.DB) UserAuthRepository {
	return &userAuthRepositoryImpl{db: db}
}

// Does login and generates jwt token
func (repo *userAuthRepositoryImpl) Login(user *LoginRequest) (string, error) {

	query := "select * from super_user where email=? "
	row := repo.db.QueryRow(query, user.Email)

	tempUser := &LoginRequest{}
	err := row.Scan(&tempUser.ID, &tempUser.Email, &tempUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get user: %v", err)
	}
	fmt.Println(tempUser.Password)
	fmt.Println(user.Password)

	match := util.ComparePasswords(tempUser.Password, user.Password)

	if !match {
		return "", fmt.Errorf("password did not match")
	}

	jwt, _ := util.GenerateJWT(tempUser.Email)
	fmt.Println("jwt", jwt)
	return jwt, nil
}

// Does login and generates jwt token
func (repo *userAuthRepositoryImpl) Register(user *LoginRequest) error {

	query := "select * from super_user where email=? "
	rows, err := repo.db.Query(query, user.Email)
	if err != nil {
		return fmt.Errorf("failed to check user: %v", err)
	}
	if rows.Next() {
		return fmt.Errorf("user already exists")
	}
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	fmt.Println(user.Email)
	fmt.Println(hashedPassword)
	result, err := repo.db.Exec("INSERT INTO super_user(email, password) VALUES (?, ?)", user.Email, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to add user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %v", err)
	}
	fmt.Println(id)
	return nil
}
