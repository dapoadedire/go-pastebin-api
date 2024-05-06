package model


type Paste struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}


type CreatePasteRequest struct {
	Content string `json:"content" binding:"required"`
}
