// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cache/mongodbcache"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/emailer/mailgun"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/openai"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/paymentprocessor/stripe"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/storage/s3"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/templatedemailer"
	controller18 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/controller"
	datastore18 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	httptransport16 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/httptransport"
	controller7 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/controller"
	datastore5 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/datastore"
	httptransport7 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/httptransport"
	controller20 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/biometric/controller"
	httptransport18 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/biometric/httptransport"
	controller17 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/controller"
	datastore17 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/datastore"
	httptransport15 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/httptransport"
	datastore6 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/equipment/datastore"
	datastore13 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/eventlog/datastore"
	controller5 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/controller"
	datastore7 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	httptransport5 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/httptransport"
	controller16 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitapp/controller"
	datastore19 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitapp/datastore"
	datastore20 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitdatum/datastore"
	controller13 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/controller"
	datastore14 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/gateway/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/gateway/httptransport"
	controller15 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/crontab"
	datastore16 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	httptransport14 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/httptransport"
	datastore21 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/datastore"
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
	controller19 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	httptransport17 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/httptransport"
	controller4 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/tag/controller"
	datastore4 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/tag/datastore"
	httptransport4 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/tag/httptransport"
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
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	crontab2 "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/crontab"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/http"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/http/fitbitapp"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/http/fitnessplan"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/http/middleware"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/jwt"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/logger"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/mongodb"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/password"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/time"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

import (
	_ "github.com/google/wire"
	_ "go.uber.org/automaxprocs"
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
	emailer := mailgun.NewEmailer(conf, slogLogger, provider)
	templatedEmailer := templatedemailer.NewTemplatedEmailer(conf, slogLogger, provider, emailer)
	paymentProcessor := stripe.NewPaymentProcessor(conf, slogLogger, provider)
	rankPointStorer := datastore.NewDatastore(conf, slogLogger, client)
	userStorer := datastore2.NewDatastore(conf, slogLogger, client)
	organizationStorer := datastore3.NewDatastore(conf, slogLogger, client)
	gatewayController := controller.NewController(conf, slogLogger, provider, jwtProvider, passwordProvider, cacher, s3Storager, client, kmutexProvider, templatedEmailer, paymentProcessor, rankPointStorer, userStorer, organizationStorer)
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
	fitnessPlanController := controller13.NewController(conf, slogLogger, provider, s3Storager, emailer, client, kmutexProvider, openAIConnector, fitnessPlanStorer, exerciseStorer, userStorer)
	fitnessplanHandler := fitnessplan.NewHandler(slogLogger, fitnessPlanController)
	nutritionPlanStorer := datastore15.NewDatastore(conf, slogLogger, client)
	nutritionPlanController := controller14.NewController(conf, slogLogger, provider, s3Storager, emailer, client, kmutexProvider, openAIConnector, nutritionPlanStorer, exerciseStorer, userStorer)
	handler12 := httptransport13.NewHandler(slogLogger, nutritionPlanController)
	googleCloudPlatformAdapter := google.NewAdapter(conf, slogLogger, client)
	googleFitAppStorer := datastore16.NewDatastore(conf, slogLogger, client)
	dataPointStorer := datastore17.NewDatastore(conf, slogLogger, client)
	aggregatePointStorer := datastore18.NewDatastore(conf, slogLogger, client)
	googleFitAppController := controller15.NewController(conf, slogLogger, provider, client, cacher, kmutexProvider, googleCloudPlatformAdapter, organizationStorer, googleFitAppStorer, userStorer, dataPointStorer, aggregatePointStorer)
	handler13 := httptransport14.NewHandler(slogLogger, googleFitAppController)
	fitBitAppStorer := datastore19.NewDatastore(conf, slogLogger, client)
	fitBitDatumStorer := datastore20.NewDatastore(conf, slogLogger, client)
	fitBitAppController := controller16.NewController(conf, slogLogger, provider, client, cacher, kmutexProvider, organizationStorer, fitBitAppStorer, userStorer, fitBitDatumStorer, dataPointStorer, aggregatePointStorer)
	fitbitappHandler := fitbitapp.NewHandler(slogLogger, fitBitAppController)
	dataPointController := controller17.NewController(conf, slogLogger, provider, client, cacher, kmutexProvider, organizationStorer, fitBitAppStorer, userStorer, fitBitDatumStorer, dataPointStorer)
	handler14 := httptransport15.NewHandler(slogLogger, dataPointController)
	aggregatePointController := controller18.NewController(conf, slogLogger, provider, client, cacher, kmutexProvider, organizationStorer, fitBitAppStorer, userStorer, fitBitDatumStorer, dataPointStorer, aggregatePointStorer)
	handler15 := httptransport16.NewHandler(slogLogger, aggregatePointController)
	rankPointController := controller19.NewController(conf, slogLogger, provider, client, cacher, kmutexProvider, s3Storager, organizationStorer, fitBitAppStorer, userStorer, fitBitDatumStorer, dataPointStorer, aggregatePointStorer, rankPointStorer)
	handler16 := httptransport17.NewHandler(slogLogger, rankPointController)
	biometricController := controller20.NewController(conf, slogLogger, provider, client, cacher, kmutexProvider, s3Storager, organizationStorer, fitBitAppStorer, userStorer, fitBitDatumStorer, dataPointStorer, aggregatePointStorer, rankPointStorer)
	handler17 := httptransport18.NewHandler(slogLogger, biometricController)
	inputPortServer := http.NewInputPort(conf, slogLogger, middlewareMiddleware, handler, httptransportHandler, handler2, handler3, handler4, handler5, handler6, handler7, handler8, handler9, handler10, handler11, stripeHandler, fitnessplanHandler, handler12, handler13, fitbitappHandler, handler14, handler15, handler16, handler17)
	googleFitDataPointStorer := datastore21.NewDatastore(conf, slogLogger, client)
	googleFitAppCrontaber := crontab.NewCrontab(slogLogger, kmutexProvider, googleCloudPlatformAdapter, googleFitDataPointStorer, googleFitAppStorer, googleFitAppController, userStorer)
	crontabInputPortServer := crontab2.NewInputPort(conf, slogLogger, userController, fitBitAppController, aggregatePointController, rankPointController, googleFitAppCrontaber)
	application := NewApplication(slogLogger, inputPortServer, crontabInputPortServer)
	return application
}