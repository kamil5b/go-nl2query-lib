package domains

type Vector struct {
	ID        string
	TenantID  string
	Embedding []float32
	Metadata  map[string]string
	Content   string
}
