package team

type Contact struct {
	Email string `json:"email"`
	Slack string `json:"slack"`
	Phone string `json:"phone"`
}

type Team struct {
	TeamId  int     `json:"teamId"`
	Slug    string  `json:"slug"`
	Name    string  `json:"name"`
	Contact Contact `json:"contact"`
	Hours   []int   `json:"hours"`
}
