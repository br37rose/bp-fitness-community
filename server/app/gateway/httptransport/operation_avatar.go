package httptransport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	gateway_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/gateway/controller"
	u_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/user/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UnmarshalOperationAvatarRequest(ctx context.Context, r *http.Request) (*gateway_c.AccountOperationAvatarRequest, error) {
	defer r.Body.Close()

	// Parse the multipart form data
	err := r.ParseMultipartForm(32 << 20) // Limit the maximum memory used for parsing to 32MB
	if err != nil {
		log.Println("UnmarshalCreateRequest:ParseMultipartForm:err:", err)
		return nil, err
	}

	// Get the values of form fields
	userIDStr := r.FormValue("user_id")

	// Get the uploaded file from the request
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println("UnmarshalCmsImageCreateRequest:FormFile:err:", err)
		// return nil, err, http.StatusInternalServerError
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		log.Println("UnmarshalCmsImageCreateRequest: primitive.ObjectIDFromHex:err:", err)
	}

	// Initialize our array which will store all the results from the remote server.
	requestData := &gateway_c.AccountOperationAvatarRequest{
		AccountID: userID,
	}
	if header != nil {
		// Extract filename and filetype from the file header
		requestData.FileName = header.Filename
		requestData.FileType = header.Header.Get("Content-Type")
		requestData.File = file
	}
	return requestData, nil
}

func (h *Handler) OperationAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqData, err := UnmarshalOperationAvatarRequest(ctx, r)
	if err != nil {
		log.Println("OperationAvatar | UnmarshalOperationAvatarRequest | err:", err)
		httperror.ResponseError(w, err)
		return
	}
	data, err := h.Controller.Avatar(ctx, reqData)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalOperationAvatarResponse(data, w)
}

func MarshalOperationAvatarResponse(res *u_s.User, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
