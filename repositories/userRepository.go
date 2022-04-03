package repositories

import (
	"github.com/VanjaRo/web-chat/models"
	"gorm.io/gorm"
)

type User struct {
	ID   string `gorm:"primaryKey; not null" json:"id"`
	Name string `gorm:"not null" json:"name"`
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

func (repo *UserRepository) AddUser(user models.User) error {
	return repo.Db.Create(&User{user.GetId(), user.GetName()}).Error
}

func (repo *UserRepository) RemoveUser(user models.User) error {
	return repo.Db.Delete(&User{}, user.GetId()).Error
}

func (repo *UserRepository) FindUserById(id string) (models.User, error) {
	var user User
	if err := repo.Db.First(&user, id).Error; err != nil {
		return &user, err
	}
	return &user, nil
}

func (repo *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []User
	if err := repo.Db.Find(&users).Error; err != nil {
		return nil, err
	}
	// convert user struct to user interface
	usersList := make([]models.User, len(users))
	for i := range usersList {
		usersList[i] = &users[i]

	}
	return usersList, nil
}
