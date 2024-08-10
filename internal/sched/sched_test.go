package sched

import (
	"context"
	"testing"
	"time"

	"github.com/Ja7ad/meilisitemap/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestSched(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	interval := 10 * time.Millisecond
	s := New(ctx, logger.DefaultLogger)

	var count int
	job := func() {
		count++
	}

	s.AddJob(job, interval)
	go s.Start()

	time.Sleep(50 * time.Millisecond)

	assert.Greater(t, count, 0, "Expected job to have run at least once")
}

func TestSchedStop(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	interval := 10 * time.Millisecond
	s := New(ctx, logger.DefaultLogger)

	var count int
	job := func() {
		count++
	}

	s.AddJob(job, interval)
	go s.Start()

	time.Sleep(30 * time.Millisecond)
}

func TestSchedNoJobs(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	s := New(ctx, logger.DefaultLogger)

	// Start and stop scheduler with no jobs
	go s.Start()
	time.Sleep(30 * time.Millisecond)
}

func TestSchedMultipleJobs(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	interval := 10 * time.Millisecond
	s := New(ctx, logger.DefaultLogger)

	var count1, count2 int
	job1 := func() {
		count1++
	}
	job2 := func() {
		count2++
	}

	s.AddJob(job1, interval)
	s.AddJob(job2, interval)
	go s.Start()

	// Let the jobs run for a few intervals
	time.Sleep(50 * time.Millisecond)

	assert.Greater(t, count1, 0, "Expected job1 to have run at least once")
	assert.Greater(t, count2, 0, "Expected job2 to have run at least once")
}
