package model

type CreateUserRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required"`
	Country string `json:"country" binding:"required"`
}

type CreateUserResponse struct {
	Status string `json:"status" binding:"required"`
}

type User struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Country   string `json:"country" binding:"required"`
	UserId    string `json:"userId" binding:"required"`
	CreatedAt string `json:"createdAt" binding:"required"`
	UpdatedAt string `json:"updatedAt" binding:"required"`
}
