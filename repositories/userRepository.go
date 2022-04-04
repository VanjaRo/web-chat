package repositories

import (
	"github.com/VanjaRo/web-chat/interfaces"
	"gorm.io/gorm"
)

type User struct {
	ID       string `gorm:"primaryKey; not null" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (user *User) GetId() string {
	return user.ID
}

func (user *User) GetName() string {
	return user.Name
}

type UserRepository struct {
	Db *gorm.DB
}

func (repo *UserRepository) AddUser(user interfaces.User) error {
	var newUser User
	newUser.ID = user.GetId()
	newUser.Name = user.GetName()
	return repo.Db.Create(&newUser).Error
}

func (repo *UserRepository) RemoveUser(user interfaces.User) error {
	return repo.Db.Delete(&User{}, user.GetId()).Error
}

func (repo *UserRepository) FindUserById(id string) (interfaces.User, error) {
	var user User
	if err := repo.Db.First(&user, id).Error; err != nil {
		return &user, err
	}
	return &user, nil
}

func (repo *UserRepository) FindUserByUsername(username string) (*User, error) {
	var user User
	if err := repo.Db.First(&user, "username = ?", username).Error; err != nil {
		return &user, err
	}
	return &user, nil
}

func (repo *UserRepository) GetAllUsers() ([]interfaces.User, error) {
	var users []User
	if err := repo.Db.Find(&users).Error; err != nil {
		return nil, err
	}
	// convert user struct to user interface
	usersList := make([]interfaces.User, len(users))
	for i := range usersList {
		usersList[i] = &users[i]

	}
	return usersList, nil
}
