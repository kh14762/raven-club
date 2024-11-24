package types

type User struct {
	Id       uint16 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    Token  `json:"token"`
}

type Token struct {
	AccessToken  string `json:"-"`
	RefreshToken string `json:"-"`
	ExpiresAt    int64  `json:"-"`
	Type         string `json:"-"`
}
