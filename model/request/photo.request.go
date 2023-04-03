package request

type PhotoCreateRequest struct {
	CategoryID uint `json:"category_id" form:"category_id" validate:"required"`
	// Image      []byte `form:"image" validate:"required"`
}
