package interfaces

type Room interface {
	GetId() string
	GetName() string
	GetPrivate() bool
}

type RoomRepository interface {
	AddRoom(room Room) error
	FindRoomByName(name string) (Room, error)
}
