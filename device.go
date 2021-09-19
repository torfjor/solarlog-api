package solarlog

// Device represents a peripheral connected to a Solarlog that reports metrics.
type Device struct {
	SerialNumber string             `json:"serialnumber"`
	Classes      []string           `json:"classes"`
	Type         string             `json:"type"`
	Model        string             `json:"model"`
	Name         string             `json:"name"`
	ID           string             `json:"id"`
	UID          string             `json:"uid"`
	Channels     map[string]Channel `json:"channels"`
}

// Channel returns a Channel with identifier c and a boolean signaling whether
// the channel exists on d.
func (d Device) Channel(c string) (found Channel, ok bool) {
	found, ok = d.Channels[c]
	return
}

// Channel represents a channel on a Device for reporting metrics.
type Channel struct {
	Channel     int    `json:"channel"`
	Position    int    `json:"position"`
	BaseName    string `json:"basename"`
	BaseChannel int    `json:"basechannel"`
	SubChannel  int    `json:"subchannel"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
}

// Value is a tuple representing a channel and a metric value.
type Value [2]string

// Channel returns v's channel element.
func (v Value) Channel() string {
	return v[0]
}

// Value returns v's value element.
func (v Value) Value() string {
	return v[1]
}

// DayValues represents a Device and a map of Values which are readings done by
// that device aggregated over a full day.
type DayValues struct {
	*Device
	// Values is a map of date strings to Value. The date layout used is
	// documented in DateLayout.
	Values map[string]Value
}
