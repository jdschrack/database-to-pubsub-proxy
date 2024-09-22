package models

type SqlCommandEvent struct {
	SqlCommand string `json:"sql_command"`
	Database   string `json:"database"`
	Timestamp  string `json:"timestamp"`
	User       string `json:"user"`
}
