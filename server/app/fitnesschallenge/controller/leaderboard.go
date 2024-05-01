package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	gcp_a "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	rank "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
)

func (c *FitnessChallengeControllerImpl) GetChallengeLeaderBoard(ctx context.Context, id primitive.ObjectID) (*rank.RankPointPaginationListResult, error) {
	chalenge, err := c.Storer.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	users := chalenge.UserIDs
	metritype := []string{}
	for _, v := range chalenge.Rules {
		// add more types here as and when implemented
		switch v.Type {
		case 1:
			metritype = append(metritype, gcp_a.DataTypeShortNameCaloriesBurned)
		case 2:
			metritype = append(metritype, gcp_a.DataTypeShortNameStepCountDelta)
		case 3:
			metritype = append(metritype, gcp_a.DataTypeShortNameSleep)
		case 4:
			metritype = append(metritype, gcp_a.DataTypeShortNameOxygenSaturation)
		case 5:
			metritype = append(metritype, gcp_a.DataTypeShortNameActivitySegment)
		case 6:
			metritype = append(metritype, gcp_a.DataTypeShortNameHeartRateBPM)
		case 7:
			metritype = append(metritype, gcp_a.DataTypeShortNameOxygenSaturation)
		}
	}
	response := new(rank.RankPointPaginationListResult)
	for _, v := range users {
		result, err := c.RankPoint.ListByFilter(ctx, &rank.RankPointPaginationListFilter{
			UserID:              v,
			StartGTE:            chalenge.StartTime,
			MetricDataTypeNames: metritype,
			SortField:           "metric_id",
			SortOrder:           -1,
			PageSize:            1,
		})
		if err != nil {
			return nil, err
		}
		response.Results = append(response.Results, result.Results...)
	}

	return response, err
}
