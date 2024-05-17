export const PAGINATION_LIMIT = 250;

export const ROOT_ROLE_ID = 1;
export const ADMIN_ROLE_ID = 2;
export const TRAINER_ROLE_ID = 3;
export const MEMBER_ROLE_ID = 4;
export const ANONYMOUS_ROLE_ID = 0;

export const EXERCISE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER = 1;
export const EXERCISE_VIDEO_TYPE_YOUTUBE = 2;
export const EXERCISE_VIDEO_TYPE_VIMEO = 3;

export const EXERCISE_THUMBNAIL_TYPE_SIMPLE_STORAGE_SERVER = 1;
export const EXERCISE_THUMBNAIL_TYPE_EXTERNAL_URL = 2;

export const EXERCISE_TYPE_SYSTEM = 1;
export const EXERCISE_TYPE_CUSTOM = 2;

export const VIDEO_COLLECTION_TYPE_MANY_VIDEOS = 1;
export const VIDEO_COLLECTION_TYPE_SINGLE_VIDEO = 2;

export const VIDEO_CONTENTE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER = 1;
export const VIDEO_CONTENTE_VIDEO_TYPE_YOUTUBE = 2;
export const VIDEO_CONTENTE_VIDEO_TYPE_VIMEO = 3;

export const VIDEO_CONTENTE_THUMBNAIL_TYPE_SIMPLE_STORAGE_SERVER = 1;
export const VIDEO_CONTENTE_THUMBNAIL_TYPE_EXTERNAL_URL = 2;

export const VIDEO_CONTENTE_TYPE_SYSTEM = 1;
export const VIDEO_CONTENTE_TYPE_CUSTOM = 2;

export const GENDER_OTHER = 1;
export const GENDER_MALE = 2;
export const GENDER_FEMALE = 3;

export const PHYSICAL_ACTIVITY_SEDENTARY = 1;
export const PHYSICAL_ACTIVITY_LIGHTLY_ACTIVE = 2;
export const PHYSICAL_ACTIVITY_MODERATELY_ACTIVE = 3;
export const PHYSICAL_ACTIVITY_VERY_ACTIVE = 4;

export const WORKOUT_INTENSITY_LOW = 1;
export const WORKOUT_INTENSITY_MEDIUM = 2;
export const WORKOUT_INTENSITY_HIGH = 3;

export const FITNESS_GOAL_STATUS_QUEUED = 1;
export const FITNESS_GOAL_STATUS_ACTIVE = 2;
export const FITNESS_GOAL_STATUS_ARCHIVED = 3;
export const FITNESS_GOAL_STATUS_ERROR = 4;
export const FITNESS_GOAL_STATUS_IN_PROGRESS = 5;
export const FITNESS_GOAL_STATUS_PENDING = 6;

export const FITNESS_GOAL_WEIGHT_LOSS = 1;
export const FITNESS_GOAL_MUSCLE_MASS_OR_STRENGTH = 2;
export const FITNESS_GOAL_CARDIOVASCULAR = 3;
export const FITNESS_GOAL_INJURY_RECOVERY = 4;
export const FITNESS_GOAL_SPORT_SPECIFIC = 5;
export const FITNESS_GOAL_AGILITY = 6;
export const FITNESS_GOAL_POWER = 7;
export const FITNESS_GOAL_SPEED = 8;
export const FITNESS_GOAL_BALANCE_OR_MOVEMENT = 9;

export const WORKOUT_PREFERENCE_FULL_BODY_WORKOUT = 1;
export const WORKOUT_PREFERENCE_FULL_HIIT = 2;

export const RANK_POINT_PERIOD_DAY = 1;
export const RANK_POINT_PERIOD_WEEK = 2;
export const RANK_POINT_PERIOD_MONTH = 3;
export const RANK_POINT_PERIOD_YEAR = 4;

export const RANK_POINT_FUNCTION_AVERAGE = 1;
export const RANK_POINT_FUNCTION_SUM = 2;
export const RANK_POINT_FUNCTION_COUNT = 3;
export const RANK_POINT_FUNCTION_MIN = 4;
export const RANK_POINT_FUNCTION_MAX = 5;

// Please see `DataTypeKey` prefixed constants via https://github.com/bci-innovation-labs/bp8fitnesscommunity/blob/main/server/adapter/cloudprovider/google/constants.go
export const RANK_POINT_METRIC_TYPE_HEART_RATE = "com.google.heart_rate.bpm";
export const RANK_POINT_METRIC_TYPE_STEP_COUNTER = "com.google.step_count.delta";
export const RANK_POINT_METRIC_TYPE_CALORIES_BURNED = "com.google.calories.expended";
export const RANK_POINT_METRIC_TYPE_DISTANCE_DELTA = "com.google.distance.delta";
//TODO: Add more health sensors here...

export const WORKOUT_ALL_VISIBLE = 1;
export const WORKOUT_PERSONAL_VISIBLE = 2;
