package requests

type AnimeCreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// type AnimeUpdateRequest struct {
// 	ID          uint   `json:"id" validate:"required"`
// 	Name        string `json:"name" validate:"required"`
// 	Description string `json:"description"`
// }

// type AnimeDeleteRequest struct {
// 	ID uint `json:"id" validate:"required"`
// }
