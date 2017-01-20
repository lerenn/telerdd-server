package config

const (
	ConfigFile = "config/nightwall.conf"

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

	MessagesModerationWithImg    = "moderate_msg_with_img"
	MessagesModerationWithoutImg = "moderate_msg_without_img"

	ClientSectionToken          = "client"
	ClientAuthorizedOriginToken = "authorized_origin"

	ImageSectionToken   = "image"
	ImageMaxWidthToken  = "max_width"
	ImageMaxHeightToken = "max_height"
)
