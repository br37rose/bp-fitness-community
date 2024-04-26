package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	_ "time/tzdata"

	_ "github.com/google/wire"   // Add dependency on this package from our app.
	_ "go.uber.org/automaxprocs" // Automatically set GOMAXPROCS to match Linux container CPU quota.

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/crontab"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/eventscheduler"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/http"
)

type Application struct {
	Logger               *slog.Logger
	HttpServer           http.InputPortServer
	CrontabServer        crontab.InputPortServer
	EventSchedulerServer eventscheduler.InputPortServer
}

// NewApplication is application construction function which is automatically called by `Google Wire` dependency injection library.
func NewApplication(
	loggerp *slog.Logger,
	httpServer http.InputPortServer,
	crontabServer crontab.InputPortServer,
	schedulerServer eventscheduler.InputPortServer,
) Application {
	return Application{
		Logger:               loggerp,
		HttpServer:           httpServer,
		CrontabServer:        crontabServer,
		EventSchedulerServer: schedulerServer,
	}
}

func (a Application) Execute() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGUSR1)

	// Run in background the HTTP server.
	go a.HttpServer.Run()

	// Run in background the Crontab server.
	go a.CrontabServer.Run()

	// Run in background the event scheduler server.
	go a.EventSchedulerServer.Run()

	a.Logger.Info("Application started")

	// Run the main loop blocking code while other input ports run in background.
	<-done

	a.Shutdown()
}

func (a Application) Shutdown() {
	a.HttpServer.Shutdown()
	a.CrontabServer.Shutdown()
	a.EventSchedulerServer.Shutdown()
	a.Logger.Info("Application shutdown")
}

// main function is the main entry point into the code.
func main() {
	// Call the `InitializeEvent` function which will call `Google Wire` dependency injection package to load up all this projects dependencies together.
	Application := InitializeEvent()

	// Start the application!
	Application.Execute()
}
