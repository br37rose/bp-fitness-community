package constants

const (
	OuraRingAuthorizationURL                       = "https://cloud.ouraring.com/oauth/authorize"
	OuraRingExchangeURL                            = "https://cloud.ouraring.com/oauth/token"
	OuraRingGetDailyActivityURL                    = "https://api.ouraring.com/v2/usercollection/daily_activity"
	OuraRingGetDailyReadinessURL                   = "https://api.ouraring.com/v2/usercollection/daily_readiness"
	OuraRingGetDailySleepURL                       = "https://api.ouraring.com/v2/usercollection/daily_sleep"
	OuraRingGetHeartRateURL                        = "https://api.ouraring.com/v2/usercollection/heartrate"
	OuraRingGetSleepPeriodsURL                     = "https://api.ouraring.com/v2/usercollection/sleep"
	OuraRingGetSessionsURL                         = "https://api.ouraring.com/v2/usercollection/session"
	OuraRingGetWorkoutsURL                         = "https://api.ouraring.com/v2/usercollection/workout"
	FitBitAuthorizationURL                         = "https://www.fitbit.com/oauth2/authorize"
	FitBitExchangeURL                              = "https://api.fitbit.com/oauth2/token"
	FitBitGetActivityRateIntradayByDateURL         = "https://api.fitbit.com/1/user/%s/activities/%s/date/%s/1d/%s.json?timezone=UTC"
	FitBitGetHeartRateIntradayByDateURL            = "https://api.fitbit.com/1/user/%s/activities/heart/date/%s/1d/%s.json?timezone=UTC"
	FitBitGetBreathingRateIntradayByDateURL        = "https://api.fitbit.com/1/user/%s/br/date/%s/all.json?timezone=UTC"
	FitBitGetHeartRateVariabilityIntradayByDateURL = "https://api.fitbit.com/1/user/%s/hrv/date/%s/all.json?timezone=UTC"
	FitBitGetSP02IntradayByDateURL                 = "https://api.fitbit.com/1/user/%s/spo2/date/%s/all.json?timezone=UTC"
	FitBitGetSleepLogByDateRangeURL                = "https://api.fitbit.com/1.2/user/%s/sleep/date/%s/%s.json"
	FitBitGetTemperatureSummaryByIntervalURL       = "https://api.fitbit.com/1/user/%s/temp/core/date/%s/%s.json"
	FitBitGetVO2MaxSummaryByIntervalURL            = "https://api.fitbit.com/1/user/%s/cardioscore/date/%s/%s.json"
	FitBitGetECGLogListURL                         = "https://api.fitbit.com/1/user/%s/ecg/list.json?afterDate=%s&sort=asc&limit=1&offset=0"
)
