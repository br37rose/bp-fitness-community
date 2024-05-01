package datastore

const (
	StatusQueued   = 1
	StatusActive   = 2
	StatusArchived = 3
	StatusError    = 4
	StatusIgnored  = 5
)

const (
	MetricDataTypeNameHeartRate            = 1
	MetricDataTypeNameActivitySteps        = 2
	MetricDataTypeNameActivityCalaries     = 3
	MetricDataTypeNameActivityDistance     = 4
	MetricDataTypeNameActivityElevation    = 5
	MetricDataTypeNameActivityFloors       = 6
	MetricDataTypeNameBreathingRate        = 7
	MetricDataTypeNameHeartRateVariability = 8
	MetricDataTypeNameOxygenSaturation     = 9
	MetricDataTypeNameSleep                = 10
	MetricDataTypeNameTemperature          = 11
	MetricDataTypeNameCardioFitnessScore   = 12 // (a.k.a. V02 Max)
	MetricDataTypeNameElectrocardiogram    = 13
)

var MetricDataTypeNameToFitBitResource = map[int]string{
	MetricDataTypeNameActivityCalaries:  "calories",
	MetricDataTypeNameActivityDistance:  "distance",
	MetricDataTypeNameActivityElevation: "elevation",
	MetricDataTypeNameActivityFloors:    "floors",
	MetricDataTypeNameActivitySteps:     "steps",
	MetricDataTypeNameHeartRate:         "distance",
	// MetricDataTypeNameBreathingRate:        "xxx",
	// MetricDataTypeNameHeartRateVariability: "xxx",
	// MetricDataTypeNameOxygenSaturation:     "xxx",
	// MetricDataTypeNameSleep:                "xxx",
	// MetricDataTypeNameTemperature:          "xxx",
	// MetricDataTypeNameCardioFitnessScore:   "xxx",
	// MetricDataTypeNameElectrocardiogram:    "xxx",
}

var FitBitResourceToMetricDataTypeName = map[string]int64{
	"calories":  MetricDataTypeNameActivityCalaries,
	"distance":  MetricDataTypeNameActivityDistance,
	"elevation": MetricDataTypeNameActivityElevation,
	"floors":    MetricDataTypeNameActivityFloors,
	"steps":     MetricDataTypeNameActivitySteps,
	"heartrate": MetricDataTypeNameHeartRate,
	// MetricDataTypeNameActivityFloors:       "floors",
	// MetricDataTypeNameActivitySteps:        "steps",
	// MetricDataTypeNameHeartRate:            "distance",
	// MetricDataTypeNameBreathingRate:        "xxx",
	// MetricDataTypeNameHeartRateVariability: "xxx",
	// MetricDataTypeNameOxygenSaturation:     "xxx",
	// MetricDataTypeNameSleep:                "xxx",
	// MetricDataTypeNameTemperature:          "xxx",
	// MetricDataTypeNameCardioFitnessScore:   "xxx",
	// MetricDataTypeNameElectrocardiogram:    "xxx",
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
