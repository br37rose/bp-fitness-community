package google

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"google.golang.org/api/fitness/v1"
)

// ParseActivity function converts the `Google Fit` activity data into usable format for our app.
func ParseActivitySegment(datasets []*fitness.Dataset) []ActivitySegmentStruct {
	var data []ActivitySegmentStruct
	// fmt.Println("datasets:", datasets)

	for _, ds := range datasets {
		// fmt.Println("---> ds:", ds)
		// fmt.Println("---> ds.DataSourceId:", ds.DataSourceId)
		// fmt.Println("---> ds.MaxEndTimeNs:", ds.MaxEndTimeNs)
		// fmt.Println("---> ds.MinStartTimeNs:", ds.MinStartTimeNs)
		// fmt.Println("---> ds.NextPageToken:", ds.NextPageToken)
		// fmt.Println("---> ds.Point:", ds.Point)
		// fmt.Println("---> ds.ForceSendFields:", ds.ForceSendFields)
		// fmt.Println("---> ds.NullFields:", ds.NullFields)

		var value int64
		for _, p := range ds.Point {
			// // For debugging purposes only.
			// fmt.Println("--- ---> ComputationTimeMillis:", p.ComputationTimeMillis)
			// fmt.Println("--- ---> DataTypeName:", p.DataTypeName)
			// fmt.Println("--- ---> EndTimeNanos:", p.EndTimeNanos)
			// fmt.Println("--- ---> ModifiedTimeMillis:", p.ModifiedTimeMillis)
			// fmt.Println("--- ---> OriginDataSourceId:", p.OriginDataSourceId)
			// fmt.Println("--- ---> StartTimeNanos:", p.StartTimeNanos)
			// fmt.Println("--- ---> Value[parent]:", p.Value)
			// fmt.Println("--- ---> ForceSendFields:", p.ForceSendFields)
			// fmt.Println("--- ---> NullFields:", p.NullFields)
			// fmt.Println()

			for _, v := range p.Value {
				// // For debugging purposes only.
				// fmt.Println("--- --- ---> v:FpVal:", v.FpVal)
				// fmt.Println("--- --- ---> v:IntVal:", v.IntVal)
				// fmt.Println("--- --- ---> v:MapVal:", v.MapVal)
				// fmt.Println("--- --- ---> v:StringVal:", v.StringVal)
				// fmt.Println("--- --- ---> v:ForceSendFields:", v.ForceSendFields)
				// fmt.Println("--- --- ---> v:NullFields:", v.NullFields)
				// // valueString := fmt.Sprintf("%.3f", v.FpVal)
				// // value, _ = strconv.ParseFloat(valueString, 64)
				value = v.IntVal
			}
			var row ActivitySegmentStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			row.ActivityTypeID = value
			row.ActivityTypeLabel = ActivitySegmentMap[value]

			// Calculate the duration between the two dates
			duration := row.EndTime.Sub(row.StartTime)

			// Convert the duration to minutes
			minutes := int(duration.Minutes())

			// Convert the duration to hours
			hours := int(duration.Hours())

			row.DurationInMinutes = minutes
			row.DurationInHours = hours

			// // For debugging purposes only.
			// fmt.Println("--- --- --- ---> StartTime:", row.StartTime)
			// fmt.Println("--- --- --- ---> EndTime:", row.EndTime)
			// fmt.Println("--- --- --- ---> ActivityTypeID:", row.ActivityTypeID)
			// fmt.Println("--- --- --- ---> ActivityTypeLabel:", row.ActivityTypeLabel)
			// fmt.Println("--- --- --- ---> Minutes:", minutes)
			// fmt.Println("--- --- --- ---> Hours:", hours)

			data = append(data, row)
		}
	}
	return data
}

// ParseBasalMetabolicRate function converts the `Google Fit` power data into usable format for our app.
func ParseBasalMetabolicRate(datasets []*fitness.Dataset) []BasalMetabolicRateStruct {
	var data []BasalMetabolicRateStruct
	// fmt.Println("datasets:", datasets)

	for _, ds := range datasets {
		// fmt.Println("ds:", ds)
		// fmt.Println("ds.DataSourceId:", ds.DataSourceId)
		// fmt.Println("ds.MaxEndTimeNs:", ds.MaxEndTimeNs)
		// fmt.Println("ds.MinStartTimeNs:", ds.MinStartTimeNs)
		// fmt.Println("ds.NextPageToken:", ds.NextPageToken)
		// fmt.Println("ds.Point:", ds.Point)
		// fmt.Println("ds.ForceSendFields:", ds.ForceSendFields)
		// fmt.Println("ds.NullFields:", ds.NullFields)

		var value float64
		for _, p := range ds.Point {
			// // For debugging purposes only.
			// fmt.Println("ComputationTimeMillis:", p.ComputationTimeMillis)
			// fmt.Println("DataTypeName:", p.DataTypeName)
			// fmt.Println("EndTimeNanos:", p.EndTimeNanos)
			// fmt.Println("ModifiedTimeMillis:", p.ModifiedTimeMillis)
			// fmt.Println("OriginDataSourceId:", p.OriginDataSourceId)
			// fmt.Println("StartTimeNanos:", p.StartTimeNanos)
			// fmt.Println("Value[parent]:", p.Value)
			// for _, v := range p.Value {
			// 	fmt.Println("Value[child]:", v)
			// }
			// fmt.Println("ForceSendFields:", p.ForceSendFields)
			// fmt.Println("NullFields:", p.NullFields)
			// fmt.Println()

			for _, v := range p.Value {
				// // For debugging purposes only.
				// fmt.Println("v:FpVal:", v.FpVal)
				// fmt.Println("v:IntVal:", v.IntVal)
				// fmt.Println("v:MapVal:", v.MapVal)
				// fmt.Println("v:StringVal:", v.StringVal)
				// fmt.Println("v:ForceSendFields:", v.ForceSendFields)
				// fmt.Println("v:NullFields:", v.NullFields)
				valueString := fmt.Sprintf("%.3f", v.FpVal)
				value, _ = strconv.ParseFloat(valueString, 64)
			}
			var row BasalMetabolicRateStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			row.Calories = float64(value)
			data = append(data, row)
		}
	}
	return data
}

// ParseCaloriesBurned function converts the `Google Fit` calories burned data into usable format for our app.
func ParseCaloriesBurned(datasets []*fitness.Dataset) []CaloriesBurnedStruct {
	var data []CaloriesBurnedStruct

	for _, ds := range datasets {
		var value float64
		for _, p := range ds.Point {
			// // For debugging purposes only.
			// fmt.Println("ComputationTimeMillis:", p.ComputationTimeMillis)
			// fmt.Println("DataTypeName:", p.DataTypeName)
			// fmt.Println("EndTimeNanos:", p.EndTimeNanos)
			// fmt.Println("ModifiedTimeMillis:", p.ModifiedTimeMillis)
			// fmt.Println("OriginDataSourceId:", p.OriginDataSourceId)
			// fmt.Println("StartTimeNanos:", p.StartTimeNanos)
			// fmt.Println("Value[parent]:", p.Value)
			// for _, v := range p.Value {
			// 	fmt.Println("Value[child]:", v)
			// }
			// fmt.Println("ForceSendFields:", p.ForceSendFields)
			// fmt.Println("NullFields:", p.NullFields)
			// fmt.Println()

			for _, v := range p.Value {
				// // For debugging purposes only.
				// fmt.Println("v:FpVal:", v.FpVal)
				// fmt.Println("v:IntVal:", v.IntVal)
				// fmt.Println("v:MapVal:", v.MapVal)
				// fmt.Println("v:StringVal:", v.StringVal)
				// fmt.Println("v:ForceSendFields:", v.ForceSendFields)
				// fmt.Println("v:NullFields:", v.NullFields)
				valueString := fmt.Sprintf("%.3f", v.FpVal)
				value, _ = strconv.ParseFloat(valueString, 64)
			}
			var row CaloriesBurnedStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			row.Calories = float64(value)
			data = append(data, row)
		}
	}
	return data
}

// ParseCyclingPedalingCadence function converts the `Google Fit` power data into usable format for our app.
func ParseCyclingPedalingCadence(datasets []*fitness.Dataset) []CyclingPedalingCadenceStruct {
	var data []CyclingPedalingCadenceStruct
	// fmt.Println("datasets:", datasets)

	for _, ds := range datasets {
		// fmt.Println("ds:", ds)
		// fmt.Println("ds.DataSourceId:", ds.DataSourceId)
		// fmt.Println("ds.MaxEndTimeNs:", ds.MaxEndTimeNs)
		// fmt.Println("ds.MinStartTimeNs:", ds.MinStartTimeNs)
		// fmt.Println("ds.NextPageToken:", ds.NextPageToken)
		// fmt.Println("ds.Point:", ds.Point)
		// fmt.Println("ds.ForceSendFields:", ds.ForceSendFields)
		// fmt.Println("ds.NullFields:", ds.NullFields)

		var value float64
		for _, p := range ds.Point {
			// // For debugging purposes only.
			// fmt.Println("ComputationTimeMillis:", p.ComputationTimeMillis)
			// fmt.Println("DataTypeName:", p.DataTypeName)
			// fmt.Println("EndTimeNanos:", p.EndTimeNanos)
			// fmt.Println("ModifiedTimeMillis:", p.ModifiedTimeMillis)
			// fmt.Println("OriginDataSourceId:", p.OriginDataSourceId)
			// fmt.Println("StartTimeNanos:", p.StartTimeNanos)
			// fmt.Println("Value[parent]:", p.Value)
			// for _, v := range p.Value {
			// 	fmt.Println("Value[child]:", v)
			// }
			// fmt.Println("ForceSendFields:", p.ForceSendFields)
			// fmt.Println("NullFields:", p.NullFields)
			// fmt.Println()

			for _, v := range p.Value {
				// // For debugging purposes only.
				// fmt.Println("v:FpVal:", v.FpVal)
				// fmt.Println("v:IntVal:", v.IntVal)
				// fmt.Println("v:MapVal:", v.MapVal)
				// fmt.Println("v:StringVal:", v.StringVal)
				// fmt.Println("v:ForceSendFields:", v.ForceSendFields)
				// fmt.Println("v:NullFields:", v.NullFields)
				valueString := fmt.Sprintf("%.3f", v.FpVal)
				value, _ = strconv.ParseFloat(valueString, 64)
			}
			var row CyclingPedalingCadenceStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			row.RPM = float64(value)
			data = append(data, row)
		}
	}
	return data
}

// ParseCyclingPedalingCumulative function converts the `Google Fit` power data into usable format for our app.
func ParseCyclingPedalingCumulative(datasets []*fitness.Dataset) []CyclingPedalingCumulativeStruct {
	var data []CyclingPedalingCumulativeStruct
	// fmt.Println("datasets:", datasets)

	for _, ds := range datasets {
		// fmt.Println("ds:", ds)
		// fmt.Println("ds.DataSourceId:", ds.DataSourceId)
		// fmt.Println("ds.MaxEndTimeNs:", ds.MaxEndTimeNs)
		// fmt.Println("ds.MinStartTimeNs:", ds.MinStartTimeNs)
		// fmt.Println("ds.NextPageToken:", ds.NextPageToken)
		// fmt.Println("ds.Point:", ds.Point)
		// fmt.Println("ds.ForceSendFields:", ds.ForceSendFields)
		// fmt.Println("ds.NullFields:", ds.NullFields)

		var value float64
		for _, p := range ds.Point {
			// // For debugging purposes only.
			// fmt.Println("ComputationTimeMillis:", p.ComputationTimeMillis)
			// fmt.Println("DataTypeName:", p.DataTypeName)
			// fmt.Println("EndTimeNanos:", p.EndTimeNanos)
			// fmt.Println("ModifiedTimeMillis:", p.ModifiedTimeMillis)
			// fmt.Println("OriginDataSourceId:", p.OriginDataSourceId)
			// fmt.Println("StartTimeNanos:", p.StartTimeNanos)
			// fmt.Println("Value[parent]:", p.Value)
			// for _, v := range p.Value {
			// 	fmt.Println("Value[child]:", v)
			// }
			// fmt.Println("ForceSendFields:", p.ForceSendFields)
			// fmt.Println("NullFields:", p.NullFields)
			// fmt.Println()

			for _, v := range p.Value {
				// // For debugging purposes only.
				// fmt.Println("v:FpVal:", v.FpVal)
				// fmt.Println("v:IntVal:", v.IntVal)
				// fmt.Println("v:MapVal:", v.MapVal)
				// fmt.Println("v:StringVal:", v.StringVal)
				// fmt.Println("v:ForceSendFields:", v.ForceSendFields)
				// fmt.Println("v:NullFields:", v.NullFields)
				valueString := fmt.Sprintf("%.3f", v.FpVal)
				value, _ = strconv.ParseFloat(valueString, 64)
			}
			var row CyclingPedalingCumulativeStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			row.Revolutions = int(value)
			data = append(data, row)
		}
	}
	return data
}

// ParseHeartPoints function converts the `Google Fit` power data into usable format for our app.
func ParseHeartPoints(datasets []*fitness.Dataset) []HeartPointsStruct {
	var data []HeartPointsStruct
	// fmt.Println("datasets:", datasets)

	for _, ds := range datasets {
		// fmt.Println("ds:", ds)
		// fmt.Println("ds.DataSourceId:", ds.DataSourceId)
		// fmt.Println("ds.MaxEndTimeNs:", ds.MaxEndTimeNs)
		// fmt.Println("ds.MinStartTimeNs:", ds.MinStartTimeNs)
		// fmt.Println("ds.NextPageToken:", ds.NextPageToken)
		// fmt.Println("ds.Point:", ds.Point)
		// fmt.Println("ds.ForceSendFields:", ds.ForceSendFields)
		// fmt.Println("ds.NullFields:", ds.NullFields)

		var value float64
		for _, p := range ds.Point {
			// // For debugging purposes only.
			// fmt.Println("ComputationTimeMillis:", p.ComputationTimeMillis)
			// fmt.Println("DataTypeName:", p.DataTypeName)
			// fmt.Println("EndTimeNanos:", p.EndTimeNanos)
			// fmt.Println("ModifiedTimeMillis:", p.ModifiedTimeMillis)
			// fmt.Println("OriginDataSourceId:", p.OriginDataSourceId)
			// fmt.Println("StartTimeNanos:", p.StartTimeNanos)
			// fmt.Println("Value[parent]:", p.Value)
			// for _, v := range p.Value {
			// 	fmt.Println("Value[child]:", v)
			// }
			// fmt.Println("ForceSendFields:", p.ForceSendFields)
			// fmt.Println("NullFields:", p.NullFields)
			// fmt.Println()

			for _, v := range p.Value {
				// // For debugging purposes only.
				// fmt.Println("v:FpVal:", v.FpVal)
				// fmt.Println("v:IntVal:", v.IntVal)
				// fmt.Println("v:MapVal:", v.MapVal)
				// fmt.Println("v:StringVal:", v.StringVal)
				// fmt.Println("v:ForceSendFields:", v.ForceSendFields)
				// fmt.Println("v:NullFields:", v.NullFields)
				valueString := fmt.Sprintf("%.3f", v.FpVal)
				value, _ = strconv.ParseFloat(valueString, 64)
			}
			var row HeartPointsStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			row.Intensity = float64(value)
			data = append(data, row)
		}
	}
	return data
}

// ParseMoveMinutes function converts the `Google Fit` power data into usable format for our app.
func ParseMoveMinutes(datasets []*fitness.Dataset) []MoveMinutesStruct {
	var data []MoveMinutesStruct
	// fmt.Println("datasets:", datasets)

	for _, ds := range datasets {
		// fmt.Println("ds:", ds)
		// fmt.Println("ds.DataSourceId:", ds.DataSourceId)
		// fmt.Println("ds.MaxEndTimeNs:", ds.MaxEndTimeNs)
		// fmt.Println("ds.MinStartTimeNs:", ds.MinStartTimeNs)
		// fmt.Println("ds.NextPageToken:", ds.NextPageToken)
		// fmt.Println("ds.Point:", ds.Point)
		// fmt.Println("ds.ForceSendFields:", ds.ForceSendFields)
		// fmt.Println("ds.NullFields:", ds.NullFields)

		var value float64
		for _, p := range ds.Point {
			// // For debugging purposes only.
			// fmt.Println("ComputationTimeMillis:", p.ComputationTimeMillis)
			// fmt.Println("DataTypeName:", p.DataTypeName)
			// fmt.Println("EndTimeNanos:", p.EndTimeNanos)
			// fmt.Println("ModifiedTimeMillis:", p.ModifiedTimeMillis)
			// fmt.Println("OriginDataSourceId:", p.OriginDataSourceId)
			// fmt.Println("StartTimeNanos:", p.StartTimeNanos)
			// fmt.Println("Value[parent]:", p.Value)
			// for _, v := range p.Value {
			// 	fmt.Println("Value[child]:", v)
			// }
			// fmt.Println("ForceSendFields:", p.ForceSendFields)
			// fmt.Println("NullFields:", p.NullFields)
			// fmt.Println()

			for _, v := range p.Value {
				// // For debugging purposes only.
				// fmt.Println("v:FpVal:", v.FpVal)
				// fmt.Println("v:IntVal:", v.IntVal)
				// fmt.Println("v:MapVal:", v.MapVal)
				// fmt.Println("v:StringVal:", v.StringVal)
				// fmt.Println("v:ForceSendFields:", v.ForceSendFields)
				// fmt.Println("v:NullFields:", v.NullFields)
				valueString := fmt.Sprintf("%.3f", v.FpVal)
				value, _ = strconv.ParseFloat(valueString, 64)
			}
			var row MoveMinutesStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			row.Duration = int(value)
			data = append(data, row)
		}
	}
	return data
}

// ParsePower function converts the `Google Fit` power data into usable format for our app.
func ParsePower(datasets []*fitness.Dataset) []PowerStruct {
	var data []PowerStruct
	// fmt.Println("datasets:", datasets)

	for _, ds := range datasets {
		// fmt.Println("ds:", ds)
		// fmt.Println("ds.DataSourceId:", ds.DataSourceId)
		// fmt.Println("ds.MaxEndTimeNs:", ds.MaxEndTimeNs)
		// fmt.Println("ds.MinStartTimeNs:", ds.MinStartTimeNs)
		// fmt.Println("ds.NextPageToken:", ds.NextPageToken)
		// fmt.Println("ds.Point:", ds.Point)
		// fmt.Println("ds.ForceSendFields:", ds.ForceSendFields)
		// fmt.Println("ds.NullFields:", ds.NullFields)

		var value float64
		for _, p := range ds.Point {
			// // For debugging purposes only.
			// fmt.Println("ComputationTimeMillis:", p.ComputationTimeMillis)
			// fmt.Println("DataTypeName:", p.DataTypeName)
			// fmt.Println("EndTimeNanos:", p.EndTimeNanos)
			// fmt.Println("ModifiedTimeMillis:", p.ModifiedTimeMillis)
			// fmt.Println("OriginDataSourceId:", p.OriginDataSourceId)
			// fmt.Println("StartTimeNanos:", p.StartTimeNanos)
			// fmt.Println("Value[parent]:", p.Value)
			// for _, v := range p.Value {
			// 	fmt.Println("Value[child]:", v)
			// }
			// fmt.Println("ForceSendFields:", p.ForceSendFields)
			// fmt.Println("NullFields:", p.NullFields)
			// fmt.Println()

			for _, v := range p.Value {
				// // For debugging purposes only.
				// fmt.Println("v:FpVal:", v.FpVal)
				// fmt.Println("v:IntVal:", v.IntVal)
				// fmt.Println("v:MapVal:", v.MapVal)
				// fmt.Println("v:StringVal:", v.StringVal)
				// fmt.Println("v:ForceSendFields:", v.ForceSendFields)
				// fmt.Println("v:NullFields:", v.NullFields)
				valueString := fmt.Sprintf("%.3f", v.FpVal)
				value, _ = strconv.ParseFloat(valueString, 64)
			}
			var row PowerStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			row.Watts = float64(value)
			data = append(data, row)
		}
	}
	return data
}

// ParseStepCountDelta function converts the `Google Fit` steps count (per minute) data into usable format for our app.
func ParseStepCountDelta(datasets []*fitness.Dataset) []StepCountDeltaStruct {
	var data []StepCountDeltaStruct

	for _, ds := range datasets {
		var value float64
		for _, p := range ds.Point {
			// // For debugging purposes only.
			// fmt.Println("ComputationTimeMillis:", p.ComputationTimeMillis)
			// fmt.Println("DataTypeName:", p.DataTypeName)
			// fmt.Println("EndTimeNanos:", p.EndTimeNanos)
			// fmt.Println("ModifiedTimeMillis:", p.ModifiedTimeMillis)
			// fmt.Println("OriginDataSourceId:", p.OriginDataSourceId)
			// fmt.Println("StartTimeNanos:", p.StartTimeNanos)
			// fmt.Println("Value[parent]:", p.Value)
			// for _, v := range p.Value {
			// 	fmt.Println("Value[child]:", v)
			// }
			// fmt.Println("ForceSendFields:", p.ForceSendFields)
			// fmt.Println("NullFields:", p.NullFields)
			// fmt.Println()

			for _, v := range p.Value {
				// // For debugging purposes only.
				// fmt.Println("v:FpVal:", v.FpVal)
				// fmt.Println("v:IntVal:", v.IntVal)
				// fmt.Println("v:MapVal:", v.MapVal)
				// fmt.Println("v:StringVal:", v.StringVal)
				// fmt.Println("v:ForceSendFields:", v.ForceSendFields)
				// fmt.Println("v:NullFields:", v.NullFields)
				valueString := fmt.Sprintf("%v", v.IntVal)
				value, _ = strconv.ParseFloat(valueString, 64)
			}
			var row StepCountDeltaStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			row.Steps = int(value)
			data = append(data, row)
		}
	}
	return data
}

// ParseStepCountCadence function converts the `Google Fit` steps count (per minute) data into usable format for our app.
func ParseStepCountCadence(datasets []*fitness.Dataset) []StepCountCadenceStruct {
	var data []StepCountCadenceStruct

	for _, ds := range datasets {
		var value float64
		for _, p := range ds.Point {
			// // For debugging purposes only.
			// fmt.Println("ComputationTimeMillis:", p.ComputationTimeMillis)
			// fmt.Println("DataTypeName:", p.DataTypeName)
			// fmt.Println("EndTimeNanos:", p.EndTimeNanos)
			// fmt.Println("ModifiedTimeMillis:", p.ModifiedTimeMillis)
			// fmt.Println("OriginDataSourceId:", p.OriginDataSourceId)
			// fmt.Println("StartTimeNanos:", p.StartTimeNanos)
			// fmt.Println("Value[parent]:", p.Value)
			// for _, v := range p.Value {
			// 	fmt.Println("Value[child]:", v)
			// }
			// fmt.Println("ForceSendFields:", p.ForceSendFields)
			// fmt.Println("NullFields:", p.NullFields)
			// fmt.Println()

			for _, v := range p.Value {
				// // For debugging purposes only.
				// fmt.Println("v:FpVal:", v.FpVal)
				// fmt.Println("v:IntVal:", v.IntVal)
				// fmt.Println("v:MapVal:", v.MapVal)
				// fmt.Println("v:StringVal:", v.StringVal)
				// fmt.Println("v:ForceSendFields:", v.ForceSendFields)
				// fmt.Println("v:NullFields:", v.NullFields)
				valueString := fmt.Sprintf("%v", v.IntVal)
				value, _ = strconv.ParseFloat(valueString, 64)
			}
			var row StepCountCadenceStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			row.RPM = int(value)
			data = append(data, row)
		}
	}
	return data
}

// ParseWorkout function which is deprecated. Note: https://9to5google.com/2020/11/30/google-fit-latest-update-removes-advanced-weight-training-features-on-wear-os/.
func ParseWorkout(datasets []*fitness.Dataset) []WorkoutStruct {
	var data []WorkoutStruct
	// fmt.Println("datasets:", datasets)

	for _, ds := range datasets {
		// fmt.Println("---> ds:", ds)
		// fmt.Println("---> ds.DataSourceId:", ds.DataSourceId)
		// fmt.Println("---> ds.MaxEndTimeNs:", ds.MaxEndTimeNs)
		// fmt.Println("---> ds.MinStartTimeNs:", ds.MinStartTimeNs)
		// fmt.Println("---> ds.NextPageToken:", ds.NextPageToken)
		// fmt.Println("---> ds.Point:", ds.Point)
		// fmt.Println("---> ds.ForceSendFields:", ds.ForceSendFields)
		// fmt.Println("---> ds.NullFields:", ds.NullFields)

		var value int64
		for _, p := range ds.Point {
			// // For debugging purposes only.
			// fmt.Println("--- ---> ComputationTimeMillis:", p.ComputationTimeMillis)
			// fmt.Println("--- ---> DataTypeName:", p.DataTypeName)
			// fmt.Println("--- ---> EndTimeNanos:", p.EndTimeNanos)
			// fmt.Println("--- ---> ModifiedTimeMillis:", p.ModifiedTimeMillis)
			// fmt.Println("--- ---> OriginDataSourceId:", p.OriginDataSourceId)
			// fmt.Println("--- ---> StartTimeNanos:", p.StartTimeNanos)
			// fmt.Println("--- ---> Value[parent]:", p.Value)
			// fmt.Println("--- ---> ForceSendFields:", p.ForceSendFields)
			// fmt.Println("--- ---> NullFields:", p.NullFields)
			// fmt.Println()

			for _, v := range p.Value {
				// // For debugging purposes only.
				// fmt.Println("--- --- ---> v:FpVal:", v.FpVal)
				// fmt.Println("--- --- ---> v:IntVal:", v.IntVal)
				// fmt.Println("--- --- ---> v:MapVal:", v.MapVal)
				// fmt.Println("--- --- ---> v:StringVal:", v.StringVal)
				// fmt.Println("--- --- ---> v:ForceSendFields:", v.ForceSendFields)
				// fmt.Println("--- --- ---> v:NullFields:", v.NullFields)
				// // valueString := fmt.Sprintf("%.3f", v.FpVal)
				// // value, _ = strconv.ParseFloat(valueString, 64)
				value = v.IntVal
			}
			var row WorkoutStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			log.Println("--->", value)
			// row.ActivityTypeID = value
			// row.ActivityTypeLabel = WorkoutMap[value]

			// Calculate the duration between the two dates
			duration := row.EndTime.Sub(row.StartTime)

			// Convert the duration to minutes
			minutes := int(duration.Minutes())

			// Convert the duration to hours
			hours := int(duration.Hours())

			row.DurationInMinutes = minutes
			row.DurationInHours = hours

			// // For debugging purposes only.
			// fmt.Println("--- --- --- ---> StartTime:", row.StartTime)
			// fmt.Println("--- --- --- ---> EndTime:", row.EndTime)
			// fmt.Println("--- --- --- ---> ActivityTypeID:", row.ActivityTypeID)
			// fmt.Println("--- --- --- ---> ActivityTypeLabel:", row.ActivityTypeLabel)
			// fmt.Println("--- --- --- ---> Minutes:", minutes)
			// fmt.Println("--- --- --- ---> Hours:", hours)

			data = append(data, row)
		}
	}
	return data
}

func ParseCyclingWheelRevolutionRPM(datasets []*fitness.Dataset) []CyclingWheelRevolutionRPMStruct {
	var data []CyclingWheelRevolutionRPMStruct
	// fmt.Println("datasets:", datasets)

	for _, ds := range datasets {
		// fmt.Println("---> ds:", ds)
		// fmt.Println("---> ds.DataSourceId:", ds.DataSourceId)
		// fmt.Println("---> ds.MaxEndTimeNs:", ds.MaxEndTimeNs)
		// fmt.Println("---> ds.MinStartTimeNs:", ds.MinStartTimeNs)
		// fmt.Println("---> ds.NextPageToken:", ds.NextPageToken)
		// fmt.Println("---> ds.Point:", ds.Point)
		// fmt.Println("---> ds.ForceSendFields:", ds.ForceSendFields)
		// fmt.Println("---> ds.NullFields:", ds.NullFields)

		var value float64
		for _, p := range ds.Point {
			// // For debugging purposes only.
			// fmt.Println("--- ---> ComputationTimeMillis:", p.ComputationTimeMillis)
			// fmt.Println("--- ---> DataTypeName:", p.DataTypeName)
			// fmt.Println("--- ---> EndTimeNanos:", p.EndTimeNanos)
			// fmt.Println("--- ---> ModifiedTimeMillis:", p.ModifiedTimeMillis)
			// fmt.Println("--- ---> OriginDataSourceId:", p.OriginDataSourceId)
			// fmt.Println("--- ---> StartTimeNanos:", p.StartTimeNanos)
			// fmt.Println("--- ---> Value[parent]:", p.Value)
			// fmt.Println("--- ---> ForceSendFields:", p.ForceSendFields)
			// fmt.Println("--- ---> NullFields:", p.NullFields)
			// fmt.Println()

			for _, v := range p.Value {
				// // For debugging purposes only.
				// fmt.Println("--- --- ---> v:FpVal:", v.FpVal)
				// fmt.Println("--- --- ---> v:IntVal:", v.IntVal)
				// fmt.Println("--- --- ---> v:MapVal:", v.MapVal)
				// fmt.Println("--- --- ---> v:StringVal:", v.StringVal)
				// fmt.Println("--- --- ---> v:ForceSendFields:", v.ForceSendFields)
				// fmt.Println("--- --- ---> v:NullFields:", v.NullFields)
				valueString := fmt.Sprintf("%.3f", v.FpVal)
				value, _ = strconv.ParseFloat(valueString, 64)
			}
			var row CyclingWheelRevolutionRPMStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			row.RPM = value

			data = append(data, row)
		}
	}
	return data
}

func ParseCyclingWheelRevolutionCumulative(datasets []*fitness.Dataset) []CyclingWheelRevolutionCumulativeStruct {
	var data []CyclingWheelRevolutionCumulativeStruct
	// fmt.Println("datasets:", datasets)

	for _, ds := range datasets {
		// fmt.Println("---> ds:", ds)
		// fmt.Println("---> ds.DataSourceId:", ds.DataSourceId)
		// fmt.Println("---> ds.MaxEndTimeNs:", ds.MaxEndTimeNs)
		// fmt.Println("---> ds.MinStartTimeNs:", ds.MinStartTimeNs)
		// fmt.Println("---> ds.NextPageToken:", ds.NextPageToken)
		// fmt.Println("---> ds.Point:", ds.Point)
		// fmt.Println("---> ds.ForceSendFields:", ds.ForceSendFields)
		// fmt.Println("---> ds.NullFields:", ds.NullFields)

		var value float64
		for _, p := range ds.Point {
			// // For debugging purposes only.
			// fmt.Println("--- ---> ComputationTimeMillis:", p.ComputationTimeMillis)
			// fmt.Println("--- ---> DataTypeName:", p.DataTypeName)
			// fmt.Println("--- ---> EndTimeNanos:", p.EndTimeNanos)
			// fmt.Println("--- ---> ModifiedTimeMillis:", p.ModifiedTimeMillis)
			// fmt.Println("--- ---> OriginDataSourceId:", p.OriginDataSourceId)
			// fmt.Println("--- ---> StartTimeNanos:", p.StartTimeNanos)
			// fmt.Println("--- ---> Value[parent]:", p.Value)
			// fmt.Println("--- ---> ForceSendFields:", p.ForceSendFields)
			// fmt.Println("--- ---> NullFields:", p.NullFields)
			// fmt.Println()

			for _, v := range p.Value {
				// // For debugging purposes only.
				// fmt.Println("--- --- ---> v:FpVal:", v.FpVal)
				// fmt.Println("--- --- ---> v:IntVal:", v.IntVal)
				// fmt.Println("--- --- ---> v:MapVal:", v.MapVal)
				// fmt.Println("--- --- ---> v:StringVal:", v.StringVal)
				// fmt.Println("--- --- ---> v:ForceSendFields:", v.ForceSendFields)
				// fmt.Println("--- --- ---> v:NullFields:", v.NullFields)
				valueString := fmt.Sprintf("%.3f", v.FpVal)
				value, _ = strconv.ParseFloat(valueString, 64)
			}
			var row CyclingWheelRevolutionCumulativeStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			row.Revolutions = value

			data = append(data, row)
		}
	}
	return data
}

func ParseDistanceDelta(datasets []*fitness.Dataset) []DistanceDeltaStruct {
	var data []DistanceDeltaStruct
	// fmt.Println("datasets:", datasets)

	for _, ds := range datasets {
		// fmt.Println("---> ds:", ds)
		// fmt.Println("---> ds.DataSourceId:", ds.DataSourceId)
		// fmt.Println("---> ds.MaxEndTimeNs:", ds.MaxEndTimeNs)
		// fmt.Println("---> ds.MinStartTimeNs:", ds.MinStartTimeNs)
		// fmt.Println("---> ds.NextPageToken:", ds.NextPageToken)
		// fmt.Println("---> ds.Point:", ds.Point)
		// fmt.Println("---> ds.ForceSendFields:", ds.ForceSendFields)
		// fmt.Println("---> ds.NullFields:", ds.NullFields)

		var value float64
		for _, p := range ds.Point {
			// // For debugging purposes only.
			// fmt.Println("--- ---> ComputationTimeMillis:", p.ComputationTimeMillis)
			// fmt.Println("--- ---> DataTypeName:", p.DataTypeName)
			// fmt.Println("--- ---> EndTimeNanos:", p.EndTimeNanos)
			// fmt.Println("--- ---> ModifiedTimeMillis:", p.ModifiedTimeMillis)
			// fmt.Println("--- ---> OriginDataSourceId:", p.OriginDataSourceId)
			// fmt.Println("--- ---> StartTimeNanos:", p.StartTimeNanos)
			// fmt.Println("--- ---> Value[parent]:", p.Value)
			// fmt.Println("--- ---> ForceSendFields:", p.ForceSendFields)
			// fmt.Println("--- ---> NullFields:", p.NullFields)
			// fmt.Println()

			for _, v := range p.Value {
				// // For debugging purposes only.
				// fmt.Println("--- --- ---> v:FpVal:", v.FpVal)
				// fmt.Println("--- --- ---> v:IntVal:", v.IntVal)
				// fmt.Println("--- --- ---> v:MapVal:", v.MapVal)
				// fmt.Println("--- --- ---> v:StringVal:", v.StringVal)
				// fmt.Println("--- --- ---> v:ForceSendFields:", v.ForceSendFields)
				// fmt.Println("--- --- ---> v:NullFields:", v.NullFields)
				valueString := fmt.Sprintf("%.3f", v.FpVal)
				value, _ = strconv.ParseFloat(valueString, 64)
			}
			var row DistanceDeltaStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			row.Distance = value

			data = append(data, row)
		}
	}
	return data
}

func getLocationValue(locationValue []*fitness.ValueMapValEntry, key string) float64 {
	for _, entry := range locationValue {
		if entry.Key == key {
			return entry.Value.FpVal
		}
	}
	return 0.0
}

func ParseLocationSample(datasets []*fitness.Dataset) []LocationSampleStruct {
	var data []LocationSampleStruct

	for _, dataset := range datasets {
		for _, point := range dataset.Point {
			// Extract relevant fields from the dataset point
			startTime := time.Unix(0, point.StartTimeNanos*int64(time.Millisecond))
			endTime := time.Unix(0, point.EndTimeNanos*int64(time.Millisecond))

			// Loop over the values in the dataset point
			for _, value := range point.Value {
				if locationValue := value.MapVal; locationValue != nil {
					// Extract latitude, longitude, accuracy, and altitude from the locationValue
					latitude := getLocationValue(locationValue, "latitude")
					longitude := getLocationValue(locationValue, "longitude")
					accuracy := getLocationValue(locationValue, "accuracy")
					altitude := getLocationValue(locationValue, "altitude")

					// Create a new LocationSampleStruct and append it to the data slice
					data = append(data, LocationSampleStruct{
						StartTime: startTime,
						EndTime:   endTime,
						Latitude:  latitude,
						Longitude: longitude,
						Accuracy:  accuracy,
						Altitude:  altitude,
					})
				}
			}
		}
	}

	return data
}

func ParseSpeed(datasets []*fitness.Dataset) []SpeedStruct {
	var data []SpeedStruct

	// fmt.Println("Speed: datasets:", datasets)

	for _, dataset := range datasets {

		// fmt.Println("Speed: dataset:", dataset)
		// fmt.Println("Speed: datasets: DataSourceId:", dataset.DataSourceId)
		// fmt.Println("Speed: datasets: MaxEndTimeNs:", dataset.MaxEndTimeNs)
		// fmt.Println("Speed: datasets: MinStartTimeNs:", dataset.MinStartTimeNs)
		// fmt.Println("Speed: datasets: NextPageToken:", dataset.NextPageToken)
		// fmt.Println("Speed: datasets: Point:", dataset.Point)
		// fmt.Println("Speed: datasets: ForceSendFields:", dataset.ForceSendFields)
		// fmt.Println("Speed: datasets: NullFields:", dataset.NullFields)

		for _, point := range dataset.Point {
			// Extract relevant fields from the dataset point
			startTime := time.Unix(0, point.StartTimeNanos*int64(time.Millisecond))
			endTime := time.Unix(0, point.EndTimeNanos*int64(time.Millisecond))

			// fmt.Println("Speed: point:", point)
			// fmt.Println("Speed: len(point.Value):", len(point.Value))

			// Loop over the values in the dataset point
			for _, value := range point.Value {
				// fmt.Println("Speed: Value:", value)
				// fmt.Println("Speed: Value.FpVal:", value.FpVal)
				// fmt.Println("Speed: Value.IntVal:", value.IntVal)
				// fmt.Println("Speed: Value.MapVal:", value.MapVal)
				// fmt.Println("Speed: Value.StringVal:", value.StringVal)
				// fmt.Println("Speed: Value.ForceSendFields:", value.ForceSendFields)
				// fmt.Println("Speed: Value.NullFields:", value.NullFields)

				// Create a new WeightStruct and append it to the data slice
				data = append(data, SpeedStruct{
					StartTime: startTime,
					EndTime:   endTime,
					Speed:     value.FpVal,
				})
			}
		}
	}

	return data
}

// ParseHeartRateBPM function converts the `Google Fit` heart rate (bpm) data into usable format for our app.
func ParseHeartRateBPM(datasets []*fitness.Dataset) []HeartRateBPMStruct {
	var data []HeartRateBPMStruct

	for _, ds := range datasets {
		var value float64
		for _, p := range ds.Point {
			// // For debugging purposes only.
			// fmt.Println("ComputationTimeMillis:", p.ComputationTimeMillis)
			// fmt.Println("DataTypeName:", p.DataTypeName)
			// fmt.Println("EndTimeNanos:", p.EndTimeNanos)
			// fmt.Println("ModifiedTimeMillis:", p.ModifiedTimeMillis)
			// fmt.Println("OriginDataSourceId:", p.OriginDataSourceId)
			// fmt.Println("StartTimeNanos:", p.StartTimeNanos)
			// fmt.Println("Value[parent]:", p.Value)
			// for _, v := range p.Value {
			// 	fmt.Println("Value[child]:", v)
			// }
			// fmt.Println("ForceSendFields:", p.ForceSendFields)
			// fmt.Println("NullFields:", p.NullFields)
			// fmt.Println()

			for _, v := range p.Value {
				// // For debugging purposes only.
				// fmt.Println("v:FpVal:", v.FpVal)
				// fmt.Println("v:IntVal:", v.IntVal)
				// fmt.Println("v:MapVal:", v.MapVal)
				// fmt.Println("v:StringVal:", v.StringVal)
				// fmt.Println("v:ForceSendFields:", v.ForceSendFields)
				// fmt.Println("v:NullFields:", v.NullFields)

				valueString := fmt.Sprintf("%.3f", v.FpVal)
				value, _ = strconv.ParseFloat(valueString, 64)
			}
			var row HeartRateBPMStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			row.BPM = int(value)
			data = append(data, row)
		}
	}
	return data
}

// ParseHydration function converts the `Google Fit` hydration data into usable format for our app. Special thanks to: https://github.com/bronnika/devto-google-fit/blob/main/google-api/parse.go#L103
func ParseHydration(datasets []*fitness.Dataset) []HydrationStruct {
	var data []HydrationStruct

	for _, ds := range datasets {
		var value float64
		for _, p := range ds.Point {
			// fmt.Println(p.DataTypeName) // For debugging purposes only.
			for _, v := range p.Value {
				valueString := fmt.Sprintf("%.3f", v.FpVal)
				value, _ = strconv.ParseFloat(valueString, 64)
			}
			var row HydrationStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			// liters to milliliters
			row.Volume = int(value * 1000)
			data = append(data, row)
		}
	}
	return data
}

// Special thanks to: https://github.com/bronnika/devto-google-fit/blob/main/google-api/parse.go#L124
func ParseNutrition(datasets []*fitness.Dataset) []NutritionStruct {
	var data []NutritionStruct

	for _, ds := range datasets {
		for _, p := range ds.Point {
			var row NutritionStruct
			for _, mapVal := range p.Value[0].MapVal {
				// there we can get more data (s.t. fat, carbs, protein, etc.) if it exists
				if mapVal.Key == NutrientCalories {
					row.Calories = int(mapVal.Value.FpVal)
				}
			}
			row.Type = MealType[int(p.Value[1].IntVal)]
			row.Name = p.Value[2].StringVal
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)

			data = append(data, row)
		}
	}
	return data
}

func ParseBloodGlucose(datasets []*fitness.Dataset) []BloodGlucoseStruct {
	var data []BloodGlucoseStruct

	for _, dataset := range datasets {
		// fmt.Println("OxygenSaturation: DataSourceId:", dataset.DataSourceId)
		// fmt.Println("OxygenSaturation: MaxEndTimeNs:", dataset.MaxEndTimeNs)
		// fmt.Println("OxygenSaturation: MinStartTimeNs:", dataset.MinStartTimeNs)
		// fmt.Println("OxygenSaturation: NextPageToken:", dataset.NextPageToken)
		// fmt.Println("OxygenSaturation: ForceSendFields:", dataset.ForceSendFields)
		// fmt.Println("OxygenSaturation: NullFields:", dataset.NullFields)

		for _, point := range dataset.Point {
			// Extract relevant fields from the dataset point
			startTime := time.Unix(0, point.StartTimeNanos*int64(time.Millisecond))
			endTime := time.Unix(0, point.EndTimeNanos*int64(time.Millisecond))
			// fmt.Println("OxygenSaturation: Point:", point)

			if len(point.Value) > 4 {
				bloodGlucoseLevel := point.Value[0]
				// fmt.Println("BloodGlucose: point.Value[0]: bloodGlucoseLevel:", bloodGlucoseLevel)
				mealType := point.Value[1]
				// fmt.Println("BloodGlucose: point.Value[1]: mealType:", mealType)
				temporalRelationToSleep := point.Value[2]
				// fmt.Println("BloodGlucose: point.Value[2]: temporalRelationToSleep:", temporalRelationToSleep)
				specimenSource := point.Value[3]
				// fmt.Println("OxygenSaturation: point.Value[3]: specimenSource:", specimenSource)

				//
				// Create a new BloodGlucoseStruct and append it to the data slice
				//
				data = append(data, BloodGlucoseStruct{
					StartTime:               startTime,
					EndTime:                 endTime,
					BloodGlucoseLevel:       bloodGlucoseLevel.FpVal,
					TemporalRelationToMeal:  0, //TODO - See https://developers.google.com/fit/datatypes/health#blood_glucose
					MealType:                mealType.IntVal,
					TemporalRelationToSleep: temporalRelationToSleep.IntVal,
					SpecimenSource:          specimenSource.IntVal,
				})
			}
		}
	}

	return data
}

func ParseBloodPressure(datasets []*fitness.Dataset) []BloodPressureStruct {
	var data []BloodPressureStruct

	for _, dataset := range datasets {
		// fmt.Println("OxygenSaturation: DataSourceId:", dataset.DataSourceId)
		// fmt.Println("OxygenSaturation: MaxEndTimeNs:", dataset.MaxEndTimeNs)
		// fmt.Println("OxygenSaturation: MinStartTimeNs:", dataset.MinStartTimeNs)
		// fmt.Println("OxygenSaturation: NextPageToken:", dataset.NextPageToken)
		// fmt.Println("OxygenSaturation: ForceSendFields:", dataset.ForceSendFields)
		// fmt.Println("OxygenSaturation: NullFields:", dataset.NullFields)

		for _, point := range dataset.Point {
			// Extract relevant fields from the dataset point
			startTime := time.Unix(0, point.StartTimeNanos*int64(time.Millisecond))
			endTime := time.Unix(0, point.EndTimeNanos*int64(time.Millisecond))
			// fmt.Println("OxygenSaturation: Point:", point)

			if len(point.Value) > 3 {
				systolic := point.Value[0]
				// fmt.Println("BloodPressure: point.Value[0]: systolic:", systolic)
				diastolic := point.Value[1]
				// fmt.Println("BloodPressure: point.Value[1]: diastolic:", diastolic)
				bodyPosition := point.Value[2]
				// fmt.Println("BloodPressure: point.Value[2]: bodyPosition:", bodyPosition)
				measurementLocation := point.Value[3]
				// fmt.Println("OxygenSaturation: point.Value[3]: measurementLocation:", measurementLocation)

				//
				// Create a new BloodPressureStruct and append it to the data slice
				//
				data = append(data, BloodPressureStruct{
					StartTime:           startTime,
					EndTime:             endTime,
					Systolic:            systolic.FpVal,
					Diastolic:           diastolic.FpVal,
					BodyPosition:        bodyPosition.IntVal,
					MeasurementLocation: measurementLocation.IntVal,
					// OxygenSaturationMeasurementMethod: oxygenSaturationMeasurementMethod.IntVal,
				})
			}
		}
	}

	return data
}

func ParseBodyFatPercentage(datasets []*fitness.Dataset) []BodyFatPercentageStruct {
	var data []BodyFatPercentageStruct

	for _, dataset := range datasets {
		for _, point := range dataset.Point {
			// Extract relevant fields from the dataset point
			startTime := time.Unix(0, point.StartTimeNanos*int64(time.Millisecond))
			endTime := time.Unix(0, point.EndTimeNanos*int64(time.Millisecond))

			for _, point := range dataset.Point {
				// // Extract relevant fields from the dataset point
				// startTime := time.Unix(0, point.StartTimeNanos*int64(time.Millisecond))
				// endTime := time.Unix(0, point.EndTimeNanos*int64(time.Millisecond))
				// // fmt.Println("BodyTemperature: Point:", point)

				if len(point.Value) > 0 {
					percentage := point.Value[0]

					//
					// Create a new BodyTemperature and append it to the data slice
					//
					data = append(data, BodyFatPercentageStruct{
						StartTime:  startTime,
						EndTime:    endTime,
						Percentage: percentage.FpVal,
					})
				}
			}
		}
	}

	return data
}

func ParseBodyTemperature(datasets []*fitness.Dataset) []BodyTemperatureStruct {
	var data []BodyTemperatureStruct

	for _, dataset := range datasets {
		for _, point := range dataset.Point {
			// Extract relevant fields from the dataset point
			startTime := time.Unix(0, point.StartTimeNanos*int64(time.Millisecond))
			endTime := time.Unix(0, point.EndTimeNanos*int64(time.Millisecond))

			for _, point := range dataset.Point {
				// // Extract relevant fields from the dataset point
				// startTime := time.Unix(0, point.StartTimeNanos*int64(time.Millisecond))
				// endTime := time.Unix(0, point.EndTimeNanos*int64(time.Millisecond))
				// // fmt.Println("BodyTemperature: Point:", point)

				if len(point.Value) > 1 {
					bodyTemperature := point.Value[0]
					// fmt.Println("BodyTemperature: point.Value[0]:", bodyTemperature)
					measurementLocation := point.Value[1]
					// // fmt.Println("BodyTemperature: point.Value[1]:", measurementLocation)

					//
					// Create a new BodyTemperature and append it to the data slice
					//
					data = append(data, BodyTemperatureStruct{
						StartTime:           startTime,
						EndTime:             endTime,
						BodyTemperature:     bodyTemperature.FpVal,
						MeasurementLocation: measurementLocation.IntVal,
					})
				}
			}
		}
	}

	return data
}

func ParseHeight(datasets []*fitness.Dataset) []HeightStruct {
	var data []HeightStruct

	for _, dataset := range datasets {
		for _, point := range dataset.Point {
			// Extract relevant fields from the dataset point
			startTime := time.Unix(0, point.StartTimeNanos*int64(time.Millisecond))
			endTime := time.Unix(0, point.EndTimeNanos*int64(time.Millisecond))

			// Loop over the values in the dataset point
			for _, value := range point.Value {
				// fmt.Println("weight: Value:", value)
				// fmt.Println("weight: Value.FpVal:", value.FpVal)
				// fmt.Println("weight: Value.IntVal:", value.IntVal)
				// fmt.Println("weight: Value.MapVal:", value.MapVal)
				// fmt.Println("weight: Value.StringVal:", value.StringVal)
				// fmt.Println("weight: Value.ForceSendFields:", value.ForceSendFields)
				// fmt.Println("weight: Value.NullFields:", value.NullFields)

				valueString := fmt.Sprintf("%.2f", value.FpVal)
				v, _ := strconv.ParseFloat(valueString, 64)

				// Conver to centimeters
				heightInCentimeters := float64(v * 100)

				// Create a new WeightStruct and append it to the data slice
				data = append(data, HeightStruct{
					StartTime: startTime,
					EndTime:   endTime,
					Height:    heightInCentimeters,
				})
			}
		}
	}

	return data
}

func ParseOxygenSaturation(datasets []*fitness.Dataset) []OxygenSaturationStruct {
	var data []OxygenSaturationStruct

	for _, dataset := range datasets {
		// fmt.Println("OxygenSaturation: DataSourceId:", dataset.DataSourceId)
		// fmt.Println("OxygenSaturation: MaxEndTimeNs:", dataset.MaxEndTimeNs)
		// fmt.Println("OxygenSaturation: MinStartTimeNs:", dataset.MinStartTimeNs)
		// fmt.Println("OxygenSaturation: NextPageToken:", dataset.NextPageToken)
		// fmt.Println("OxygenSaturation: ForceSendFields:", dataset.ForceSendFields)
		// fmt.Println("OxygenSaturation: NullFields:", dataset.NullFields)

		for _, point := range dataset.Point {
			// Extract relevant fields from the dataset point
			startTime := time.Unix(0, point.StartTimeNanos*int64(time.Millisecond))
			endTime := time.Unix(0, point.EndTimeNanos*int64(time.Millisecond))
			// fmt.Println("OxygenSaturation: Point:", point)

			if len(point.Value) > 3 {
				oxygenSaturation := point.Value[0]
				// fmt.Println("OxygenSaturation: point.Value[0]:", oxygenSaturation)
				supplementalOxygenFlowRate := point.Value[1]
				// fmt.Println("OxygenSaturation: point.Value[1]:", supplementalOxygenFlowRate)
				oxygenTherapyAdministrationMode := point.Value[2]
				// fmt.Println("OxygenSaturation: point.Value[2]:", oxygenTherapyAdministrationMode)
				oxygenSaturationSystem := point.Value[3]
				// fmt.Println("OxygenSaturation: point.Value[3]:", oxygenSaturationSystem)
				oxygenSaturationMeasurementMethod := point.Value[4]
				// fmt.Println("OxygenSaturation: point.Value[4]:", oxygenSaturationMeasurementMethod)

				//
				// Create a new OxygenSaturationStruct and append it to the data slice
				//
				data = append(data, OxygenSaturationStruct{
					StartTime:                         startTime,
					EndTime:                           endTime,
					OxygenSaturation:                  oxygenSaturation.FpVal,
					SupplementalOxygenFlowRate:        supplementalOxygenFlowRate.FpVal,
					OxygenTherapyAdministrationMode:   oxygenTherapyAdministrationMode.IntVal,
					OxygenSaturationSystem:            oxygenSaturationSystem.IntVal,
					OxygenSaturationMeasurementMethod: oxygenSaturationMeasurementMethod.IntVal,
				})
			}
		}
	}

	return data
}

func ParseSleep(datasets []*fitness.Dataset) []SleepStruct {
	var data []SleepStruct

	for _, dataset := range datasets {
		for _, point := range dataset.Point {
			// Extract relevant fields from the dataset point
			startTime := time.Unix(0, point.StartTimeNanos*int64(time.Millisecond))
			endTime := time.Unix(0, point.EndTimeNanos*int64(time.Millisecond))

			// fmt.Println("sleep: point:", point)

			// Loop over the values in the dataset point
			for _, value := range point.Value {
				// fmt.Println("sleep: Value:", value)
				// fmt.Println("sleep: Value.FpVal:", value.FpVal)
				// fmt.Println("sleep: Value.IntVal:", value.IntVal)
				// fmt.Println("sleep: Value.MapVal:", value.MapVal)
				// fmt.Println("sleep: Value.StringVal:", value.StringVal)
				// fmt.Println("sleep: Value.ForceSendFields:", value.ForceSendFields)
				// fmt.Println("sleep: Value.NullFields:", value.NullFields)

				// Create a new SleepStruct and append it to the data slice
				data = append(data, SleepStruct{
					StartTime:        startTime,
					EndTime:          endTime,
					SleepSegmentType: int(value.IntVal),
				})
			}
		}
	}

	return data
}

func ParseWeight(datasets []*fitness.Dataset) []WeightStruct {
	var data []WeightStruct

	for _, dataset := range datasets {
		for _, point := range dataset.Point {
			// Extract relevant fields from the dataset point
			startTime := time.Unix(0, point.StartTimeNanos*int64(time.Millisecond))
			endTime := time.Unix(0, point.EndTimeNanos*int64(time.Millisecond))

			// fmt.Println("weight: point:", point)

			// Loop over the values in the dataset point
			for _, value := range point.Value {
				// fmt.Println("weight: Value:", value)
				// fmt.Println("weight: Value.FpVal:", value.FpVal)
				// fmt.Println("weight: Value.IntVal:", value.IntVal)
				// fmt.Println("weight: Value.MapVal:", value.MapVal)
				// fmt.Println("weight: Value.StringVal:", value.StringVal)
				// fmt.Println("weight: Value.ForceSendFields:", value.ForceSendFields)
				// fmt.Println("weight: Value.NullFields:", value.NullFields)

				// Create a new WeightStruct and append it to the data slice
				data = append(data, WeightStruct{
					StartTime: startTime,
					EndTime:   endTime,
					Weight:    value.FpVal,
				})
			}
		}
	}

	return data
}
