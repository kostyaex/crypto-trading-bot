package app

import (
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/service/exchange"
	"log"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	c         *cron.Cron
	exchanges []exchange.Exchange
	logger    *logger.Logger
}

func NewScheduler(exchanges []exchange.Exchange, logger *logger.Logger) *Scheduler {
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(logger.Writer(), "cron: ", log.LstdFlags))))
	return &Scheduler{
		c:         c,
		exchanges: exchanges,
		logger:    logger,
	}
}

func (s *Scheduler) Start() {
	s.c.Start()
	s.logger.Infof("Scheduler started")
}

func (s *Scheduler) Stop() {
	s.c.Stop()
	s.logger.Infof("Scheduler stopped")
}

func (s *Scheduler) AddJob(spec string, job cron.Job) (cron.EntryID, error) {
	id, err := s.c.AddJob(spec, job)
	if err != nil {
		s.logger.Errorf("Failed to add job: %v", err)
		return 0, err
	}
	s.logger.Infof("Job added with spec: %s, ID: %d", spec, id)
	return id, nil
}
