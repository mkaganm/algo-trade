package domain

type Message struct {
	ID        string
	Values    map[string]interface{}
	Processed bool
	Timestamp string
}
