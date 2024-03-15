package httptransport

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	gateway_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/gateway/controller"
	user_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (h *Handler) Account(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	profile, err := h.Controller.Account(ctx)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	MarshalAccountResponse(profile, w)
}

func MarshalAccountResponse(responseData *user_s.User, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&responseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type AccountUpdateRequestIDO struct {
	FirstName                 string            `bson:"first_name" json:"first_name"`
	LastName                  string            `bson:"last_name" json:"last_name"`
	Email                     string            `json:"email"`
	Phone                     string            `bson:"phone,omitempty" json:"phone,omitempty"`
	Country                   string            `bson:"country,omitempty" json:"country,omitempty"`
	Region                    string            `bson:"region,omitempty" json:"region,omitempty"`
	City                      string            `bson:"city,omitempty" json:"city,omitempty"`
	PostalCode                string            `bson:"postal_code,omitempty" json:"postal_code,omitempty"`
	AddressLine1              string            `bson:"address_line_1,omitempty" json:"address_line_1,omitempty"`
	AddressLine2              string            `bson:"address_line_2,omitempty" json:"address_line_2,omitempty"`
	HowDidYouHearAboutUs      int8              `bson:"how_did_you_hear_about_us,omitempty" json:"how_did_you_hear_about_us,omitempty"`
	HowDidYouHearAboutUsOther string            `bson:"how_did_you_hear_about_us_other,omitempty" json:"how_did_you_hear_about_us_other,omitempty"`
	AgreePromotionsEmail      bool              `bson:"agree_promotions_email,omitempty" json:"agree_promotions_email,omitempty"`
	Tags                      []*user_s.UserTag `bson:"tags" json:"tags"`
}

func UnmarshalAccountUpdateRequest(ctx context.Context, r *http.Request) (*user_s.User, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData AccountUpdateRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Defensive Code: For security purposes we need to remove all whitespaces from the email and lower the characters.
	requestData.Email = strings.ToLower(requestData.Email)
	requestData.Email = strings.ReplaceAll(requestData.Email, " ", "")

	// Perform our validation and return validation error on any issues detected.
	if err = ValidateAccountUpdateRequest(&requestData); err != nil {
		return nil, err
	}

	// Convert to the user collection.
	return &user_s.User{
		FirstName:                 requestData.FirstName,
		LastName:                  requestData.LastName,
		Email:                     requestData.Email,
		Phone:                     requestData.Phone,
		Country:                   requestData.Country,
		Region:                    requestData.Region,
		City:                      requestData.City,
		PostalCode:                requestData.PostalCode,
		AddressLine1:              requestData.AddressLine1,
		AddressLine2:              requestData.AddressLine2,
		HowDidYouHearAboutUs:      requestData.HowDidYouHearAboutUs,
		HowDidYouHearAboutUsOther: requestData.HowDidYouHearAboutUsOther,
		AgreePromotionsEmail:      requestData.AgreePromotionsEmail,
		Tags:                      requestData.Tags,
	}, nil
}

func ValidateAccountUpdateRequest(dirtyData *AccountUpdateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.FirstName == "" {
		e["first_name"] = "missing value"
	}
	if dirtyData.LastName == "" {
		e["last_name"] = "missing value"
	}
	if dirtyData.Email == "" {
		e["email"] = "missing value"
	}
	if len(dirtyData.Email) > 255 {
		e["email"] = "too long"
	}
	if dirtyData.Phone == "" {
		e["phone"] = "missing value"
	}
	if dirtyData.Country == "" {
		e["country"] = "missing value"
	}
	if dirtyData.Region == "" {
		e["region"] = "missing value"
	}
	if dirtyData.City == "" {
		e["city"] = "missing value"
	}
	if dirtyData.PostalCode == "" {
		e["postal_code"] = "missing value"
	}
	if dirtyData.AddressLine1 == "" {
		e["address_line_1"] = "missing value"
	}
	if len(dirtyData.Tags) > 0 {
		for _, tag := range dirtyData.Tags {
			if tag.Text == "" {
				e["text"] = "missing value"
			}
		}
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) AccountUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := UnmarshalAccountUpdateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	if err := h.Controller.AccountUpdate(ctx, data); err != nil {
		httperror.ResponseError(w, err)
		return
	}

	// Get the request
	h.Account(w, r)
}

func UnmarshalAccountChangePasswordRequest(ctx context.Context, r *http.Request) (*gateway_c.AccountChangePasswordRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData gateway_c.AccountChangePasswordRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Return our result
	return &requestData, nil
}

func (h *Handler) AccountChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := UnmarshalAccountChangePasswordRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	if err := h.Controller.AccountChangePassword(ctx, data); err != nil {
		httperror.ResponseError(w, err)
		return
	}

	// Get the request
	h.Account(w, r)
}

func (h *Handler) AccountListLatestInvoices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Here is where you extract url parameters.
	query := r.URL.Query()

	var cursor int64 = 0
	cursorStr := query.Get("cursor")
	if cursorStr != "" {
		cursorInt64, _ := strconv.ParseInt(cursorStr, 10, 64)
		cursor = cursorInt64
	}

	var limit int64 = 25
	limitStr := query.Get("page_size")
	if limitStr != "" {
		pageSize, _ := strconv.ParseInt(limitStr, 10, 64)
		if pageSize == 0 || pageSize > 250 {
			limit = 25
		}
		limit = pageSize
	}

	m, err := h.Controller.AccountListLatestInvoices(ctx, cursor, limit)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalPublicListResponse(m, w)
}

func MarshalPublicListResponse(res *user_s.StripeInvoiceListResult, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
