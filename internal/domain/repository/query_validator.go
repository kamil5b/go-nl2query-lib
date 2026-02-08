package repository

type QueryValidatorRepository interface {
	IsSafe(query string) (bool, error)
	ContainsDDLDML(query string) bool
}
