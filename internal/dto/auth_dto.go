package dto

import "time"

// RegisterRequest correspond a la creation d'utilisateur par un SuperAdmin.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=120"`
	Email    string `json:"email" binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Role     string `json:"role" binding:"required,oneof=SuperAdmin Admin"`
	ShopID   string `json:"shopID" binding:"required,uuid"`
}

// LoginRequest inclut shopID pour lever l'ambiguite email multi-shop.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	ShopID   string `json:"shopID" binding:"required,uuid"`
}

// AuthUserResponse est la representation utilisateur renvoyee au client.
type AuthUserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	ShopID    string    `json:"shopID"`
	CreatedAt time.Time `json:"createdAt"`
}

// LoginResponse expose le token JWT et les informations utilisateur.
type LoginResponse struct {
	Token string           `json:"token"`
	User  AuthUserResponse `json:"user"`
}
