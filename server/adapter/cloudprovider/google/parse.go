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
			// liters to milliliters
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
			// liters to milliliters
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
			// liters to milliliters
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
			// liters to milliliters
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
			// liters to milliliters
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
			// liters to milliliters
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
			// liters to milliliters
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
			// liters to milliliters
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
			// liters to milliliters
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
			var row SpeedStruct
			row.StartTime = NanosToTime(p.StartTimeNanos)
			row.EndTime = NanosToTime(p.EndTimeNanos)
			// liters to milliliters
			row.Speed = int(value)
			data = append(data, row)
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
			// liters to milliliters
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

func getIntValue(locationValue []*fitness.ValueMapValEntry, key string) int {
	for _, entry := range locationValue {
		if entry.Key == key {
			return int(entry.Value.FpVal)
		}
	}
	return 0
}

func getFloat64Value(locationValue []*fitness.ValueMapValEntry, key string) float64 {
	for _, entry := range locationValue {
		if entry.Key == key {
			return entry.Value.FpVal
		}
	}
	return 0.0
}

func ParseBloodGlucose(datasets []*fitness.Dataset) []BloodGlucoseStruct {
	var data []BloodGlucoseStruct

	for _, dataset := range datasets {
		for _, point := range dataset.Point {
			// Extract relevant fields from the dataset point
			startTime := time.Unix(0, point.StartTimeNanos*int64(time.Millisecond))
			endTime := time.Unix(0, point.EndTimeNanos*int64(time.Millisecond))

			// Loop over the values in the dataset point
			for _, value := range point.Value {
				if bloodGlucoseValue := value.MapVal; bloodGlucoseValue != nil {
					//TODO: IMPLEMENT THE CODE BELOW...

					fmt.Println("bloodGlucoseValue:", bloodGlucoseValue)

					// Extract latitude, longitude, accuracy, and altitude from the locationValue
					bloodGlucoseLevel := getIntValue(bloodGlucoseValue, "blood glucose level ")
					// longitude := getLocationValue(locationValue, "longitude")
					// accuracy := getLocationValue(locationValue, "accuracy")
					// altitude := getLocationValue(locationValue, "altitude")

					// Create a new BloodGlucoseStruct and append it to the data slice
					data = append(data, BloodGlucoseStruct{
						StartTime:         startTime,
						EndTime:           endTime,
						BloodGlucoseLevel: bloodGlucoseLevel,
						// Longitude: longitude,
						// Accuracy:  accuracy,
						// Altitude:  altitude,
					})
				}
			}
		}
	}

	return data
}

func ParseBloodPressure(datasets []*fitness.Dataset) []BloodPressureStruct {
	var data []BloodPressureStruct

	for _, dataset := range datasets {
		for _, point := range dataset.Point {
			// Extract relevant fields from the dataset point
			startTime := time.Unix(0, point.StartTimeNanos*int64(time.Millisecond))
			endTime := time.Unix(0, point.EndTimeNanos*int64(time.Millisecond))

			// Loop over the values in the dataset point
			for _, value := range point.Value {
				if bloodPressureValue := value.MapVal; bloodPressureValue != nil {
					//TODO: IMPLEMENT THE CODE BELOW...

					fmt.Println("bloodPressureValue:", bloodPressureValue)

					// Extract latitude, longitude, accuracy, and altitude from the locationValue
					systolic := getFloat64Value(bloodPressureValue, "systolic")
					// longitude := getLocationValue(locationValue, "longitude")
					// accuracy := getLocationValue(locationValue, "accuracy")
					// altitude := getLocationValue(locationValue, "altitude")

					// Create a new BloodPressureStruct and append it to the data slice
					data = append(data, BloodPressureStruct{
						StartTime: startTime,
						EndTime:   endTime,
						Systolic:  systolic,
						// Longitude: longitude,
						// Accuracy:  accuracy,
						// Altitude:  altitude,
					})
				}
			}
		}
	}

	return data
}
