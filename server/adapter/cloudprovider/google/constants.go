package google

const (
	DataTypeShortNameActivity                  = "activity.segment"            // https://developers.google.com/fit/datatypes/activity
	DataTypeShortNameBasalMetabolicRate        = "calories.bmr"                // https://developers.google.com/fit/datatypes/activity#basal_metabolic_rate_bmr
	DataTypeShortNameCaloriesBurned            = "calories.expended"           // https://developers.google.com/fit/datatypes/activity#rest_8
	DataTypeShortNameCyclingPedalingCadence    = "cycling.pedaling.cadence"    // https://developers.google.com/fit/datatypes/activity#cycling_pedaling_cadence
	DataTypeShortNameCyclingPedalingCumulative = "cycling.pedaling.cumulative" // https://developers.google.com/fit/datatypes/activity#cycling_pedaling_cumulative
	DataTypeShortNameHeartPoints               = "heart_minutes"               // https://developers.google.com/fit/datatypes/activity#heart_points
	DataTypeShortNameeMoveMinutes              = "active_minutes"              // https://developers.google.com/fit/datatypes/activity#move_minutes
	DataTypeShortNamePower                     = "power.sample"                // https://developers.google.com/fit/datatypes/activity#power
	DataTypeShortNameStepCountCadence          = "step_count.cadence"          //https://developers.google.com/fit/datatypes/activity#step_count_cadence
	DataTypeShortNameStepCountDelta            = "step_count.delta"            // https://developers.google.com/fit/datatypes/activity#step_count_delta
	DataTypeShortNameWorkout                   = "activity.exercise"           //https://developers.google.com/fit/datatypes/activity#workout

	DataTypeShortNameCyclingWheelRevolution           = "cycling.wheel_revolution.rpm"        // https://developers.google.com/fit/datatypes/location#cycling_wheel_revolution_rpm
	DataTypeShortNameCyclingWheelRevolutionCumulative = "cycling.wheel_revolution.cumulative" // https://developers.google.com/fit/datatypes/location#cycling_wheel_revolution_cumulative
	DataTypeShortNameDistanceDelta                    = "distance.delta"                      // https://developers.google.com/fit/datatypes/location#distance_delta
	DataTypeShortNameLocationSample                   = "location.sample"                     // https://developers.google.com/fit/datatypes/location#location_sample
	DataTypeShortNameSpeed                            = "speed"                               // https://developers.google.com/fit/datatypes/location#speed

	DataTypeShortNameHydration = "hydration" // https://developers.google.com/fit/datatypes/nutrition
	DataTypeShortNameNutrition = "nutrition" // https://developers.google.com/fit/datatypes/nutrition

	DataTypeShortNameBloodGlucose      = "blood_glucose"       // https://developers.google.com/fit/datatypes/health#blood_glucose
	DataTypeShortNameBloodPressure     = "blood_pressure"      // https://developers.google.com/fit/datatypes/health#blood_pressure
	DataTypeShortNameBodyFatPercentage = "body.fat.percentage" // https://developers.google.com/fit/datatypes/health#body_fat_percentage
	DataTypeShortNameBodyTemperature   = "body.temperature"    // https://developers.google.com/fit/datatypes/health#body_temperature
	DataTypeShortNameCervicalMucus     = "cervical_mucus"      // https://developers.google.com/fit/datatypes/health#cervical_mucus
	DataTypeShortNameCervicalPosition  = "cervical_position"   // https://developers.google.com/fit/datatypes/health#cervical_position
	DataTypeShortNameHeartRateBPM      = "heart_rate.bpm"      // https://developers.google.com/fit/datatypes/health#heart_rate
	DataTypeShortNameHeight            = "height"              // https://developers.google.com/fit/datatypes/health#height
	DataTypeShortNameMenstruation      = "menstruation"        // https://developers.google.com/fit/datatypes/health#menstruation
	DataTypeShortNameOvulationTest     = "ovulation_test"      // https://developers.google.com/fit/datatypes/health#ovulation_test
	DataTypeShortNameOxygenSaturation  = "oxygen_saturation"   // https://developers.google.com/fit/datatypes/health#oxygen_saturation
	DataTypeShortNameSleep             = "sleep.segment"       // https://developers.google.com/fit/datatypes/health#sleep
	DataTypeShortNameVaginalSpotting   = "vaginal_spotting"    // https://developers.google.com/fit/datatypes/health#vaginal_spotting
	DataTypeShortNameWeight            = "weight"              // https://developers.google.com/fit/datatypes/health#weight
)

const (
	DataTypeNameActivity                  = "com.google.activity.segment"            // https://developers.google.com/fit/datatypes/activity
	DataTypeNameBasalMetabolicRate        = "com.google.calories.bmr"                // https://developers.google.com/fit/datatypes/activity#basal_metabolic_rate_bmr
	DataTypeNameCaloriesBurned            = "com.google.calories.expended"           // https://developers.google.com/fit/datatypes/activity#rest_8
	DataTypeNameCyclingPedalingCadence    = "com.google.cycling.pedaling.cadence"    // https://developers.google.com/fit/datatypes/activity#cycling_pedaling_cadence
	DataTypeNameCyclingPedalingCumulative = "com.google.cycling.pedaling.cumulative" // https://developers.google.com/fit/datatypes/activity#cycling_pedaling_cumulative
	DataTypeNameHeartPoints               = "com.google.heart_minutes"               // https://developers.google.com/fit/datatypes/activity#heart_points
	DataTypeNameeMoveMinutes              = "com.google.active_minutes"              // https://developers.google.com/fit/datatypes/activity#move_minutes
	DataTypeNamePower                     = "com.google.power.sample"                // https://developers.google.com/fit/datatypes/activity#power
	DataTypeNameStepCountCadence          = "com.google.step_count.cadence"          // https://developers.google.com/fit/datatypes/activity#step_count_cadence
	DataTypeNameStepCountDelta            = "com.google.step_count.delta"            // https://developers.google.com/fit/datatypes/activity#step_count_delta
	DataTypeNameWorkout                   = "com.google.activity.exercise"           //https://developers.google.com/fit/datatypes/activity#workout

	DataTypeNameCyclingWheelRevolution           = "com.google.cycling.wheel_revolution.rpm"        // https://developers.google.com/fit/datatypes/location#cycling_wheel_revolution_rpm
	DataTypeNameCyclingWheelRevolutionCumulative = "com.google.cycling.wheel_revolution.cumulative" // https://developers.google.com/fit/datatypes/location#cycling_wheel_revolution_cumulative
	DataTypeNameDistanceDelta                    = "com.google.distance.delta"                      // https://developers.google.com/fit/datatypes/location#distance_delta
	DataTypeNameLocationSample                   = "com.google.location.sample"                     // https://developers.google.com/fit/datatypes/location#location_sample
	DataTypeNameSpeed                            = "com.google.speed"                               // https://developers.google.com/fit/datatypes/location#speed

	DataTypeNameHydration = "com.google.hydration" // https://developers.google.com/fit/datatypes/nutrition
	DataTypeNameNutrition = "com.google.nutrition" // https://developers.google.com/fit/datatypes/nutrition

	DataTypeNameBloodGlucose      = "com.google.blood_glucose"       // https://developers.google.com/fit/datatypes/health#blood_glucose
	DataTypeNameBloodPressure     = "com.google.blood_pressure"      // https://developers.google.com/fit/datatypes/health#blood_pressure
	DataTypeNameBodyFatPercentage = "com.google.body.fat.percentage" // https://developers.google.com/fit/datatypes/health#body_fat_percentage
	DataTypeNameBodyTemperature   = "com.google.body.temperature"    // https://developers.google.com/fit/datatypes/health#body_temperature
	DataTypeNameCervicalMucus     = "com.google.cervical_mucus"      // https://developers.google.com/fit/datatypes/health#cervical_mucus
	DataTypeNameCervicalPosition  = "com.google.cervical_position"   // https://developers.google.com/fit/datatypes/health#cervical_position
	DataTypeNameHeartRateBPM      = "com.google.heart_rate.bpm"      // https://developers.google.com/fit/datatypes/health#heart_rate
	DataTypeNameHeight            = "com.google.height"              // https://developers.google.com/fit/datatypes/health#height
	DataTypeNameMenstruation      = "com.google.menstruation"        // https://developers.google.com/fit/datatypes/health#menstruation
	DataTypeNameOvulationTest     = "com.google.ovulation_test"      // https://developers.google.com/fit/datatypes/health#ovulation_test
	DataTypeNameOxygenSaturation  = "com.google.oxygen_saturation"   // https://developers.google.com/fit/datatypes/health#oxygen_saturation
	DataTypeNameSleep             = "com.google.sleep.segment"       // https://developers.google.com/fit/datatypes/health#sleep
	DataTypeNameVaginalSpotting   = "com.google.vaginal_spotting"    // https://developers.google.com/fit/datatypes/health#vaginal_spotting
	DataTypeNameWeight            = "com.google.weight"              // https://developers.google.com/fit/datatypes/health#weight
)
