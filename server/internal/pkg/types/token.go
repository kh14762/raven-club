package types

type Token struct {
	AccessToken      string `json:"-"`
	RefreshToken     string `json:"-"`
	AccessExpiresAt  int64  `json:"-"`
	RefreshExpiresAt int64  `json:"-"`
	Type             string `json:"-"`
}
