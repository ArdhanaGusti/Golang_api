package validation

type CreateArticlePayload struct {
	Title string `json:"Title" form:"Title" binding:"required"`
	Desc  string `json:"Desc" form:"Desc" binding:"required"`
	Tag   string `json:"Tag" form:"Tag" binding:"required"`
}
