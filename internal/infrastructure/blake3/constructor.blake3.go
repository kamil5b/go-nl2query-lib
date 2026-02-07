package hash

type Blake3HashGenerator struct{}

func New() *Blake3HashGenerator {
	return &Blake3HashGenerator{}
}
