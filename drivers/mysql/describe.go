package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/PumpkinSeed/sqlfuzz/drivers/types"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

const (
	asd = `SELECT CONSTRAINT_NAME,TABLE_NAME,COLUMN_NAME,REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME 
							   from INFORMATION_SCHEMA.KEY_COLUMN_USAGE 
                               where REFERENCED_TABLE_NAME <> 'NULL' and REFERENCED_COLUMN_NAME <> 'NULL' and TABLE_NAME = '%s'`
)

func (MySQL) Describe(table string, db *sql.DB) ([]types.FieldDescriptor, error) {
	//goqu.New("mysql", db).From("INFORMATION_SCHEMA.KEY_COLUMN_USAGE").
	//	Select("CONSTRAINT_CATALOG", "CONSTRAINT_SCHEMA","CONSTRAINT_NAME", "TABLE_CATALOG", "TABLE_SCHEMA", "TABLE_NAME", )
	describeQuery := fmt.Sprintf("DESCRIBE %s;", table)
	results, err := db.Query(describeQuery)
	if err != nil {
		return nil, err
	}
	fkRows, err := db.Query(fmt.Sprintf(mysqlFKQuery, strings.ToLower(table)))
	if err != nil {
		return nil, err
	}
	return parseMySQLFields(results, fkRows)
}
