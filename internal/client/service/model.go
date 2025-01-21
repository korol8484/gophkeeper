package service

type auth struct {
	token   string
	passKey string
}

func newAuthState(token, passKey string) *auth {
	return &auth{
		token:   token,
		passKey: passKey,
	}
}

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
