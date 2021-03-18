package fuzzer

import (
	"database/sql"
	"log"
	"sync"

	"github.com/PumpkinSeed/sqlfuzz/drivers"
	"github.com/PumpkinSeed/sqlfuzz/pkg/action"
	"github.com/PumpkinSeed/sqlfuzz/pkg/flags"
	_ "github.com/lib/pq"
)

// Run the commands in a worker pool
func Run(db *sql.DB, fields []drivers.FieldDescriptor, f flags.Flags) error {
	numJobs := f.Num
	workers := f.Workers
	jobs := make(chan struct{}, numJobs)
	wg := &sync.WaitGroup{}
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go worker(db, jobs, fields, wg, f)
	}

	for j := 0; j < numJobs; j++ {
		jobs <- struct{}{}
	}
	close(jobs)
	wg.Wait()

	return action.Insert(db, fields, drivers.New(f.Driver), f.Table)
}

// worker of the worker pool, executing the command, logging if fails
func worker(db *sql.DB, jobs <-chan struct{}, fields []drivers.FieldDescriptor, wg *sync.WaitGroup, f flags.Flags) {
	for range jobs {
		if err := action.Insert(db, fields, drivers.New(f.Driver), f.Table); err != nil {
			log.Println(err)
		}
	}

	wg.Done()
}
