package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	bio_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/biometric/controller"
	rp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (h *Handler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f := &bio_c.LeaderboardRequest{
		Cursor:   "",
		PageSize: 250,
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

	metricDataTypeNamesStr := query.Get("metric_data_type_names")
	metricDataTypeNamesArr := strings.Split(metricDataTypeNamesStr, ",")
	f.MetricDataTypeNames = metricDataTypeNamesArr
	// log.Println("--->", f.MetricDataTypeNames)

	functionStr := query.Get("function")
	if functionStr != "" {
		function, _ := strconv.ParseInt(functionStr, 10, 64)
		if function == 0 || function > 250 {
			function = 1
		}
		f.Function = int8(function)
	}

	periodStr := query.Get("period")
	if periodStr != "" {
		period, _ := strconv.ParseInt(periodStr, 10, 64)
		if period == 0 || period > 250 {
			period = 1
		}
		f.Period = int8(period)
	}

	// Perform our database operation.
	res, err := h.Controller.Leaderboard(ctx, f)
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
