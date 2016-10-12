package api

import (
	"sync"
	"time"

	config "github.com/lerenn/go-config"
	cst "github.com/lerenn/telerdd-server/constants"
)

type Data struct {
	// Request
	authorizedOrigin string
	// Message
	msgLimit  int
	msgIP     map[string]*time.Time
	msgIPLock *sync.Mutex
	// Image
	imgMaxWidth int
	imgMaxHeight int
}

func newData(c *config.Config) (*Data, error) {
	var d Data
	var err error

	// Get msg limit
	if d.msgLimit, err = c.GetInt(cst.MESSAGES_SECTION_TOKEN, cst.MESSAGES_LIMIT_TOKEN); err != nil {
		return nil, err
	}

	// Get authorized URL for client
	if d.authorizedOrigin, err = c.GetString(cst.CLIENT_SECTION_TOKEN, cst.CLIENT_AUTHORIZED_ORIGIN_TOKEN); err != nil {
		return nil, err
	}

	// Get max size for picture
	if d.imgMaxWidth, err = c.GetInt(cst.IMAGE_SECTION_TOKEN, cst.IMAGE_MAX_WIDTH_TOKEN); err != nil {
		return nil, err
	}
	if d.imgMaxHeight, err = c.GetInt(cst.IMAGE_SECTION_TOKEN, cst.IMAGE_MAX_HEIGHT_TOKEN); err != nil {
		return nil, err
	}

	d.msgIP = make(map[string]*time.Time)
	d.msgIPLock = &sync.Mutex{}

	return &d, nil
}

func (d *Data) ProceedMessageLimit(ip string) (int, error) {
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

func (d *Data) AuthorizedOrigin() string {
	return d.authorizedOrigin
}

// Private methods
////////////////////////////////////////////////////////////////////////////////

// UNSAFE : You have to lock use
func (d *Data) addMsgIP(ip string) {
	t := time.Now()
	d.msgIP[ip] = &t
}
