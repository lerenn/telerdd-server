package api

import (
	"strings"
	"sync"
	"time"

	libConfig "github.com/lerenn/go-config"
	appConfig "github.com/nightwall/nightwall-server/config"
)

type data struct {
	// Request
	authorizedOrigin string
	// Message
	msgLimit     int
	MsgLimitSize int
	msgIP        map[string]*time.Time
	msgIPLock    *sync.Mutex
	// Moderation
	MsgModerationWithImg    bool
	MsgModerationWithoutImg bool
	// Image
	imgMaxWidth, imgMaxHeight int
}

func newData(c *libConfig.Config) (*data, error) {
	var d data
	var err error
	var str string

	// Get msg limit
	if d.msgLimit, err = c.GetInt(appConfig.MessagesSectionToken, appConfig.MessagesLimitToken); err != nil {
		return nil, err
	}

	// Get msg limit size
	if d.MsgLimitSize, err = c.GetInt(appConfig.MessagesSectionToken, appConfig.MessagesLimitSizeToken); err != nil {
		return nil, err
	}

	// Get authorized URL for client
	if d.authorizedOrigin, err = c.GetString(appConfig.ClientSectionToken, appConfig.ClientAuthorizedOriginToken); err != nil {
		return nil, err
	}

	// Get max size for picture
	if d.imgMaxWidth, err = c.GetInt(appConfig.ImageSectionToken, appConfig.ImageMaxWidthToken); err != nil {
		return nil, err
	}
	if d.imgMaxHeight, err = c.GetInt(appConfig.ImageSectionToken, appConfig.ImageMaxHeightToken); err != nil {
		return nil, err
	}

	// Get message moderation
	if str, err = c.GetString(appConfig.MessagesSectionToken, appConfig.MessagesModerationWithImg); err != nil {
		return nil, err
	}
	d.MsgModerationWithImg = strings.EqualFold(str, "true")
	if str, err = c.GetString(appConfig.MessagesSectionToken, appConfig.MessagesModerationWithoutImg); err != nil {
		return nil, err
	}
	d.MsgModerationWithoutImg = strings.EqualFold(str, "true")

	d.msgIP = make(map[string]*time.Time)
	d.msgIPLock = &sync.Mutex{}

	return &d, nil
}

func (d *data) ProceedMessageLimit(ip string) (int, error) {
	// Lock msgIP
	d.msgIPLock.Lock()
	defer d.msgIPLock.Unlock()

	// Get last time from ip
	last := d.msgIP[ip]

	// Check if it exists
	if last == nil {
		d.addMsgIP(ip)
		return -1, nil
	}

	// Compare time and return error if too soon
	limit := time.Duration(d.msgLimit)
	if time.Now().Before(last.Add(limit * time.Second)) {
		return int(time.Since(*last).Seconds()), nil
	}

	d.addMsgIP(ip)
	return -1, nil
}

// Accessors
////////////////////////////////////////////////////////////////////////////////

func (d *data) AuthorizedOrigin() string {
	return d.authorizedOrigin
}

// Private methods
////////////////////////////////////////////////////////////////////////////////

// UNSAFE : You have to lock use
func (d *data) addMsgIP(ip string) {
	t := time.Now()
	d.msgIP[ip] = &t
}
