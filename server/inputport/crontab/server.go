package crontab

import (
	"log/slog"

	"github.com/mileusna/crontab"

	oai_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/openai"
	agg_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/controller"
	fp_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	googlefitapp_cron "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/crontab"
	googlefitdp_cron "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/crontab"
	rank_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/controller"
	user_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

type InputPortServer interface {
	Run()
	Shutdown()
}

type crontabInputPort struct {
	Config                    *config.Conf
	Logger                    *slog.Logger
	Crontab                   *crontab.Crontab
	UserController            user_c.UserController
	AggregatePointController  agg_c.AggregatePointController
	RankPointController       rank_c.RankPointController
	GoogleFitDataPointCrontab googlefitdp_cron.GoogleFitDataPointCrontaber
	GoogleFitAppCrontab       googlefitapp_cron.GoogleFitAppCrontaber
	FitnessPlanStorer         fp_d.FitnessPlanStorer
	OpenAIConnector           oai_c.OpenAIConnector
}

func NewInputPort(
	configp *config.Conf,
	loggerp *slog.Logger,
	usrContr user_c.UserController,
	aggContr agg_c.AggregatePointController,
	rankContr rank_c.RankPointController,
	gfdb googlefitdp_cron.GoogleFitDataPointCrontaber,
	gfaCron googlefitapp_cron.GoogleFitAppCrontaber,
	fitnessPlanStorer fp_d.FitnessPlanStorer,
	openAI oai_c.OpenAIConnector,
) InputPortServer {
	// Initialize.

	ctab := crontab.New() // create cron table

	// Create our HTTP server controller.
	p := &crontabInputPort{
		Config:                    configp,
		Logger:                    loggerp,
		Crontab:                   ctab,
		UserController:            usrContr,
		AggregatePointController:  aggContr,
		RankPointController:       rankContr,
		GoogleFitDataPointCrontab: gfdb,
		GoogleFitAppCrontab:       gfaCron,
		FitnessPlanStorer:         fitnessPlanStorer,
		OpenAIConnector:           openAI,
	}

	return p
}

func (port *crontabInputPort) Run() {
	// // (For debugging purposes only)
	// // Run the following code on startup of the application.
	// // port.GoogleFitAppCrontab.RefreshTokensFromGoogleJob()
	// // port.GoogleFitAppCrontab.PullDataFromGoogleJob()
	// // port.GoogleFitAppCrontab.ProcessAllQueuedDataTask()
	// // port.GoogleFitDataPointCrontab.DeleteAllAnomalousData()
	// // port.processAllActiveSimulators() //TODO: IMPLEMENT.
	//
	// // (For debugging purposes only)
	// // Run the following code on startup of the application.
	// // port.AggregateThisHour()
	// // port.AggregateLastHour()
	// // port.AggregateToday()
	// // port.AggregateYesterday()
	// // port.AggregateThisISOWeek()
	// // port.AggregateLastISOWeek()
	// // port.AggregateThisMonth()
	// // port.AggregateLastMonth()
	// // port.AggregateThisYear()
	// // port.AggregateLastYear()
	// // port.RankToday()
	// // port.RankThisISOWeek()
	// // port.RankThisMonth()
	// // port.RankThisYear()
	//
	// port.Logger.Info("Crontab server running")
	//
	// // MustAddJob is like AddJob but panics on wrong syntax or problems with func/args
	// // This aproach is similar to regexp.Compile and regexp.MustCompile from go's standard library,  used for easier initialization on startup
	//
	// // For debugging purposes
	// // port.Crontab.MustAddJob("* * * * *", port.ping)                 // every minute
	//
	// // Google Fit data.
	// // The following section will include Google Fit web-services interaction
	// // related background tasks that are important for fetching or simulating
	// port.Crontab.MustAddJob("* * * * *", port.GoogleFitAppCrontab.RefreshTokensFromGoogleJob)   // every minute TOOD: ADD BACK WHEN READY
	// port.Crontab.MustAddJob("*/5 * * * *", port.GoogleFitAppCrontab.PullDataFromGoogleJob)      // every 5 minutes TOOD: ADD BACK WHEN READY
	// port.Crontab.MustAddJob("* * * * *", port.GoogleFitAppCrontab.ProcessAllQueuedDataTask)     // every minute
	// port.Crontab.MustAddJob("* * * * *", port.GoogleFitDataPointCrontab.DeleteAllAnomalousData) // every minute
	// // port.Crontab.MustAddJob("* * * * *", port.processAllActiveSimulators) // every minute //TODO: IMPLEMENT
	//
	// // Aggregation.
	// // The following section will include code that takes the raw data points
	// port.Crontab.MustAddJob("* * * * *", port.AggregateThisHour)    // every minute
	// port.Crontab.MustAddJob("* * * * *", port.AggregateLastHour)    // every minute
	// port.Crontab.MustAddJob("* * * * *", port.AggregateToday)       // every minute
	// port.Crontab.MustAddJob("* * * * *", port.AggregateYesterday)   // every minute
	// port.Crontab.MustAddJob("* * * * *", port.AggregateThisISOWeek) // every minute
	// port.Crontab.MustAddJob("* * * * *", port.AggregateLastISOWeek) // every minute
	//
	// // // TODO: Comment the following when going live.
	// port.Crontab.MustAddJob("* * * * *", port.AggregateThisMonth) // every minute
	// port.Crontab.MustAddJob("* * * * *", port.AggregateLastMonth) // every minute
	// port.Crontab.MustAddJob("* * * * *", port.AggregateThisYear)  // every minute
	// port.Crontab.MustAddJob("* * * * *", port.AggregateLastYear)  // every minute
	// // // TODO: Uncomment the following when going live.
	// // // port.Crontab.MustAddJob("0 * * * *", port.AggregateThisMonth)   // every hour
	// // // port.Crontab.MustAddJob("0 * * * *", port.AggregateLastMonth)   // every hour
	// // // port.Crontab.MustAddJob("0 0 * * 0", port.AggregateThisYear)    // every sunday (Code via https://www.linuxshelltips.com/cron-run-every-sunday-at-midnight/)
	// // // port.Crontab.MustAddJob("0 0 * * 0", port.AggregateLastYear)    // every sunday (Code via https://www.linuxshelltips.com/cron-run-every-sunday-at-midnight/)
	//
	// // Leaderboard Ranking
	// // // The following section will enable the ranking system for the different periods of the year.
	// port.Crontab.MustAddJob("* * * * *", port.RankToday)       // every minute
	// port.Crontab.MustAddJob("* * * * *", port.RankThisISOWeek) // every minute
	// port.Crontab.MustAddJob("* * * * *", port.RankThisMonth)   // every minute
	// port.Crontab.MustAddJob("* * * * *", port.RankThisYear)    // every minute
	//
	// // Fitness Plans
	// port.Crontab.MustAddJob("* * * * *", port.updateFitnessPlans) // every minute
}

func (port *crontabInputPort) Shutdown() {
	port.Logger.Info("Crontab server shutdown")
	port.Crontab.Clear()
}

func (port *crontabInputPort) ping() {
	port.Logger.Info("Crontab ping")
}
