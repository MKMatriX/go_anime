package requests

type IdParamRequest struct {
	ID uint `path:"id" param:"id" binding:"required"`
}
