package models

import (
	"database/sql"
	"time"

	"github.com/aZ4ziL/blogs_api/auth"
)

// User is implement for user model
type User struct {
	ID         int          `gorm:"primaryKey" json:"id"`
	FirstName  string       `gorm:"size:50" json:"first_name"`
	LastName   string       `gorm:"size:50" json:"last_name"`
	Username   string       `gorm:"size:30;unique" json:"username"`
	Email      string       `gorm:"size:50;unique" json:"email"`
	Password   string       `gorm:"size:128" json:"-"`
	IsAdmin    bool         `gorm:"default:false" json:"is_admin"`
	IsActive   bool         `gorm:"default:true" json:"is_active"`
	LastLogin  sql.NullTime `gorm:"null" json:"last_login"`
	DateJoined time.Time    `gorm:"autoCreateTime" json:"date_joined"`
	Articles   []*Article   `gorm:"foreignKey:UserID" json:"articles,omitempty"`
}

// CreateNewUser is function to create new user.
func CreateNewUser(user *User) error {
	user.Password = auth.EncryptionPassword(user.Password)
	return db.Create(user).Error
}

// GetUserByUsername is function to get user by username.
func GetUserByUsername(username string) (User, error) {
	var user User
	err := db.Model(&User{}).Where("username = ?", username).Preload("Articles").First(&user).Error
	return user, err
}

// GetUserByID is function to get user by id.
func GetUserByID(id int) (User, error) {
	var user User
	err := db.Model(&User{}).Where("id = ?", id).Preload("Articles").First(&user).Error
	return user, err
}
