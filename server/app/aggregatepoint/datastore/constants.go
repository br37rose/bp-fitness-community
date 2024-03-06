package datastore

const (
	StatusQueued   = 1
	StatusActive   = 2
	StatusArchived = 3
	StatusError    = 4
	StatusIgnored  = 5
)

const (
	TypeActivityCalaries     = 1000
	TypeActivityDistance     = 1001
	TypeActivityElevation    = 1002
	TypeActivityFloors       = 1003
	TypeActivitySteps        = 1004
	TypeHeartRate            = 1005
	TypeBreathingRate        = 1006
	TypeHeartRateVariability = 1007
	TypeOxygenSaturation     = 1008
	TypeSleep                = 1009
	TypeTemperature          = 1010
	TypeCardioFitnessScore   = 1011 // (a.k.a. V02 Max)
	TypeElectrocardiogram    = 1012
)

var TypeToFitBitResource = map[int]string{
	TypeActivityCalaries:  "calories",
	TypeActivityDistance:  "distance",
	TypeActivityElevation: "elevation",
	TypeActivityFloors:    "floors",
	TypeActivitySteps:     "steps",
	TypeHeartRate:         "distance",
	// TypeBreathingRate:        "xxx",
	// TypeHeartRateVariability: "xxx",
	// TypeOxygenSaturation:     "xxx",
	// TypeSleep:                "xxx",
	// TypeTemperature:          "xxx",
	// TypeCardioFitnessScore:   "xxx",
	// TypeElectrocardiogram:    "xxx",
}

var FitBitResourceToType = map[string]int64{
	"calories":  TypeActivityCalaries,
	"distance":  TypeActivityDistance,
	"elevation": TypeActivityElevation,
	"floors":    TypeActivityFloors,
	"steps":     TypeActivitySteps,
	"heartrate": TypeHeartRate,
	// TypeActivityFloors:       "floors",
	// TypeActivitySteps:        "steps",
	// TypeHeartRate:            "distance",
	// TypeBreathingRate:        "xxx",
	// TypeHeartRateVariability: "xxx",
	// TypeOxygenSaturation:     "xxx",
	// TypeSleep:                "xxx",
	// TypeTemperature:          "xxx",
	// TypeCardioFitnessScore:   "xxx",
	// TypeElectrocardiogram:    "xxx",
}

const (
	PeriodDay   = 1
	PeriodWeek  = 2
	PeriodMonth = 3
	PeriodYear  = 4
	PeriodHour  = 5
)
