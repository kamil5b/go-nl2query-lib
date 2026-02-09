package ports

type QueryValidatorPort interface {
	IsSafe(query string) (bool, error)
	ContainsDDLDML(query string) bool
}
