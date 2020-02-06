package sns

// UploadPayload describes the shape of relevant to upload functionality message.Body contents of sns message
type UploadPayload struct {
	Message string `json:"Message"`
}
