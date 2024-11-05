package dtos

import "github.com/Eggi19/simple-social-media/utils"

type ResponseMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Login struct {
	Token *utils.JwtToken  `json:"token"`
}
