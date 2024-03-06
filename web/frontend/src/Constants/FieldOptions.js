import {
    GENDER_OTHER, GENDER_MALE, GENDER_FEMALE,
    PHYSICAL_ACTIVITY_SEDENTARY, PHYSICAL_ACTIVITY_LIGHTLY_ACTIVE, PHYSICAL_ACTIVITY_MODERATELY_ACTIVE, PHYSICAL_ACTIVITY_VERY_ACTIVE,
    WORKOUT_INTENSITY_LOW, WORKOUT_INTENSITY_MEDIUM, WORKOUT_INTENSITY_HIGH,
    FITNESS_GOAL_WEIGHT_LOSS, FITNESS_GOAL_MUSCLE_MASS_OR_STRENGTH, FITNESS_GOAL_CARDIOVASCULAR, FITNESS_GOAL_INJURY_RECOVERY, FITNESS_GOAL_SPORT_SPECIFIC, FITNESS_GOAL_AGILITY, FITNESS_GOAL_POWER, FITNESS_GOAL_SPEED, FITNESS_GOAL_BALANCE_OR_MOVEMENT,
    WORKOUT_PREFERENCE_FULL_BODY_WORKOUT, WORKOUT_PREFERENCE_FULL_HIIT,
    FITNESS_GOAL_STATUS_QUEUED, FITNESS_GOAL_STATUS_ACTIVE, FITNESS_GOAL_STATUS_ARCHIVED, FITNESS_GOAL_STATUS_ERROR
} from "./App";

export const HOW_DID_YOU_HEAR_ABOUT_US_OPTIONS = [
    { value: 5, label: 'Friend' },
    { value: 6, label: 'Social media' },
    { value: 7, label: 'Blog post article' },
    { value: 1, label: 'Other (Please specify)' },
];

export const HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...HOW_DID_YOU_HEAR_ABOUT_US_OPTIONS
];

export const MEMBER_STATUS_OPTIONS = [
    { value: 1, label: 'Active' },
    { value: 2, label: 'Archived' },
];

export const MEMBER_STATUS_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...MEMBER_STATUS_OPTIONS
];

export const TRAINER_STATUS_OPTIONS = [
    { value: 1, label: 'Active' },
    { value: 100, label: 'Archived' },
];

export const TRAINER_STATUS_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...TRAINER_STATUS_OPTIONS
];

export const DURATION_IN_MINUTES_OPTIONS = [
    { value: 5, label: '5 Minutes' },
    { value: 10, label: '10 Minutes' },
    { value: 15, label: '15 Minutes' },
    { value: 20, label: '20 Minutes' },
    { value: 25, label: '25 Minutes' },
    { value: 30, label: '30 Minutes' },
    { value: 35, label: '35 Minutes' },
    { value: 40, label: '40 Minutes' },
    { value: 45, label: '45 Minutes' },
    { value: 50, label: '50 Minutes' },
    { value: 55, label: '55 Minutes' },
    { value: 60, label: '60 Minutes' },
    { value: 65, label: '65 Minutes' },
    { value: 70, label: '70 Minutes' },
    { value: 75, label: '75 Minutes' },
    { value: 85, label: '80 Minutes' },
    { value: 85, label: '85 Minutes' },
    { value: 90, label: '90 Minutes' },
    { value: 95, label: '95 Minutes' },
    { value: 100, label: '100 Minutes' },
    { value: 120, label: '120 Minutes' },
    { value: 240, label: '240 Minutes' },
];

export const DURATION_IN_MINUTES_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...DURATION_IN_MINUTES_OPTIONS
];

export const ATENDEE_LIMIT_OPTIONS = [
    { value: 1, label: '1 Person' },
    { value: 2, label: '2 People' },
    { value: 3, label: '3 People' },
    { value: 4, label: '4 People' },
    { value: 5, label: '5 People' },
    { value: 6, label: '6 People' },
    { value: 7, label: '7 People' },
    { value: 8, label: '8 People' },
    { value: 9, label: '9 People' },
    { value: 10, label: '10 People' },
    { value: 11, label: '11 People' },
    { value: 12, label: '12 People' },
    { value: 13, label: '13 People' },
    { value: 14, label: '14 People' },
    { value: 15, label: '15 People' },
    { value: 16, label: '16 People' },
    { value: 17, label: '17 People' },
    { value: 18, label: '18 People' },
    { value: 19, label: '19 People' },
    { value: 20, label: '20 People' },
];

export const ATTENDEE_LIMIT_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...ATENDEE_LIMIT_OPTIONS
];

export const SELF_CANCEL_MINIMUM_HOUR_OPTIONS = [
    { value: 1, label: '1 Hour before Session' },
    { value: 2, label: '2 Hours before Session' },
    { value: 3, label: '3 Hours before Session' },
    { value: 3, label: '4 Hours before Session' },
    { value: 3, label: '5 Hours before Session' },
];

export const SELF_CANCEL_MINIMUM_HOUR_OPTIONS_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...SELF_CANCEL_MINIMUM_HOUR_OPTIONS
];

export const PAGE_SIZE_OPTIONS = [
    { value: 2, label: '2 Rows' },
    { value: 5, label: '5 Rows' },
    { value: 10, label: '10 Rows' },
    { value: 25, label: '25 Rows' },
    { value: 50, label: '50 Rows' },
    { value: 100, label: '100 Rows' },
    { value: 250, label: '250 Rows' },
];

export const STATUS_OPTIONS = [
    { value: 1, label: "Pending", color: "warning" },
    { value: 2, label: "Active", color: "success" },
    { value: 3, label: "Error", color: "danger" },
    { value: 4, label: "Cancelled", color: "danger" },
    { value: 5, label: "Archived", color: "dark" },
  ];

export const PAYMENT_OPTIONS = [
    { value: "Credit Card", label: "Credit Card" },
    { value: "Debit Card", label: "Debit Card" },
    { value: "PayPal", label: "PayPal" },
    { value: "Bank Transfer", label: "Bank Transfer" },
    { value: "Cash", label: "Cash" },
    { value: "Mobile Payment", label: "Mobile Payment" },
  ];

export const CLASSES_AND_EVENTS = [
    {
      "id":1,
      "title":"7 Day Trial (Local Residents Only)",
      "subTitle": "Join us for 7-day sweat dripping, heart-pumping fun!",
      "expirationDate": "7 days from first use",
      "price": "$35.00",
      "category": "Functional Team Training"
    },
    {
      "id":2,
      "title":"Drop In",
      "subTitle": "Any day. Any time.",
      "price": "$30.00",
      "category": "Functional Team Training"
    },
    {
      "id":3,
      "title":"10 Classes",
      "expirationDate": "2023-12-05",
      "price": "$265.00",
      "category": "Other Team"
    },
    {
      "id":4,
      "title":"20 Classes",
      "subTitle": "Train with us when it works for you! Passes are valid for 6 months from date of purchase.",
      "expirationDate": "2023-08-15",
      "price": "$495.00",
      "category": "Other Team"
    }
  ];

export const SUBSCRIPTION_STATUS_OPTIONS = [
    { value: "all", label: 'All' },
    { value: "canceled", label: 'Canceled' },
    { value: "incomplete", label: 'Incomplete' },
    { value: "incomplete_expired", label: 'Incomplete Expired' },
    { value: "past_due", label: 'Past Due' },
    { value: "trialing", label: 'Trialing' },
    { value: "unpaid", label: 'Unpaid' },
];

export const SUBSCRIPTION_STATUS_WITH_EMPTY_OPTIONS = [
    { value: "", label: "Please select" }, // EMPTY OPTION
    ...SUBSCRIPTION_STATUS_OPTIONS
];

export const SUBSCRIPTION_INTERVAL_OPTIONS = [
    { value: "day", label: 'Day' },
    { value: "week", label: 'Week' },
    { value: "month", label: 'Month' },
    { value: "year", label: 'Year' },
];

export const SUBSCRIPTION_TINTERVAL_WITH_EMPTY_OPTIONS = [
    { value: "", label: "Please select" }, // EMPTY OPTION
    ...SUBSCRIPTION_INTERVAL_OPTIONS
];

export const EXERCISE_TYPE_MAP = {
    1: "System",
    2: "Custom",
};

export const EXERCISE_TYPE_OPTIONS = [
    { value: 1, label: 'System' },
    { value: 2, label: 'Custom' },
];

export const EXERCISE_TYPE_WITH_EMPTY_OPTIONS = [
    { value: "", label: "Please select" }, // EMPTY OPTION
    ...EXERCISE_TYPE_OPTIONS
];

export const EXERCISE_CATEGORY_MAP = {
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
};

export const EXERCISE_CATEGORY_OPTIONS = [
    { value: 1, label: "Anterior" },
    { value: 2, label: "Anti-Extension" },
    { value: 3, label: "Anti-Lateral Flexion" },
    { value: 4, label: "Anti-Rotation" },
    { value: 5, label: "Band" },
    { value: 6, label: "Bar" },
    { value: 7, label: "Barbell" },
    { value: 8, label: "Biceps" },
    { value: 9, label: "Bodyweight" },
    { value: 10, label: "Cable" },
    { value: 11, label: "Cable or band" },
    { value: 12, label: "Carries" },
    { value: 13, label: "Combination" },
    { value: 14, label: "Dumbbell" },
    { value: 15, label: "Dumbell" },
    { value: 16, label: "Dynamic" },
    { value: 17, label: "Flexion" },
    { value: 18, label: "Ground-Based Exercises" },
    { value: 19, label: "Hip Belt" },
    { value: 20, label: "Lunge" },
    { value: 21, label: "Plate" },
    { value: 22, label: "Posterior" },
    { value: 23, label: "Pushup" },
    { value: 24, label: "Ring" },
    { value: 25, label: "Static" },
    { value: 26, label: "Triceps" },
];

export const EXERCISE_CATEGORY_OPTIONS_WITH_EMPTY_OPTION = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...EXERCISE_CATEGORY_OPTIONS
];

export const EXERCISE_MOMENT_TYPE_MAP = {
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

export const EXERCISE_MOMENT_TYPE_OPTIONS = [
    { value: 1, label: "Arms" },
    { value: 2, label: "Core Work" },
    { value: 3, label: "Corrective Work" },
    { value: 4, label: "Hip-Hinge" },
    { value: 5, label: "Horizontal Pressing" },
    { value: 6, label: "Horizontal Pulling" },
    { value: 7, label: "Jumps" },
    { value: 8, label: "Single-Leg" },
    { value: 9, label: "Squats" },
    { value: 10, label: "Vertical Pressing" },
    { value: 11, label: "Vertical Pulling" },
    { value: 12, label: "Warmups & Mobility Fillers" },
    { value: 13, label: "Work" },
];

export const EXERCISE_MOMENT_TYPE_OPTIONS_WITH_EMPTY_OPTION = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...EXERCISE_MOMENT_TYPE_OPTIONS
];

export const EXERCISE_VIDEO_FILE_TYPE_MAP = {
    1: "S3",
    2: "YouTube",
    3: "Vimeo",
}

export const EXERCISE_VIDEO_FILE_TYPE_OPTIONS = [
    { value: 1, label: "S3" },
    { value: 2, label: "YouTube" },
    { value: 3, label: "Vimeo" },
];

export const EXERCISE_VIDEO_FILE_TYPE_OPTIONS_WITH_EMPTY_OPTION = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...EXERCISE_VIDEO_FILE_TYPE_OPTIONS
];

export const EXERCISE_THUMBNAIL_FILE_TYPE_MAP = {
	1: "S3",
	2: "External URL",
	3: "Local",
}

export const EXERCISE_STATUS_MAP = {
    0: "Archived",
    1: "Active",
    2: "Archived",
}

export const EXERCISE_STATUS_OPTIONS = [
    { value: 1, label: "Active" },
    { value: 2, label: "Archived" },
];

export const EXERCISE_STATUS_OPTIONS_WITH_EMPTY_OPTION = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...EXERCISE_STATUS_OPTIONS
];

export const EXERCISE_GENDER_OPTIONS = [
    { value: "Male", label: "Male" },
    { value: "Female", label: "Female" },
    { value: "Other", label: "Other" },
];

export const EXERCISE_GENDER_OPTIONS_WITH_EMPTY_OPTION = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...EXERCISE_GENDER_OPTIONS
];

export const VIDEO_COLLECTION_TYPE_OPTIONS = [
    { value: 1, label: "Many Videos" },
    { value: 2, label: "Single Video" },
];

export const VIDEO_COLLECTION_TYPE_OPTIONS_WITH_EMPTY_OPTION = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...VIDEO_COLLECTION_TYPE_OPTIONS
];

export const VIDEO_COLLECTION_TYPE_MAP = {
    1: "Many Videos",
    2: "Single Videos",
};

export const VIDEO_COLLECTION_STATUS_MAP = {
    1:  "Active",
    2:  "Archived",
};

export const VIDEO_COLLECTION_STATUS_OPTIONS = [
    { value: 1, label: "Active" },
    { value: 2, label: "Archived" },
];

export const VIDEO_COLLECTION_STATUS_OPTIONS_WITH_EMPTY_OPTION = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...VIDEO_COLLECTION_STATUS_OPTIONS
];

export const VIDEO_CONTENT_TYPE_MAP = {
    1: "System",
    2: "Custom",
};

export const VIDEO_CONTENT_TYPE_OPTIONS = [
    { value: 1, label: 'System' },
    { value: 2, label: 'Custom' },
];

export const VIDEO_CONTENT_TYPE_WITH_EMPTY_OPTIONS = [
    { value: "", label: "Please select" }, // EMPTY OPTION
    ...VIDEO_CONTENT_TYPE_OPTIONS
];

export const VIDEO_CONTENT_VIDEO_TYPE_OPTIONS = [
    { value: 1, label: 'S3' },
    { value: 2, label: 'YouTube' },
    { value: 3, label: 'Vimeo' },
];

export const VIDEO_CONTENT_VIDEO_TYPE_WITH_EMPTY_OPTIONS = [
    { value: "", label: "Please select" }, // EMPTY OPTION
    ...VIDEO_CONTENT_VIDEO_TYPE_OPTIONS
];

export const VIDEO_CONTENT_VIDEO_TYPE_MAP = {
    1: "S3",
    2: "YouTube",
    3: "Vimeo",
};

export const PAY_FREQUENCY = [
    { id: 1, value: 1, label: 'Single Pay', identifier: '' },
    { id: 2, value: 2, label: 'Subscription', identifier: 'month' },
    { id: 3, value: 3, label: 'Annual', identifier: 'year' }
]

const OFFER_PAY_FREQUENCY_OPTIONS = [
    { value: 1, label: 'One-Time' },
    { value: 2, label: 'Day' },
    { value: 3, label: 'Week' },
    { value: 4, label: 'Month' },
    { value: 5, label: 'Year' },
];

export const OFFER_PAY_FREQUENCY_WITH_EMPTY_OPTIONS = [
    { value: "", label: "Please select" }, // EMPTY OPTION
    ...OFFER_PAY_FREQUENCY_OPTIONS
];

export const BUSINESS_FUNCTION_OPTIONS = [
    { value: 1, label: 'Provide Access to Content based on Membership Rank' },
];

export const BUSINESS_FUNCTION_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...BUSINESS_FUNCTION_OPTIONS
];

const OFFER_MEMBERSHIP_RANK_OPTIONS = [  // Membership Rank No
    { value: 4, label: 'Regular' },
    { value: 5, label: 'Elite' },
];

export const OFFER_MEMBERSHIP_RANK_WITH_EMPTY_OPTIONS = [
    { value: "", label: "Please select" }, // EMPTY OPTION
    ...OFFER_MEMBERSHIP_RANK_OPTIONS
];

export const OFFER_STATUS_OPTIONS = [  // Membership Rank No
    { value: 1, label: 'Pending' },
    { value: 2, label: 'Active' },
    { value: 3, label: 'Archive' },
];

export const OFFER_STATUS_WITH_EMPTY_OPTIONS = [
    { value: "", label: "Please select" }, // EMPTY OPTION
    ...OFFER_STATUS_OPTIONS
];

export const TIMED_LOCK_DURATION_OPTIONS = [
    { value: '24h', label: '1 Day' },
    { value: '48h', label: '2 Days' },
    { value: '168h', label: '1 Week' },
    { value: '336h', label: '2 Weeks' },
    { value: '504h', label: '3 Weeks' },
    { value: '672h', label: '4 Weeks' },
    { value: '720h', label: '1 Month' },
    { value: '1440h', label: '2 Months' },
    { value: '2160h', label: '3 Months' },
    { value: '2880h', label: '4 Months' },
    { value: '3600h', label: '5 Months' },
    { value: '4320h', label: '6 Months' },
    { value: '8760h', label: '1 Year' },
];

export const TIMED_LOCK_DURATION_WITH_EMPTY_OPTIONS = [
    { value: "", label: "Please select" }, // EMPTY OPTION
    ...TIMED_LOCK_DURATION_OPTIONS
];

// --- FITNESS PLAN --- //

export const HOME_GYM_EQUIPMENT_MAP = {
    2: "Bench/Boxes/Floor Mat",
    3: "Free Weights",
    4: "Barbell",
    5: "Cables",
    6: "Rower",
    7: "Stationary Bike",
    8: "Treadmill",
    9: "Resistant Bands",
    10: "Skipping Ropex",
    11: "Pull Up Bar",
    12: "Kettle Bells",
};

export const HOME_GYM_EQUIPMENT_OPTIONS = [
    { value: 2, label: 'Bench/Boxes/Floor Mat' },
    { value: 3, label: 'Free Weights' },
    { value: 4, label: 'Barbell' },
    { value: 5, label: 'Cables' },
    { value: 6, label: 'Rower' },
    { value: 7, label: 'Stationary Bike' },
    { value: 8, label: 'Treadmill' },
    { value: 9, label: 'Resistant Bands' },
    { value: 10, label: 'Skipping Ropex' },
    { value: 11, label: 'Pull Up Bar' },
    { value: 12, label: 'Kettle Bells' },
];

export const HOME_GYM_EQUIPMENT_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...HOME_GYM_EQUIPMENT_OPTIONS
];

export const FEET_OPTIONS = [
    { value: 4, label: '4\'' },
    { value: 5, label: '5\'' },
    { value: 6, label: '6\'' },
    { value: 7, label: '7\'' },
];

export const FEET_WITH_EMPTY_OPTIONS = [
    { value: -1, label: "ft" }, // EMPTY OPTION
    ...FEET_OPTIONS
];

export const INCHES_OPTIONS = [
    { value: 0, label: '0\"' },
    { value: 1, label: '1\"' },
    { value: 2, label: '2\"' },
    { value: 3, label: '3\"' },
    { value: 4, label: '4\"' },
    { value: 5, label: '5\"' },
    { value: 6, label: '6\"' },
    { value: 7, label: '7\"' },
    { value: 8, label: '8\"' },
    { value: 9, label: '9\"' },
    { value: 10, label: '10\"' },
    { value: 11, label: '11\"' },
    { value: 12, label: '12\"' },
];

export const INCHES_WITH_EMPTY_OPTIONS = [
    { value: -1, label: "in" }, // EMPTY OPTION
    ...INCHES_OPTIONS
];

export const GENDER_OPTIONS = [
    { value: GENDER_MALE, label: "Male" },
    { value: GENDER_FEMALE, label: "Female" },
    { value: GENDER_OTHER, label: "Other" },
];

export const GENDER_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...GENDER_OPTIONS
];

export const PHYSICAL_ACTIVITY_MAP = {
    [PHYSICAL_ACTIVITY_SEDENTARY]: "Sedentary",
    [PHYSICAL_ACTIVITY_LIGHTLY_ACTIVE]: "Lightly Active",
    [PHYSICAL_ACTIVITY_MODERATELY_ACTIVE]: "Moderately Active",
    [PHYSICAL_ACTIVITY_VERY_ACTIVE]: "Very Active",
};

export const PHYSICAL_ACTIVITY_OPTIONS = [
    { value: PHYSICAL_ACTIVITY_SEDENTARY, label: "Sedentary" },
    { value: PHYSICAL_ACTIVITY_LIGHTLY_ACTIVE, label: "Lightly Active" },
    { value: PHYSICAL_ACTIVITY_MODERATELY_ACTIVE, label: "Moderately Active" },
    { value: PHYSICAL_ACTIVITY_VERY_ACTIVE, label: "Very Active" },
];

export const PHYSICAL_ACTIVITY_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...PHYSICAL_ACTIVITY_OPTIONS
];

export const WORKOUT_INTENSITY_OPTIONS = [
    { value: WORKOUT_INTENSITY_LOW, label: "Low" },
    { value: WORKOUT_INTENSITY_MEDIUM, label: "Medium" },
    { value: WORKOUT_INTENSITY_HIGH, label: "High" },
];

export const WORKOUT_INTENSITY_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...WORKOUT_INTENSITY_OPTIONS
];

export const DAYS_PER_WEEK_MAP = {
    1: "1",
    2: "2",
    3: "3",
    4: "4",
    5: "5",
    6: "6",
    7: "7",
}

export const DAYS_PER_WEEK_OPTIONS = [
    { value: 1, label: "1" },
    { value: 2, label: "2" },
    { value: 3, label: "3" },
    { value: 4, label: "4" },
    { value: 5, label: "5" },
    { value: 6, label: "6" },
    { value: 7, label: "7" },
];

export const DAYS_PER_WEEK_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...DAYS_PER_WEEK_OPTIONS
];

export const TIME_PER_DAY_MAP = {
    30: "30 mins",
    60: "60 min",
    90: "90 min",
};

export const TIME_PER_DAY_OPTIONS = [
    { value: 30, label: "30 mins" },
    { value: 60, label: "60 min" },
    { value: 90, label: "90 min" },
];

export const TIME_PER_DAY_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...TIME_PER_DAY_OPTIONS
];

export const MAX_WEEK_MAP = {
    1: "1 week",
    2: "2 weeks",
    3: "3 weeks",
    4: "4 weeks",
    5: "5 weeks",
    6: "6 weeks",
};

export const MAX_WEEK_OPTIONS = [
    { value: 1, label: "1 week" },
    { value: 2, label: "2 weeks" },
    { value: 3, label: "3 weeks" },
    { value: 4, label: "4 weeks" },
    { value: 5, label: "5 weeks" },
    { value: 6, label: "6 weeks" },
];

export const MAX_WEEK_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...MAX_WEEK_OPTIONS
];

export const FITNESS_GOAL_MAP = {
    [FITNESS_GOAL_WEIGHT_LOSS]: "Weight Loss",
    [FITNESS_GOAL_MUSCLE_MASS_OR_STRENGTH]: "Muscle Mass/Strength",
    [FITNESS_GOAL_CARDIOVASCULAR]: "Cardiovascular",
    [FITNESS_GOAL_INJURY_RECOVERY]: "Injury Recovery",
    [FITNESS_GOAL_SPORT_SPECIFIC]: "Sport Specific",
    [FITNESS_GOAL_AGILITY]: "Goal Agility",
    [FITNESS_GOAL_POWER]: "Power",
    [FITNESS_GOAL_SPEED]: "Speed",
    [FITNESS_GOAL_BALANCE_OR_MOVEMENT]: "Balance/Movement",
};

export const FITNESS_GOAL_OPTIONS = [
    { value: FITNESS_GOAL_WEIGHT_LOSS, label: "Weight Loss" },
    { value: FITNESS_GOAL_MUSCLE_MASS_OR_STRENGTH, label: "Muscle Mass/Strength" },
    { value: FITNESS_GOAL_CARDIOVASCULAR, label: "Cardiovascular" },
    { value: FITNESS_GOAL_INJURY_RECOVERY, label: "Injury Recovery" },
    { value: FITNESS_GOAL_SPORT_SPECIFIC, label: "Sport Specific" },
    { value: FITNESS_GOAL_AGILITY, label: "Goal Agility" },
    { value: FITNESS_GOAL_POWER, label: "Power" },
    { value: FITNESS_GOAL_SPEED, label: "Speed" },
    { value: FITNESS_GOAL_BALANCE_OR_MOVEMENT, label: "Balance/Movement" },
];

export const FITNESS_GOAL_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...FITNESS_GOAL_OPTIONS
];

export const WORKOUT_PREFERENCE_MAP = {
    [WORKOUT_PREFERENCE_FULL_BODY_WORKOUT]: "Full-Body Workout",
    [WORKOUT_PREFERENCE_FULL_HIIT]: "HIIT",
};

export const WORKOUT_PREFERENCE_OPTIONS = [
    { value: WORKOUT_PREFERENCE_FULL_BODY_WORKOUT, label: "Full-Body Workout" },
    { value: WORKOUT_PREFERENCE_FULL_HIIT, label: "HIIT" },
];

export const WORKOUT_PREFERENCE_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...WORKOUT_PREFERENCE_OPTIONS
];

export const FITNESS_PLAN_STATUS_MAP = {
    [FITNESS_GOAL_STATUS_QUEUED]: "Not Ready",
    [FITNESS_GOAL_STATUS_ACTIVE]: "Ready",
    [FITNESS_GOAL_STATUS_ARCHIVED]: "Archived",
    [FITNESS_GOAL_STATUS_ERROR]: "Problem",
};

export const MEALS_PER_DAY_MAP = {
    1: "1 meal per day",
    2: "2 meals per day",
    3: "3 meals per day",
    4: "4 meals per day",
    5: "5 meals per day",
    6: "6 meals per day",
    7: "7 meals per day",
    8: "8 meals per day",
    9: "9 meals per day",
    10: "10 meals per day",
};

export const MEALS_PER_DAY_OPTIONS = [
    { value: 1, label: "1 meal per day" },
    { value: 2, label: "2 meals per day" },
    { value: 3, label: "3 meals per day" },
    { value: 4, label: "4 meals per day" },
    { value: 5, label: "5 meals per day" },
    { value: 6, label: "6 meals per day" },
    { value: 7, label: "7 meals per day" },
    { value: 8, label: "8 meals per day" },
    { value: 9, label: "9 meals per day" },
    { value: 10, label: "10 meals per day" },
];

export const MEALS_PER_DAY_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...MEALS_PER_DAY_OPTIONS
];

export const CONSUME_FREQUENCY_OPTIONS = [
    { value: 1, label: "never" },
    { value: 2, label: "once a year" },
    { value: 3, label: "once a month" },
    { value: 4, label: "once a week" },
    { value: 5, label: "3-4 times a week" },
    { value: 6, label: "daily" },
];

export const CONSUME_FREQUENCY_MAP = {
    1: "never",
    2: "once a year",
    3: "once a month",
    4: "once a week",
    5: "3-4 times a week",
    6: "daily",
};

export const CONSUME_FREQUENCY_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...CONSUME_FREQUENCY_OPTIONS
];

export const NUTRITIONAL_GOAL_MAP = {
    1: "Weight loss",
    2: "Weight gain",
    3: "Get lean",
    4: "Vegan diet",
    5: "Vegiterian diet",
    6: "Fruitarian diet",
    7: "Carnivore diet",
};

export const NUTRITIONAL_GOAL_OPTIONS = [
    { value: 1, label: "Weight loss" },
    { value: 2, label: "Weight gain" },
    { value: 3, label: "Get lean" },
    { value: 4, label: "Vegan diet" },
    { value: 5, label: "Vegiterian diet" },
    { value: 6, label: "Fruitarian diet" },
    { value: 7, label: "Carnivore diet" },
];

export const NUTRITIONAL_GOAL_WITH_EMPTY_OPTIONS = [
    { value: 0, label: "Please select" }, // EMPTY OPTION
    ...NUTRITIONAL_GOAL_OPTIONS
];
