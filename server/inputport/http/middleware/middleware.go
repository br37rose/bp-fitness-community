package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"

	"go.uber.org/ratelimit"

	gateway_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/gateway/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/jwt"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/time"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/provider/uuid"
)

type Middleware interface {
	Attach(fn http.HandlerFunc) http.HandlerFunc
}

type middleware struct {
	Config            *config.Conf
	Logger            *slog.Logger
	Time              time.Provider
	JWT               jwt.Provider
	UUID              uuid.Provider
	GatewayController gateway_c.GatewayController
}

func NewMiddleware(
	configp *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	timep time.Provider,
	jwtp jwt.Provider,
	gatewayController gateway_c.GatewayController,
) Middleware {
	return &middleware{
		Logger:            loggerp,
		UUID:              uuidp,
		Time:              timep,
		JWT:               jwtp,
		GatewayController: gatewayController,
	}
}

// Attach function attaches to HTTP router to apply for every API call.
func (mid *middleware) Attach(fn http.HandlerFunc) http.HandlerFunc {
	// Attach our middleware handlers here. Please note that all our middleware
	// will start from the bottom and proceed upwards.
	// Ex: `URLProcessorMiddleware` will be executed first and
	//     `PostJWTProcessorMiddleware` will be executed last.
	fn = mid.ProtectedURLsMiddleware(fn)
	fn = mid.IPAddressMiddleware(fn)
	fn = mid.PostJWTProcessorMiddleware(fn) // Note: Must be above `JWTProcessorMiddleware`.
	fn = mid.JWTProcessorMiddleware(fn)     // Note: Must be above `PreJWTProcessorMiddleware`.
	fn = mid.PreJWTProcessorMiddleware(fn)  // Note: Must be above `URLProcessorMiddleware`.
	fn = mid.URLProcessorMiddleware(fn)
	fn = mid.RateLimitMiddleware(fn)

	return func(w http.ResponseWriter, r *http.Request) {
		// Flow to the next middleware.
		fn(w, r)
	}
}

func (mid *middleware) RateLimitMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	// Special thanks: https://ubogdan.com/2021/09/ip-based-rate-limit-middleware-using-go.uber.org/ratelimit/
	var lmap sync.Map

	return func(w http.ResponseWriter, r *http.Request) {
		// Open our program's context based on the request and save the
		// slash-seperated array from our URL path.
		ctx := r.Context()

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			mid.Logger.Error("invalid RemoteAddr", slog.Any("err", err))
			http.Error(w, fmt.Sprintf("invalid RemoteAddr: %s", err), http.StatusInternalServerError)
			return
		}

		lif, ok := lmap.Load(host)
		if !ok {
			lif = ratelimit.New(50) // per second.
		}

		lm, ok := lif.(ratelimit.Limiter)
		if !ok {
			mid.Logger.Error("internal middleware error: typecast failed")
			http.Error(w, "internal middleware error: typecast failed", http.StatusInternalServerError)
			return
		}

		lm.Take()
		lmap.Store(host, lm)

		// Flow to the next middleware.
		fn(w, r.WithContext(ctx))
	}
}

// URLProcessorMiddleware Middleware will split the full URL path into slash-sperated parts and save to
// the context to flow downstream in the app for this particular request.
func (mid *middleware) URLProcessorMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Split path into slash-separated parts, for example, path "/foo/bar"
		// gives p==["foo", "bar"] and path "/" gives p==[""]. Our API starts with
		// "/api", as a result we will start the array slice at "1".
		p := strings.Split(r.URL.Path, "/")[1:]

		// log.Println(p) // For debugging purposes only.

		// Open our program's context based on the request and save the
		// slash-seperated array from our URL path.
		ctx := r.Context()
		ctx = context.WithValue(ctx, "url_split", p)

		// Flow to the next middleware.
		fn(w, r.WithContext(ctx))
	}
}

// PreJWTProcessorMiddleware checks to see if we are visiting an unprotected URL and if so then
// let the system know we need to skip authorization handling.
func (mid *middleware) PreJWTProcessorMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Open our program's context based on the request and save the
		// slash-seperated array from our URL path.
		ctx := r.Context()

		// The following code will lookup the URL path in a whitelist and
		// if the visited path matches then we will skip URL protection.
		// We do this because a majority of API endpoints are protected
		// by authorization.

		urlSplit := ctx.Value("url_split").([]string)
		skipPath := map[string]bool{
			"health-check":    true,
			"version":         true,
			"greeting":        true,
			"login":           true,
			"refresh-token":   true,
			"register-member": true,
			"verify":          true,
			"forgot-password": true,
			"password-reset":  true,
			"select-options":  true,
			"public":          true,
			"callback":        true,
		}

		// DEVELOPERS NOTE:
		// If the URL cannot be split into the size then do not skip authorization.
		if len(urlSplit) < 3 {
			// mid.Logger.Warn("Skipping authorization | len less then 3")
			ctx = context.WithValue(ctx, constants.SessionSkipAuthorization, false)
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			return
		}

		// Skip authorization if the URL matches the whitelist else we need to
		// run authorization check.
		if skipPath[urlSplit[2]] {
			// mid.Logger.Warn("Skipping authorization | skipPath found")
			ctx = context.WithValue(ctx, constants.SessionSkipAuthorization, true)
		} else {
			// For debugging purposes only.
			// log.Println("PreJWTProcessorMiddleware | Protected URL detected")
			// log.Println("PreJWTProcessorMiddleware | urlSplit:", urlSplit)
			// log.Println("PreJWTProcessorMiddleware | urlSplit[2]:", urlSplit[2])
			ctx = context.WithValue(ctx, constants.SessionSkipAuthorization, false)
		}

		// Flow to the next middleware.
		fn(w, r.WithContext(ctx))
	}
}

func (mid *middleware) JWTProcessorMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		skipAuthorization, ok := ctx.Value(constants.SessionSkipAuthorization).(bool)
		if ok && skipAuthorization {
			// mid.Logger.Warn("Skipping authorization")
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			return
		}

		// Extract our auth header array.
		reqToken := r.Header.Get("Authorization")

		// For debugging purposes.
		// log.Println("JWTProcessorMiddleware | reqToken:", reqToken)

		// Before running our JWT middleware we need to confirm there is an
		// an `Authorization` header to run our middleware. This is an important
		// step!
		if reqToken != "" && strings.Contains(reqToken, "undefined") == false {

			// Special thanks to "poise" via https://stackoverflow.com/a/44700761
			splitToken := strings.Split(reqToken, "JWT ")
			if len(splitToken) < 2 {
				mid.Logger.Warn("not properly formatted authorization header")
				http.Error(w, "not properly formatted authorization header", http.StatusBadRequest)
				return
			}

			reqToken = splitToken[1]
			// log.Println("JWTProcessorMiddleware | reqToken:", reqToken) // For debugging purposes only.

			sessionID, err := mid.JWT.ProcessJWTToken(reqToken)
			// log.Println("JWTProcessorMiddleware | sessionUUID:", sessionUUID) // For debugging purposes only.

			if err == nil {
				// Update our context to save our JWT token content information.
				ctx = context.WithValue(ctx, constants.SessionIsAuthorized, true)
				ctx = context.WithValue(ctx, constants.SessionID, sessionID)

				// Flow to the next middleware with our JWT token saved.
				fn(w, r.WithContext(ctx))
				return
			}

			// The following code will lookup the URL path in a whitelist and
			// if the visited path matches then we will skip any token errors.
			// We do this because a majority of API endpoints are protected
			// by authorization.

			urlSplit := ctx.Value("url_split").([]string)
			skipPath := map[string]bool{
				"health-check":    true,
				"version":         true,
				"greeting":        true,
				"login":           true,
				"refresh-token":   true,
				"register-member": true,
				"verify":          true,
				"forgot-password": true,
				"password-reset":  true,
				"select-options":  true,
				"public":          true,
				"callback":        true,
			}

			// DEVELOPERS NOTE:
			// If the URL cannot be split into the size we want then skip running
			// this middleware.
			if len(urlSplit) >= 3 {
				if skipPath[urlSplit[2]] {
					mid.Logger.Warn("Skipping expired or error token")
				} else {
					// For debugging purposes only.
					// log.Println("JWTProcessorMiddleware | ProcessJWT | err", err, "for reqToken:", reqToken)
					// log.Println("JWTProcessorMiddleware | ProcessJWT | urlSplit:", urlSplit)
					// log.Println("JWTProcessorMiddleware | ProcessJWT | urlSplit[2]:", urlSplit[2])
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
			}
		}

		// Flow to the next middleware without anything done.
		ctx = context.WithValue(ctx, constants.SessionIsAuthorized, false)
		fn(w, r.WithContext(ctx))
	}
}

func (mid *middleware) PostJWTProcessorMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Skip this middleware if user is on a whitelisted URL path.
		skipAuthorization, ok := ctx.Value(constants.SessionSkipAuthorization).(bool)
		if ok && skipAuthorization {
			// mid.Logger.Warn("Skipping authorization")
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			return
		}

		// Get our authorization information.
		isAuthorized, ok := ctx.Value(constants.SessionIsAuthorized).(bool)
		if ok && isAuthorized {
			sessionID := ctx.Value(constants.SessionID).(string)

			// Lookup our user profile in the session or return 500 error.
			user, err := mid.GatewayController.GetUserBySessionID(ctx, sessionID) //TODO: IMPLEMENT.
			if err != nil {
				mid.Logger.Warn("GetUserBySessionID error", slog.Any("err", err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// If no user was found then that means our session expired and the
			// user needs to login or use the refresh token.
			if user == nil {
				mid.Logger.Warn("Session expired - please log in again")
				http.Error(w, "Session expired - please log in again", http.StatusUnauthorized)
				return
			}

			// // If system administrator disabled the user account then we need
			// // to generate a 403 error letting the user know their account has
			// // been disabled and you cannot access the protected API endpoint.
			// if user.Status == 0 {
			// 	http.Error(w, "Account disabled - please contact admin", http.StatusForbidden)
			// 	return
			// }

			// The following session verification code enforces the user submitted
			// their 2FA code after login to get access to the session. If the
			// user did not verify 2FA code then block their access.
			if user.OTPEnabled && user.OTPVerified && !user.OTPValidated {
				// Check the current URL, split it up and skip this 2FA code
				// validation security if the user is making a call to the
				// `/api/v1/otp/validate` API endpoint. Please note these
				// URLs are dependent to what you are using in the server file.
				urlSplit := ctx.Value("url_split").([]string)

				// Check to see if the user is calling the `/api/v1/otp/validate`.
				if len(urlSplit) > 3 && urlSplit[2] == "otp" && urlSplit[3] == "validate" {
					// We skip validation so proceed in this function. Provide
					// the following log for debugging purposes only.
					mid.Logger.Debug("skipping session requires 2fa validation after login",
						slog.Any("url_split", urlSplit),
					)
				} else {
					// For debuggin purposes only.
					mid.Logger.Warn("session requires 2fa validation after login",
						slog.Any("url_split", urlSplit),
					)

					// Halt proceeding further.
					http.Error(w, "attempting to access a protected endpoint without validating 2fa code after login", http.StatusForbidden)
					return
				}
			}

			// Save our user information to the context.
			// Save our user.
			ctx = context.WithValue(ctx, constants.SessionUser, user)

			// // For debugging purposes only.
			// mid.Logger.Debug("Fetched session record",
			// 	slog.Any("ID", user.ID),
			// 	slog.String("SessionID", sessionID),
			// 	slog.String("Name", user.Name),
			// 	slog.String("FirstName", user.FirstName),
			// 	slog.String("Email", user.Email),
			// 	slog.Any("OrganizationID", user.OrganizationID),
			// 	slog.String("OrganizationName", user.OrganizationName))

			// Save individual pieces of the user profile.
			ctx = context.WithValue(ctx, constants.SessionID, sessionID)
			ctx = context.WithValue(ctx, constants.SessionUserID, user.ID)
			ctx = context.WithValue(ctx, constants.SessionUserRole, user.Role)
			ctx = context.WithValue(ctx, constants.SessionUserName, user.Name)
			ctx = context.WithValue(ctx, constants.SessionUserLexicalName, user.LexicalName)
			ctx = context.WithValue(ctx, constants.SessionUserFirstName, user.FirstName)
			ctx = context.WithValue(ctx, constants.SessionUserLastName, user.LastName)
			ctx = context.WithValue(ctx, constants.SessionUserOrganizationID, user.OrganizationID)
			ctx = context.WithValue(ctx, constants.SessionUserOrganizationName, user.OrganizationName)
			ctx = context.WithValue(ctx, constants.SessionUserOTPValidated, user.OTPValidated)
		}

		fn(w, r.WithContext(ctx))
	}
}

func (mid *middleware) IPAddressMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the IPAddress. Code taken from: https://stackoverflow.com/a/55738279
		IPAddress := r.Header.Get("X-Real-Ip")
		if IPAddress == "" {
			IPAddress = r.Header.Get("X-Forwarded-For")
		}
		if IPAddress == "" {
			IPAddress = r.RemoteAddr
		}

		// Save our IP address to the context.
		ctx := r.Context()
		ctx = context.WithValue(ctx, constants.SessionIPAddress, IPAddress)
		fn(w, r.WithContext(ctx)) // Flow to the next middleware.
	}
}

// ProtectedURLsMiddleware The purpose of this middleware is to return a `401 unauthorized` error if
// the user is not authorized when visiting a protected URL.
func (mid *middleware) ProtectedURLsMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Skip this middleware if user is on a whitelisted URL path.
		skipAuthorization, ok := ctx.Value(constants.SessionSkipAuthorization).(bool)
		if ok && skipAuthorization {
			// mid.Logger.Warn("Skipping authorization")
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			return
		}

		// The following code will lookup the URL path in a whitelist and
		// if the visited path matches then we will skip URL protection.
		// We do this because a majority of API endpoints are protected
		// by authorization.

		urlSplit := ctx.Value("url_split").([]string)
		skipPath := map[string]bool{
			"health-check":    true,
			"version":         true,
			"greeting":        true,
			"login":           true,
			"register-member": true,
			"refresh-token":   true,
			"verify":          true,
			"forgot-password": true,
			"password-reset":  true,
			"select-options":  true,
			"public":          true,
			"callback":        true,
		}

		// DEVELOPERS NOTE:
		// If the URL cannot be split into the size we want then skip running
		// this middleware.
		if len(urlSplit) < 3 {
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			return
		}

		if skipPath[urlSplit[2]] {
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
		} else {
			// Get our authorization information.
			isAuthorized, ok := ctx.Value(constants.SessionIsAuthorized).(bool)

			// Either accept continuing execution or return 401 error.
			if ok && isAuthorized {
				fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			} else {
				mid.Logger.Warn("unauthorized api call")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}
	}
}
