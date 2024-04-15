package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	w_d "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/workout/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f := &w_d.WorkoutListFilter{
		Cursor:          primitive.NilObjectID,
		PageSize:        25,
		SortField:       "_id",
		SortOrder:       1,
		ExcludeArchived: true,
		// StatusList:      []int8{1},
	}

	// Here is where you extract URL parameters.
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

	visibilityStr := query.Get("visibility")
	if visibilityStr != "" {
		visibility, _ := strconv.ParseBool(visibilityStr)
		f.Visibility = visibility
	}

	createdByUserID := query.Get("created_by_user_id")
	if createdByUserID != "" {
		createdByUserID, err := primitive.ObjectIDFromHex(createdByUserID)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.CreatedByUserID = createdByUserID
	}

	typesStr := query.Get("types")
	if typesStr != "" {
		types := []int8{}
		typesStrings := strings.Split(typesStr, ",")
		for _, typeStr := range typesStrings {
			typ, _ := strconv.ParseInt(typeStr, 10, 64)
			types = append(types, int8(typ))
		}
		f.Types = types
	}

	excludeArchivedStr := query.Get("exclude_archived")
	if excludeArchivedStr != "" {
		excludeArchived, _ := strconv.ParseBool(excludeArchivedStr)
		f.ExcludeArchived = excludeArchived
	}

	searchKeyword := query.Get("search")
	if searchKeyword != "" {
		f.SearchText = searchKeyword
	}

	userId := query.Get("user_id")
	if userId != "" {
		userId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.UserId = userId
	}
	getExc := query.Get("get_exercise")
	if getExc != "" {
		if val, err := strconv.ParseBool(getExc); err != nil && val {
			f.GetExcercise = val
		}
	}

	m, err := h.Controller.ListByFilter(ctx, f)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalListResponse(m, w)
}

func MarshalListResponse(res *w_d.WorkoutistResult, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
