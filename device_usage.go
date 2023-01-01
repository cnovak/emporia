package emporia

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type DeviceListUsages struct {
	Devices    []DeviceUsage `json:"devices"`
	EnergyUnit EnergyUnit    `json:"energyUnit"`
	Instant    time.Time     `json:"instant"`
	Scale      Scale         `json:"scale"`
}

type DeviceUsage struct {
	ChannelUsages []ChannelUsage `json:"channelUsages"`
	DeviceGid     uint64         `json:"deviceGid"`
}

func (d *DeviceUsage) getChannel(channelUsageType ChannelUsageType) (*ChannelUsage, error) {
	channelName := channelUsageType.String()
	for _, channel := range d.ChannelUsages {
		if channel.Name == channelName {
			return &channel, nil
		}
	}
	return nil, fmt.Errorf("%d is not a valid ChannelUsageType", channelUsageType)
}

// Special channels in the device usage response
type ChannelUsageType uint8

const (
	NetUsage ChannelUsageType = iota
	TotalUsage
	Balance
)

func (c ChannelUsageType) String() string {
	switch c {
	case NetUsage:
		return "Main" // Amount coming/going to grid
	case TotalUsage:
		return "TotalUsage" // Energy being used by location
	case Balance:
		return "Balance" // Energy not being tracked in other channels
	}
	return "unknown"
}

type ChannelUsage struct {
	Name          string        `json:"name"`
	Percentage    float64       `json:"percentage"`
	NestedDevices []DeviceUsage `json:"nestedDevices"`
	Usage         float64       `json:"usage"`
	DeviceGid     uint64        `json:"deviceGid"`
	ChannelNum    string        `json:"channelNum"`
}

// Enum for EnergyUnit to be sent to GetDeviceListUsages
type EnergyUnit uint8

// Valid EnergyUnit values
const (
	KilowattHours EnergyUnit = iota
	Dollars
	AmpHours
	Trees
	GallonsOfGas
	MilesDriven
	CarbonEnergyUnit
)

func (e EnergyUnit) String() string {
	switch e {
	case KilowattHours:
		return "KilowattHours"
	case Dollars:
		return "Dollars"
	case AmpHours:
		return "AmpHours"
	case Trees:
		return "Trees"
	case GallonsOfGas:
		return "GallonsOfGas"
	case MilesDriven:
		return "MilesDriven"
	case CarbonEnergyUnit:
		return "CarbonEnergyUnit"
	}
	return "unknown"
}

func ParseEnergyUnit(s string) (EnergyUnit, error) {
	s = strings.TrimSpace(s)
	switch s {
	case "KilowattHours":
		return KilowattHours, nil
	case "Dollars":
		return Dollars, nil
	case "AmpHours":
		return AmpHours, nil
	case "Trees":
		return Trees, nil
	case "GallonsOfGas":
		return GallonsOfGas, nil
	case "MilesDriven":
		return MilesDriven, nil
	case "CarbonEnergyUnit":
		return CarbonEnergyUnit, nil
	}
	return 0, fmt.Errorf("%q is not a valid EnergyUnit", s)
}

func (e *EnergyUnit) UnmarshalJSON(byte []byte) error {
	var energyUnit string
	err := json.Unmarshal(byte, &energyUnit)
	if err != nil {
		return err
	}
	*e, err = ParseEnergyUnit(energyUnit)
	if err != nil {
		return err
	}
	return nil
}

// Enum for Scale to be sent to GetDeviceListUsages
type Scale int

// Valid Instant values
const (
	Second Scale = iota
	Minute
	Hour
	Day
	Week
	Month
	Year
)

// Format Scale into string the API expects
func (i Scale) String() string {
	switch i {
	case Second:
		return "1S"
	case Minute:
		return "1MIN"
	case Hour:
		return "1H"
	case Day:
		return "1D"
	case Week:
		return "1W"
	case Month:
		return "1MON"
	case Year:
		return "1Y"
	}
	return "unknown"
}

func ParseScale(s string) (Scale, error) {
	s = strings.TrimSpace(s)
	switch s {
	case "1S":
		return Second, nil
	case "1MIN":
		return Minute, nil
	case "1H":
		return Hour, nil
	case "1D":
		return Day, nil
	case "1W":
		return Week, nil
	case "1MON":
		return Month, nil
	case "1Y":
		return Year, nil
	}
	return 0, fmt.Errorf("%q is not a valid Scale", s)
}

func (s *Scale) UnmarshalJSON(byte []byte) error {
	var scale string
	err := json.Unmarshal(byte, &scale)
	if err != nil {
		return err
	}
	*s, err = ParseScale(scale)
	if err != nil {
		return err
	}
	return nil
}

// GetDeviceListUsages returns usage data for a list of devices
func (c *Client) GetDeviceListUsages(deviceGids []uint64, scale Scale, energyUnit EnergyUnit, instant time.Time) (*DeviceListUsages, error) {
	url, err := url.Parse(fmt.Sprintf("%s/AppAPI", apiBaseURL))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	// appending to existing query args
	q := req.URL.Query()
	q.Add("apiMethod", "getDeviceListUsages")
	q.Add("deviceGids", SplitToString(deviceGids, "+"))
	q.Add("instant", instant.Format("2006-01-02T15:04:05Z"))
	q.Add("scale", scale.String())
	q.Add("energyUnit", energyUnit.String())
	req.URL.RawQuery = q.Encode()

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get device usage: %v", err)
	}
	defer resp.Body.Close()

	response := struct {
		DeviceListUsages DeviceListUsages `json:"deviceListUsages"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode device usages: %v", err)
	}
	return &response.DeviceListUsages, nil
}

func SplitToString(a []uint64, sep string) string {
	if len(a) == 0 {
		return ""
	}

	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.FormatUint(v, 10)
	}
	return strings.Join(b, sep)
}
