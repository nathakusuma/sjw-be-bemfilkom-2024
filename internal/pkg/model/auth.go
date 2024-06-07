package model

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token          string `json:"token"`
	NIM            string `json:"nim"`
	Email          string `json:"email"`
	FullName       string `json:"full_name"`
	Role           string `json:"role"`
	Angkatan       string `json:"angkatan"`
	ProgramStudi   string `json:"program_studi"`
	ProfilePicture string `json:"profile_picture"`
}
