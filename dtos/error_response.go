package dtos

type ValidationErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrResponse struct {
	Status  int                  `json:"status,omitempty"`
	Message string               `json:"message"`
	Data    *any                 `json:"data"`
	Details []ValidationErrorMsg `json:"details,omitempty"`
}
