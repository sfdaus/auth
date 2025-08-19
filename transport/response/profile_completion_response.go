package response

type ProfileCompletionResponse struct {
	BasicResponse
	Data ProfileCompletionResponseData `json:"data,omitempty"`
}

type ProfileCompletionResponseData struct {
	CompletionStatus bool `json:"completion_status"`
}
