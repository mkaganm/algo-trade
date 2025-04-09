package scheduler

import (
	"context"
	"github.com/robfig/cron/v3"
)

type CronScheduler struct {
	cron *cron.Cron
}

func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		cron: cron.New(),
	}
}

func (s *CronScheduler) Schedule(interval string, job func()) (cron.EntryID, error) {
	return s.cron.AddFunc(interval, job)
}

func (s *CronScheduler) Start() {
	s.cron.Start()
}

func (s *CronScheduler) Stop() context.Context {
	return s.cron.Stop()
}
