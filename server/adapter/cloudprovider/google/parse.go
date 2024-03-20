package google

import (
	"fmt"
	"strconv"

	"google.golang.org/api/fitness/v1"
)

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
			row.Amount = float64(value)
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
			row.Amount = int(value)
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
			row.Amount = int(value)
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
			row.Amount = int(value * 1000)
			data = append(data, row)
		}
	}
	return data
}
