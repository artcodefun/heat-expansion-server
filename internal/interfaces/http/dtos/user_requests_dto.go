package dtos

// userRegisterBody represents the JSON payload for user registration.
type userRegisterBody struct {
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserRegisterRequest = Request[None, None, userRegisterBody]

// userLoginBody represents the JSON payload for user login.
type userLoginBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserLoginRequest = Request[None, None, userLoginBody]

// LoginResponse represents the response for successful user login.
type LoginResponse struct {
	Token string `json:"token"`
}

// UserCrystalBalanceResponse represents the current crystal balance for the authenticated user.
type UserCrystalBalanceResponse struct {
	Crystals int `json:"crystals"`
}
