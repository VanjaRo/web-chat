package repositories

import (
	"github.com/VanjaRo/web-chat/interfaces"
	"gorm.io/gorm"
)

type Room struct {
	ID      string `gorm:"primaryKey; not null" json:"id"`
	Name    string `gorm:"not null" json:"name"`
	Private bool   `json:"private"`
}

func (room *Room) GetId() string {
	return room.ID
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) GetPrivate() bool {
	return room.Private
}

type RoomRepository struct {
	Db *gorm.DB
}

func (repo *RoomRepository) AddRoom(room interfaces.Room) error {
	return repo.Db.Create(&Room{room.GetId(), room.GetName(), room.GetPrivate()}).Error
}

func (repo *RoomRepository) FindRoomByName(name string) (interfaces.Room, error) {
	var room Room
	if err := repo.Db.First(&room, "name = ?", name).Error; err != nil {
		return &room, err
	}
	return &room, nil
}
