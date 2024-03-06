package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	gateway_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/gateway/controller"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func UnmarshalMemberRegisterRequest(ctx context.Context, r *http.Request) (*gateway_c.MemberRegisterRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData gateway_c.MemberRegisterRequestIDO

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
	requestData.EmailRepeated = strings.ToLower(requestData.EmailRepeated)
	requestData.EmailRepeated = strings.ReplaceAll(requestData.EmailRepeated, " ", "")

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateRegisterRequest(&requestData); err != nil {
		return nil, err
	}

	return &requestData, nil
}

func ValidateRegisterRequest(dirtyData *gateway_c.MemberRegisterRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.OrganizationID.IsZero() {
		e["organization_id"] = "missing value"
	}
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
	if dirtyData.EmailRepeated == "" {
		e["email_repeated"] = "missing value"
	}
	if len(dirtyData.EmailRepeated) > 255 {
		e["email_repeated"] = "too long"
	}
	if dirtyData.Email != dirtyData.EmailRepeated {
		e["email"] = "does not match email repeated"
		e["email_repeated"] = "does not match email"
	}
	if dirtyData.Password == "" {
		e["password"] = "missing value"
	}
	if dirtyData.PasswordRepeated == "" {
		e["password_repeated"] = "missing value"
	}
	if dirtyData.PasswordRepeated != dirtyData.Password {
		e["password"] = "does not match"
		e["password_repeated"] = "does not match"
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
	if dirtyData.Password == "" {
		e["password"] = "missing value"
	}
	if dirtyData.AddressLine1 == "" {
		e["address_line_1"] = "missing value"
	}
	// if dirtyData.HowDidYouHearAboutUs == 0 {
	// 	e["how_did_you_hear_about_us"] = "missing value"
	// }
	if dirtyData.AgreeTOS == false {
		e["agree_tos"] = "you must agree to the terms before proceeding"
	}
	if dirtyData.HowDidYouHearAboutUs > 7 || dirtyData.HowDidYouHearAboutUs < 1 {
		e["how_did_you_hear_about_us"] = "missing value"
	} else {
		if dirtyData.HowDidYouHearAboutUs == 1 && dirtyData.HowDidYouHearAboutUsOther == "" {
			e["how_did_you_hear_about_us_other"] = "missing value"
		}
	}

	// The following logic will enforce shipping address input validation.
	if dirtyData.HasShippingAddress == true {
		if dirtyData.ShippingName == "" {
			e["shipping_name"] = "missing value"
		}
		if dirtyData.ShippingPhone == "" {
			e["shipping_phone"] = "missing value"
		}
		if dirtyData.ShippingCountry == "" {
			e["shipping_country"] = "missing value"
		}
		if dirtyData.ShippingRegion == "" {
			e["shipping_region"] = "missing value"
		}
		if dirtyData.ShippingCity == "" {
			e["shipping_city"] = "missing value"
		}
		if dirtyData.ShippingPostalCode == "" {
			e["shipping_postal_code"] = "missing value"
		}
		if dirtyData.ShippingAddressLine1 == "" {
			e["shipping_address_line1"] = "missing value"
		}
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) MemberRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := UnmarshalMemberRegisterRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	res, err := h.Controller.MemberRegister(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	MarshalLoginResponse(res, w)
}
