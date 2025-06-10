package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func addRoutes(api huma.API) {

	huma.Register(api, huma.Operation{
		OperationID: "device-setup",
		Path:        "/api/setup/",
		Summary:     "Setup a new device",
		Method:      http.MethodGet,
	}, setupHandler)

	huma.Register(api, huma.Operation{
		OperationID: "display",
		Path:        "/api/display",
		Summary:     "Update display image",
		Method:      http.MethodGet,
	}, displayHandler)

	huma.Register(api, huma.Operation{
		OperationID: "log",
		Path:        "/api/log",
		Summary:     "Log from device",
		Method:      http.MethodPost,
	}, logHandler)

}

type SetupRequest struct {
	Host  string
	Proto string
	ID    string `header:"ID" doc:"Device Mac Address"`
}

func (m *SetupRequest) Resolve(ctx huma.Context) []error {
	m.Host = ctx.Host()
	m.Proto = "http"
	if m.Host == "trmnl.gmar.dev" {
		m.Proto = "https"
	}
	return nil
}

type SetupResponse struct {
	Status int
	Body   struct {
		Status     int    `json:"status" doc:"Setup status"`
		Message    string `json:"message" doc:"Message"`
		APIKey     string `json:"api_key" doc:"API Key"`
		FriendlyID string `json:"friendly_id" doc:"Friendly ID"`
		ImageURL   string `json:"image_url" doc:"Image URL"`
		Filename   string `json:"filename" doc:"Filename"`
	}
}

// {"log":{"logs_array":[{"creation_timestamp":1749404319,"device_status_stamp":{"wifi_rssi_level":-66,"wifi_status":"connected","refresh_rate":900,"time_since_last_sleep_start":280,"current_fw_version":"1.5.5","special_function":"none","battery_voltage":3.844,"wakeup_reason":"button","free_heap_size":215512,"max_alloc_size":192500},"log_id":29,"log_message":"Failed to resolve hostname after 5 attempts, continuing...","log_codeline":573,"log_sourcefile":"src/bl.cpp","additional_info":{"filename_current":"","filename_new":"","retry_attempt":1}}]}}
func setupHandler(ctx context.Context, input *SetupRequest) (*SetupResponse, error) {
	fmt.Printf("Received log: %v\n", input)
	resp := &SetupResponse{}
	resp.Status = 200
	resp.Body.Status = 200
	resp.Body.Message = "Setup successful"
	resp.Body.APIKey = "sk-123456789013456789"
	resp.Body.FriendlyID = "ABCDEF"
	resp.Body.ImageURL = input.Proto + "://" + input.Host + "/static/placeholder.png"
	resp.Body.Filename = "empty_state"
	return resp, nil
}

type DisplayRequest struct {
	Host           string
	Proto          string
	ID             string  `header:"ID" doc:"Device Mac Address"`
	AccessToken    string  `header:"Access-Token" doc:"Access Token"`
	RefreshRate    int     `header:"Refresh-Rate" doc:"Refresh Rate"`
	BatteryVoltage float64 `header:"Battery-Voltage" doc:"Battery Voltage"`
	FWVersion      string  `header:"FW-Version" doc:"Firmware Version"`
	RSSI           int     `header:"RSSI" doc:"Received Signal Strength Indicator"`
}

func (m *DisplayRequest) Resolve(ctx huma.Context) []error {
	m.Host = ctx.Host()
	m.Proto = "http"
	if m.Host == "trmnl.gmar.dev" {
		m.Proto = "https"
	}
	return nil
}

type DisplayResponse struct {
	Status int
	Body   struct {
		Status         int    `json:"status" doc:"Display"`
		ImageURL       string `json:"image_url" doc:"Image URL"`
		Filename       string `json:"filename" doc:"Filename"`
		RefreshRate    string `json:"refresh_rate" doc:"Refresh Rate"`
		UpdateFirmware bool   `json:"update_firmware" doc:"Update Firmware"`
		FirmwareURL    string `json:"firmware_url" doc:"Firmware URL"`
		ResetFirmware  bool   `json:"reset_firmware" doc:"Reset Firmware"`
	}
}

// Received log: &{192.168.1.51:8888 DC:06:75:B8:89:2C 7Bi1rqZlnDFg16dG9ZIKK4 900 3.85 1.5.5 -45}
func displayHandler(ctx context.Context, input *DisplayRequest) (*DisplayResponse, error) {
	fmt.Printf("Received log: %v\n", input)
	resp := &DisplayResponse{}
	resp.Status = 200
	resp.Body.Status = 200
	resp.Body.ImageURL = input.Proto + "://" + input.Host + "/output/placeholder.png"
	resp.Body.Filename = "2025-06-08T00:00:00"
	resp.Body.RefreshRate = "1800"
	resp.Body.UpdateFirmware = false
	resp.Body.FirmwareURL = ""
	resp.Body.ResetFirmware = false
	return resp, nil
}

type LogRequest struct {
	RawBody []byte
}

type LogResponse struct {
	Status int
}

func logHandler(ctx context.Context, input *LogRequest) (*LogResponse, error) {
	resp := &LogResponse{}
	resp.Status = 200
	return resp, nil
}
