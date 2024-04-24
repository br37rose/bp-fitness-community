package datastore

var typeMap = map[int8]string{
	ExerciseTypeSystem: "System",
	ExerciseTypeCustom: "Custom",
}

var categoryMap = map[int8]string{
	1:  "Anterior",
	2:  "Anti-Extension",
	3:  "Anti-Lateral Flexion",
	4:  "Anti-Rotation",
	5:  "Band",
	6:  "Bar",
	7:  "Barbell",
	8:  "Biceps",
	9:  "Bodyweight",
	10: "Cable",
	11: "Cable or band",
	12: "Carries",
	13: "Combination",
	14: "Dumbbell",
	15: "Dumbell",
	16: "Dynamic",
	17: "Flexion",
	18: "Ground-Based Exercises",
	19: "Hip Belt",
	20: "Lunge",
	21: "Plate",
	22: "Posterior",
	23: "Pushup",
	24: "Ring",
	25: "Static",
	26: "Triceps",
}

var movementTypeMap = map[int8]string{
	1:  "Arms",
	2:  "Core Work",
	3:  "Corrective Work",
	4:  "Hip-Hinge",
	5:  "Horizontal Pressing",
	6:  "Horizontal Pulling",
	7:  "Jumps",
	8:  "Single-Leg",
	9:  "Squats",
	10: "Vertical Pressing",
	11: "Vertical Pulling",
	12: "Warmups & Mobility Fillers",
	13: "Work",
}

// videoFileTypeMap := map[int8]string{
//     models.ExerciseS3VideoType:      "S3",
//     models.ExerciseYouTubeVideoType: "YouTube",
//     models.ExerciseVimeoVideoType:   "Vimeo",
// }
//
// // thumbnailFileTypeMap := map[int8]string{
// // 	models.ExerciseS3ThumbnailType:          "S3",
// // 	models.ExerciseExternalURLThumbnailType: "External URL",
// // 	models.ExerciseLocalThumbnailType:       "Local",
// // }
//
// stateMap := map[int8]string{
//     models.ExerciseActiveState:   "Active",
//     models.ExerciseArchivedState: "Archived",
// }

//
