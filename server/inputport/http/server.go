package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rs/cors"

	aggregatepoint_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/httptransport"
	attachment_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/httptransport"
	biometric_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/biometric/httptransport"
	datapoint_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/datapoint/httptransport"
	exercise_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/exercise/httptransport"
	gateway_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/gateway/httptransport"
	googlefitapp_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitapp/httptransport"
	googlefitdp_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/googlefitdatapoint/httptransport"
	invoice_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/invoice/httptransport"
	member_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/member/httptransport"
	nutritionplan_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/nutritionplan/httptransport"
	offer_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/offer/httptransport"
	organization_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/organization/httptransport"
	strpp "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/paymentprocessor/httptransport/stripe"
	rankpoint_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/httptransport"
	tag_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/tag/httptransport"
	tp_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/httptransport"
	user_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/httptransport"
	videocategory_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocategory/httptransport"
	videocollection_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocollection/httptransport"
	videocontent_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/videocontent/httptransport"
	w_http "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/httptransport"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/http/fitnessplan"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/inputport/http/middleware"
)

type InputPortServer interface {
	Run()
	Shutdown()
}

type httpInputPort struct {
	Config                   *config.Conf
	Logger                   *slog.Logger
	Server                   *http.Server
	Middleware               middleware.Middleware
	Gateway                  *gateway_http.Handler
	User                     *user_http.Handler
	Organization             *organization_http.Handler
	Tag                      *tag_http.Handler
	Exercise                 *exercise_http.Handler
	Member                   *member_http.Handler
	Attachment               *attachment_http.Handler
	VideoCategory            *videocategory_http.Handler
	VideoCollection          *videocollection_http.Handler
	VideoContent             *videocontent_http.Handler
	Offer                    *offer_http.Handler
	Invoice                  *invoice_http.Handler
	StripePaymentProcessor   *strpp.Handler
	FitnessPlan              *fitnessplan.Handler
	NutritionPlan            *nutritionplan_http.Handler
	GoogleFitDataPointStorer *googlefitdp_http.Handler
	GoogleFitApp             *googlefitapp_http.Handler
	DataPoint                *datapoint_http.Handler
	AggregatePoint           *aggregatepoint_http.Handler
	RankPoint                *rankpoint_http.Handler
	Biometric                *biometric_http.Handler
	TrainingProgram          *tp_http.Handler
	Workout                  *w_http.Handler
}

func NewInputPort(
	configp *config.Conf,
	loggerp *slog.Logger,
	mid middleware.Middleware,
	gh *gateway_http.Handler,
	cu *user_http.Handler,
	org *organization_http.Handler,
	tag *tag_http.Handler,
	exc *exercise_http.Handler,
	mem *member_http.Handler,
	att *attachment_http.Handler,
	vc *videocategory_http.Handler,
	vcol *videocollection_http.Handler,
	vcon *videocontent_http.Handler,
	off *offer_http.Handler,
	inv *invoice_http.Handler,
	strpp *strpp.Handler,
	ff *fitnessplan.Handler,
	np *nutritionplan_http.Handler,
	gfdp *googlefitdp_http.Handler,
	gfa *googlefitapp_http.Handler,
	dp *datapoint_http.Handler,
	ap *aggregatepoint_http.Handler,
	rp *rankpoint_http.Handler,
	bio *biometric_http.Handler,
	tp *tp_http.Handler,
	w *w_http.Handler,
) InputPortServer {
	// Initialize the ServeMux.
	mux := http.NewServeMux()

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation via `https://github.com/rs/cors` for more options.
	handler := cors.AllowAll().Handler(mux)

	// Bind the HTTP server to the assigned address and port.
	addr := fmt.Sprintf("%s:%s", configp.AppServer.IP, configp.AppServer.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Create our HTTP server controller.
	p := &httpInputPort{
		Config:                   configp,
		Logger:                   loggerp,
		Middleware:               mid,
		Gateway:                  gh,
		User:                     cu,
		Organization:             org,
		Tag:                      tag,
		Exercise:                 exc,
		Member:                   mem,
		Attachment:               att,
		VideoCategory:            vc,
		VideoCollection:          vcol,
		VideoContent:             vcon,
		Server:                   srv,
		Offer:                    off,
		Invoice:                  inv,
		StripePaymentProcessor:   strpp,
		FitnessPlan:              ff,
		NutritionPlan:            np,
		GoogleFitApp:             gfa,
		GoogleFitDataPointStorer: gfdp,
		DataPoint:                dp,
		AggregatePoint:           ap,
		RankPoint:                rp,
		Biometric:                bio,
		TrainingProgram:          tp,
		Workout:                  w,
	}

	// Attach the HTTP server controller to the ServerMux.
	mux.HandleFunc("/", mid.Attach(p.HandleRequests))

	return p
}

func (port *httpInputPort) Run() {
	port.Logger.Info("HTTP server running")
	if err := port.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		port.Logger.Error("listen failed", slog.Any("error", err))

		// DEVELOPERS NOTE: We terminate app here b/c dependency injection not allowed to fail, so fail here at startup of app.
		panic("failed running")
	}
}

func (port *httpInputPort) Shutdown() {
	port.Logger.Info("HTTP server shutdown")
}

func (port *httpInputPort) HandleRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get our URL paths which are slash-seperated.
	ctx := r.Context()
	p := ctx.Value("url_split").([]string)
	n := len(p)
	port.Logger.Debug("Handling request",
		slog.Int("n", n),
		slog.String("m", r.Method),
		slog.Any("p", p),
	)

	switch {
	// --- GATEWAY & PROFILE & DASHBOARD --- //
	case n == 3 && p[1] == "v1" && p[2] == "health-check" && r.Method == http.MethodGet:
		port.Gateway.HealthCheck(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "version" && r.Method == http.MethodGet:
		port.Gateway.Version(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "greeting" && r.Method == http.MethodPost:
		port.Gateway.Greet(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "login" && r.Method == http.MethodPost:
		port.Gateway.Login(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "register-member" && r.Method == http.MethodPost:
		port.Gateway.MemberRegister(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "refresh-token" && r.Method == http.MethodPost:
		port.Gateway.RefreshToken(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "verify" && r.Method == http.MethodPost:
		port.Gateway.Verify(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "logout" && r.Method == http.MethodPost:
		port.Gateway.Logout(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "account" && r.Method == http.MethodGet:
		port.Gateway.Account(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "account" && r.Method == http.MethodPut:
		port.Gateway.AccountUpdate(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "account" && p[3] == "change-password" && r.Method == http.MethodPut:
		port.Gateway.AccountChangePassword(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "account" && p[3] == "invoices" && r.Method == http.MethodGet:
		port.Gateway.AccountListLatestInvoices(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "forgot-password" && r.Method == http.MethodPost:
		port.Gateway.ForgotPassword(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "password-reset" && r.Method == http.MethodPost:
		port.Gateway.PasswordReset(w, r)
	case n == 5 && p[1] == "v1" && p[2] == "account" && p[3] == "operation" && p[4] == "avatar" && r.Method == http.MethodPost:
		port.Gateway.OperationAvatar(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "otp" && p[3] == "generate" && r.Method == http.MethodPost:
		port.Gateway.GenerateOTP(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "otp" && p[3] == "generate-qr-code" && r.Method == http.MethodPost:
		port.Gateway.GenerateOTPAndQRCodePNGImage(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "otp" && p[3] == "verify" && r.Method == http.MethodPost:
		port.Gateway.VerifyOTP(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "otp" && p[3] == "validate" && r.Method == http.MethodPost:
		port.Gateway.ValidateOTP(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "otp" && p[3] == "disable" && r.Method == http.MethodPost:
		port.Gateway.DisableOTP(w, r)

	// --- ORGANIZATION --- //
	case n == 3 && p[1] == "v1" && p[2] == "organizations" && r.Method == http.MethodGet:
		port.Organization.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "organizations" && r.Method == http.MethodPost:
		port.Organization.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "organization" && r.Method == http.MethodGet:
		port.Organization.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "organization" && r.Method == http.MethodPut:
		port.Organization.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "organization" && r.Method == http.MethodDelete:
		port.Organization.DeleteByID(w, r, p[3])

	// --- MEMBERS --- //
	case n == 3 && p[1] == "v1" && p[2] == "members" && r.Method == http.MethodGet:
		port.Member.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "members" && r.Method == http.MethodPost:
		port.Member.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "member" && r.Method == http.MethodGet:
		port.Member.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "member" && r.Method == http.MethodPut:
		port.Member.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "member" && r.Method == http.MethodDelete:
		port.Member.DeleteByID(w, r, p[3])
	case n == 5 && p[1] == "v1" && p[2] == "members" && p[3] == "operation" && p[4] == "create-comment" && r.Method == http.MethodPost:
		port.Member.OperationCreateComment(w, r)
	case n == 5 && p[1] == "v1" && p[2] == "members" && p[3] == "operation" && p[4] == "avatar" && r.Method == http.MethodPost:
		port.Member.OperationAvatar(w, r)
	case n == 5 && p[1] == "v1" && p[2] == "select-options" && p[4] == "members" && r.Method == http.MethodGet:
		port.Member.ListAsSelectOptionsByOrganization(w, r, p[3])

	// --- EXERCISES --- //
	case n == 3 && p[1] == "v1" && p[2] == "exercises" && r.Method == http.MethodGet:
		port.Exercise.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "exercises" && r.Method == http.MethodPost:
		port.Exercise.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "exercise" && r.Method == http.MethodGet:
		port.Exercise.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "exercise" && r.Method == http.MethodPut:
		port.Exercise.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "exercise" && r.Method == http.MethodDelete:
		port.Exercise.DeleteByID(w, r, p[3])

	// --- ATTACHMENTS --- //
	case n == 3 && p[1] == "v1" && p[2] == "attachments" && r.Method == http.MethodGet:
		port.Attachment.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "attachments" && r.Method == http.MethodPost:
		port.Attachment.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "attachment" && r.Method == http.MethodGet:
		port.Attachment.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "attachment" && r.Method == http.MethodPut:
		port.Attachment.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "attachment" && r.Method == http.MethodDelete:
		port.Attachment.DeleteByID(w, r, p[3])

	// --- VIDEO CATEGORY --- //
	case n == 3 && p[1] == "v1" && p[2] == "video-categories" && r.Method == http.MethodGet:
		port.VideoCategory.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "video-categories" && r.Method == http.MethodPost:
		port.VideoCategory.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "video-category" && r.Method == http.MethodGet:
		port.VideoCategory.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "video-category" && r.Method == http.MethodPut:
		port.VideoCategory.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "video-category" && r.Method == http.MethodDelete:
		port.VideoCategory.DeleteByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "video-categories" && p[3] == "select-options" && r.Method == http.MethodGet:
		port.VideoCategory.ListAsSelectOptionByFilter(w, r)

	// --- VIDEO COLLECTION --- //
	case n == 3 && p[1] == "v1" && p[2] == "video-collections" && r.Method == http.MethodGet:
		port.VideoCollection.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "video-collections" && r.Method == http.MethodPost:
		port.VideoCollection.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "video-collection" && r.Method == http.MethodGet:
		port.VideoCollection.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "video-collection" && r.Method == http.MethodPut:
		port.VideoCollection.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "video-collection" && r.Method == http.MethodDelete:
		port.VideoCollection.DeleteByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "video-collections" && p[3] == "select-options" && r.Method == http.MethodGet:
		port.VideoCollection.ListAsSelectOptionByFilter(w, r)

	// --- VIDEO CONTENT --- //
	case n == 3 && p[1] == "v1" && p[2] == "video-contents" && r.Method == http.MethodGet:
		port.VideoContent.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "video-contents" && r.Method == http.MethodPost:
		port.VideoContent.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "video-content" && r.Method == http.MethodGet:
		port.VideoContent.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "video-content" && r.Method == http.MethodPut:
		port.VideoContent.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "video-content" && r.Method == http.MethodDelete:
		port.VideoContent.DeleteByID(w, r, p[3])

	// --- PAYMENT PROCESSOR --- //
	case n == 4 && p[1] == "v1" && p[2] == "stripe" && p[3] == "create-checkout-session" && r.Method == http.MethodPost:
		port.StripePaymentProcessor.CreateStripeCheckoutSession(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "stripe" && p[3] == "complete-checkout-session" && r.Method == http.MethodGet:
		port.StripePaymentProcessor.CompleteStripeCheckoutSession(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "stripe" && p[3] == "cancel-subscription" && r.Method == http.MethodPost:
		port.StripePaymentProcessor.CancelStripeSubscription(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "stripe" && p[3] == "invoices" && r.Method == http.MethodGet:
		port.StripePaymentProcessor.ListLatestStripeInvoices(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "public" && p[3] == "stripe-webhook":
		port.StripePaymentProcessor.Webhook(w, r)

	// --- OFFERS --- //
	case n == 3 && p[1] == "v1" && p[2] == "offers" && r.Method == http.MethodGet:
		port.Offer.List(w, r)
	// case n == 3 && p[1] == "v1" && p[2] == "offers" && r.Method == http.MethodPost:
	// 	port.Offer.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "offer" && r.Method == http.MethodGet:
		port.Offer.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "offer" && r.Method == http.MethodPut:
		port.Offer.UpdateByID(w, r, p[3])
	// case n == 4 && p[1] == "v1" && p[2] == "offer" && r.Method == http.MethodDelete:
	// 	port.Offer.DeleteByID(w, r, p[3])
	// case n == 5 && p[1] == "v1" && p[2] == "offer" && p[3] == "operation" && p[4] == "create-comment" && r.Method == http.MethodPost:
	// 	port.Offer.OperationCreateComment(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "offers" && p[3] == "select-options" && r.Method == http.MethodGet:
		port.Offer.ListAsSelectOptions(w, r)

	// --- FITNESS PLAN --- //
	case n == 3 && p[1] == "v1" && p[2] == "fitness-plans" && r.Method == http.MethodGet:
		port.FitnessPlan.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "fitness-plans" && r.Method == http.MethodPost:
		port.FitnessPlan.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "fitness-plan" && r.Method == http.MethodGet:
		port.FitnessPlan.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "fitness-plan" && r.Method == http.MethodPut:
		port.FitnessPlan.UpdateByID(w, r, p[3])

	// --- NUTRITION PLAN --- //
	case n == 3 && p[1] == "v1" && p[2] == "nutrition-plans" && r.Method == http.MethodGet:
		port.NutritionPlan.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "nutrition-plans" && r.Method == http.MethodPost:
		port.NutritionPlan.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "nutrition-plan" && r.Method == http.MethodGet:
		port.NutritionPlan.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "nutrition-plan" && r.Method == http.MethodPut:
		port.NutritionPlan.UpdateByID(w, r, p[3])

		// --- GOOGLE FIT --- //
	case n == 3 && p[1] == "v1" && p[2] == "google-login":
		port.GoogleFitApp.GetGoogleLoginURL(w, r)
	case n == 5 && p[1] == "v1" && p[2] == "callback" && p[3] == "google" && p[4] == "auth":
		port.GoogleFitApp.GoogleCallback(w, r)

	// case n == 4 && p[1] == "v1" && p[2] == "googlefit" && p[3] == "simulators":
	// 	port.FitBitApp.CreateSimulator(w, r)
	// case n == 5 && p[1] == "v1" && p[2] == "callback" && p[3] == "googlefit" && p[4] == "auth":
	// 	port.FitBitApp.Auth(w, r)
	// case n == 5 && p[1] == "v1" && p[2] == "callback" && p[3] == "googlefit" && p[4] == "subscriber":
	// 	port.FitBitApp.Subscriber(w, r)

	// --- DATA POINT --- //
	case n == 3 && p[1] == "v1" && p[2] == "data-points" && r.Method == http.MethodGet:
		port.DataPoint.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "google-fit-data-points" && r.Method == http.MethodGet:
		port.GoogleFitDataPointStorer.List(w, r)

	// --- AGGREGATE POINT --- //
	case n == 4 && p[1] == "v1" && p[2] == "aggregate-points" && p[3] == "summary" && r.Method == http.MethodGet:
		port.AggregatePoint.GetSummary(w, r)

	// --- RANK POINT --- //
	case n == 3 && p[1] == "v1" && p[2] == "rank-points" && r.Method == http.MethodGet:
		port.RankPoint.List(w, r)

	// --- BIOMETRIC --- //
	case n == 3 && p[1] == "v1" && p[2] == "leaderboard" && r.Method == http.MethodGet: // Deprecated URL.
		port.Biometric.Leaderboard(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "biometrics" && p[3] == "leaderboard" && r.Method == http.MethodGet:
		port.Biometric.Leaderboard(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "biometrics" && p[3] == "summary" && r.Method == http.MethodGet:
		port.Biometric.GetSummary(w, r)

		// --- TRAINING PROGRAM --- //
	case n == 3 && p[1] == "v1" && p[2] == "training-program" && r.Method == http.MethodPost:
		port.TrainingProgram.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "training-program" && r.Method == http.MethodGet:
		port.TrainingProgram.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "training-program" && r.Method == http.MethodDelete:
		port.TrainingProgram.DeleteByID(w, r, p[3])
	case n == 3 && p[1] == "v1" && p[2] == "training-program" && r.Method == http.MethodGet:
		port.TrainingProgram.List(w, r)

	// --- TAG --- //
	// case n == 3 && p[1] == "v1" && p[2] == "tags" && r.Method == http.MethodGet:
	// 	port.Tag.List(w, r)
	// case n == 3 && p[1] == "v1" && p[2] == "tags" && r.Method == http.MethodPost:
	// 	port.Tag.Create(w, r)
	// case n == 4 && p[1] == "v1" && p[2] == "tag" && r.Method == http.MethodGet:
	// 	port.Tag.GetByID(w, r, p[3])
	// case n == 4 && p[1] == "v1" && p[2] == "tag" && r.Method == http.MethodPut:
	// 	port.Tag.UpdateByID(w, r, p[3])
	// case n == 4 && p[1] == "v1" && p[2] == "tag" && r.Method == http.MethodDelete:
	// 	port.Tag.DeleteByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "tags" && p[3] == "select-options" && r.Method == http.MethodGet:
		port.Tag.ListAsSelectOptionByFilter(w, r)

	// --- Workout --- //
	case n == 3 && p[1] == "v1" && p[2] == "workouts" && r.Method == http.MethodGet:
		port.Workout.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "workouts" && r.Method == http.MethodPost:
		port.Workout.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "workouts" && r.Method == http.MethodGet:
		port.Workout.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "workouts" && r.Method == http.MethodPut:
		port.Workout.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "workouts" && r.Method == http.MethodDelete:
		port.Workout.DeleteByID(w, r, p[3])

	// --- CATCH ALL: D.N.E. ---
	default:
		http.NotFound(w, r)
	}
}
