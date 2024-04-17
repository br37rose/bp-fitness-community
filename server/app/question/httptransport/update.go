package httptransport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	q_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/question/controller"
	q_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/question/datastore"

	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func UnmarshalUpdateRequest(ctx context.Context, r *http.Request) (*q_c.QuestionUpdateRequest, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData *q_c.QuestionUpdateRequest

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	if err := ValidateQuestionUpdateRequest(requestData); err != nil {
		return nil, err
	}

	return requestData, nil
}

func ValidateQuestionUpdateRequest(dirtyData *q_c.QuestionUpdateRequest) error {
	e := make(map[string]string)

	if dirtyData.ID.IsZero() {
		e["id"] = "missing value"
	}

	if dirtyData.Question == "" {
		e["question"] = "missing value"

	}
	if len(dirtyData.Content) == 0 {
		e["content"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) UpdateByID(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	data, err := UnmarshalUpdateRequest(ctx, r)
	if err != nil {
		log.Println("Create | member | err:", err)
		httperror.ResponseError(w, err)
		return
	}

	q, err := h.Controller.UpdateByID(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalUpdateResponse(q, w)
}

func MarshalUpdateResponse(res *q_s.Question, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
