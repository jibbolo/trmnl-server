package main

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func addRoutes(api huma.API) {

	huma.Register(api, huma.Operation{
		OperationID: "device-setup",
		Path:        "/api/setup",
		Summary:     "Setup a new device",
		Method:      http.MethodGet,
	}, setupHandler)

	huma.Register(api, huma.Operation{
		OperationID: "disaply",
		Path:        "/api/display",
		Summary:     "Update display image",
		Method:      http.MethodGet,
	}, displayHandler)

}

type SetupRequest struct {
	ID string `header:"ID" doc:"Device Mac Address"`
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

func setupHandler(ctx context.Context, input *SetupRequest) (*SetupResponse, error) {
	resp := &SetupResponse{}
	resp.Status = 200
	resp.Body.Status = 200
	resp.Body.Message = "Setup successful"
	resp.Body.APIKey = "sk-123456789013456789"
	resp.Body.FriendlyID = "ABCDEF"
	resp.Body.ImageURL = "http://localhost:8888/static/placeholder.png"
	resp.Body.Filename = "empty_state"
	return resp, nil
}

type DisplayRequest struct {
	ID             string  `header:"ID" doc:"Device Mac Address"`
	AccessToken    string  `header:"Access-Token" doc:"Access Token"`
	RefreshRate    int     `header:"Refresh-Rate" doc:"Refresh Rate"`
	BatteryVoltage float64 `header:"Battery-Voltage" doc:"Battery Voltage"`
	FWVersion      string  `header:"FW-Version" doc:"Firmware Version"`
	RSSI           int     `header:"RSSI" doc:"Received Signal Strength Indicator"`
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

func displayHandler(ctx context.Context, input *DisplayRequest) (*DisplayResponse, error) {
	resp := &DisplayResponse{}
	resp.Status = 200
	resp.Body.Status = 200
	resp.Body.ImageURL = "http://localhost:8888/static/placeholder.png"
	resp.Body.Filename = "2025-06-08T00:00:00"
	resp.Body.RefreshRate = "1800"
	resp.Body.UpdateFirmware = false
	resp.Body.FirmwareURL = ""
	resp.Body.ResetFirmware = false
	return resp, nil
}
