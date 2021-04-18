package fuzzer

import (
	"database/sql"
	"log"
	"sync"

	"github.com/PumpkinSeed/sqlfuzz/drivers"
	"github.com/PumpkinSeed/sqlfuzz/pkg/action"
	"github.com/PumpkinSeed/sqlfuzz/pkg/connector"
	"github.com/PumpkinSeed/sqlfuzz/pkg/flags"
	_ "github.com/lib/pq"
)

func getDriverAndDB(f flags.Flags) (drivers.Driver, *sql.DB) {
	driver := drivers.New(f.Driver)
	db := connector.Connection(driver, f)
	return driver, db
}

func runHelper(f flags.Flags, input action.SQLInsertInput) error {
	numJobs := f.Num
	workers := f.Workers
	jobs := make(chan struct{}, numJobs)
	wg := &sync.WaitGroup{}
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go worker(jobs, wg, f, input)
	}

	for j := 0; j < numJobs; j++ {
		jobs <- struct{}{}
	}
	close(jobs)
	wg.Wait()

	return nil
}

func worker(jobs <-chan struct{}, wg *sync.WaitGroup, f flags.Flags, input action.SQLInsertInput) {
	defer wg.Done()
	driver := drivers.New(f.Driver)
	db := connector.Connection(driver, f)
	defer func() {
		if err := db.Close(); err != nil {
			log.Print(err)
		}
	}()
	for range jobs {
		if err := input.Insert(); err != nil {
			log.Println(err)
		}
	}
}

// Run the commands in a worker pool
func Run(fields []drivers.FieldDescriptor, f flags.Flags) error {
	driver, db := getDriverAndDB(f)
	defer func() {
		if err := db.Close(); err != nil {
			log.Print(err)
		}
	}()
	sqlInsertInput := action.SQLInsertInput{
		SingleInsertParams: &action.SingleInsertParams{
			DB:     db,
			Driver: driver,
			Table:  f.Table,
			Fields: fields,
		},
	}
	return runHelper(f, sqlInsertInput)
}

func RunMulti(tableToFieldsMap map[string][]drivers.FieldDescriptor, insertionOrder []string, f flags.Flags) error {
	driver, db := getDriverAndDB(f)
	defer func() {
		if err := db.Close(); err != nil {
			log.Print(err)
		}
	}()
	sqlInsertInput := action.SQLInsertInput{MultiInsertParams: &action.MultiInsertParams{
		DB:               db,
		Driver:           driver,
		InsertionOrder:   insertionOrder,
		TableToFieldsMap: tableToFieldsMap,
	}}
	return runHelper(f, sqlInsertInput)
}
