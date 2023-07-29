package ginserver

type getUserRequest struct {
	Password string `json:"password"`
}

type postUserRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}
