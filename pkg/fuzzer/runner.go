package fuzzer

import (
	"github.com/PumpkinSeed/sqlfuzz/drivers"
	"github.com/PumpkinSeed/sqlfuzz/pkg/action"
	"github.com/PumpkinSeed/sqlfuzz/pkg/connector"
	"github.com/PumpkinSeed/sqlfuzz/pkg/flags"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

// Run the commands in a worker pool
func Run(fields []drivers.FieldDescriptor, f flags.Flags) error {
	numJobs := f.Num
	workers := f.Workers
	jobs := make(chan struct{}, numJobs)
	wg := &sync.WaitGroup{}
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go worker(jobs, fields, wg, f)
	}

	for j := 0; j < numJobs; j++ {
		jobs <- struct{}{}
	}
	close(jobs)
	wg.Wait()

	return nil
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
