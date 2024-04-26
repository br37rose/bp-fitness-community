package controller

import (
	"context"

	metric "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/adapter/cloudprovider/google"
	rank "github.com/bci-innovation-labs/bp8fitnesscommunity-backend/app/rankpoint/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *FitnessChallengeControllerImpl) GetChallengeLeaderBoard(ctx context.Context, id primitive.ObjectID) (*rank.RankPointPaginationListResult, error) {
	chalenge, err := c.Storer.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	users := chalenge.UserIDs
	metritype := []int8{}
	for _, v := range chalenge.Rules {
		// add more types here as and when implemented
		switch v.Type {
		case 1:
			metritype = append(metritype, metric.DataTypeKeyCaloriesBurned)
		case 2:
			metritype = append(metritype, metric.DataTypeKeyStepCountDelta)
		case 3:
			metritype = append(metritype, metric.DataTypeKeySleep)
		case 4:
			metritype = append(metritype, metric.DataTypeKeyOxygenSaturation)
		case 5:
			metritype = append(metritype, metric.DataTypeKeyActivitySegment)
		case 6:
			metritype = append(metritype, metric.DataTypeKeyHeartRateBPM)
		case 7:
			metritype = append(metritype, metric.DataTypeKeyOxygenSaturation)
		}
	}
	response := new(rank.RankPointPaginationListResult)
	for _, v := range users {
		result, err := c.RankPoint.ListByFilter(ctx, &rank.RankPointPaginationListFilter{
			UserID:      v,
			StartGTE:    chalenge.StartTime,
			MetricTypes: metritype,
			SortField:   "metric_id",
			SortOrder:   -1,
			PageSize:    1,
		})
		if err != nil {
			return nil, err
		}
		response.Results = append(response.Results, result.Results...)
	}

	return response, err
}
