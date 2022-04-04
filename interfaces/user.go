package interfaces

type User interface {
	GetId() string
	GetName() string
}

type UserAuth interface {
	User
	GetPassword() string
}

type UserRepository interface {
	AddUser(userAuth UserAuth) error
	RemoveUser(user User) error
	FindUserById(id string) (User, error)
	GetAllUsers() ([]User, error)
}
