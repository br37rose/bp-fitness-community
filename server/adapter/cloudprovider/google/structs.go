package google

import "time"

type Measurement struct {
	Type      string    `json:"type"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Value     float64   `json:"value"`
	AvValue   float64   `json:"av_value"`
	MaxValue  float64   `json:"max_value"`
	MinValue  float64   `json:"min_value"`
}

type ActiveMinute struct {
	Type      string    `json:"type"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Value     int64     `json:"value"`
}
type Activity struct {
	Amount       int64     `json:"amount"`
	Date         time.Time `json:"date"`
	ActivityType string    `json:"activity_type"`
	Tracker      string    `json:"tracker"`
	Additional   string    `json:"additional"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
}

type ActivitySegments struct {
	ActivityType string    `json:"activity_type"`
	Minutes      float64   `json:"minutes"`
	SessionNum   int64     `json:"session_num"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
}

// HydrationStruct defines the hydration data type provided by `Google Fit`. Special thanks to https://github.com/bronnika/devto-google-fit/blob/main/models/models.go#L52C1-L56C2
type Calories struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Amount    float64   `json:"amount"`
}

type HeartRateBPMStruct struct {
	Amount    int       `json:"amount"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// HydrationStruct defines the hydration data type provided by `Google Fit`. Special thanks to https://github.com/bronnika/devto-google-fit/blob/main/models/models.go#L52C1-L56C2
type HydrationStruct struct {
	Amount    int       `json:"amount"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// NutritionStruct defines the nutritional data type provided by `Google Fit`. Special thanks to https://github.com/bronnika/devto-google-fit/blob/main/models/models.go#L58C1-L77C2
type NutritionStruct struct {
	Name               string    `json:"name" default:" "`
	Type               string    `json:"type"`
	Calories           int       `json:"calories" default:"0"`
	TotalFat           float64   `json:"total_fat" default:"0"`
	SaturatedFat       float64   `json:"saturated_fat" default:"0"`
	UnsaturatedFat     float64   `json:"unsaturated_fat" default:"0"`
	PolyunsaturatedFat float64   `json:"polyunsaturated_fat" default:"0"`
	TransFat           float64   `json:"trans_fat" default:"0"`
	Cholesterol        float64   `json:"cholesterol" default:"0"`
	Sodium             float64   `json:"sodium" default:"0"`
	Potassium          float64   `json:"potassium" default:"0"`
	Carbohydrates      float64   `json:"carbohydrates" default:"0"`
	Fibre              float64   `json:"fibre" default:"0"`
	Sugar              float64   `json:"sugar" default:"0"`
	Protein            float64   `json:"protein" default:"0"`
	StartTime          time.Time `json:"start_time" default:"0"`
	EndTime            time.Time `json:"end_time" default:"0"`
	Amount             int       `json:"amount" default:"0"`
}
