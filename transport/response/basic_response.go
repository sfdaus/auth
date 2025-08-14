package response

type BasicResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    any         `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}
