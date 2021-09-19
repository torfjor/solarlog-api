package solarlog

type DayValuesResponse struct {
	Header struct {
		Epoch        string                       `json:"epoch"`
		Timezone     string                       `json:"timezone"`
		UTCOffset    string                       `json:"utf_offset"`
		EpochDevices map[string]map[string]Device `json:"epoch_devices"`
	} `json:"header"`
	Readings map[string]interface{} `json:"body"`
}
