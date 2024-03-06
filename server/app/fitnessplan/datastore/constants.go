package datastore

const (
	StatusQueued   = 1
	StatusActive   = 2
	StatusArchived = 3
	StatusError    = 4

	GenderOther  = 1
	GenderMale   = 2
	GenderFemale = 3

	// OwnershipTypeTemporary indicates file has been uploaded and saved in our system but not assigned ownership to anything. As a result, if this fitnessplan is not assigned within 24 hours then the crontab will delete this fitnessplan record and the uploaded file.
	OwnershipTypeTemporary         = 1
	OwnershipTypeExerciseVideo     = 2
	OwnershipTypeExerciseThumbnail = 3
	OwnershipTypeUser              = 4
	OwnershipTypeOrganization      = 5
	ContentTypeFile                = 6
	ContentTypeImage               = 7

	EquipmentAccessNoEquipmen    = 1
	EquipmentAccessFullGymAccess = 2
	EquipmentAccessHomeGym       = 3

	HasWorkoutsAtHomeYes = 1
	HasWorkoutsAtHomeNo  = 2

	PhysicalActivitySedentary        = 1
	PhysicalActivityLightlyActive    = 2
	PhysicalActivityModeratelyActive = 3
	PhysicalActivityVeryActive       = 4

	WorkoutIntensityLow    = 1
	WorkoutIntensityMedium = 2
	WorkoutIntensityHigh   = 2
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

var HasIntermittentFastingMap = map[int8]string{
	1: "Yes",
	2: "No",
}

var HomeGymEquipmentMap = map[int8]string{
	2:  "Bench/Boxes/Floor Mat",
	3:  "Free Weights",
	4:  "Barbell",
	5:  "Cables",
	6:  "Rower",
	7:  "Stationary Bike",
	8:  "Treadmill",
	9:  "Resistant Bands",
	10: "Skipping Ropex",
	11: "Pull Up Bar",
	12: "Kettle Bells",
}

var HasWorkoutsAtHomeMap = map[int8]string{
	1: "Yes",
	2: "No",
}

var DaysPerWeekMap = map[int8]string{
	1: "1",
	2: "2",
	3: "3",
	4: "4",
	5: "5",
	6: "6",
	7: "7",
}

var TimePerDayMap = map[int8]string{
	30: "30 mins",
	60: "60 min",
	90: "90 min",
}

var MaxWeekMap = map[int8]string{
	1: "1 week",
	2: "2 weeks",
	3: "3 weeks",
	4: "4 weeks",
	5: "5 weeks",
	6: "6 weeks",
}

var GoalMap = map[int8]string{
	1: "Weight Loss",
	2: "Muscle Mass/Strength",
	3: "Cardiovascular",
	4: "Injury Recovery",
	5: "Sport Specific",
	6: "Goal Agility",
	7: "Power",
	8: "Speed",
	9: "Balance/Movement",
}

var WorkoutPreferenceMap = map[int8]string{
	1: "Full-Body Workout",
	2: "HIIT",
}
