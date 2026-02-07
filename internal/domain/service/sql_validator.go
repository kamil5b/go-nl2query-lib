package service

type SQLValidator interface {
	IsSafe(sql string) (bool, error)
	ContainsDDLDML(sql string) bool
}
