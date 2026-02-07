package model

type Vector struct {
	ID        string
	Embedding []float32
	Metadata  map[string]string
	Content   string
}
