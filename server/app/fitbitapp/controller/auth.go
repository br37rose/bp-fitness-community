package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	fitbitapp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitbitapp/datastore"
	u_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/config/constants"
)

type AuthResponse struct {
	URL string `bson:"url" json:"url"`
}

func (c *FitBitAppControllerImpl) Auth(ctx context.Context, code string) (*AuthResponse, error) {
	c.Logger.Debug("fitbit authentication beginning",
		slog.Any("code", code))
	defer c.Logger.Debug("fitbit authentication finished",
		slog.Any("code", code))

	if code == "" {
		err := errors.New("code is missing")
		c.Logger.Error("failed starting auth",
			slog.Any("error", err))
		return nil, err
	}

	// DEVELOPERS NOTE:
	// For more information on how to process authorization then please read the following:
	// ---> https://dev.fitbit.com/build/reference/web-api/developer-guide/authorization/

	// For convinence.
	config := c.Config

	// Get our client ID and client secret.
	clientID := config.FitBitApp.ClientID
	clientSecret := config.FitBitApp.ClientSecret

	// STEP 1
	// Generate the bearer.
	bearerValueRaw := fmt.Sprintf("%v:%v", clientID, clientSecret)
	bearerValue := Base64EncodeStrippedFromString(bearerValueRaw)
	bearer := fmt.Sprintf("Basic %s", string(bearerValue))

	// STEP 2: Exchange Code For Access Token

	userIDs := make([]primitive.ObjectID, 0, len(c.CodeVerifierMap))
	for k := range c.CodeVerifierMap {
		userIDs = append(userIDs, k)
	}

	for _, userID := range userIDs {
		codeVerifier := c.CodeVerifierMap[userID]
		if err := c.attemptAuthorizationForKey(ctx, userID, clientID, code, codeVerifier, bearer); err != nil {
			c.Logger.Error("fitbit authentication failed attempt authorization",
				slog.Any("code", code),
				slog.Any("error", err))
			return nil, err
		}
	}
	return &AuthResponse{URL: c.Config.FitBitApp.RegistrationSuccessRedirectURL}, nil
}

func (c *FitBitAppControllerImpl) attemptAuthorizationForKey(
	ctx context.Context,
	userID primitive.ObjectID,
	clientID string,
	code string,
	codeVerifier string,
	bearer string) error {

	url1 := fmt.Sprintf("%s?client_id=%s&code=%s&code_verifier=%s&grant_type=%s",
		constants.FitBitExchangeURL,
		clientID,
		code,
		codeVerifier,
		"authorization_code",
	)

	c.Logger.Debug("calling fitbit remote api for authorization",
		slog.Any("client_id", clientID),
		slog.Any("code", code),
		slog.Any("code_verifier", codeVerifier),
		slog.Any("grant_type", "authorization_code"),
		slog.Any("url", url1))

	req, err := http.NewRequest("POST", url1, nil)
	if err != nil {
		c.Logger.Error("failed setting up fitbit remote api request",
			slog.Any("code", code),
			slog.Any("error", err))
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", bearer)

	// Set client timeout
	client := &http.Client{Timeout: time.Second * 25}

	// Send req using http Client
	resp, err := client.Do(req)
	if err != nil {
		c.Logger.Error("failed calling fitbit remote api",
			slog.Any("code", code),
			slog.Any("error", err))
		return err
	}

	defer resp.Body.Close()

	// Read the response body
	data := FitBitAuthReponse{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		c.Logger.Error("failed decoding response from remote api",
			slog.Any("code", code),
			slog.Any("response", resp.Body),
			slog.Any("error", err))
		return err
	}

	dataBin, err := json.Marshal(data)
	if err != nil {
		c.Logger.Error("failed unmarshalling response from remote api",
			slog.Any("code", code),
			slog.Any("response", resp.Body),
			slog.Any("error", err))
		return err
	}

	c.Logger.Debug("response data from calling fitbit authorization api",
		slog.Any("code", code),
		slog.Any("response", resp.Body),
		slog.Any("data", data),
		slog.Any("dataBin", dataBin))

	// Clear memory for the user id because we already used it.
	defer delete(c.CodeVerifierMap, userID)
	c.Logger.Debug("cleared userID",
		slog.Any("user_id", userID),
		slog.Any("data", dataBin))

	u, err := c.UserStorer.GetByID(ctx, userID)
	if err != nil {
		c.Logger.Error("database get error",
			slog.Any("code", code),
			slog.Any("user_id", userID),
			slog.Any("error", err))
		return err
	}
	if u == nil {
		c.Logger.Error("user does not exist",
			slog.Any("code", code),
			slog.Any("user_id", userID),
			slog.Any("error", err))
		return errors.New("user does not exist")
	}

	//
	// STEP 4: Create our objects (or get them if they previously existed).
	//

	// Get to modify or create remote device (or get previously existed one).

	fba, err := c.FitBitAppStorer.GetByUserID(ctx, u.ID)
	if err != nil {
		c.Logger.Error("failed getting fitbit app",
			slog.Any("code", code),
			slog.Any("fba", fba),
			slog.Any("user_id", userID),
			slog.Any("error", err))
		return err
	}
	if fba == nil {
		c.Logger.Debug("fitbit does not exist, proceeding to create one for user",
			slog.Any("code", code),
			slog.Any("user_id", userID))
		fba = &fitbitapp_s.FitBitApp{
			ID:                 primitive.NewObjectID(),
			UserFirstName:      u.FirstName,
			UserLastName:       u.LastName,
			UserName:           u.Name,
			UserLexicalName:    u.LexicalName,
			UserID:             u.ID,
			Status:             fitbitapp_s.StatusActive,
			CreatedAt:          time.Now(),
			ModifiedAt:         time.Now(),
			OrganizationID:     u.OrganizationID,
			OrganizationName:   u.OrganizationName,
			FitBitUserID:       data.UserID,
			AuthType:           fitbitapp_s.AuthTypeOAuth2,
			Errors:             "",
			Scope:              data.Scope,
			TokenType:          data.TokenType,
			AccessToken:        data.AccessToken,
			ExpiresIn:          data.ExpiresIn,
			RefreshToken:       data.RefreshToken,
			ExpireTime:         time.Now().Add(time.Second * time.Duration(data.ExpiresIn)),
			LastFetchedAt:      time.Date(2013, 1, 1, 00, 00, 00, 000000000, time.UTC), // 2013-01-01 00:00:00.00 UTC
			HeartRateMetricID:  primitive.NewObjectID(),
			StepsCountMetricID: primitive.NewObjectID(),
			IsTestMode:         false,
			SimulatorAlgorithm: "",
		}
	} else {
		c.Logger.Debug("previous fitbit exist, proceeding to update for user",
			slog.Any("code", code),
			slog.Any("user_id", userID))
		fba.UserFirstName = u.FirstName
		fba.UserLastName = u.LastName
		fba.UserName = u.Name
		fba.UserLexicalName = u.LexicalName
		fba.Status = fitbitapp_s.StatusActive
		fba.ModifiedAt = time.Now()
		fba.OrganizationID = u.OrganizationID
		fba.OrganizationName = u.OrganizationName
		fba.FitBitUserID = data.UserID
		fba.AuthType = fitbitapp_s.AuthTypeOAuth2
		fba.Errors = ""
		fba.Scope = data.Scope
		fba.TokenType = data.TokenType
		fba.AccessToken = data.AccessToken
		fba.ExpiresIn = data.ExpiresIn
		fba.RefreshToken = data.RefreshToken
		fba.ExpireTime = time.Now().Add(time.Second * time.Duration(data.ExpiresIn))
		fba.LastFetchedAt = time.Date(2013, 1, 1, 00, 00, 00, 000000000, time.UTC) // 2013-01-01 00:00:00.00 UTC
	}

	// Essentially run create or update function.
	if err := c.FitBitAppStorer.UpsertByUserID(ctx, fba); err != nil {
		c.Logger.Error("database upsert error",
			slog.Any("code", code),
			slog.Any("fba", fba),
			slog.Any("user_id", userID),
			slog.Any("error", err))
		return err
	}

	// Update our user with our new device.
	u.PrimaryHealthTrackingDeviceType = u_s.UserPrimaryHealthTrackingDeviceTypeFitBit
	u.PrimaryHealthTrackingDeviceHeartRateMetricID = fba.HeartRateMetricID
	u.PrimaryHealthTrackingDeviceStepsCountMetricID = fba.StepsCountMetricID
	u.FitBitAppID = fba.ID
	u.ModifiedAt = time.Now()
	if err := c.UserStorer.UpdateByID(ctx, u); err != nil {
		c.Logger.Error("failed updating user by id",
			slog.Any("code", code),
			slog.Any("user_id", userID),
			slog.Any("error", err))
		return err
	}

	return nil
}

type FitBitAuthReponse struct {
	Scope        string `json:"scope,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"` // Measured in seconds.
	RefreshToken string `json:"refresh_token,omitempty"`
	UserID       string `json:"user_id,omitempty"`
}
