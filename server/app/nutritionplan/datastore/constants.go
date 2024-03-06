package datastore

const (
	StatusQueued   = 1
	StatusActive   = 2
	StatusArchived = 3
	StatusError    = 4

	GenderOther  = 1
	GenderMale   = 2
	GenderFemale = 3

	// OwnershipTypeTemporary indicates file has been uploaded and saved in our system but not assigned ownership to anything. As a result, if this nutritionplan is not assigned within 24 hours then the crontab will delete this nutritionplan record and the uploaded file.
	OwnershipTypeTemporary         = 1
	OwnershipTypeExerciseVideo     = 2
	OwnershipTypeExerciseThumbnail = 3
	OwnershipTypeUser              = 4
	OwnershipTypeOrganization      = 5
	ContentTypeFile                = 6
	ContentTypeImage               = 7

	PhysicalActivitySedentary        = 1
	PhysicalActivityLightlyActive    = 2
	PhysicalActivityModeratelyActive = 3
	PhysicalActivityVeryActive       = 4

	WorkoutIntensityLow    = 1
	WorkoutIntensityMedium = 2
	WorkoutIntensityHigh   = 2

	HasAllergiesYes = 1
	HasAllergiesNo  = 2
)

var PhysicalActivityMap = map[int8]string{
	1: "Sedentary",
	2: "Lightly Active",
	3: "Moderately Active",
	4: "Very Active",
}

var WorkoutIntensityMap = map[int8]string{
	1: "Low",
	2: "Medium",
	3: "High",
}

var GoalMap = map[int8]string{
	1: "Weight loss",
	2: "Weight gain",
	3: "Get lean",
	4: "Vegan diet",
	5: "Vegiterian diet",
	6: "Fruitarian diet",
	7: "Carnivore diet",
}

var MealsPerDayMap = map[int8]string{
	1:  "1 meal per day",
	2:  "2 meals per day",
	3:  "3 meals per day",
	4:  "4 meals per day",
	5:  "5 meals per day",
	6:  "6 meals per day",
	7:  "7 meals per day",
	8:  "8 meals per day",
	9:  "9 meals per day",
	10: "10 meals per day",
}

var ConsumeFrequencyMap = map[int8]string{
	1: "never",
	2: "once a year",
	3: "once a month",
	4: "once a week",
	5: "3-4 times a week",
	6: "daily",
}

var HasIntermittentFastingMap = map[int8]string{
	1: "Yes",
	2: "No",
}

var MaxWeekMap = map[int8]string{
	1: "1 week",
	2: "2 weeks",
	3: "3 weeks",
	4: "4 weeks",
	5: "5 weeks",
	6: "6 weeks",
}
