package google

import "time"

// type ActivityStruct struct {
// 	Amount       int64     `json:"amount"`
// 	Date         time.Time `json:"date"`
// 	ActivityType string    `json:"activity_type"`
// 	Tracker      string    `json:"tracker"`
// 	Additional   string    `json:"additional"`
// 	StartTime    time.Time `json:"start_time"`
// 	EndTime      time.Time `json:"end_time"`
// }
//
// type Measurement struct {
// 	Type      string    `json:"type"`
// 	StartTime time.Time `json:"start_time"`
// 	EndTime   time.Time `json:"end_time"`
// 	Value     float64   `json:"value"`
// 	AvValue   float64   `json:"av_value"`
// 	MaxValue  float64   `json:"max_value"`
// 	MinValue  float64   `json:"min_value"`
// }
//
// type ActiveMinute struct {
// 	Type      string    `json:"type"`
// 	StartTime time.Time `json:"start_time"`
// 	EndTime   time.Time `json:"end_time"`
// 	Value     int64     `json:"value"`
// }

type ActivitySegmentStruct struct {
	ActivityTypeID    int64     `bson:"activity_type_id,omitempty" json:"activity_type_id"`
	ActivityTypeLabel string    `bson:"activity_type_label,omitempty" json:"activity_type_label"`
	DurationInMinutes int       `bson:"duration_in_minutes,omitempty" json:"duration_in_minutes"`
	DurationInHours   int       `bson:"duration_in_hours,omitempty" json:"duration_in_hours"`
	StartTime         time.Time `bson:"start_time,omitempty" json:"start_time"`
	EndTime           time.Time `bson:"end_time,omitempty" json:"end_time"`
}

// BasalMetabolicRate defines the power type provided by `Google Fit`.
type BasalMetabolicRateStruct struct {
	StartTime time.Time `bson:"start_time,omitempty" json:"start_time"`
	EndTime   time.Time `bson:"end_time,omitempty" json:"end_time"`
	Calories  float64   `bson:"calories,omitempty" json:"calories"` // kcal per day
}

// Calories defines the calories burned type provided by `Google Fit`.
type CaloriesBurnedStruct struct {
	StartTime time.Time `bson:"start_time,omitempty" json:"start_time"`
	EndTime   time.Time `bson:"end_time,omitempty" json:"end_time"`
	Calories  float64   `bson:"calories,omitempty" json:"calories"` // kcal
}

// CyclingPedalingCadence defines the power type provided by `Google Fit`.
type CyclingPedalingCadenceStruct struct {
	StartTime time.Time `bson:"start_time,omitempty" json:"start_time"`
	EndTime   time.Time `bson:"end_time,omitempty" json:"end_time"`
	RPM       float64   `bson:"rpm,omitempty" json:"rpm"` // RPM
}

// CyclingPedalingCumulative defines the power type provided by `Google Fit`.
type CyclingPedalingCumulativeStruct struct {
	StartTime   time.Time `bson:"start_time,omitempty" json:"start_time"`
	EndTime     time.Time `bson:"end_time,omitempty" json:"end_time"`
	Revolutions int       `bson:"revolutions,omitempty" json:"revolutions"` // count
}

// HeartPoints defines the power type provided by `Google Fit`.
type HeartPointsStruct struct {
	StartTime time.Time `bson:"start_time,omitempty" json:"start_time"`
	EndTime   time.Time `bson:"end_time,omitempty" json:"end_time"`
	Intensity float64   `bson:"intensity,omitempty" json:"intensity"` // Heart Points
}

// MoveMinutesStruct defines the power type provided by `Google Fit`.
type MoveMinutesStruct struct {
	StartTime time.Time `bson:"start_time,omitempty" json:"start_time"`
	EndTime   time.Time `bson:"end_time,omitempty" json:"end_time"`
	Duration  int       `bson:"duration,omitempty" json:"duration"` // milliseconds
}

// Power defines the power type provided by `Google Fit`.
type PowerStruct struct {
	StartTime time.Time `bson:"start_time,omitempty" json:"start_time"`
	EndTime   time.Time `bson:"end_time,omitempty" json:"end_time"`
	Watts     float64   `bson:"watts,omitempty" json:"watts"` // watts
}

type StepCountDeltaStruct struct {
	Steps     int       `bson:"steps,omitempty" json:"steps"` // Count
	StartTime time.Time `bson:"start_time,omitempty" json:"start_time"`
	EndTime   time.Time `bson:"end_time,omitempty" json:"end_time"`
}

type StepCountCadenceStruct struct {
	RPM       int       `bson:"rpm,omitempty" json:"rpm"` // steps/minute
	StartTime time.Time `bson:"start_time,omitempty" json:"start_time"`
	EndTime   time.Time `bson:"end_time,omitempty" json:"end_time"`
}

// WorkoutStruct is deprecated. Note: https://9to5google.com/2020/11/30/google-fit-latest-update-removes-advanced-weight-training-features-on-wear-os/.
type WorkoutStruct struct {
	Exercise          int       `bson:"exercise,omitempty" json:"exercise"`
	Repetitions       int       `bson:"repetitions,omitempty" json:"repetitions"`
	ResistanceType    int       `bson:"resistance_type ,omitempty" json:"resistance_type"`
	Resistance        float64   `bson:"resistance,omitempty" json:"resistance"`
	DurationInMinutes int       `bson:"duration_in_minutes,omitempty" json:"duration_in_minutes"`
	DurationInHours   int       `bson:"duration_in_hours,omitempty" json:"duration_in_hours"`
	StartTime         time.Time `bson:"start_time,omitempty" json:"start_time"`
	EndTime           time.Time `bson:"end_time,omitempty" json:"end_time"`
}

type HeartRateBPMStruct struct {
	BPM       int       `bson:"bpm,omitempty" json:"bpm"` // bpm
	StartTime time.Time `bson:"start_time,omitempty" json:"start_time"`
	EndTime   time.Time `bson:"end_time,omitempty" json:"end_time"`
}

// HydrationStruct defines the hydration data type provided by `Google Fit`. Special thanks to https://github.com/bronnika/devto-google-fit/blob/main/models/models.go#L52C1-L56C2
type HydrationStruct struct {
	Volume    int       `bson:"volume,omitempty" json:"volume"` // liters
	StartTime time.Time `bson:"start_time,omitempty" json:"start_time"`
	EndTime   time.Time `bson:"end_time,omitempty" json:"end_time"`
}

// NutritionStruct defines the nutritional data type provided by `Google Fit`. Special thanks to https://github.com/bronnika/devto-google-fit/blob/main/models/models.go#L58C1-L77C2
type NutritionStruct struct {
	Name               string    `bson:"name,omitempty" json:"name" default:" "`
	Type               string    `bson:"type,omitempty" json:"type"`
	Calories           int       `bson:"calories,omitempty" json:"calories" default:"0"`
	TotalFat           float64   `bson:"total_fat,omitempty" json:"total_fat" default:"0"`
	SaturatedFat       float64   `bson:"saturated_fat,omitempty" json:"saturated_fat" default:"0"`
	UnsaturatedFat     float64   `bson:"unsaturated_fat,omitempty" json:"unsaturated_fat" default:"0"`
	PolyunsaturatedFat float64   `bson:"polyunsaturated_fat,omitempty" json:"polyunsaturated_fat" default:"0"`
	TransFat           float64   `bson:"trans_fat,omitempty" json:"trans_fat" default:"0"`
	Cholesterol        float64   `bson:"cholesterol,omitempty" json:"cholesterol" default:"0"`
	Sodium             float64   `bson:"sodium,omitempty" json:"sodium" default:"0"`
	Potassium          float64   `bson:"potassium,omitempty" json:"potassium" default:"0"`
	Carbohydrates      float64   `bson:"carbohydrates,omitempty" json:"carbohydrates" default:"0"`
	Fibre              float64   `bson:"fibre,omitempty" json:"fibre" default:"0"`
	Sugar              float64   `bson:"sugar,omitempty" json:"sugar" default:"0"`
	Protein            float64   `bson:"protein,omitempty" json:"protein" default:"0"`
	StartTime          time.Time `bson:"start_time,omitempty" json:"start_time" default:"0"`
	EndTime            time.Time `bson:"end_time,omitempty" json:"end_time" default:"0"`
	Amount             int       `bson:"amount,omitempty" json:"amount" default:"0"`
}
