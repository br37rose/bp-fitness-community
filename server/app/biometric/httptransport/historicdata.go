package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"

	bio_c "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/biometric/controller"
	rp_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
)

func (h *Handler) HistoricData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f := &bio_c.HistoricDataRequest{
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

	metricType := int8(0)
	metricTypeStr := query.Get("metric_type")
	if metricTypeStr != "" {
		mt, _ := strconv.ParseInt(metricTypeStr, 10, 64)
		if mt != 0 {
			metricType = int8(mt)
		}
	}
	f.MetricType = metricType

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
	res, err := h.Controller.HistoricData(ctx, f)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalHistoricDataListResponse(res, w)
}

func MarshalHistoricDataListResponse(res *rp_s.RankPointPaginationListResult, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
