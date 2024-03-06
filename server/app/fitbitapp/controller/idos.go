package controller

import (
	"encoding/json"
	"strconv"
	"strings"
)

//------------- Activities Calorie Intraday ---------------//

type FitBitActivitiesCalorieIntradayDataRequestIDO struct {
	Value float64 `json:"value,omitempty"`
	Time  string  `json:"time,omitempty"`
	Level int64   `json:"level,omitempty"`
	METS  int64   `json:"mets,omitempty"`
}

type FitBitActivitiesCalorieIntradayRequestIDO struct {
	Dataset         []FitBitActivitiesCalorieIntradayDataRequestIDO `json:"dataset,omitempty"`
	DatasetInterval int64                                           `json:"datasetInterval,omitempty"`
	DatasetType     string                                          `json:"datasetType,omitempty"`
}

type FitBitActivitiesCalorieRequestIDO struct {
	Value    string `json:"value,omitempty"`
	DateTime string `json:"dateTime,omitempty"`
}

type FitBitActivityCalariesRawData struct {
	ActivitiesCalories          []FitBitActivitiesCalorieRequestIDO       `json:"activities-calories,omitempty"`
	ActivitiesCaloriesIntradays FitBitActivitiesCalorieIntradayRequestIDO `json:"activities-calories-intraday,omitempty"`
}

func UnmarshalFitBitActivityCalariesRawData(jsonData string) (*FitBitActivityCalariesRawData, error) {
	d := &FitBitActivityCalariesRawData{}

	err := json.Unmarshal([]byte(jsonData), &d) // Read the JSON string and convert it into our golang stuct else we get an error.
	if err != nil {
		return nil, err
	}

	return d, nil
}

//------------- Activities Steps Intraday ---------------//

type FitBitActivitiesStepIntradayDataRequestIDO struct {
	Value float64 `json:"value,omitempty"`
	Time  string  `json:"time,omitempty"`
	Level int64   `json:"level,omitempty"`
	METS  int64   `json:"mets,omitempty"`
}

type FitBitActivitiesStepIntradayRequestIDO struct {
	Dataset         []FitBitActivitiesStepIntradayDataRequestIDO `json:"dataset,omitempty"`
	DatasetInterval int64                                        `json:"datasetInterval,omitempty"`
	DatasetType     string                                       `json:"datasetType,omitempty"`
}

type FitBitActivitiesStepRequestIDO struct {
	Value    string `json:"value,omitempty"`
	DateTime string `json:"dateTime,omitempty"`
}

type FitBitStepStepsRawData struct {
	ActivitiesSteps          []FitBitActivitiesStepRequestIDO       `json:"activities-steps,omitempty"`
	ActivitiesStepsIntradays FitBitActivitiesStepIntradayRequestIDO `json:"activities-steps-intraday,omitempty"`
}

func UnmarshalFitBitStepStepsRawData(jsonData string) (*FitBitStepStepsRawData, error) {
	d := &FitBitStepStepsRawData{}

	err := json.Unmarshal([]byte(jsonData), &d) // Read the JSON string and convert it into our golang stuct else we get an error.
	if err != nil {
		return nil, err
	}

	return d, nil
}

//------------- Activities HeartRates Intraday ---------------//

type FitBitActivitiesHeartRateIntradayDataRequestIDO struct {
	Value float64 `json:"value,omitempty"`
	Time  string  `json:"time,omitempty"`
	Level int64   `json:"level,omitempty"`
	METS  int64   `json:"mets,omitempty"`
}

type FitBitActivitiesHeartRateIntradayRequestIDO struct {
	Dataset         []FitBitActivitiesHeartRateIntradayDataRequestIDO `json:"dataset,omitempty"`
	DatasetInterval int64                                             `json:"datasetInterval,omitempty"`
	DatasetType     string                                            `json:"datasetType,omitempty"`
}

type FitBitActivitiesHeartRateRequestIDO struct { // Not interested
	// Value    float64 `json:"value,omitempty"`
	DateTime string `json:"dateTime,omitempty"`
}

type FitBitHeartRatesRawData struct {
	ActivitiesHeartRates          []FitBitActivitiesHeartRateRequestIDO        `json:"activities-heart,omitempty"` // Not interested
	ActivitiesHeartRatesIntradays *FitBitActivitiesHeartRateIntradayRequestIDO `json:"activities-heart-intraday,omitempty"`
}

func UnmarshalFitBitHeartRatesRawData(jsonData string) (*FitBitHeartRatesRawData, error) {
	d := &FitBitHeartRatesRawData{}

	err := json.Unmarshal([]byte(jsonData), &d) // Read the JSON string and convert it into our golang stuct else we get an error.
	if err != nil {
		return nil, err
	}

	return d, nil
}

//------------- Activities Breathing Rates Intraday ---------------//

type FitBitActivitiesBreathingRateIntradayDataInternalRequestIDO struct {
	BreathingRate float64 `json:"breathingRate,omitempty"`
}

type FitBitActivitiesBreathingRateIntradayDataRequestIDO struct {
	DeepSleepSummary  FitBitActivitiesBreathingRateIntradayDataInternalRequestIDO `json:"deepSleepSummary,omitempty"`
	RemSleepSummary   FitBitActivitiesBreathingRateIntradayDataInternalRequestIDO `json:"remSleepSummary,omitempty"`
	FullSleepSummary  FitBitActivitiesBreathingRateIntradayDataInternalRequestIDO `json:"fullSleepSummary,omitempty"`
	LightSleepSummary FitBitActivitiesBreathingRateIntradayDataInternalRequestIDO `json:"lightSleepSummary,omitempty"`
}

type FitBitActivitiesBreathingRateIntradayRequestIDO struct {
	Value    FitBitActivitiesBreathingRateIntradayDataRequestIDO `json:"value,omitempty"`
	DateTime string                                              `json:"dateTime,omitempty"`
}

// type FitBitActivitiesBreathingRateRequestIDO struct { // Not interested
// 	Value    string `json:"value,omitempty"`
// 	DateTime string `json:"dateTime,omitempty"`
// }

type FitBitBreathingRatesRawData struct {
	BR []FitBitActivitiesBreathingRateIntradayRequestIDO `json:"br,omitempty"`
}

func UnmarshalFitBitBreathingRatesRawData(jsonData string) (*FitBitBreathingRatesRawData, error) {
	d := &FitBitBreathingRatesRawData{}

	err := json.Unmarshal([]byte(jsonData), &d) // Read the JSON string and convert it into our golang stuct else we get an error.
	if err != nil {
		return nil, err
	}

	return d, nil
}

//------------- Heart Rate Variability Intraday ---------------//

type FitBitActivitiesHeartRateVariabilityIntradayDataInternalRequestIDO struct {
	RMSSD    float64 `json:"rmssd,omitempty"` //  Root Mean Square of Successive Differences (RMSSD)
	Coverage float64 `json:"coverage,omitempty"`
	HF       float64 `json:"hf,omitempty"`
	LF       float64 `json:"lf,omitempty"`
}

type FitBitActivitiesHeartRateVariabilityIntradayDataRequestIDO struct {
	Value  FitBitActivitiesHeartRateVariabilityIntradayDataInternalRequestIDO `json:"value,omitempty"`
	Minute string                                                             `json:"minute,omitempty"`
}

type FitBitActivitiesHeartRateVariabilityIntradayRequestIDO struct {
	Minutes  []FitBitActivitiesHeartRateVariabilityIntradayDataRequestIDO `json:"minutes,omitempty"`
	DateTime string                                                       `json:"dateTime,omitempty"`
}

type FitBitHeartRateVariabilitysRawData struct {
	HRV []FitBitActivitiesHeartRateVariabilityIntradayRequestIDO `json:"hrv,omitempty"`
}

func UnmarshalFitBitHeartRateVariabilitysRawData(jsonData string) (*FitBitHeartRateVariabilitysRawData, error) {
	d := &FitBitHeartRateVariabilitysRawData{}

	err := json.Unmarshal([]byte(jsonData), &d) // Read the JSON string and convert it into our golang stuct else we get an error.
	if err != nil {
		return nil, err
	}

	return d, nil
}

//------------- OxygenSaturation Intraday ---------------//

type FitBitActivitiesOxygenSaturationIntradayRequestIDO struct {
	Value  float64 `json:"value,omitempty"`
	Minute string  `json:"minute,omitempty"`
}

type FitBitOxygenSaturationsRawData struct {
	Minutes  []FitBitActivitiesOxygenSaturationIntradayRequestIDO `json:"minutes,omitempty"`
	DateTime string                                               `json:"dateTime,omitempty"`
}

func UnmarshalFitBitOxygenSaturationsRawData(jsonData string) (*FitBitOxygenSaturationsRawData, error) {
	d := &FitBitOxygenSaturationsRawData{}

	err := json.Unmarshal([]byte(jsonData), &d) // Read the JSON string and convert it into our golang stuct else we get an error.
	if err != nil {
		return nil, err
	}

	return d, nil
}

//------------- Sleep ---------------//

type FitBitSleepDayLevelSummaryItemIDO struct {
	Count   int64 `json:"count,omitempty"`
	Minutes int64 `json:"minutes,omitempty"`
}

type FitBitSleepDayLevelSummaryIDO struct {
	Asleep   *FitBitSleepDayLevelSummaryItemIDO `json:"asleep,omitempty"`
	Awake    *FitBitSleepDayLevelSummaryItemIDO `json:"awake,omitempty"`
	Restless *FitBitSleepDayLevelSummaryItemIDO `json:"restless,omitempty"`
}

type FitBitSleepDayLevelDatumIDO struct {
	DateTime string `json:"dateTime,omitempty"`
	Level    string `json:"level,omitempty"`
	Seconds  int64  `json:"seconds,omitempty"`
}

type FitBitSleepDayLevelIDO struct {
	Data      []*FitBitSleepDayLevelDatumIDO `json:"data,omitempty"`
	ShortData []*FitBitSleepDayLevelDatumIDO `json:"shortData,omitempty"`
	Summary   *FitBitSleepDayLevelSummaryIDO `json:"summary,omitempty"`
}

type FitBitSleepDayIDO struct {
	DateOfSleep         string                  `json:"dateOfSleep,omitempty"`
	Duration            int64                   `json:"duration,omitempty"`
	Efficiency          int64                   `json:"efficiency,omitempty"`
	EndTime             string                  `json:"endTime,omitempty"`
	InfoCode            int64                   `json:"infoCode,omitempty"`
	IsMainSleep         bool                    `json:"isMainSleep,omitempty"`
	Levels              *FitBitSleepDayLevelIDO `json:"levels,omitempty"`
	LogId               int64                   `json:"logId,omitempty"`
	MinutesAfterWakeup  int64                   `json:"minutesAfterWakeup,omitempty"`
	MinutesAsleep       int64                   `json:"minutesAsleep,omitempty"`
	MinutesAwake        int64                   `json:"minutesAwake,omitempty"`
	MinutesToFallAsleep int64                   `json:"minutesToFallAsleep,omitempty"`
	LogType             string                  `json:"logType,omitempty"`
	StartTime           string                  `json:"startTime,omitempty"`
	TimeInBed           int64                   `json:"timeInBed,omitempty"`
	Type                string                  `json:"type,omitempty"`
}

type FitBitSleepRawData struct {
	Sleep []FitBitSleepDayIDO `json:"sleep,omitempty"`
}

func UnmarshalFitBitSleepRawData(jsonData string) (*FitBitSleepRawData, error) {
	d := &FitBitSleepRawData{}

	err := json.Unmarshal([]byte(jsonData), &d) // Read the JSON string and convert it into our golang stuct else we get an error.
	if err != nil {
		return nil, err
	}

	return d, nil
}

//------------- Temperature ---------------//

type FitBitTempCoreIntervalIDO struct {
	DateTime string  `json:"dateTime,omitempty"`
	Value    float64 `json:"value,omitempty"`
}

type FitBitTemperatureRawData struct {
	TempCore []FitBitTempCoreIntervalIDO `json:"tempCore,omitempty"`
}

func UnmarshalFitBitTemperatureRawData(jsonData string) (*FitBitTemperatureRawData, error) {
	d := &FitBitTemperatureRawData{}

	err := json.Unmarshal([]byte(jsonData), &d) // Read the JSON string and convert it into our golang stuct else we get an error.
	if err != nil {
		return nil, err
	}

	return d, nil
}

//------------- Cardio Fitness Score (VO2 Max) ---------------//

type FitBitVO2MaxIntervalValueIDO struct {
	VO2Max string `json:"vo2Max,omitempty"`
}

type FitBitVO2MaxIntervalIDO struct {
	DateTime string                       `json:"dateTime,omitempty"`
	Value    FitBitVO2MaxIntervalValueIDO `json:"value,omitempty"`
	AvgValue float64                      `json:"avg_value,omitempty"`
}

type FitBitVO2MaxRawData struct {
	CardioScore []FitBitVO2MaxIntervalIDO `json:"cardioScore,omitempty"`
}

func UnmarshalFitBitV02MaxRawData(jsonData string) (*FitBitVO2MaxRawData, error) {
	d := &FitBitVO2MaxRawData{}

	err := json.Unmarshal([]byte(jsonData), &d) // Read the JSON string and convert it into our golang stuct else we get an error.
	if err != nil {
		return nil, err
	}

	cardioScore := []FitBitVO2MaxIntervalIDO{} // Array with modified values.

	for _, cs := range d.CardioScore {
		splitArr := strings.Split(cs.Value.VO2Max, "-")
		if len(splitArr) == 0 {
			v, _ := strconv.ParseFloat(cs.Value.VO2Max, 64)
			cs.AvgValue = v
		} else {
			v1, _ := strconv.ParseFloat(splitArr[0], 64)
			v2, _ := strconv.ParseFloat(splitArr[1], 64)
			cs.AvgValue = (v1 + v2) / 2
		}

		cardioScore = append(cardioScore, cs)
	}

	d.CardioScore = cardioScore // Replace with our updated value.

	return d, nil
}

// //------------- Electrocardiogram ---------------//
//
// type ElectrocardiogramIntervalIDO struct {
// 	StartTime string `json:"startTime,omitempty"`
// }
//
// type FitBitElectrocardiogramRawData struct {
// 	ECGReadings []ElectrocardiogramIntervalIDO `json:"ecgReadings,omitempty"`
// }
//
// func UnmarshalFitBitElectrocardiogramRawData(jsonData string) (*FitBitTemperatureRawData, error) {
// 	d := &FitBitTemperatureRawData{}
//
// 	err := json.Unmarshal([]byte(jsonData), &d) // Read the JSON string and convert it into our golang stuct else we get an error.
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return d, nil
// }
