package action

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/PumpkinSeed/sqlfuzz/drivers"
	"github.com/brianvoe/gofakeit/v5"
	_ "github.com/lib/pq"
	"github.com/rs/xid"
)

func InsertMulti(args ...interface{}) error {
	db := args[0].(*sql.DB)
	driver := args[1].(drivers.Driver)
	tableToFieldsMap := args[2].(map[string][]drivers.FieldDescriptor)
	insertionOrder := args[3].([]string)
	//func InsertMulti(db *sql.DB, driver drivers.Driver, tableToFieldsMap map[string][]drivers.FieldDescriptor, insertionOrder []string) error {
	tableFieldValuesMap := make(map[string]map[string]interface{})
	for _, table := range insertionOrder {
		if fields, ok := tableToFieldsMap[table]; ok {
			var f = make([]string, 0, len(fields))
			var values []interface{}
			for _, field := range fields {
				f = append(f, field.Field)
				if field.HasDefaultValue {
					continue
				}
				var data interface{}
				if field.ForeignKeyDescriptor != nil {
					if foreignTableFields, ok := tableFieldValuesMap[field.ForeignKeyDescriptor.ForeignTableName]; ok {
						if val, ok := foreignTableFields[field.ForeignKeyDescriptor.ForeignColumnName]; ok {
							data = val
							continue
						}
					}
					val, err := getLatestColumnValue(field.ForeignKeyDescriptor.ForeignTableName, field.ForeignKeyDescriptor.ForeignColumnName, db)
					if err != nil {
						return err
					}
					data = val
					// Get from table. If no value present in table as well, throw error.
				} else {
					data = generateData(driver, field)
				}
				values = append(values, data)
			}
			query := driver.Insert(f, table)
			_, err := db.Exec(query, values...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getLatestColumnValue(table, column string, db *sql.DB) (interface{}, error) {
	query := fmt.Sprintf("select %v from %v order by %v desc limit 1", column, table, column)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	var val interface{}
	for rows.Next() {
		rows.Scan(&val)
	}
	return val, nil
}

// Insert is inserting a random generated data into the chosen table
func Insert(args ...interface{}) error {
	//func Insert(db *sql.DB, fields []drivers.FieldDescriptor, driver drivers.Driver, table string) error {
	db := args[0].(*sql.DB)
	fields := args[1].([]drivers.FieldDescriptor)
	driver := args[2].(drivers.Driver)
	table := args[3].(string)
	var f = make([]string, 0, len(fields))
	var values = make([]interface{}, 0, len(fields))
	for _, field := range fields {
		// Has default value. No need to insert this field manually.
		if field.HasDefaultValue {
			continue
		}
		f = append(f, field.Field)
		values = append(values, generateData(driver, field))
	}
	query := driver.Insert(f, table)

	_, err := db.Exec(query, values...)
	return err
}

// generateData generates random data based on the field
func generateData(driver drivers.Driver, fieldDescriptor drivers.FieldDescriptor) interface{} {
	field := driver.MapField(fieldDescriptor)
	switch field.Type {
	case drivers.String:
		if field.Length > 19 {
			return xid.New().String()
		}
		if field.Length > 0 {
			return randomString(field.Length)
		}
		return randomString(20)
	case drivers.Int16:
		return gofakeit.Number(1, 32766)
	case drivers.Int32:
		return gofakeit.Number(1, 2147483647)
	case drivers.Float:
		max := 2147483647
		if fieldDescriptor.Precision.Valid && fieldDescriptor.Scale.Valid {
			max = int(math.Pow10(fieldDescriptor.Precision.Int - fieldDescriptor.Scale.Int))
		}
		return gofakeit.Number(1, max)
	case drivers.Blob:
		return base64.StdEncoding.EncodeToString([]byte(randomString(12)))
	case drivers.Text:
		return randomString(12)
	case drivers.Enum:
		return field.Enum[gofakeit.Number(0, len(field.Enum)-1)]
	case drivers.Bool:
		if gofakeit.Number(1, 200)%2 == 0 {
			return true
		}
		return false
	case drivers.Json:
		return fmt.Sprintf(
			`{"%s": "%s", "%s": "%s"}`,
			gofakeit.Password(true, true, false, false, false, 6),
			gofakeit.Password(true, true, false, false, false, 6),
			gofakeit.Password(true, true, false, false, false, 6),
			gofakeit.Password(true, true, false, false, false, 6),
		)
	case drivers.Time:
		return time.Date(
			gofakeit.Number(1970, 2038),
			time.Month(gofakeit.Number(0, 12)),
			gofakeit.Day(),
			gofakeit.Hour(),
			gofakeit.Minute(),
			gofakeit.Second(),
			gofakeit.NanoSecond(),
			time.UTC)
	case drivers.Year:
		return gofakeit.Number(1901, 2155)
	case drivers.XML:
		xml, err := gofakeit.XML(&gofakeit.XMLOptions{
			Type:          "single",
			RootElement:   "xml",
			RecordElement: "record",
			RowCount:      2,
			Indent:        true,
			Fields: []gofakeit.Field{
				{Name: "first_name", Function: "firstname"},
				{Name: "last_name", Function: "lastname"},
				{Name: "password", Function: "password", Params: map[string][]string{"special": {"false"}}},
			},
		})
		if err != nil {
			return nil
		}
		return string(xml)
	case drivers.UUID:
		return gofakeit.UUID()
	case drivers.BinaryString:
		return binaryString(int(field.Length))
	case drivers.Unknown:
		log.Printf("unknown field type: %s\n", fieldDescriptor.Field)
		return nil
	}

	return nil
}

// randomString generates a length size random string
func randomString(length int16) string {
	var charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func binaryString(length int) string {
	var str []string
	for i := 0; i < length; i++ {
		str = append(str, strconv.Itoa(gofakeit.Number(0, 1)))
	}
	return strings.Join(str, "")
}
