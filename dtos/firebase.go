package dtos

type SubscribeTopicRequest struct {
	FcmToken string `json:"fcm_token" binding:"required"`
	Topic    string `json:"topic" binding:"required"`
}

type UnsubscribeTopicRequest struct {
	FcmToken string `json:"fcm_token" binding:"required"`
	Topic    string `json:"topic" binding:"required"`
}
