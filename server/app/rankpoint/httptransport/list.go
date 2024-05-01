package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bartmika/timekit"
	"go.mongodb.org/mongo-driver/bson/primitive"

	rp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f := &rp_s.RankPointPaginationListFilter{
		Cursor:    "",
		PageSize:  250,
		SortField: "timestamp",
		SortOrder: rp_s.OrderDescending,
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
		f.SortOrder = rp_s.OrderAscending
	}
	if sortOrder == "DESC" {
		f.SortOrder = rp_s.OrderDescending
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

	metricDataTypeNames := make([]string, 0)
	MetricDataTypeNamesStr := query.Get("metric_data_type_names")
	if MetricDataTypeNamesStr != "" {
		arr := strings.Split(MetricDataTypeNamesStr, ",")
		for _, mtStr := range arr {
			metricDataTypeNames = append(metricDataTypeNames, mtStr)
		}
	}
	f.MetricDataTypeNames = metricDataTypeNames

	periodStr := query.Get("period")
	if periodStr != "" {
		period, _ := strconv.ParseInt(periodStr, 10, 64)
		if period == 0 || period > 250 {
			period = rp_s.PeriodDay
		}
		f.Period = int8(period)
	}

	funcStr := query.Get("function")
	if periodStr != "" {
		funcNumb, _ := strconv.ParseInt(funcStr, 10, 64)
		if funcNumb == 0 || funcNumb > 250 {
			funcNumb = rp_s.FunctionAverage
		}
		f.Function = int8(funcNumb)
	}

	isTodayOnly := query.Get("is_today_only")
	if isTodayOnly == "true" {
		f.StartGTE = timekit.Midnight(time.Now)
		f.EndLTE = timekit.MidnightTomorrow(time.Now)
	}

	// Perform our database operation.
	res, err := h.Controller.ListByFilter(ctx, f)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalListResponse(res, w)
}

func MarshalListResponse(res *rp_s.RankPointPaginationListResult, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
