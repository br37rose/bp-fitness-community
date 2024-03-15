package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"

	sub_s "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/attachment/datastore"
	"github.com/bci-innovation-labs/bp8fitnesscommunity-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f := &sub_s.AttachmentListFilter{
		Cursor:          primitive.NilObjectID,
		OwnershipID:     primitive.NilObjectID,
		PageSize:        25,
		SortField:       "_id",
		SortOrder:       1, // 1=ascending | -1=descending
		ExcludeArchived: true,
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

	ownershipID := query.Get("ownership_id")
	if ownershipID != "" {
		ownershipID, err := primitive.ObjectIDFromHex(ownershipID)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.OwnershipID = ownershipID
	}

	pageSize := query.Get("page_size")
	if pageSize != "" {
		pageSize, _ := strconv.ParseInt(pageSize, 10, 64)
		if pageSize == 0 || pageSize > 250 {
			pageSize = 250
		}
		f.PageSize = pageSize
	}

	m, err := h.Controller.ListByFilter(ctx, f)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalListResponse(m, w)
}

func MarshalListResponse(res *sub_s.AttachmentListResult, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListAsSelectOptionByFilter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f := &sub_s.AttachmentListFilter{
		Cursor:          primitive.NilObjectID,
		PageSize:        6,
		SortField:       "_id",
		SortOrder:       1, // 1=ascending | -1=descending
		ExcludeArchived: true,
	}

	m, err := h.Controller.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalListAsSelectOptionResponse(m, w)
}

func MarshalListAsSelectOptionResponse(res []*sub_s.AttachmentAsSelectOption, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
