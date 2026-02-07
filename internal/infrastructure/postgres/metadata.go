package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kamil5b/go-nl-sql/internal/domain/model"
	_ "github.com/lib/pq"
)

type MetadataExtractor struct{}

func NewMetadataExtractor() *MetadataExtractor {
	return &MetadataExtractor{}
}

func (e *MetadataExtractor) Extract(ctx context.Context, dbURL string) (*model.DatabaseMetadata, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("open connection: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	tables, err := e.extractTables(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("extract tables: %w", err)
	}

	for i := range tables {
		columns, err := e.extractColumns(ctx, db, tables[i].Name)
		if err != nil {
			return nil, fmt.Errorf("extract columns for %s: %w", tables[i].Name, err)
		}
		tables[i].Columns = columns

		indexes, err := e.extractIndexes(ctx, db, tables[i].Name)
		if err != nil {
			return nil, fmt.Errorf("extract indexes for %s: %w", tables[i].Name, err)
		}
		tables[i].Indexes = indexes

		constraints, err := e.extractConstraints(ctx, db, tables[i].Name)
		if err != nil {
			return nil, fmt.Errorf("extract constraints for %s: %w", tables[i].Name, err)
		}
		tables[i].Constraints = constraints
	}

	return &model.DatabaseMetadata{
		Tables: tables,
	}, nil
}

func (e *MetadataExtractor) extractTables(ctx context.Context, db *sql.DB) ([]model.Table, error) {
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []model.Table
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, model.Table{Name: tableName})
	}

	return tables, rows.Err()
}

func (e *MetadataExtractor) extractColumns(ctx context.Context, db *sql.DB, tableName string) ([]model.Column, error) {
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable,
			column_default,
			(SELECT EXISTS (
				SELECT 1 FROM information_schema.table_constraints tc
				JOIN information_schema.key_column_usage kcu ON tc.constraint_name = kcu.constraint_name
				WHERE tc.table_name = $1 AND kcu.column_name = c.column_name
				AND tc.constraint_type = 'PRIMARY KEY'
			)) as is_primary_key
		FROM information_schema.columns c
		WHERE table_name = $1
		ORDER BY ordinal_position
	`

	rows, err := db.QueryContext(ctx, query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []model.Column
	for rows.Next() {
		var col model.Column
		var nullable string
		var defaultVal sql.NullString

		if err := rows.Scan(&col.Name, &col.Type, &nullable, &defaultVal, &col.IsPrimaryKey); err != nil {
			return nil, err
		}

		col.Nullable = nullable == "YES"
		if defaultVal.Valid {
			col.Default = defaultVal.String
		}

		columns = append(columns, col)
	}

	return columns, rows.Err()
}

func (e *MetadataExtractor) extractIndexes(ctx context.Context, db *sql.DB, tableName string) ([]model.Index, error) {
	query := `
		SELECT indexname
		FROM pg_indexes
		WHERE tablename = $1 AND schemaname = 'public'
	`

	rows, err := db.QueryContext(ctx, query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []model.Index
	for rows.Next() {
		var indexName string
		if err := rows.Scan(&indexName); err != nil {
			return nil, err
		}
		indexes = append(indexes, model.Index{Name: indexName})
	}

	return indexes, rows.Err()
}

func (e *MetadataExtractor) extractConstraints(ctx context.Context, db *sql.DB, tableName string) ([]model.Constraint, error) {
	query := `
		SELECT constraint_name, constraint_type
		FROM information_schema.table_constraints
		WHERE table_name = $1 AND table_schema = 'public'
	`

	rows, err := db.QueryContext(ctx, query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var constraints []model.Constraint
	for rows.Next() {
		var name, ctype string
		if err := rows.Scan(&name, &ctype); err != nil {
			return nil, err
		}
		constraints = append(constraints, model.Constraint{
			Name: name,
			Type: ctype,
		})
	}

	return constraints, rows.Err()
}
