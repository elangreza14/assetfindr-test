package dto

// implement this https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/

type CreateOrUpdatePostRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Tags    []string `json:"tags" binding:"required,gt=0,dive,required"`
}

type UriPostRequest struct {
	ID int `uri:"id" binding:"required,gt=0"`
}

type GetPostResponse struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}
