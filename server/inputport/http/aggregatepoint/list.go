package aggregatepoint

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bartmika/timekit"
	dp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/aggregatepoint/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f := &dp_s.AggregatePointPaginationListFilter{
		Cursor:    "",
		PageSize:  250,
		SortField: "timestamp",
		SortOrder: dp_s.SortOrderDescending,
	}

	// Here is where you extract url parameters.
	query := r.URL.Query()

	cursor := query.Get("cursor")
	if cursor != "" {
		f.Cursor = cursor
	}

	pageSize := query.Get("page_size")
	if pageSize != "" {
		pageSize, _ := strconv.ParseInt(pageSize, 10, 64)
		if pageSize == 0 || pageSize > 250 {
			pageSize = 250
		}
		f.PageSize = pageSize
	}

	sortField := query.Get("sort_field")
	if sortField != "" {
		f.SortField = sortField
	}
	sortOrder := query.Get("sort_order")
	if sortOrder == "ASC" {
		f.SortOrder = dp_s.SortOrderAscending
	}
	if sortOrder == "DESC" {
		f.SortOrder = dp_s.SortOrderDescending
	}
	createdAtGTEStr := query.Get("created_at_gte")
	if createdAtGTEStr != "" {
		createdAtGTE, err := timekit.ParseJavaScriptTimeString(createdAtGTEStr)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.CreatedAtGTE = createdAtGTE
	}
	metricIDs := make([]primitive.ObjectID, 0)
	heartRateID := query.Get("heart_rate_id")
	if heartRateID != "" {
		heartRateID, err := primitive.ObjectIDFromHex(heartRateID)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		metricIDs = append(metricIDs, heartRateID)
	}

	stepsCounterID := query.Get("steps_counter_id")
	if stepsCounterID != "" {
		stepsCounterID, err := primitive.ObjectIDFromHex(stepsCounterID)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		metricIDs = append(metricIDs, stepsCounterID)
	}
	f.MetricIDs = metricIDs

	// Perform our database operation.
	res, err := h.Controller.ListByFilter(ctx, f)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalListResponse(res, w)
}

func MarshalListResponse(res *dp_s.AggregatePointPaginationListResult, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
