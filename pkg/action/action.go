package action

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/PumpkinSeed/sqlfuzz/drivers/types"
	"github.com/brianvoe/gofakeit/v5"
	_ "github.com/lib/pq"
	"github.com/rs/xid"
)

type SingleInsertParams struct {
	DB     *sql.DB
	Driver types.Driver
	Table  string
	Fields []types.FieldDescriptor
}

type MultiInsertParams struct {
	DB               *sql.DB
	Driver           types.Driver
	InsertionOrder   []string
	TableToFieldsMap map[string][]types.FieldDescriptor
}

type SQLInsertInput struct {
	SingleInsertParams *SingleInsertParams
	MultiInsertParams  *MultiInsertParams
}

func (sqlInsertInput SQLInsertInput) Insert() error {
	if sqlInsertInput.SingleInsertParams != nil {
		return sqlInsertInput.singleInsert()
	} else if sqlInsertInput.MultiInsertParams != nil {
		return sqlInsertInput.multiInsert()
	}
	return errors.New("action: error in sql insert input. Both single and multi insert arguments are not initialised")
}

func (sqlInsertInput SQLInsertInput) multiInsert() error {
	multiInsertParams := sqlInsertInput.MultiInsertParams
	if multiInsertParams == nil {
		return errors.New("action : error during multi insert. Could not find necessary arguments")
	}
	tableFieldValuesMap := make(map[string]map[string]interface{})
	for _, table := range multiInsertParams.InsertionOrder {
		if fields, ok := multiInsertParams.TableToFieldsMap[table]; ok {
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
					val, err := multiInsertParams.Driver.GetLatestColumnValue(field.ForeignKeyDescriptor.ForeignTableName, field.ForeignKeyDescriptor.ForeignColumnName, multiInsertParams.DB)
					if err != nil {
						return err
					}
					data = val
					// Get from table. If no value present in table as well, throw error.
				} else {
					data = generateData(multiInsertParams.Driver, field)
				}
				values = append(values, data)
			}
			query := multiInsertParams.Driver.Insert(f, table)
			_, err := multiInsertParams.DB.Exec(query, values...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// singleInsert is inserting a random generated data into the chosen table
func (sqlInsertInput SQLInsertInput) singleInsert() error {
	insertParams := sqlInsertInput.SingleInsertParams
	if insertParams == nil {
		return errors.New("action : error during insert. Could not find necessary arguments")
	}
	var f = make([]string, 0, len(insertParams.Fields))
	var values = make([]interface{}, 0, len(insertParams.Fields))
	for _, field := range insertParams.Fields {
		// Has default value. No need to insert this field manually.
		if field.HasDefaultValue {
			continue
		}
		f = append(f, field.Field)
		values = append(values, generateData(insertParams.Driver, field))
	}
	query := insertParams.Driver.Insert(f, insertParams.Table)

	_, err := insertParams.DB.Exec(query, values...)
	return err
}

// generateData generates random data based on the field
func generateData(driver types.Driver, fieldDescriptor types.FieldDescriptor) interface{} {
	field := driver.MapField(fieldDescriptor)
	switch field.Type {
	case types.String:
		if field.Length > 19 {
			return xid.New().String()
		}
		if field.Length > 0 {
			return randomString(field.Length)
		}
		return randomString(20)
	case types.Int16:
		return gofakeit.Number(1, 32766)
	case types.Int32:
		return gofakeit.Number(1, 2147483647)
	case types.Float:
		max := 2147483647
		if fieldDescriptor.Precision.Valid && fieldDescriptor.Scale.Valid {
			max = int(math.Pow10(fieldDescriptor.Precision.Int - fieldDescriptor.Scale.Int))
		}
		return gofakeit.Number(1, max)
	case types.Blob:
		return base64.StdEncoding.EncodeToString([]byte(randomString(12)))
	case types.Text:
		return randomString(12)
	case types.Enum:
		return field.Enum[gofakeit.Number(0, len(field.Enum)-1)]
	case types.Bool:
		if gofakeit.Number(1, 200)%2 == 0 {
			return true
		}
		return false
	case types.Json:
		return fmt.Sprintf(
			`{"%s": "%s", "%s": "%s"}`,
			gofakeit.Password(true, true, false, false, false, 6),
			gofakeit.Password(true, true, false, false, false, 6),
			gofakeit.Password(true, true, false, false, false, 6),
			gofakeit.Password(true, true, false, false, false, 6),
		)
	case types.Time:
		return time.Date(
			gofakeit.Number(1980, 2028),
			time.Month(gofakeit.Number(0, 12)),
			gofakeit.Day(),
			gofakeit.Hour(),
			gofakeit.Minute(),
			gofakeit.Second(),
			gofakeit.NanoSecond(),
			time.UTC)
	case types.Year:
		return gofakeit.Number(1901, 2155)
	case types.XML:
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
	case types.UUID:
		return gofakeit.UUID()
	case types.BinaryString:
		return binaryString(int(field.Length))
	case types.Unknown:
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
