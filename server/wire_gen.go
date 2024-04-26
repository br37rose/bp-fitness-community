// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cache/mongodbcache"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/distributedlocker"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/emailer/mailgun"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/eventscheduler"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/openai"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/paymentprocessor/stripe"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/storage/s3"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/templatedemailer"
	controller18 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/controller"
	datastore19 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	httptransport17 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/httptransport"
	scheduler3 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/scheduler"
	controller7 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/controller"
	datastore5 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/datastore"
	httptransport7 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/httptransport"
	controller20 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/biometric/controller"
	httptransport19 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/biometric/httptransport"
	controller17 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/controller"
	datastore18 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/datastore"
	httptransport16 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/httptransport"
	datastore6 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/equipment/datastore"
	datastore13 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/eventlog/datastore"
	controller5 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/controller"
	datastore7 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	httptransport5 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/httptransport"
	controller24 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnesschallenge/controller"
	datastore23 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnesschallenge/datastore"
	httptransport23 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnesschallenge/httptransport"
	controller13 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/controller"
	datastore14 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	scheduler5 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/scheduler"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/gateway/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/gateway/httptransport"
	controller16 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/controller"
	crontab2 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/crontab"
	datastore17 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	httptransport15 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/httptransport"
	scheduler2 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/scheduler"
	controller15 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/crontab"
	datastore16 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/datastore"
	httptransport14 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/httptransport"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/scheduler"
	controller12 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/invoice/controller"
	datastore12 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/invoice/datastore"
	httptransport12 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/invoice/httptransport"
	controller6 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/member/controller"
	httptransport6 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/member/httptransport"
	controller14 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/controller"
	datastore15 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/datastore"
	httptransport13 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/httptransport"
	controller11 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/controller"
	datastore8 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/datastore"
	httptransport11 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/httptransport"
	controller3 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/controller"
	datastore3 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	httptransport3 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/httptransport"
	stripe2 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/paymentprocessor/controller/stripe"
	stripe3 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/paymentprocessor/httptransport/stripe"
	controller23 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/question/controller"
	datastore22 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/question/datastore"
	httptransport22 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/question/httptransport"
	controller19 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	httptransport18 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/httptransport"
	scheduler4 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/scheduler"
	controller4 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/tag/controller"
	datastore4 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/tag/datastore"
	httptransport4 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/tag/httptransport"
	controller22 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/controller"
	datastore20 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/datastore"
	httptransport20 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/httptransport"
	controller2 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/controller"
	datastore2 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	httptransport2 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/httptransport"
	controller8 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocategory/controller"
	datastore9 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocategory/datastore"
	httptransport8 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocategory/httptransport"
	controller9 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocollection/controller"
	datastore10 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocollection/datastore"
	httptransport9 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocollection/httptransport"
	controller10 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocontent/controller"
	datastore11 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocontent/datastore"
	httptransport10 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocontent/httptransport"
	controller21 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/controller"
	datastore21 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/datastore"
	httptransport21 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/httptransport"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	crontab3 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/crontab"
	eventscheduler2 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/eventscheduler"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/http"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/http/fitnessplan"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/http/middleware"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/jwt"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/logger"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/mongodb"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/password"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/redis"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/time"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

import (
	_ "github.com/google/wire"
	_ "go.uber.org/automaxprocs"
	_ "time/tzdata"
)

// Injectors from wire.go:

func InitializeEvent() Application {
	slogLogger := logger.NewProvider()
	conf := config.New()
	provider := uuid.NewProvider()
	timeProvider := time.NewProvider()
	jwtProvider := jwt.NewProvider(conf)
	passwordProvider := password.NewProvider()
	client := mongodb.NewProvider(conf, slogLogger)
	cacher := mongodbcache.NewCache(conf, slogLogger, client)
	s3Storager := s3.NewStorage(conf, slogLogger, provider)
	kmutexProvider := kmutex.NewProvider()
	universalClient := redis.NewProvider(conf, slogLogger)
	adapter := distributedlocker.NewAdapter(slogLogger, universalClient)
	emailer := mailgun.NewEmailer(conf, slogLogger, provider)
	templatedEmailer := templatedemailer.NewTemplatedEmailer(conf, slogLogger, provider, emailer)
	paymentProcessor := stripe.NewPaymentProcessor(conf, slogLogger, provider)
	rankPointStorer := datastore.NewDatastore(conf, slogLogger, client)
	userStorer := datastore2.NewDatastore(conf, slogLogger, client)
	organizationStorer := datastore3.NewDatastore(conf, slogLogger, client)
	gatewayController := controller.NewController(conf, slogLogger, provider, jwtProvider, passwordProvider, cacher, s3Storager, client, kmutexProvider, adapter, templatedEmailer, paymentProcessor, rankPointStorer, userStorer, organizationStorer)
	middlewareMiddleware := middleware.NewMiddleware(conf, slogLogger, provider, timeProvider, jwtProvider, gatewayController)
	handler := httptransport.NewHandler(slogLogger, gatewayController)
	userController := controller2.NewController(conf, slogLogger, provider, client, passwordProvider, userStorer)
	httptransportHandler := httptransport2.NewHandler(slogLogger, userController)
	organizationController := controller3.NewController(conf, slogLogger, provider, passwordProvider, s3Storager, client, templatedEmailer, organizationStorer, userStorer)
	handler2 := httptransport3.NewHandler(slogLogger, organizationController)
	tagStorer := datastore4.NewDatastore(conf, slogLogger, client)
	tagController := controller4.NewController(conf, slogLogger, provider, s3Storager, passwordProvider, kmutexProvider, client, templatedEmailer, userStorer, tagStorer)
	handler3 := httptransport4.NewHandler(slogLogger, tagController)
	attachmentStorer := datastore5.NewDatastore(conf, slogLogger, client)
	equipmentStorer := datastore6.NewDatastore(conf, slogLogger, client)
	exerciseStorer := datastore7.NewDatastore(conf, slogLogger, client)
	offerStorer := datastore8.NewDatastore(conf, slogLogger, client)
	exerciseController := controller5.NewController(conf, slogLogger, provider, kmutexProvider, s3Storager, client, attachmentStorer, equipmentStorer, exerciseStorer, offerStorer, userStorer)
	handler4 := httptransport5.NewHandler(slogLogger, exerciseController)
	memberController := controller6.NewController(conf, slogLogger, provider, s3Storager, passwordProvider, paymentProcessor, client, organizationStorer, rankPointStorer, userStorer)
	handler5 := httptransport6.NewHandler(slogLogger, memberController)
	attachmentController := controller7.NewController(conf, slogLogger, provider, s3Storager, client, emailer, attachmentStorer, exerciseStorer, userStorer)
	handler6 := httptransport7.NewHandler(attachmentController)
	videoCategoryStorer := datastore9.NewDatastore(conf, slogLogger, client)
	videoCategoryController := controller8.NewController(conf, slogLogger, provider, s3Storager, emailer, client, videoCategoryStorer, exerciseStorer, userStorer)
	handler7 := httptransport8.NewHandler(slogLogger, videoCategoryController)
	videoCollectionStorer := datastore10.NewDatastore(conf, slogLogger, client)
	videoCollectionController := controller9.NewController(conf, slogLogger, provider, kmutexProvider, s3Storager, client, attachmentStorer, videoCategoryStorer, videoCollectionStorer)
	handler8 := httptransport9.NewHandler(slogLogger, videoCollectionController)
	videoContentStorer := datastore11.NewDatastore(conf, slogLogger, client)
	videoContentController := controller10.NewController(conf, slogLogger, provider, kmutexProvider, s3Storager, attachmentStorer, client, videoCategoryStorer, videoCollectionStorer, videoContentStorer, offerStorer, userStorer)
	handler9 := httptransport10.NewHandler(slogLogger, videoContentController)
	offerontroller := controller11.NewController(conf, slogLogger, provider, client, organizationStorer, offerStorer, userStorer, exerciseStorer, videoContentStorer)
	handler10 := httptransport11.NewHandler(slogLogger, offerontroller)
	invoiceStorer := datastore12.NewDatastore(conf, slogLogger, client)
	invoiceController := controller12.NewController(conf, slogLogger, provider, client, organizationStorer, invoiceStorer)
	handler11 := httptransport12.NewHandler(slogLogger, invoiceController)
	eventLogStorer := datastore13.NewDatastore(conf, slogLogger, client)
	stripePaymentProcessorController := stripe2.NewController(conf, slogLogger, provider, s3Storager, passwordProvider, emailer, templatedEmailer, paymentProcessor, organizationStorer, userStorer, invoiceStorer, offerStorer, eventLogStorer)
	stripeHandler := stripe3.NewHandler(slogLogger, stripePaymentProcessorController)
	openAIConnector := openai.NewOpenAIConnector(conf, slogLogger, client)
	fitnessPlanStorer := datastore14.NewDatastore(conf, slogLogger, client)
	fitnessPlanController := controller13.NewController(conf, slogLogger, provider, s3Storager, emailer, client, kmutexProvider, openAIConnector, fitnessPlanStorer, exerciseStorer, userStorer, exerciseController)
	fitnessplanHandler := fitnessplan.NewHandler(slogLogger, fitnessPlanController)
	nutritionPlanStorer := datastore15.NewDatastore(conf, slogLogger, client)
	nutritionPlanController := controller14.NewController(conf, slogLogger, provider, s3Storager, emailer, client, kmutexProvider, openAIConnector, nutritionPlanStorer, exerciseStorer, userStorer)
	handler12 := httptransport13.NewHandler(slogLogger, nutritionPlanController)
	googleFitDataPointStorer := datastore16.NewDatastore(conf, slogLogger, client)
	googleFitDataPointController := controller15.NewController(conf, slogLogger, provider, client, cacher, kmutexProvider, organizationStorer, userStorer, googleFitDataPointStorer)
	handler13 := httptransport14.NewHandler(slogLogger, googleFitDataPointController)
	googleCloudPlatformAdapter := google.NewAdapter(conf, slogLogger, client)
	googleFitAppStorer := datastore17.NewDatastore(conf, slogLogger, client)
	dataPointStorer := datastore18.NewDatastore(conf, slogLogger, client)
	aggregatePointStorer := datastore19.NewDatastore(conf, slogLogger, client)
	googleFitAppController := controller16.NewController(conf, slogLogger, provider, client, cacher, kmutexProvider, googleCloudPlatformAdapter, organizationStorer, googleFitDataPointStorer, googleFitAppStorer, userStorer, dataPointStorer, aggregatePointStorer)
	handler14 := httptransport15.NewHandler(slogLogger, googleFitAppController)
	dataPointController := controller17.NewController(conf, slogLogger, provider, client, cacher, kmutexProvider, organizationStorer, userStorer, dataPointStorer)
	handler15 := httptransport16.NewHandler(slogLogger, dataPointController)
	aggregatePointController := controller18.NewController(conf, slogLogger, provider, client, cacher, kmutexProvider, organizationStorer, userStorer, googleFitAppStorer, googleFitDataPointStorer, dataPointStorer, aggregatePointStorer)
	handler16 := httptransport17.NewHandler(slogLogger, aggregatePointController)
	rankPointController := controller19.NewController(conf, slogLogger, provider, client, cacher, kmutexProvider, s3Storager, organizationStorer, userStorer, googleFitAppStorer, googleFitDataPointStorer, dataPointStorer, aggregatePointStorer, rankPointStorer)
	handler17 := httptransport18.NewHandler(slogLogger, rankPointController)
	biometricController := controller20.NewController(conf, slogLogger, provider, client, cacher, kmutexProvider, s3Storager, organizationStorer, userStorer, dataPointStorer, aggregatePointStorer, rankPointStorer)
	handler18 := httptransport19.NewHandler(slogLogger, biometricController)
	trainingProgramStorer := datastore20.NewDatastore(conf, slogLogger, client)
	workoutStorer := datastore21.NewDatastore(conf, slogLogger, client)
	workoutController := controller21.NewController(conf, slogLogger, provider, client, workoutStorer, exerciseStorer)
	trainingprogramController := controller22.NewController(conf, slogLogger, provider, client, trainingProgramStorer, workoutStorer, userStorer, workoutController)
	handler19 := httptransport20.NewHandler(slogLogger, trainingprogramController)
	handler20 := httptransport21.NewHandler(slogLogger, workoutController)
	questionStorer := datastore22.NewDatastore(conf, slogLogger, client)
	questionController := controller23.NewController(conf, slogLogger, provider, client, questionStorer)
	handler21 := httptransport22.NewHandler(slogLogger, questionController)
	fitnessChallengeStorer := datastore23.NewDatastore(conf, slogLogger, client)
	fitnessChallengeController := controller24.NewController(conf, slogLogger, provider, client, fitnessChallengeStorer, workoutStorer, userStorer, workoutController, rankPointController)
	handler22 := httptransport23.NewHandler(slogLogger, fitnessChallengeController)
	inputPortServer := http.NewInputPort(conf, slogLogger, middlewareMiddleware, handler, httptransportHandler, handler2, handler3, handler4, handler5, handler6, handler7, handler8, handler9, handler10, handler11, stripeHandler, fitnessplanHandler, handler12, handler13, handler14, handler15, handler16, handler17, handler18, handler19, handler20, handler21, handler22)
	googleFitDataPointCrontaber := crontab.NewCrontab(slogLogger, kmutexProvider, googleCloudPlatformAdapter, dataPointStorer, googleFitDataPointStorer, googleFitDataPointController, userStorer)
	googleFitAppCrontaber := crontab2.NewCrontab(slogLogger, kmutexProvider, googleCloudPlatformAdapter, dataPointStorer, googleFitDataPointStorer, googleFitAppStorer, googleFitAppController, userStorer)
	crontabInputPortServer := crontab3.NewInputPort(conf, slogLogger, userController, aggregatePointController, rankPointController, googleFitDataPointCrontaber, googleFitAppCrontaber, fitnessPlanStorer, openAIConnector)
	eventSchedulerAdapter := eventscheduler.NewAdapter(slogLogger, universalClient)
	googleFitDataPointScheduler := scheduler.NewScheduler(slogLogger, kmutexProvider, googleCloudPlatformAdapter, eventSchedulerAdapter, googleFitDataPointController)
	googleFitAppScheduler := scheduler2.NewScheduler(slogLogger, kmutexProvider, eventSchedulerAdapter, googleFitAppController)
	aggregatePointScheduler := scheduler3.NewScheduler(slogLogger, kmutexProvider, eventSchedulerAdapter, aggregatePointController)
	rankPointScheduler := scheduler4.NewScheduler(slogLogger, kmutexProvider, eventSchedulerAdapter, rankPointController)
	fitnessPlanScheduler := scheduler5.NewScheduler(slogLogger, kmutexProvider, eventSchedulerAdapter, fitnessPlanStorer, openAIConnector, fitnessPlanController)
	eventschedulerInputPortServer := eventscheduler2.NewInputPort(conf, slogLogger, eventSchedulerAdapter, googleFitDataPointScheduler, googleFitAppScheduler, aggregatePointScheduler, rankPointScheduler, fitnessPlanScheduler)
	application := NewApplication(slogLogger, inputPortServer, crontabInputPortServer, eventschedulerInputPortServer)
	return application
}
