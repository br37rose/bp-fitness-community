package datastore

const (
	StatusQueued   = 1
	StatusActive   = 2
	StatusArchived = 3
	StatusError    = 4
	StatusIgnored  = 5
)

const (
	MetricTypeHeartRate            = 1
	MetricTypeActivitySteps        = 2
	MetricTypeActivityCalaries     = 3
	MetricTypeActivityDistance     = 4
	MetricTypeActivityElevation    = 5
	MetricTypeActivityFloors       = 6
	MetricTypeBreathingRate        = 7
	MetricTypeHeartRateVariability = 8
	MetricTypeOxygenSaturation     = 9
	MetricTypeSleep                = 10
	MetricTypeTemperature          = 11
	MetricTypeCardioFitnessScore   = 12 // (a.k.a. V02 Max)
	MetricTypeElectrocardiogram    = 13
)

var MetricTypeToFitBitResource = map[int]string{
	MetricTypeActivityCalaries:  "calories",
	MetricTypeActivityDistance:  "distance",
	MetricTypeActivityElevation: "elevation",
	MetricTypeActivityFloors:    "floors",
	MetricTypeActivitySteps:     "steps",
	MetricTypeHeartRate:         "distance",
	// MetricTypeBreathingRate:        "xxx",
	// MetricTypeHeartRateVariability: "xxx",
	// MetricTypeOxygenSaturation:     "xxx",
	// MetricTypeSleep:                "xxx",
	// MetricTypeTemperature:          "xxx",
	// MetricTypeCardioFitnessScore:   "xxx",
	// MetricTypeElectrocardiogram:    "xxx",
}

var FitBitResourceToMetricType = map[string]int64{
	"calories":  MetricTypeActivityCalaries,
	"distance":  MetricTypeActivityDistance,
	"elevation": MetricTypeActivityElevation,
	"floors":    MetricTypeActivityFloors,
	"steps":     MetricTypeActivitySteps,
	"heartrate": MetricTypeHeartRate,
	// MetricTypeActivityFloors:       "floors",
	// MetricTypeActivitySteps:        "steps",
	// MetricTypeHeartRate:            "distance",
	// MetricTypeBreathingRate:        "xxx",
	// MetricTypeHeartRateVariability: "xxx",
	// MetricTypeOxygenSaturation:     "xxx",
	// MetricTypeSleep:                "xxx",
	// MetricTypeTemperature:          "xxx",
	// MetricTypeCardioFitnessScore:   "xxx",
	// MetricTypeElectrocardiogram:    "xxx",
}

const (
	FunctionAverage = 1
	FunctionSum     = 2
	FunctionCount   = 3
	FunctionMin     = 4
	FunctionMax     = 5
	PeriodDay       = 1
	PeriodWeek      = 2
	PeriodMonth     = 3
	PeriodYear      = 4
)
