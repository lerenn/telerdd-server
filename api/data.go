package api

import (
	"sync"
	"time"
)

type Data struct {
	// Request
	authorizedOrigin string
	// Message
	msgLimit  int
	msgIP     map[string]*time.Time
	msgIPLock *sync.Mutex
}

func newData() *Data {
	var d Data
	d.msgIP = make(map[string]*time.Time)
	d.msgIPLock = &sync.Mutex{}
	return &d
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

// Mutators
////////////////////////////////////////////////////////////////////////////////

func (d *Data) SetMsgLimit(lim int) {
	d.msgLimit = lim
}

func (d *Data) SetAuthorizedOrigin(orig string) {
	d.authorizedOrigin = orig
}

// Private methods
////////////////////////////////////////////////////////////////////////////////

// UNSAFE : You have to lock use
func (d *Data) addMsgIP(ip string) {
	t := time.Now()
	d.msgIP[ip] = &t
}
