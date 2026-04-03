package dto

type Repository struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       uint32 `json:"stars"`
	Forks       uint32 `json:"forks"`
	CreatedAt   string `json:"created_at"`
}
