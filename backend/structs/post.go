package structs

type Post struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Text_Content string `json:"text_content"`
	User         User   `json:"user"`
}
