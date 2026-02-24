package tui

import (
	"github.com/yteraoka/sbottui/internal/api"
	"github.com/yteraoka/sbottui/internal/domain"
)

// MsgLoaded is sent when devices and scenes are successfully loaded.
type MsgLoaded struct {
	Items []domain.ListItem
}

// MsgLoadError is sent when loading fails.
type MsgLoadError struct {
	Err error
}

// MsgDeviceStatus is sent when a device status is fetched.
type MsgDeviceStatus struct {
	Status *api.DeviceStatus
	Err    error
}

// MsgCommandDone is sent when a command completes.
type MsgCommandDone struct {
	Err error
}

// MsgSceneDone is sent when a scene execution completes.
type MsgSceneDone struct {
	Name string
	Err  error
}
