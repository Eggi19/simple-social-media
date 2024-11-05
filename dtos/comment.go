package dtos

type CreateCommentRequest struct {
	Comment string `json:"comment" binding:"required"`
	TweetId int64  `json:"tweet_id" binding:"required"`
}
