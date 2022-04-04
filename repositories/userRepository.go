package repositories

import (
	"github.com/VanjaRo/web-chat/interfaces"
	"gorm.io/gorm"
)

// type FriendRequest struct {
// 	gorm.Model
// 	FromUser *User
// 	ToUser   *User
// }

type User struct {
	ID       string `gorm:"primaryKey; not null" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	Password string `json:"password"`
	// Friends         []*User          `gorm:"many2many:friends;" json:"friends"`
	// FRSended []*FriendRequest
}

func (user *User) GetId() string {
	return user.ID
}

func (user *User) GetName() string {
	return user.Name
}

func (user *User) GetPassword() string {
	return user.Password
}

type UserRepository struct {
	Db *gorm.DB
}

func (repo *UserRepository) AddUser(user interfaces.UserAuth) error {
	var newUser User
	newUser.ID = user.GetId()
	newUser.Name = user.GetName()
	newUser.Password = user.GetPassword()
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

func (repo *UserRepository) FindUserByName(username string) (*User, error) {
	var user User
	if err := repo.Db.First(&user, "name = ?", username).Error; err != nil {
		return &user, err
	}
	return &user, nil
}

// func (repo *UserRepository) GetAllFriends

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
