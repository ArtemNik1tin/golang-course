package dto

type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       uint32 `json:"stars"`
	Forks       uint32 `json:"forks"`
	CreatedAt   string `json:"created_at"`
}
