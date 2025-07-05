package parser

import (
	"database/sql"
	"fmt"
	"tables/pkg/schema"
)

type SchemaParser struct {
	db *sql.DB
}

func NewSchemaParser(db *sql.DB) *SchemaParser {
	return &SchemaParser{db: db}
}

func (si *SchemaParser) GetTables() ([]schema.Table, error) {
	query := `
		SELECT 
			t.table_name,
			c.column_name,
			c.data_type,
			c.is_nullable
		FROM 
			information_schema.tables t
		JOIN 
			information_schema.columns c ON t.table_name = c.table_name
		WHERE 
			t.table_schema = 'public'
			AND t.table_type = 'BASE TABLE'
		ORDER BY 
			t.table_name, c.ordinal_position
	`

	rows, err := si.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query schema: %w", err)
	}
	defer rows.Close()

	tablesMap := make(map[string]*schema.Table)
	var tables []schema.Table

	for rows.Next() {
		var tableName, columnName, dataType, nullable string

		if err := rows.Scan(&tableName, &columnName, &dataType, &nullable); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Get or create table
		table, exists := tablesMap[tableName]
		if !exists {
			table = &schema.Table{Name: tableName, Columns: []schema.Column{}}
			tablesMap[tableName] = table
		}

		// Add column to table
		column := schema.Column{
			Name:     columnName,
			Type:     dataType,
			Nullable: nullable == "YES",
		}
		table.Columns = append(table.Columns, column)
	}

	// Convert map to slice
	for _, table := range tablesMap {
		tables = append(tables, *table)
	}

	return tables, nil
}
