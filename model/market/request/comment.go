package request

type CreateEventComment struct {
	Content string `json:"content" form:"content" binding:"required"`
	ReplyId uint   `json:"reply_id" form:"reply_id"`
	Code    string `json:"code" form:"code" binding:"required"`
}

type GetEventComment struct {
	Code string `json:"code" form:"code" binding:"required"`
}

type RemoveEventComment struct {
	CommentId uint `json:"comment_id" form:"comment_id" binding:"required"`
}
