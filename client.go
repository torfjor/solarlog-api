package solarlog

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	baseURL = "https://stromberg.solarlog-web.eu/api"
)

type doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client provides methods to interact with the Solarlog API.
type Client struct {
	user     string
	password string
	solarlog int
	doer
}

// NewClient returns a configured Client.
func NewClient(user, password string, solarlogID int, opts ...ClientOption) Client {
	c := Client{user: user, password: password, solarlog: solarlogID}

	for _, o := range opts {
		o(&c)
	}

	return c
}

// ClientOption is passed to NewClient to configure a Client.
type ClientOption func(*Client)

// WithDoer sets the Client's doer to d.
func WithDoer(d doer) ClientOption {
	return func(c *Client) {
		c.doer = d
	}
}

const (
	// DateLayout describes the date layout used by the Solarlog API as used in
	// time.Format.
	DateLayout = "2006-01-02"
)

// CurrentDayValues returns DayValues for all devices connected to the Solarlog
// for the current day.
func (c Client) CurrentDayValues(ctx context.Context) (map[string]DayValues, error) {
	now := time.Now()
	return c.dayValues(ctx, now, now)
}

// DayValues returns DayValues for all devices connected to the Solarlog for the
// time span between from and to.
func (c Client) DayValues(ctx context.Context, from, to time.Time) (map[string]DayValues, error) {
	if from.After(to) || to.IsZero() || from.IsZero() {
		return nil, fmt.Errorf("from must be before to")
	}

	return c.dayValues(ctx, from, to)
}

func (c Client) dayValues(ctx context.Context, from, to time.Time) (map[string]DayValues, error) {
	q := url.Values{
		"username":  {c.user},
		"password":  {c.password},
		"format":    {"json"},
		"function":  {"getDayValues"},
		"solarlog":  {strconv.Itoa(c.solarlog)},
		"date_from": {from.Format(DateLayout)},
		"date_to":   {to.Format(DateLayout)},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s?%s", baseURL, q.Encode()), nil)
	if err != nil {
		return nil, err
	}

	var resp DayValuesResponse
	if err := do(c.doer, req, &resp); err != nil {
		return nil, err
	}

	return decodeDayValuesResponse(resp)
}

func decodeDayValuesResponse(resp DayValuesResponse) (map[string]DayValues, error) {
	dayValueMap := newDeviceValueMap(resp.Header.EpochDevices[resp.Header.Epoch])

	for date, dm := range resp.Readings {
		deviceMap, ok := dm.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid device map")
		}
		for device, vm := range deviceMap {
			valueMap, ok := vm.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid value map")
			}
			for channel, value := range valueMap {
				dayValueMap[device].Values[date] = Value{channel, value.(string)}
			}
		}
	}

	return dayValueMap, nil
}

func newDeviceValueMap(dMap map[string]Device) map[string]DayValues {
	m := map[string]DayValues{}
	for deviceID, device := range dMap {
		d := device
		m[deviceID] = DayValues{
			Device: &d,
			Values: map[string]Value{},
		}
	}

	return m
}

func do(d doer, req *http.Request, v interface{}) error {
	req.Header.Set("Accept", "application/json")

	if d == nil {
		d = http.DefaultClient
	}

	resp, err := d.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// The Solarlog API returns http.StatusOK even for requests with invalid
	// credentials. The only reliable way to figure out if a request failed is
	// to peek at the first three bytes of the response.
	buf := bufio.NewReader(resp.Body)
	p, err := buf.Peek(3)
	if err != nil {
		return err
	}

	if bytes.Equal(p, []byte("ERR")) {
		b, err := ioutil.ReadAll(buf)
		if err != nil {
			return err
		}
		return fmt.Errorf("errenous response: %q", b)
	}

	if code := resp.StatusCode; code != http.StatusOK {
		return fmt.Errorf("invalid status: %d", code)
	}

	if v == nil {
		return nil
	}

	if err := json.NewDecoder(buf).Decode(v); err != nil {
		return err
	}

	return nil
}
