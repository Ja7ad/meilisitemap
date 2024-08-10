package sched

import (
	"context"
	"sync"
	"time"

	"github.com/Ja7ad/meilisitemap/internal/logger"
)

type Sched struct {
	jobs []*job
	ctx  context.Context
	wg   sync.WaitGroup
	log  logger.Logger
}

type job struct {
	interval time.Duration
	fn       func()
}

func New(ctx context.Context, log logger.Logger) *Sched {
	s := new(Sched)
	s.ctx = ctx
	s.log = log
	s.jobs = make([]*job, 0)
	return s
}

func (s *Sched) AddJob(jobFunc func(), interval time.Duration) {
	s.jobs = append(s.jobs, &job{fn: jobFunc, interval: interval})
}

func (s *Sched) Len() int {
	return len(s.jobs)
}

func (s *Sched) Start() {
	if len(s.jobs) == 0 {
		s.log.Warn("sched: no job functions defined, scheduler will not start")
		return
	}

	s.log.Info("starting scheduler jobs...", "total_jobs", len(s.jobs))

	for _, j := range s.jobs {
		s.wg.Add(1)
		go func(jFunc func(), interval time.Duration) {
			defer s.wg.Done()
			ticker := time.NewTicker(interval)

			defer func() {
				ticker.Stop()
				ticker = nil
			}()

			for {
				select {
				case <-s.ctx.Done():
					return
				case <-ticker.C:
					jFunc()
				}
			}
		}(j.fn, j.interval)
	}

	s.wg.Wait()
}
