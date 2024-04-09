package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	tp_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/trainingprogram/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f := &tp_d.TrainingProgramListFilter{
		Cursor:     primitive.NilObjectID,
		PageSize:   25,
		SortField:  "_id",
		SortOrder:  1,
		StatusList: []int8{1},
	}

	// Extract URL parameters.
	query := r.URL.Query()

	cursor := query.Get("cursor")
	if cursor != "" {
		cursor, err := primitive.ObjectIDFromHex(cursor)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
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

	sortOrderStr := query.Get("sort_order")
	if sortOrderStr != "" {
		sortOrder, _ := strconv.ParseInt(sortOrderStr, 10, 64)
		if sortOrder != 1 && sortOrder != -1 {
			sortOrder = 1
		}
		f.SortOrder = int8(sortOrder)
	}

	userID := query.Get("user_id")
	if userID != "" {
		userID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.UserID = userID
	}

	organizationID := query.Get("organization_id")
	if organizationID != "" {
		organizationID, err := primitive.ObjectIDFromHex(organizationID)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.OrganizationID = organizationID
	}

	name := query.Get("name")
	if name != "" {
		f.Name = name
	}

	phases := query.Get("phases")
	if phases != "" {
		phases, _ := strconv.ParseInt(phases, 10, 64)
		f.Phases = phases
	}

	weeks := query.Get("weeks")
	if weeks != "" {
		weeks, _ := strconv.ParseInt(weeks, 10, 64)
		f.Weeks = weeks
	}

	durationInWeeks := query.Get("duration_in_weeks")
	if durationInWeeks != "" {
		durationInWeeks, _ := strconv.ParseInt(durationInWeeks, 10, 64)
		f.DurationInWeeks = durationInWeeks
	}

	startTimeStr := query.Get("start_time")
	if startTimeStr != "" {
		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.StartTime = startTime
	}

	endTimeStr := query.Get("end_time")
	if endTimeStr != "" {
		endTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.EndTime = endTime
	}

	statusListStr := query.Get("status_list")
	if statusListStr != "" {
		statusList := []int8{}
		statusListStrings := strings.Split(statusListStr, ",")
		for _, statusStr := range statusListStrings {
			status, _ := strconv.ParseInt(statusStr, 10, 64)
			statusList = append(statusList, int8(status))
		}
		f.StatusList = statusList
	}

	m, err := h.Controller.ListByFilter(ctx, f)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalListResponse(m, w)
}

func MarshalListResponse(res *tp_d.TrainingProgramListResult, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
