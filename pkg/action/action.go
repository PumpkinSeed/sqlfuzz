package action

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/PumpkinSeed/sqlfuzz/drivers"
	"github.com/PumpkinSeed/sqlfuzz/pkg/descriptor"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/rs/xid"
)

// Insert is inserting a random generated data into the chosen table
func Insert(db *sql.DB, fields []descriptor.FieldDescriptor, driver drivers.Driver, table string) error {
	var f = make([]string, 0, len(fields))
	var values = make([]interface{}, 0, len(fields))
	for _, field := range fields {
		f = append(f, field.Field)

		values = append(values, generateData(driver, field.Type))
	}
	query := driver.Insert(f, table)

	ins, err := db.Prepare(query)
	if err != nil {
		log.Printf("invalid preparing query: %s\n", query)
		return fmt.Errorf("error preparing query: %w", err)
	}

	_, err = ins.Exec(values...)
	return err
}

// generateData generates random data based on the field
func generateData(driver drivers.Driver, t string) interface{} {
	field := driver.MapField(t)
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
		return gofakeit.Number(1, 2147483647)
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
		return gofakeit.Date()
	case drivers.Unknown:
		log.Printf("unknown field type: %s\n", t)
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
