package models

// define the Todo model
type Todo struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Status  bool   `json:"status"`
	Deleted bool   `json:"deleted"`
}
