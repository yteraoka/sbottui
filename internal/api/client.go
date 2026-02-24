package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://api.switch-bot.com/v1.1"

// Client is a SwitchBot API v1.1 client.
type Client struct {
	token      string
	secret     string
	httpClient *http.Client
}

// NewClient creates a new Client with the given token and secret.
func NewClient(token, secret string) *Client {
	return &Client{
		token:  token,
		secret: secret,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) do(method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	for k, v := range AuthHeaders(c.token, c.secret) {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(data))
	}

	return data, nil
}

// GetDevices fetches all physical and IR devices.
func (c *Client) GetDevices() (*DevicesResponse, error) {
	data, err := c.do("GET", "/devices", nil)
	if err != nil {
		return nil, err
	}
	var resp DevicesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal devices: %w", err)
	}
	if resp.StatusCode != 100 {
		return nil, &ErrAPI{StatusCode: resp.StatusCode, Message: resp.Message}
	}
	return &resp, nil
}

// GetScenes fetches all scenes.
func (c *Client) GetScenes() (*ScenesResponse, error) {
	data, err := c.do("GET", "/scenes", nil)
	if err != nil {
		return nil, err
	}
	var resp ScenesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal scenes: %w", err)
	}
	if resp.StatusCode != 100 {
		return nil, &ErrAPI{StatusCode: resp.StatusCode, Message: resp.Message}
	}
	return &resp, nil
}

// GetDeviceStatus fetches the current status of a device.
func (c *Client) GetDeviceStatus(deviceID string) (*DeviceStatus, error) {
	data, err := c.do("GET", "/devices/"+deviceID+"/status", nil)
	if err != nil {
		return nil, err
	}
	var resp DeviceStatusResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal device status: %w", err)
	}
	if resp.StatusCode != 100 {
		return nil, &ErrAPI{StatusCode: resp.StatusCode, Message: resp.Message}
	}
	return &resp.Body, nil
}

// SendCommand sends a command to a physical device.
func (c *Client) SendCommand(deviceID string, cmd CommandRequest) (*CommandResponse, error) {
	data, err := c.do("POST", "/devices/"+deviceID+"/commands", cmd)
	if err != nil {
		return nil, err
	}
	var resp CommandResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal command response: %w", err)
	}
	if resp.StatusCode != 100 {
		return nil, &ErrAPI{StatusCode: resp.StatusCode, Message: resp.Message}
	}
	return &resp, nil
}

// SendIRCommand sends a command to an IR device.
func (c *Client) SendIRCommand(deviceID string, cmd CommandRequest) (*CommandResponse, error) {
	cmd.CommandType = "command"
	return c.SendCommand(deviceID, cmd)
}

// ExecuteScene triggers a scene by ID.
func (c *Client) ExecuteScene(sceneID string) (*CommandResponse, error) {
	data, err := c.do("POST", "/scenes/"+sceneID+"/execute", nil)
	if err != nil {
		return nil, err
	}
	var resp CommandResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal scene response: %w", err)
	}
	if resp.StatusCode != 100 {
		return nil, &ErrAPI{StatusCode: resp.StatusCode, Message: resp.Message}
	}
	return &resp, nil
}
