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
	db := connector.Connection(driver)
	db.SetMaxOpenConns(f.Workers)
	db.SetMaxIdleConns(f.Workers)
	return driver, db
}

func runHelper(f flags.Flags, exec func(...interface{}) error, args []interface{}) error {
	numJobs := f.Num
	workers := f.Workers
	jobs := make(chan struct{}, numJobs)
	wg := &sync.WaitGroup{}
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go newWorker(jobs, wg, f, exec, args)
	}

	for j := 0; j < numJobs; j++ {
		jobs <- struct{}{}
	}
	close(jobs)
	wg.Wait()

	return nil
}

// Run the commands in a worker pool
func Run(fields []drivers.FieldDescriptor, f flags.Flags) error {
	driver, db := getDriverAndDB(f)
	defer func() {
		if err := db.Close(); err != nil {
			log.Print(err)
		}
	}()
	return runHelper(f, action.Insert, []interface{}{db, fields, driver, f.Table})
	//numJobs := f.Num
	//workers := f.Workers
	//jobs := make(chan struct{}, numJobs)
	//wg := &sync.WaitGroup{}
	//wg.Add(workers)
	//for w := 0; w < workers; w++ {
	//	go worker(jobs, fields, wg, f)
	//}
	//
	//for j := 0; j < numJobs; j++ {
	//	jobs <- struct{}{}
	//}
	//close(jobs)
	//wg.Wait()
}

func RunMulti(tableToFieldsMap map[string][]drivers.FieldDescriptor, insertionOrder []string, f flags.Flags) error {
	driver, db := getDriverAndDB(f)
	defer func() {
		if err := db.Close(); err != nil {
			log.Print(err)
		}
	}()
	return runHelper(f, action.InsertMulti, []interface{}{db, driver, tableToFieldsMap, insertionOrder})
}

func newWorker(jobs <-chan struct{}, wg *sync.WaitGroup, f flags.Flags, exec func(...interface{}) error, args []interface{}) {
	defer wg.Done()
	driver := drivers.New(f.Driver)
	db := connector.Connection(driver)
	defer func() {
		if err := db.Close(); err != nil {
			log.Print(err)
		}
	}()
	for range jobs {
		if err := exec(args...); err != nil {
			log.Println(err)
		}
	}
}

// worker of the worker pool, executing the command, logging if fails
func worker(jobs <-chan struct{}, fields []drivers.FieldDescriptor, wg *sync.WaitGroup, f flags.Flags) {
	driver := drivers.New(f.Driver)
	db := connector.Connection(driver)
	defer func() {
		if err := db.Close(); err != nil {
			log.Print(err)
		}
	}()

	for range jobs {
		if err := action.Insert(db, fields, driver, f.Table); err != nil {
			log.Println(err)
		}
	}

	wg.Done()
}
