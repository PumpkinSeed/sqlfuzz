package drivers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

const (
	DefaultTableCreateQueryKey = ""
)

func multiDescribeHelper(tables []string, processedTables map[string]struct{}, db *sql.DB, d Driver) (map[string][]FieldDescriptor, []string, error) {
	knownTables := make(map[string]bool)
	tableDescriptorMap := make(map[string][]FieldDescriptor)
	var newlyReferencedTables []string
	for _, table := range tables {
		knownTables[table] = true
	}
	for _, table := range tables {
		fields, err := d.Describe(table, db)
		if err != nil {
			return nil, nil, err
		}
		for _, field := range fields {
			if field.ForeignKeyDescriptor == nil {
				continue
			}
			foreignTableName := field.ForeignKeyDescriptor.ForeignTableName
			if _, ok := processedTables[foreignTableName]; ok && !knownTables[foreignTableName] {
				newlyReferencedTables = append(newlyReferencedTables, foreignTableName)
				knownTables[foreignTableName] = true
			}
		}
		tableDescriptorMap[table] = fields
		processedTables[table] = struct{}{}
	}
	return tableDescriptorMap, newlyReferencedTables, nil
}

func getInsertionOrder(tablesToFieldsMap map[string][]FieldDescriptor) ([]string, error) {
	var tablesVisitOrder []string
	tablesVisited := make(map[string]struct{})
	for len(tablesVisitOrder) < len(tablesToFieldsMap) {
		newInsertCount := 0
		for table, fields := range tablesToFieldsMap {
			if _, ok := tablesVisited[table]; ok {
				continue
			}
			canInsert := true
			for _, field := range fields {
				if field.ForeignKeyDescriptor == nil {
					continue
				}
				if _, ok := tablesVisited[field.ForeignKeyDescriptor.ForeignTableName]; ok {
					continue
				}
				// Necessary table is not yet visited.
				canInsert = false
				break
			}
			if canInsert {
				newInsertCount++
				tablesVisited[table] = struct{}{}
				tablesVisitOrder = append(tablesVisitOrder, table)
			}
		}
		if newInsertCount == 0 {
			return nil, errors.New("error generating insertion order. Maybe necessary dependencies are not met")
		}
	}
	return tablesVisitOrder, nil
}

func testTable(db *sql.DB, testCase, table string, d Testable) error {
	test, err := d.GetTestCase(testCase)
	if err != nil {
		return err
	}
	if test.TableCreationOrder == nil {
		if query, ok := test.TableToCreateQueryMap[DefaultTableCreateQueryKey]; ok {
			if res, err := db.ExecContext(context.Background(), fmt.Sprintf(query, table)); err != nil {
				return err
			} else if _, err := res.RowsAffected(); err != nil {
				return err
			}
		}
	} else {
		for _, table := range test.TableCreationOrder {
			createCommand := test.TableToCreateQueryMap[table]
			_, err := db.Query(strings.TrimSpace(createCommand))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
