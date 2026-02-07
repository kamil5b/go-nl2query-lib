package repository

type SQLValidatorRepository interface {
	IsSafe(sql string) (bool, error)
	ContainsDDLDML(sql string) bool
}
