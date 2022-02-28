package common_http

type Response struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}
