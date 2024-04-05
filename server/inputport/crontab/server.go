package crontab

import (
	"log/slog"

	"github.com/mileusna/crontab"

	agg_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/controller"
	googlefitapp_cron "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/crontab"
	rank_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/controller"
	user_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
)

type InputPortServer interface {
	Run()
	Shutdown()
}

type crontabInputPort struct {
	Config                   *config.Conf
	Logger                   *slog.Logger
	Crontab                  *crontab.Crontab
	UserController           user_c.UserController
	AggregatePointController agg_c.AggregatePointController
	RankPointController      rank_c.RankPointController
	GoogleFitAppCrontab      googlefitapp_cron.GoogleFitAppCrontaber
}

func NewInputPort(
	configp *config.Conf,
	loggerp *slog.Logger,
	usrContr user_c.UserController,
	aggContr agg_c.AggregatePointController,
	rankContr rank_c.RankPointController,
	gfaCron googlefitapp_cron.GoogleFitAppCrontaber,
) InputPortServer {
	// Initialize.

	ctab := crontab.New() // create cron table

	// Create our HTTP server controller.
	p := &crontabInputPort{
		Config:                   configp,
		Logger:                   loggerp,
		Crontab:                  ctab,
		UserController:           usrContr,
		AggregatePointController: aggContr,
		RankPointController:      rankContr,
		GoogleFitAppCrontab:      gfaCron,
	}

	return p
}

func (port *crontabInputPort) Run() {

	port.Logger.Info("Crontab server running")

	// MustAddJob is like AddJob but panics on wrong syntax or problems with func/args
	// This aproach is similar to regexp.Compile and regexp.MustCompile from go's standard library,  used for easier initialization on startup

	// For debugging purposes
	// port.Crontab.MustAddJob("* * * * *", port.ping)                 // every minute

	//----------------------------------------------------------------------------
	// The following section will include Google Fit web-services interaction
	// related background tasks that are important for fetching or simulating
	// Google Fit data.
	// port.Crontab.MustAddJob("*/15 * * * *", port.pullGoogleFitAppRawData) // every 15 minutes
	// port.Crontab.MustAddJob("*/5 * * * *", port.processAllQueuedData)     // every 5 minutes
	// port.Crontab.MustAddJob("* * * * *", port.processAllActiveSimulators) // every minute
	//----------------------------------------------------------------------------

	// port.Crontab.MustAddJob("* * * * *", port.GoogleFitAppCrontab.RefreshTokensFromGoogleJob) // every minute TOOD: ADD BACK WHEN READY
	// port.Crontab.MustAddJob("*/5 * * * *", port.GoogleFitAppCrontab.PullDataFromGoogleJob)    // every 5 minutes TOOD: ADD BACK WHEN READY
	port.Crontab.MustAddJob("* * * * *", port.GoogleFitAppCrontab.ProcessQueuedDataTask) // every minute

	// The following section will include code that takes the raw data points
	port.Crontab.MustAddJob("* * * * *", port.AggregateThisHour)    // every minute
	port.Crontab.MustAddJob("* * * * *", port.AggregateLastHour)    // every minute
	port.Crontab.MustAddJob("* * * * *", port.AggregateToday)       // every minute
	port.Crontab.MustAddJob("* * * * *", port.AggregateYesterday)   // every minute
	port.Crontab.MustAddJob("* * * * *", port.AggregateThisISOWeek) // every minute
	port.Crontab.MustAddJob("* * * * *", port.AggregateLastISOWeek) // every minute

	// TODO: Comment the following when going live.
	port.Crontab.MustAddJob("* * * * *", port.AggregateThisMonth) // every hour
	port.Crontab.MustAddJob("* * * * *", port.AggregateLastMonth) // every hour
	port.Crontab.MustAddJob("* * * * *", port.AggregateThisYear)  // every hour
	port.Crontab.MustAddJob("* * * * *", port.AggregateLastYear)  // every hour

	// TODO: The code below is commented out until we need to use it for performance reasons.
	// port.Crontab.MustAddJob("0 * * * *", port.AggregateThisMonth)   // every hour
	// port.Crontab.MustAddJob("0 * * * *", port.AggregateLastMonth)   // every hour
	// port.Crontab.MustAddJob("0 0 * * 0", port.AggregateThisYear)    // every sunday (Code via https://www.linuxshelltips.com/cron-run-every-sunday-at-midnight/)
	// port.Crontab.MustAddJob("0 0 * * 0", port.AggregateLastYear)    // every sunday (Code via https://www.linuxshelltips.com/cron-run-every-sunday-at-midnight/)

	// The following section will enable the ranking system for the different periods of the year.
	port.Crontab.MustAddJob("* * * * *", port.RankToday)       // every minute
	port.Crontab.MustAddJob("* * * * *", port.RankThisISOWeek) // every minute
	port.Crontab.MustAddJob("* * * * *", port.RankThisMonth)   // every minute
	port.Crontab.MustAddJob("* * * * *", port.RankThisYear)    // every minute

	// // (For debugging purposes only)
	// // Run the following code on startup of the application.
	// port.pullGoogleFitAppRawData()
	// port.processAllQueuedData()
	// port.processAllActiveSimulators()

	// (For debugging purposes only)
	// Run the following code on startup of the application.
	// port.AggregateThisHour()
	// port.AggregateLastHour()
	// port.AggregateToday()
	// port.AggregateYesterday()
	// port.AggregateThisISOWeek()
	// port.AggregateLastISOWeek()
	// port.AggregateThisMonth()
	// port.AggregateLastMonth()
	// port.AggregateThisYear()
	// port.AggregateLastYear()
	port.RankToday()
	// port.RankThisISOWeek()
	// port.RankThisMonth()
	// port.RankThisYear()

	// port.GoogleFitAppCrontab.RefreshTokensFromGoogleJob() //TODO: Comment out when ready.
	// port.GoogleFitAppCrontab.PullDataFromGoogleJob()      //TODO: Comment out when ready.
	port.GoogleFitAppCrontab.ProcessQueuedDataTask() //TODO: Comment out when ready.
}

func (port *crontabInputPort) Shutdown() {
	port.Logger.Info("Crontab server shutdown")
	port.Crontab.Clear()
}

func (port *crontabInputPort) ping() {
	port.Logger.Info("Crontab ping")
}
