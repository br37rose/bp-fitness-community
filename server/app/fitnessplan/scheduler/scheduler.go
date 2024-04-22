package scheduler

import (
	"log/slog"

	dscheduler "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/distributedscheduler"
	oai_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/openai"
	fp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/controller"
	fp_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
)

type FitnessPlanScheduler interface {
	RunEveryMinuteUpdateFitnessPlans() error
}

// Handler Creates http request handler
type fitnessPlanSchedulerImpl struct {
	Logger               *slog.Logger
	Kmutex               kmutex.Provider
	DistributedScheduler dscheduler.DistributedSchedulerAdapter
	FitnessPlanStorer    fp_d.FitnessPlanStorer
	OpenAIConnector      oai_c.OpenAIConnector
	Controller           fp_c.FitnessPlanController
}

// NewHandler Constructor
func NewScheduler(
	loggerp *slog.Logger,
	kmutexp kmutex.Provider,
	ds dscheduler.DistributedSchedulerAdapter,
	fitnessPlanStorer fp_d.FitnessPlanStorer,
	openAI oai_c.OpenAIConnector,
	c fp_c.FitnessPlanController,
) FitnessPlanScheduler {
	return &fitnessPlanSchedulerImpl{
		Logger:               loggerp,
		Kmutex:               kmutexp,
		DistributedScheduler: ds,
		FitnessPlanStorer:    fitnessPlanStorer,
		OpenAIConnector:      openAI,
		Controller:           c,
	}
}
