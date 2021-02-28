package descriptor

import (
	"database/sql"
)

// ShowTables queries the available tables of the database
func ShowTables(db *sql.DB) ([]string, error) {
	results, err := db.Query("SHOW TABLES;")
	if err != nil {
		return nil, err
	}
	defer results.Close()
	var tables []string
	for results.Next() {
		var table string
		if err := results.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}
