package domains

type DatabaseMetadata struct {
	TenantID  string // Hashed Client's DB URL
	Tables    []Table
	Relations []Relation
	Checksum  string
}

type Table struct {
	Name        string
	Columns     []Column
	Indexes     []Index
	Constraints []Constraint
	Comments    string
}

type Column struct {
	Name         string
	Type         string
	Nullable     bool
	Default      string
	IsPrimaryKey bool
	IsForeignKey bool
	Comments     string
}

type Index struct {
	Name    string
	Columns []string
	Unique  bool
}

type Constraint struct {
	Name      string
	Type      string
	Columns   []string
	Reference string
}

type Relation struct {
	SourceTable  string
	SourceColumn string
	TargetTable  string
	TargetColumn string
	RelationType string
}
