package run

import (
	"database/sql"
	"log"
	"sync"

	"github.com/PumpkinSeed/sqlfuzz/drivers"
	"github.com/PumpkinSeed/sqlfuzz/pkg/action"
	"github.com/PumpkinSeed/sqlfuzz/pkg/descriptor"
	"github.com/PumpkinSeed/sqlfuzz/pkg/flags"
)

func Run(db *sql.DB, fields []descriptor.FieldDescriptor, f flags.Flags) error {
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

	return action.Exec(db, fields, drivers.New(f.Driver), f.Table)
}

func worker(db *sql.DB, jobs <-chan struct{}, fields []descriptor.FieldDescriptor, wg *sync.WaitGroup, f flags.Flags) {
	for range jobs {
		if err := action.Exec(db, fields, drivers.New(f.Driver), f.Table); err != nil {
			log.Println(err)
		}
	}

	wg.Done()
}
