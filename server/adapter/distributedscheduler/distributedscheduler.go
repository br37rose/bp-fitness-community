package distributedscheduler

import (
	"log"
	"log/slog"
	"time"

	redislock "github.com/go-co-op/gocron-redis-lock/v2"
	"github.com/go-co-op/gocron/v2"
	_ "github.com/google/wire" // Add dependency on this package from our app.
	redis "github.com/redis/go-redis/v9"
)

// DistributedSchedulerAdapter interface provides the functions necessary for
// your application to submit tasks to be executed in the background assuming
// your application has multiple instances running concurrently.
type DistributedSchedulerAdapter interface {
	Start() error
	Shutdown() error
	ScheduleOneTimeFunc(function any, parameters ...any) error
	ScheduleEveryMinuteFunc(function any, parameters ...any) error
}

type distributedSchedulerAdapter struct {
	Logger    *slog.Logger
	Redis     redis.UniversalClient
	Scheduler gocron.Scheduler
	Locker    gocron.Locker
}

func NewAdapter(loggerp *slog.Logger, redisClient redis.UniversalClient) DistributedSchedulerAdapter {
	locker, err := redislock.NewRedisLocker(redisClient, redislock.WithTries(1)) //tODO: FIGURE THIS OUT
	if err != nil {
		log.Fatalf("failed staring redis locker with error: %v", err)
	}
	location, _ := time.LoadLocation("America/Toronto")

	s, err := gocron.NewScheduler(
		gocron.WithLocation(location),
		gocron.WithLogger(
			loggerp,
		),
		gocron.WithDistributedLocker(locker), //tODO: FIGURE THIS OUT
	)
	if err != nil {
		log.Fatalf("failed staring new scheduler with error: %v", err)
	}
	return &distributedSchedulerAdapter{
		Logger:    loggerp,
		Redis:     redisClient,
		Locker:    locker,
		Scheduler: s,
	}
}

func (adapt *distributedSchedulerAdapter) Start() error {
	adapt.Scheduler.Start()
	return nil
}

func (adapt *distributedSchedulerAdapter) Shutdown() error {
	return adapt.Scheduler.Shutdown()
}

func (adapt *distributedSchedulerAdapter) ScheduleOneTimeFunc(function any, parameters ...any) error {
	// xxx, err1 := adapt.Locker.Lock(context.Background(), "test-1-2-3")
	// if err1 != nil {
	// 	return err1
	// }
	// if xxx != nil {
	// 	defer xxx.Unlock(context.Background())
	// }

	// Create a new task with the provided function and parameters
	task := gocron.NewTask(function, parameters...)

	// run a job once, immediately - https://pkg.go.dev/github.com/go-co-op/gocron/v2#OneTimeJob
	// Create a new one-time job starting immediately
	job := gocron.OneTimeJob(gocron.OneTimeJobStartImmediately())

	// Create a new job with the task and one-time job configuration
	_, err := adapt.Scheduler.NewJob(
		job,
		task,
		// gocron.WithEventListeners(
		// 	gocron.AfterJobRuns(
		// 		func(jobID uuid.UUID, jobName string) {
		// 			// do something after the job completes
		// 			fmt.Println("done", jobID, jobName)
		// 		},
		// 	),
		// ),
	)
	if err != nil {
		adapt.Logger.Error("failed enqueuing one-time job", slog.Any("err", err))
		return err
	}
	// adapt.Logger.Debug("enqueued one-time job", slog.Any("id", j.ID()))

	return nil
}

func (adapt *distributedSchedulerAdapter) ScheduleEveryMinuteFunc(function any, parameters ...any) error {
	// xxx, err1 := adapt.Locker.Lock(context.Background(), "test-1-2-3")
	// if err1 != nil {
	// 	return err1
	// }
	// if xxx != nil {
	// 	defer xxx.Unlock(context.Background())
	// }

	// Create a new task with the provided function and parameters
	task := gocron.NewTask(function, parameters...)

	// Create a new recuring job every 1 minute.
	job := gocron.DurationJob(
		60 * time.Second,
	)

	// Create a new job with the task and one-time job configuration
	_, err := adapt.Scheduler.NewJob(
		job,
		task,
		// gocron.WithEventListeners(
		// 	gocron.AfterJobRuns(
		// 		func(jobID uuid.UUID, jobName string) {
		// 			// do something after the job completes
		// 			fmt.Println("done", jobID, jobName)
		// 		},
		// 	),
		// ),
	)
	if err != nil {
		adapt.Logger.Error("failed enqueuing one-time job", slog.Any("err", err))
		return err
	}
	// adapt.Logger.Debug("enqueued one-time job", slog.Any("id", j.ID()))

	return nil
}
