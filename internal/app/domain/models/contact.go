package models

type Contact struct {
	ID        string `json:"_id"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
	Telephone Phone  `json:"telephone"`
}

type Phone struct {
	Mobile string `json:"mobile"`
	Home   string `json:"home"`
}
