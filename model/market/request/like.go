package request

type CreateCommentLike struct {
	CommentId uint `json:"comment_id" form:"comment_id" binding:"required"`
}

type UpdateCommentLike struct {
	IsLike    bool `json:"is_like" form:"is_like" binding:"required"`
	CommentId uint `json:"comment_id" form:"comment_id" binding:"required"`
}
