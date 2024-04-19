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

func UnmarshalCreateRequest(ctx context.Context, r *http.Request) (*q_c.QuestionRequest, error) {
	var requestData *q_c.QuestionRequest

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	if err := ValidateCreateRequest(requestData); err != nil {
		return nil, err
	}

	return requestData, nil
}

func ValidateCreateRequest(dirtyData *q_c.QuestionRequest) error {
	e := make(map[string]string)

	if dirtyData.Title == "" {
		e["title"] = "missing value"

	}
	if len(dirtyData.Options) == 0 {
		e["options"] = "missing value"
	}

	if len(dirtyData.Options) > 0 {
		for _, opt := range dirtyData.Options {
			if opt == "" {
				e["options"] = "missing value"
			}
		}
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := UnmarshalCreateRequest(ctx, r)
	if err != nil {
		log.Println("Create | question | err:", err)
		httperror.ResponseError(w, err)
		return
	}

	q, err := h.Controller.Create(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalCreateResponse(q, w)
}

func MarshalCreateResponse(res *q_s.Question, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
