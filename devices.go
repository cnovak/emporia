package emporia

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type (
	DeviceResponse struct {
		CustomerGid int      `json:"customerGid"`
		Email       string   `json:"email"`
		FirstName   string   `json:"firstName"`
		LastName    string   `json:"lastName"`
		CreatedAt   string   `json:"createdAt"`
		Devices     []Device `json:"devices"`
	}

	Device struct {
		DeviceGid            int64            `json:"deviceGid"`
		ManufacturerDeviceID string           `json:"manufacturerDeviceId"`
		Model                string           `json:"model"`
		Firmware             string           `json:"firmware"`
		Channels             []Channel        `json:"channels"`
		Devices              []Device         `json:"devices"`
		LocationProperties   LocationProperty `json:"locationProperties"`
	}

	LocationProperty struct {
		DeviceGid             int64               `json:"deviceGid"`
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
		DeviceGid         int     `json:"deviceGid"`
		Name              string  `json:"name"`
		ChannelNum        string  `json:"channelNum"`
		ChannelMultiplier float64 `json:"channelMultiplier"`
		ChannelTypeGid    int     `json:"channelTypeGid"`
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

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Printf("============== \nResponse Body:\n============== \n%s\n\n", string(body))

	if err := json.NewDecoder(resp.Body).Decode(&deviceResponse); err != nil {
		return nil, fmt.Errorf("failed to decode devices: %v", err)
	}

	return &deviceResponse, nil
}
