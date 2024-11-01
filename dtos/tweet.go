package dtos

type CreateTweetRequest struct {
	Tweet string `json:"tweet" binding:"required"`
}
