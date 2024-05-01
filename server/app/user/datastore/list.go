package datastore

import (
	"context"
	"log"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl UserStorerImpl) ListByFilter(ctx context.Context, f *UserListFilter) (*UserListResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	// Create the filter based on the cursor
	filter := bson.M{}
	if !f.Cursor.IsZero() {
		filter["_id"] = bson.M{"$gt": f.Cursor} // Add the cursor condition to the filter
	}

	// Apply the conditions to the filter
	if !f.OrganizationID.IsZero() {
		filter["organization_id"] = f.OrganizationID
	}
	if f.Role > 0 {
		filter["role"] = f.Role
	}
	if f.FirstName != "" {
		filter["first_name"] = f.FirstName
	}
	if f.LastName != "" {
		filter["last_name"] = f.LastName
	}
	if f.Email != "" {
		filter["email"] = f.Email
	}
	if f.Phone != "" {
		filter["phone"] = f.Phone
	}
	if f.Status != 0 {
		filter["status"] = f.Status
	}

	// impl.Logger.Debug("listing filter:",
	// 	slog.Any("filter", filter))

	// Include additional filters for our cursor-based pagination pertaining to sorting and limit.
	options := options.Find().
		SetSort(bson.M{f.SortField: f.SortOrder}).
		SetLimit(f.PageSize)

	// Include Full-text search
	if f.SearchText != "" {
		filter["$text"] = bson.M{"$search": f.SearchText}
		options.SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}})
		options.SetSort(bson.D{{"score", bson.M{"$meta": "textScore"}}})
	}

	// Execute the query
	cursor, err := impl.Collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// var results = []*ComicSubmission{}
	// if err = cursor.All(ctx, &results); err != nil {
	// 	panic(err)
	// }

	// Retrieve the documents and check if there is a next page
	results := []*User{}
	hasNextPage := false
	for cursor.Next(ctx) {
		document := &User{}
		if err := cursor.Decode(document); err != nil {
			return nil, err
		}
		results = append(results, document)
		// Stop fetching documents if we have reached the desired page size
		if int64(len(results)) >= f.PageSize {
			hasNextPage = true
			break
		}
	}

	// Get the next cursor and encode it
	nextCursor := primitive.NilObjectID
	if int64(len(results)) == f.PageSize {
		// Remove the extra document from the current page
		results = results[:len(results)]

		// Get the last document's _id as the next cursor
		nextCursor = results[len(results)-1].ID
	}

	return &UserListResult{
		Results:     results,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}

func (impl UserStorerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *UserListFilter) ([]*UserAsSelectOption, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	// Get a reference to the collection
	collection := impl.Collection

	// Pagination parameters
	pageSize := 10
	startAfter := "" // The ID to start after, initially empty for the first page

	// Sorting parameters
	sortField := "_id"
	sortOrder := 1 // 1=ascending | -1=descending

	// Pagination query
	query := bson.M{}
	options := options.Find().
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{sortField, sortOrder}})

	// Apply the conditions to the filter
	if !f.OrganizationID.IsZero() {
		query["organization_id"] = f.OrganizationID
	}
	if f.Role > 0 {
		query["role"] = f.Role
	}
	if f.FirstName != "" {
		query["first_name"] = f.FirstName
	}
	if f.LastName != "" {
		query["last_name"] = f.LastName
	}
	if f.Email != "" {
		query["email"] = f.Email
	}
	if f.Phone != "" {
		query["phone"] = f.Phone
	}
	if f.Status != 0 {
		query["status"] = f.Status
	}

	if startAfter != "" {
		// Find the document with the given startAfter ID
		cursor, err := collection.FindOne(ctx, bson.M{"_id": startAfter}).DecodeBytes()
		if err != nil {
			log.Fatal(err)
		}
		options.SetSkip(1)
		query["_id"] = bson.M{"$gt": cursor.Lookup("_id").ObjectID()}
	}

	if f.ExcludeArchived {
		query["status"] = bson.M{"$ne": UserStatusArchived} // Do not list archived items! This code
	}

	options.SetSort(bson.D{{sortField, 1}}) // Sort in ascending order based on the specified field

	// Retrieve the list of items from the collection
	cursor, err := collection.Find(ctx, query, options)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var results = []*UserAsSelectOption{}
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}

	return results, nil
}

func (impl UserStorerImpl) ListLatestStripeInvoices(ctx context.Context, userID primitive.ObjectID, cursor int64, limit int64) (*StripeInvoiceListResult, error) {
	u, err := impl.GetByID(ctx, userID)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if u == nil {
		impl.Logger.Error("user does not exist error", slog.Any("userID", userID))
		return nil, err
	}

	// Get the current list of invoices which are sorted from latest to
	// oldest creation date.
	invoices := u.StripeInvoices

	impl.Logger.Debug("list invoices",
		slog.Any("invoices", invoices),
		slog.Int("invoices Count", len(invoices)),
		slog.Any("userID", userID),
		slog.Int64("cursor", cursor),
		slog.Int64("limit", limit))

	res, err := paginateList(invoices, int(cursor), int(limit))

	impl.Logger.Debug("list invoices",
		slog.Any("res", res))

	return res, err
}

func paginateList(list []*StripeInvoice, cursor int, itemsPerPage int) (*StripeInvoiceListResult, error) {
	if cursor < 0 {
		return &StripeInvoiceListResult{
			Results:     []*StripeInvoice{},
			NextCursor:  -1,
			HasNextPage: false,
		}, nil
	}
	if cursor >= len(list) {
		return &StripeInvoiceListResult{
			Results:     []*StripeInvoice{},
			NextCursor:  -1,
			HasNextPage: false,
		}, nil
	}

	startIndex := cursor
	endIndex := cursor + itemsPerPage

	if startIndex >= len(list) {
		// No more items after the cursor
		return nil, nil
	}

	if endIndex > len(list) {
		endIndex = len(list)
	}

	// Get the items for the current page
	items := list[startIndex:endIndex]

	// Calculate the next cursor
	var nextCursor int
	var hasNextPage bool
	if endIndex < len(list) {
		nextCursor = endIndex
		hasNextPage = nextCursor <= itemsPerPage
	} else {
		// No more items after the current page
		nextCursor = -1
		hasNextPage = false
	}

	// Return the paginated result
	return &StripeInvoiceListResult{
		Results:     items,
		NextCursor:  int64(nextCursor),
		HasNextPage: hasNextPage,
	}, nil
}
