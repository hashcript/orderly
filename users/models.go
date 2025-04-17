package users

import (
	"errors"

	"com.orderly/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	ID       uint   `gorm:"primary_key"`
	Name     string `gorm:"size:2048"`
	Phone    string `gorm:"size:2048"`
	Role     string `gorm:"size:2048"`
	Status   string `gorm:"size:2048"`
	Location string `gorm:"size:2048"`
	Password string `gorm:"column:Password;not null"`
	Email    string `gorm:"column:email;unique_index"`
}

// Migrate the schema to the database if needed
func AutoMigrate() {
	db := database.GetConnection()
	db.AutoMigrate(&UserModel{})
}

// FindSingleUser finds a single User based on the provided condition
func FindSingleUser(condition interface{}) (UserModel, error) {
	db := database.GetConnection()
	var model UserModel
	err := db.Where(condition).First(&model).Error
	return model, err
}

// SaveSingleUser saves a single User to the database
func SaveSingleUser(data interface{}) error {
	db := database.GetConnection()
	err := db.Save(data).Error
	return err
}

// UpdateSingleUser updates a User with new data
func UpdateSingleUser(model *UserModel, data interface{}) error {
	db := database.GetConnection()
	err := db.Model(model).Updates(data).Error
	return err
}

// DeleteSingleUser deletes a User from the database
func DeleteSingleUser(model *UserModel) error {
	db := database.GetConnection()
	err := db.Delete(model).Error
	return err
}

// GetAllUsers gets all Users from the database
func GetAllUsers() ([]UserModel, error) {
	db := database.GetConnection()
	var models []UserModel
	err := db.Find(&models).Error
	return models, err
}

// fix codwa
func GetUserByID(id uint) (UserModel, error) {
	db := database.GetConnection()
	var User UserModel
	err := db.First(&User, id).Error
	return User, err
}

// get user by  email
func GetUserByEmail(email string) (UserModel, error) {
	db := database.GetConnection()
	var user UserModel
	err := db.Where("Email = ?", email).First(&user).Error
	return user, err
}

func (u *UserModel) setPassword(Password string) error {
	//check Password legnth
	if len(Password) == 0 {
		return errors.New("Password should never be empty")
	}
	bytePassword := []byte(Password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	PasswordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.Password = string(PasswordHash)
	return nil
}

func (u *UserModel) checkPassword(Password string) error {
	bytePassword := []byte(Password)
	byteHashedPassword := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}
