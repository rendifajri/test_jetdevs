package structs

type Article struct {
	Id        string `json:"id,omitempty" db:"id"`
	Nickname  string `json:"nickname,omitempty" db:"nickname"`
	Title     string `json:"title,omitempty" db:"title"`
	Content   string `json:"content,omitempty" db:"content"`
	CreatedOn string `json:"created_on,omitempty" db:"created_on"`
}

type Comment struct {
	Id        string `json:"id,omitempty" db:"id"`
	ArticleId string `json:"article_id,omitempty" db:"article_id"`
	CommentId string `json:"comment_id,omitempty" db:"comment_id"`
	Nickname  string `json:"nickname,omitempty" db:"nickname"`
	Content   string `json:"content,omitempty" db:"content"`
	CreatedOn string `json:"created_on,omitempty" db:"created_on"`
}
