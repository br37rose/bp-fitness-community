package fitnessplan

import (
	"encoding/json"
	"net/http"
	"strconv"

	sub_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/fitnessplan/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f := &sub_s.FitnessPlanListFilter{
		Cursor:     primitive.NilObjectID,
		PageSize:   25,
		SortField:  "_id",
		SortOrder:  1, // 1=ascending | -1=descending
		StatusList: make([]int8, 0),
	}

	// Here is where you extract url parameters.
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

	organizationID := query.Get("organization_id")
	if organizationID != "" {
		organizationID, err := primitive.ObjectIDFromHex(organizationID)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.OrganizationID = organizationID
	}

	searchKeyword := query.Get("search")
	if searchKeyword != "" {
		f.SearchText = searchKeyword
	}

	statusStr := query.Get("status")
	if statusStr != "" {
		status, _ := strconv.ParseInt(statusStr, 10, 64)
		if status != 0 {
			f.StatusList = []int8{int8(status)}
		}
	}

	userIdStr := query.Get("user_id")
	if userIdStr != "" {
		userId, err := primitive.ObjectIDFromHex(userIdStr)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.UserID = userId
	}

	m, err := h.Controller.ListByFilter(ctx, f)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalListResponse(m, w)
}

func MarshalListResponse(res *sub_s.FitnessPlanListResult, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
