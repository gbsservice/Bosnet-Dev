package constant

const (
	LayoutISO  = "2006-01-02"
	LayoutUS   = "January 2, 2006"
	LayoutTime = "15:04"
	LayoutDay  = "02 Jan"
)

const (
	JwtClaim    = "JwtClaim"
	Postgres    = "DBPostgres"
	MysqlHPM    = "DBMysqlHPM"
	RequestTime = "RequestTime"
	Redis       = "Redis"
	Auth        = "Auth"
	Activity    = "Activity"
)

const (
	BaseUrl        = ""
	BaseUrlFcm     = "https://fcm.googleapis.com/fcm/"
	BaseUrlGeoCode = "https://maps.googleapis.com/maps/api/geocode/json"
	ApiVersion     = "0.1"
)

const (
	ErrorLogin            = "invalid credential"
	ErrorTimesUp          = "TimesUp"
	ErrorScheduleNotReady = "ScheduleNotReady"
	ErrorNotFound         = "data tidak ditemukan"
	ErrorChangePassword   = "current password does not match"
	Error2FA              = "incorrect security code"
	ErrorToken            = "invalid token"
	Error2FARequest       = "error requesting 2fa"
	TwoFAIssuer           = "Kontinum"
	ResponseOK            = "OK"
)

const (
	GeoAddress   = "administrative_area_level_5"
	GeoKelurahan = "administrative_area_level_4"
	GeoKecamatan = "administrative_area_level_3"
	GeoKota      = "administrative_area_level_2"
	GeoProvinsi  = "administrative_area_level_1"
)

const (
	UserRoleAdministrator = "c5b4792c-0d62-11ec-8c16-0e138fa03d7b"
	UserRolePartner       = "a0d9d358-0ff2-11ec-8caf-00155d585831"
	UserRoleParticipant   = "b72ab000-0ff2-11ec-8caf-00155d585831"
)
