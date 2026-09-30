package model

// Post is the document stored in Elasticsearch for each uploaded image or video.
type Post struct {
	ID      string `json:"id"`
	User    string `json:"user"`
	Message string `json:"message"`
	URL     string `json:"url"`
	Type    string `json:"type"`
}

// User is the account document stored in Elasticsearch in this course project.
// A production system should store a password hash rather than a plain password.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Age      int64  `json:"age"`
	Gender   string `json:"gender"`
}
