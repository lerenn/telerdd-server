package config

const (
	ConfigFile = "config/telerdd.conf"

	LogSectionToken = "log"
	LogFileToken    = "file"

	DbSectionToken  = "database"
	DbUserToken     = "user"
	DbPasswordToken = "password"
	DbAddrToken     = "address"
	DbPortToken     = "port"
	DbNameToken     = "name"

	MessagesSectionToken   = "messages"
	MessagesLimitToken     = "limit"
	MessagesLimitSizeToken = "limit_size"

	ClientSectionToken          = "client"
	ClientAuthorizedOriginToken = "authorized_origin"

	ImageSectionToken   = "image"
	ImageMaxWidthToken  = "max_width"
	ImageMaxHeightToken = "max_height"
)
