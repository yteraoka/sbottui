package api

// DevicesResponse is the response from GET /v1.1/devices
type DevicesResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Body       struct {
		DeviceList   []Device   `json:"deviceList"`
		InfraredRemoteList []IRDevice `json:"infraredRemoteList"`
	} `json:"body"`
}

// Device represents a physical SwitchBot device
type Device struct {
	DeviceID           string `json:"deviceId"`
	DeviceName         string `json:"deviceName"`
	DeviceType         string `json:"deviceType"`
	EnableCloudService bool   `json:"enableCloudService"`
	HubDeviceID        string `json:"hubDeviceId"`
}

// IRDevice represents an infrared remote device
type IRDevice struct {
	DeviceID       string `json:"deviceId"`
	DeviceName     string `json:"deviceName"`
	RemoteType     string `json:"remoteType"`
	HubDeviceID    string `json:"hubDeviceId"`
}

// ScenesResponse is the response from GET /v1.1/scenes
type ScenesResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Body       []Scene `json:"body"`
}

// Scene represents a SwitchBot scene
type Scene struct {
	SceneID   string `json:"sceneId"`
	SceneName string `json:"sceneName"`
}

// DeviceStatusResponse is the response from GET /v1.1/devices/{id}/status
type DeviceStatusResponse struct {
	StatusCode int          `json:"statusCode"`
	Message    string       `json:"message"`
	Body       DeviceStatus `json:"body"`
}

// DeviceStatus holds the current status of a device
type DeviceStatus struct {
	DeviceID    string  `json:"deviceId"`
	DeviceType  string  `json:"deviceType"`
	HubDeviceID string  `json:"hubDeviceId"`
	Power       string  `json:"power"`
	Brightness  int     `json:"brightness"`
	ColorTemp   int     `json:"colorTemperature"`
	Color       string  `json:"color"`
	Version     string  `json:"version"`
	Voltage     float64 `json:"voltage"`
	Weight      float64 `json:"weight"`
}

// CommandRequest is the body for POST /v1.1/devices/{id}/commands
type CommandRequest struct {
	Command     string `json:"command"`
	Parameter   string `json:"parameter,omitempty"`
	CommandType string `json:"commandType"`
}

// CommandResponse is the response from POST /v1.1/devices/{id}/commands
type CommandResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}
