package dtos

type UpdateEmployeePositionInput struct {
	PositionID     int      `json:"position_id" binding:"required"`
	JoiningAt      string   `json:"joining_at"`
	Description    string   `json:"description"`
	AttachmentUrls []string `json:"attachment_urls"`
}
