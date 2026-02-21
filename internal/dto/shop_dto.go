package dto

type ShopUpdatePayload struct {
	Name           string `json:"name" binding:"required,max=120"`
	WhatsAppNumber string `json:"whatsAppNumber" binding:"omitempty,max=20"`
}
