package user

type User struct {
	ID       uint32 `json:"ID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
