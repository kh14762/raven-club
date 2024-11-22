package user

import "go.uber.org/fx"

var Module = fx.Module("UserService", fx.Provide(fx.Annotate(MakeUserList(), fx.As(new(UserService)))))

type UserService interface {
	List() []User
	Get(id uint16) (User, bool)
	Add(User)
	Update(User) (User, error)
	Delete(id uint16)
}

type UserList struct {
	Id    uint16 `json:"id"`
	Users []User `json:"users"`
}

type User struct {
	Id         uint16              `json:"id"`
	Username   string              `json:"username"`
	Email      string              `json:"email"`
	Password   string              `json:"password"` // TODO: make this unreadable json
	Properties map[string]Property `json:"properties"`
}

type Property interface{}

func (u *UserList) List() []User { return u.Users }

func (u *UserList) Get(id uint16) (User, bool) {
	for _, user := range u.Users {
		if user.Id == id {
			return user, true
		}
	}
	return User{}, false
}

func (u *UserList) Add(user User) {
	u.Id++
	user.Id = u.Id
	if user.Properties == nil {
		user.Properties = make(map[string]Property)
	}

	u.Users = append(u.Users, user)
}

func MakeUserList() *UserList {
	return &UserList{
		Id:    0,
		Users: make([]User, 0),
	}
}
