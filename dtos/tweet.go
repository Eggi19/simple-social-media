package dtos

type CreateTweetRequest struct {
	Tweet string `json:"tweet" binding:"required"`
}

type TweetLikesResponse struct {
	Likes int64 `json:"likes"`
}

type TweetLikesRequest struct {
	TweetId int64 `json:"tweet_id" binding:"required"`
}
