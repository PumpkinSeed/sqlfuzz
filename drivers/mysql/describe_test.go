package mysql

import (
	"database/sql"
	"testing"
)

func TestDescribe(t *testing.T) {
	db, err := sql.Open("mysql", "test:test@tcp(localhost:3306)/test")
	if err != nil {
		t.Error(err)
	}
	MySQL{}.Describe("t_currency", db)
}
