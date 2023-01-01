package emporia

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type (
	DeviceResponse struct {
		CustomerGID int      `json:"customerGid"`
		Email       string   `json:"email"`
		FirstName   string   `json:"firstName"`
		LastName    string   `json:"lastName"`
		CreatedAt   string   `json:"createdAt"`
		Devices     []Device `json:"devices"`
	}

	Device struct {
		DeviceGID            uint64           `json:"deviceGid"`
		ManufacturerDeviceID string           `json:"manufacturerDeviceId"`
		Model                string           `json:"model"`
		Firmware             string           `json:"firmware"`
		Channels             []Channel        `json:"channels"`
		Devices              []Device         `json:"devices"`
		LocationProperties   LocationProperty `json:"locationProperties"`
	}

	LocationProperty struct {
		DeviceGID             uint64              `json:"deviceGid"`
		DeviceName            string              `json:"deviceName"`
		ZipCode               string              `json:"zipCode"`
		TimeZone              string              `json:"timeZone"`
		BillingCycleStartDay  int                 `json:"billingCycleStartDay"`
		UsageCentPerKwHour    float64             `json:"usageCentPerKwHour"`
		PeakDemandDollarPerKw float64             `json:"peakDemandDollarPerKw"`
		LocationInformation   LocationInformation `json:"locationInformation"`
		LatitudeLongitude     LatitudeLongitude   `json:"latitudeLongitude"`
	}

	LatitudeLongitude struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	Channel struct {
		DeviceGID         int     `json:"deviceGid"`
		Name              string  `json:"name"`
		ChannelNum        string  `json:"channelNum"`
		ChannelMultiplier float64 `json:"channelMultiplier"`
		ChannelTypeGID    int     `json:"channelTypeGid"`
	}

	LocationInformation struct {
		AirConditioning bool   `json:"airConditioning"`
		HeatSource      string `json:"heatSource"`
		LocationSqFt    string `json:"locationSqFt"`
		NumElectricCars string `json:"numElectricCars"`
		LocationType    string `json:"locationType"`
		NumPeople       string `json:"numPeople"`
		SwimmingPool    bool   `json:"swimmingPool"`
		HotTub          bool   `json:"hotTub"`
	}
)

type ChannelType uint8

const (
	AirConditioner  ChannelType = 1
	Battery         ChannelType = 2
	Boiler          ChannelType = 3
	ClothesDryer    ChannelType = 4
	ClothesWasher   ChannelType = 5
	Computer        ChannelType = 6
	Cooktop         ChannelType = 7
	Dishwasher      ChannelType = 8
	ElectricVehicle ChannelType = 9
	Fridge          ChannelType = 10
	Furnace         ChannelType = 11
	Garage          ChannelType = 12
	Solar           ChannelType = 13
	HotTub          ChannelType = 14
	Humidifier      ChannelType = 15
	Kitchen         ChannelType = 16
	Microwave       ChannelType = 17
	Other           ChannelType = 18
	Pump            ChannelType = 19
	Room            ChannelType = 20
	SubPanel        ChannelType = 21
	WaterHeater     ChannelType = 22
	HeatPump        ChannelType = 24
	Lights          ChannelType = 26
)

var ChannelTypeLookup = map[string]ChannelType{
	"Air Conditioner":          AirConditioner,
	"Battery":                  Battery,
	"Boiler":                   Boiler,
	"Clothes Dryer":            ClothesDryer,
	"Clothes Washer":           ClothesWasher,
	"Computer/Network":         Computer,
	"Cooktop/Range/Oven/Stove": Cooktop,
	"Dishwasher":               Dishwasher,
	"Electric Vehicle/RV":      ElectricVehicle,
	"Fridge/Freezer":           Fridge,
	"Furnace":                  Furnace,
	"Garage/Shop/Barn/Shed":    Garage,
	"Solar/Generation":         Solar,
	"Hot Tub/Spa":              HotTub,
	"Humidifier/Dehumidifier":  Humidifier,
	"Kitchen":                  Kitchen,
	"Microwave":                Microwave,
	"Other":                    Other,
	"Pump":                     Pump,
	"Room/Multi-use Circuit":   Room,
	"Sub Panel":                SubPanel,
	"Water Heater":             WaterHeater,
	"Heat Pump":                HeatPump,
	"Lights":                   Lights,
}

func (c *Client) GetDevices() (*DeviceResponse, error) {
	url, err := url.Parse(fmt.Sprintf("%s/customers/devices", apiBaseURL))
	if err != nil {
		return nil, err
	}

	var deviceResponse DeviceResponse
	resp, err := c.get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&deviceResponse); err != nil {
		return nil, fmt.Errorf("failed to decode devices: %v", err)
	}

	return &deviceResponse, nil
}

func (c *Client) GetSolar() (*[]Device, error) {
	devices, err := c.GetDevices()
	if err != nil {
		return nil, err
	}

	var solarDevices []Device
	for _, device := range devices.Devices {
		if device.Model == "Solar" {
			solarDevices = append(solarDevices, device)
		}
	}

	return &solarDevices, nil
}
