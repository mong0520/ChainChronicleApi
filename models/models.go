package models

type GeneralResponse struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Error   string      `json:"error_code"`
}
