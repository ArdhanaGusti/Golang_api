package validation

type RegisterUserPayload struct {
	Username string `json:"Username" form:"Username" binding:"required"`
	Fullname string `json:"Fullname" form:"Fullname" binding:"required"`
	Email    string `json:"Email" form:"Email" binding:"required,email"`
	Password string `json:"Password" form:"Password" binding:"required"`
}
