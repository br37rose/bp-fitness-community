//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cache/mongodbcache"
	gcp_a "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/emailer/mailgun"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/openai"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/paymentprocessor/stripe"
	s3_storage "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/storage/s3"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/templatedemailer"
	ap_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/controller"
	ap_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	ap_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/httptransport"
	attachment_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/controller"
	attachment_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/datastore"
	attachment_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/httptransport"
	bio_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/biometric/controller"
	bio_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/biometric/httptransport"
	dp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/controller"
	dp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/datastore"
	dp_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/httptransport"
	equip_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/equipment/datastore"
	eventlog_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/eventlog/datastore"
	exercise_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/controller"
	exercise_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/datastore"
	exercise_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/httptransport"
	fitnessplan_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/controller"
	fitnessplan_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	gateway_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/gateway/controller"
	gateway_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/gateway/httptransport"
	googlefitapp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/controller"
	googlefitapp_cron "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/crontab"
	googlefitapp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/datastore"
	googlefitapp_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/httptransport"
	googlefitdp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/controller"
	googlefitdp_cron "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/crontab"
	googlefitdp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/datastore"
	googlefitdp_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/httptransport"
	inv_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/invoice/controller"
	inv_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/invoice/datastore"
	inv_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/invoice/httptransport"
	member_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/member/controller"
	member_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/member/httptransport"
	nutritionplan_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/controller"
	nutritionplan_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/datastore"
	nutritionplan_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/httptransport"
	off_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/controller"
	off_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/datastore"
	off_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/httptransport"
	organization_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/controller"
	organization_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/datastore"
	organization_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/httptransport"
	strpayproc_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/paymentprocessor/controller/stripe"
	strpayproc_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/paymentprocessor/httptransport/stripe"
	q_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/question/controller"
	q_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/question/datastore"
	q_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/question/httptransport"
	rp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/controller"
	rp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	rp_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/httptransport"
	tag_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/tag/controller"
	tag_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/tag/datastore"
	tag_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/tag/httptransport"
	tp_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/controller"
	tp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/datastore"
	tp_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/httptransport"
	user_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/controller"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	user_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/httptransport"
	vcat_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocategory/controller"
	vcat_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocategory/datastore"
	vc_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocategory/httptransport"
	vcol_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocollection/controller"
	vcol_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocollection/datastore"
	vcol_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocollection/httptransport"
	vcon_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocontent/controller"
	vcon_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocontent/datastore"
	vcon_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocontent/httptransport"
	w_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/controller"
	w_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/datastore"
	w_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/httptransport"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/crontab"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/http"
	fitnessplan_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/http/fitnessplan"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/http/middleware"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/jwt"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/kmutex"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/logger"
	mongodb_p "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/mongodb"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/password"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/time"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

func InitializeEvent() Application {
	// Our application is dependent on the following Golang packages. We need to
	// provide them to Google wire so it can sort out the dependency injection
	// at compile time.
	wire.Build(
		config.New,
		uuid.NewProvider,
		time.NewProvider,
		logger.NewProvider,
		jwt.NewProvider,
		kmutex.NewProvider,
		mailgun.NewEmailer,
		templatedemailer.NewTemplatedEmailer,
		password.NewProvider,
		mongodb_p.NewProvider,
		mongodbcache.NewCache,
		gcp_a.NewAdapter,
		openai.NewOpenAIConnector,
		s3_storage.NewStorage,
		stripe.NewPaymentProcessor,
		eventlog_s.NewDatastore,
		user_s.NewDatastore,
		user_c.NewController,
		attachment_s.NewDatastore,
		attachment_c.NewController,
		off_s.NewDatastore,
		off_c.NewController,
		member_c.NewController,
		organization_s.NewDatastore,
		organization_c.NewController,
		equip_s.NewDatastore,
		tag_s.NewDatastore,
		tag_c.NewController,
		exercise_s.NewDatastore,
		exercise_c.NewController,
		vcat_s.NewDatastore,
		vcat_c.NewController,
		vcol_s.NewDatastore,
		vcol_c.NewController,
		vcon_s.NewDatastore,
		vcon_c.NewController,
		inv_s.NewDatastore,
		inv_c.NewController,
		fitnessplan_s.NewDatastore,
		fitnessplan_c.NewController,
		nutritionplan_s.NewDatastore,
		nutritionplan_c.NewController,
		googlefitdp_s.NewDatastore,
		googlefitdp_c.NewController,
		googlefitdp_cron.NewCrontab,
		googlefitdp_http.NewHandler,
		googlefitapp_s.NewDatastore,
		googlefitapp_c.NewController,
		strpayproc_c.NewController,
		dp_s.NewDatastore,
		dp_c.NewController,
		ap_s.NewDatastore,
		ap_c.NewController,
		rp_s.NewDatastore,
		rp_c.NewController,
		bio_c.NewController,
		gateway_c.NewController,
		gateway_http.NewHandler,
		user_http.NewHandler,
		attachment_http.NewHandler,
		off_http.NewHandler,
		organization_http.NewHandler,
		tag_http.NewHandler,
		exercise_http.NewHandler,
		member_http.NewHandler,
		vc_http.NewHandler,
		vcol_http.NewHandler,
		vcon_http.NewHandler,
		inv_http.NewHandler,
		fitnessplan_http.NewHandler,
		nutritionplan_http.NewHandler,
		googlefitapp_http.NewHandler,
		googlefitapp_cron.NewCrontab,
		dp_http.NewHandler,
		ap_http.NewHandler,
		rp_http.NewHandler,
		bio_http.NewHandler,
		// payproc_http.NewHandler,
		strpayproc_http.NewHandler,
		middleware.NewMiddleware,
		crontab.NewInputPort,
		http.NewInputPort,
		tp_s.NewDatastore,
		tp_c.NewController,
		tp_http.NewHandler,
		w_s.NewDatastore,
		w_c.NewController,
		w_http.NewHandler,
		q_s.NewDatastore,
		q_c.NewController,
		q_http.NewHandler,
		NewApplication)
	return Application{}
}
