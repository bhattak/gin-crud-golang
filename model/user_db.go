// model/user_db.go
package model

import (
	"database/sql"
	"fmt"

	"project/util"
)

// UserRepository provides an interface for interacting with the user database.
type UserRepository interface {
	AddUser(*User) error
	UpdateUser(*User) error
	DeleteUser(int64) error
	GetUserByID(int64) (*User, error)
	GetAllUsers() ([]*User, error)
	// Login(*User) (string, error)
}

type userRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository with the given database connection.
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

// AddUser adds a new user to the database.
func (repo *userRepositoryImpl) AddUser(user *User) error {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	result, err := repo.db.Exec("INSERT INTO users(name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to add user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %v", err)
	}

	user.ID = id
	user.Password = ""
	return nil
}

// model/user_db.go (continued)
// UpdateUser updates an existing user in the database.
func (repo *userRepositoryImpl) UpdateUser(user *User) error {
	_, err := repo.db.Exec("UPDATE users SET name=?, email=? WHERE id=?", user.Name, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

// DeleteUser deletes an existing user from the database.
func (repo *userRepositoryImpl) DeleteUser(id int64) error {
	_, err := repo.db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}

// GetUserByID retrieves a user with the given ID from the database.
func (repo *userRepositoryImpl) GetUserByID(id int64) (*User, error) {
	row := repo.db.QueryRow("SELECT * FROM users WHERE id=?", id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	user.Password = ""
	return user, nil
}

// GetAllUsers retrieves all users from the database.
func (repo *userRepositoryImpl) GetAllUsers() ([]*User, error) {
	rows, err := repo.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to get users: %v", err)
		}
		user.Password = ""
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Failed to get users: %v", err)
	}

	return users, nil
}

// ADoes login and generates jwt token
// func (repo *userRepositoryImpl) Login(user *User) (string, error) {

// 	// var exists bool
// 	// query := "select exists(select * from users where email=? )"
// 	// error := repo.db.QueryRow(query, user.Email).Scan(&exists)
// 	// fmt.Println(exists)
// 	// if error != nil {
// 	// 	return "", fmt.Errorf("eeor querying db: %v", err)
// 	// }
// 	// if !exists {
// 	// 	return "", fmt.Errorf("user does not exist")

// 	// }
// 	query := "select * from users where email=? "
// 	row := repo.db.QueryRow(query, user.Email)

// 	tempUser := &User{}
// 	err := row.Scan(&tempUser.ID, &tempUser.Name, &tempUser.Email, &tempUser.Password)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return "", nil
// 		}
// 		return "", fmt.Errorf("failed to get user: %v", err)
// 	}
// 	fmt.Println(tempUser.Password)
// 	fmt.Println(user.Password)

// 	match := util.ComparePasswords(tempUser.Password, user.Password)

// 	if !match {
// 		return "", fmt.Errorf("password did not match")
// 	}

// 	jwt, _ := util.GenerateJWT(tempUser.Email)

// 	return jwt, nil
// }
