const HTTP_API_SERVER =
  process.env.REACT_APP_API_PROTOCOL + "://" + process.env.REACT_APP_API_DOMAIN;
export const BP8_FITNESS_API_BASE_PATH = "/api/v1";
export const BP8_FITNESS_VERSION_ENDPOINT = "version";
export const BP8_FITNESS_LOGIN_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/login";
export const BP8_FITNESS_REGISTER_MEMBER_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/register-member";
export const BP8_FITNESS_REFRESH_TOKEN_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/refresh-token";
export const BP8_FITNESS_EMAIL_VERIFICATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/verify";
export const BP8_FITNESS_LOGOUT_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/logout";
export const BP8_FITNESS_2FA_GENERATE_OTP_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/otp/generate";
export const BP8_FITNESS_2FA_GENERATE_OTP_AND_QR_CODE_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/otp/generate-qr-code";
export const BP8_FITNESS_2FA_VERIFY_OTP_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/otp/verify";
export const BP8_FITNESS_2FA_VALIDATE_OTP_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/otp/validate";
export const BP8_FITNESS_2FA_DISABLED_OTP_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/otp/disable";
export const BP8_FITNESS_ACCOUNT_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/account";
export const BP8_FITNESS_ACCOUNT_CHANGE_PASSWORD_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/account/change-password";
export const BP8_FITNESS_ACCOUNT_AVATAR_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/account/operation/avatar";
export const BP8_FITNESS_ORGANIZATIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/organizations";
export const BP8_FITNESS_ORGANIZATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/organization/{id}";
export const BP8_FITNESS_FORGOT_PASSWORD_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/forgot-password";
export const BP8_FITNESS_PASSWORD_RESET_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/password-reset";
export const BP8_FITNESS_REGISTRY_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/cpsrn/{id}";
export const BP8_FITNESS_MEMBERS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/members";
export const BP8_FITNESS_MEMBER_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/member/{id}";
export const BP8_FITNESS_MEMBER_CREATE_COMMENT_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/members/operation/create-comment";
export const BP8_FITNESS_EXERCISES_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/exercises";
export const BP8_FITNESS_EXERCISE_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/exercise/{id}";
export const BP8_FITNESS_ATTACHMENTS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/attachments";
export const BP8_FITNESS_ATTACHMENT_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/attachment/{id}";
export const BP8_FITNESS_VIDEO_CATEGORIES_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/video-categories";
export const BP8_FITNESS_VIDEO_CATEGORY_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/video-category/{id}";
export const BP8_FITNESS_VIDEO_CATEGORY_SELECT_OPTIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/video-categories/select-options";
export const BP8_FITNESS_VIDEO_COLLECTIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/video-collections";
export const BP8_FITNESS_VIDEO_COLLECTION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/video-collection/{id}";
export const BP8_FITNESS_VIDEO_CONTENTS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/video-contents";
export const BP8_FITNESS_VIDEO_CONTENT_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/video-content/{id}";
export const BP8_FITNESS_OFFERS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/offers";
export const BP8_FITNESS_OFFER_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/offer/{id}";
export const BP8_FITNESS_OFFER_SELECT_OPTIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/offers/select-options";
export const BP8_FITNESS_FITNESS_PLANS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/fitness-plans";
export const BP8_FITNESS_FITNESS_PLAN_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/fitness-plan/{id}";
export const BP8_FITNESS_FITNESS_PLAN_SELECT_OPTIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/fitness-plans/select-options";
export const BP8_FITNESS_NUTRITION_PLANS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/nutrition-plans";
export const BP8_FITNESS_NUTRITION_PLAN_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/nutrition-plan/{id}";
export const BP8_FITNESS_NUTRITION_PLAN_SELECT_OPTIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/nutrition-plans/select-options";
export const BP8_FITNESS_WEARABLE_FITBIT_DEVICE_REGISTRATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/fitbit-app-registration";
export const BP8_FITNESS_WEARABLE_GOOGLE_FIT_REGISTRATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/google-login";
export const BP8_FITNESS_WEARABLE_FITBITAPP_CREATE_SIMULATOR_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/fitbit/simulators";
export const BP8_FITNESS_DATA_POINTS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/data-points";
export const BP8_FITNESS_GOOGLE_FIT_DATA_POINTS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/google-fit-data-points";
export const BP8_FITNESS_AGGREGATE_POINTS_SUMMARY_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/aggregate-points/summary";
export const BP8_FITNESS_RANK_POINTS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/rank-points";
export const BP8_FITNESS_LEADERBOARD_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/leaderboard";
export const BP8_FITNESS_TAGS_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/tags";
export const BP8_FITNESS_TAG_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/tag/{id}";
export const BP8_FITNESS_TAG_SELECT_OPTIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/tags/select-options";
export const BP8_FITNESS_BIOMETRICS_MY_SUMMARY_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/biometrics/summary";
export const BP8_FITNESS_BIOMETRICS_HISTORIC_DATA_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/biometrics/historic-data";
export const BP8_FITNESS_PUBLIC_WORKOUT_SESSIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/public/workout-sessions";
export const BP8_FITNESS_COMPLETE_STRIPE_CHECKOUT_SESSION_API_ENDPOINT =
  HTTP_API_SERVER +
  "/api/v1/stripe/complete-checkout-session?session_id={sessionID}";
export const BP8_FITNESS_CREATE_STRIPE_CHECKOUT_SESSION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/stripe/create-checkout-session";
export const BP8_FITNESS_CANCEL_SUBSCRIPTION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/stripe/cancel-subscription";
export const BP8_FITNESS_PAYMENT_PROCESSOR_STRIPE_INVOICES_API_ENDPOINT =
  HTTP_API_SERVER +
  "/api/v1/stripe/invoices?user_id={userID}&cursor={cursor}&page_size={pageSize}";
export const BP8_FITNESS_WORKOUTS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/workouts";
export const BP8_FITNESS_T_PROGRAM_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/training-program";
export const BP8_FITNESS_MEMBER_SELECT_OPTION_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/select-options/{bid}/members";
export const BP8_FITNESS_T_PROGRAM_PHASE_PATCH_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/training-program/{pid}/phases";
export const BP8_FITNESS_QUESTIONNAIRE_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/questions";
export const BP8_FITNESS_FITNESS_CHALLENGE_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/fitness-challenge";

// Auto-generated comment for change 4
