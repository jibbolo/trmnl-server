package main

import (
	"context"
	"net/http"
	"time"

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

func setupHandler(ctx context.Context, input *SetupRequest) (*SetupResponse, error) {
	resp := &SetupResponse{}
	resp.Status = 200
	resp.Body.Status = 200
	resp.Body.Message = "Setup successful"
	resp.Body.APIKey = "sk-123456789013456789"
	resp.Body.FriendlyID = "ABCDEF"
	resp.Body.ImageURL = input.Proto + "://" + input.Host + "/static/empty_state.bmp"
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
	UserAgent      string  `header:"User-Agent" doc:"User Agent"`
	Width          int     `header:"Width" doc:"Width"`
	Height         int     `header:"Height" doc:"Height"`
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
		Status          int    `json:"status" doc:"Display"`
		Filename        string `json:"filename" doc:"Filename"`
		ImageURL        string `json:"image_url" doc:"Image URL"`
		ImageURLTimeout int    `json:"image_url_timeout" doc:"Image URL Timeout"`
		RefreshRate     string `json:"refresh_rate" doc:"Refresh Rate"`
		UpdateFirmware  bool   `json:"update_firmware" doc:"Update Firmware"`
		FirmwareURL     string `json:"firmware_url" doc:"Firmware URL"`
		ResetFirmware   bool   `json:"reset_firmware" doc:"Reset Firmware"`
		SpecialFunction string `json:"special_function" doc:"Special Function"`
		Action          string `json:"action" doc:"Action"`
	}
}

func displayHandler(ctx context.Context, input *DisplayRequest) (*DisplayResponse, error) {
	resp := &DisplayResponse{}
	resp.Status = http.StatusOK
	resp.Body.Status = http.StatusOK
	resp.Body.ImageURL = input.Proto + "://" + input.Host + "/static/placeholder.bmp"
	resp.Body.ImageURLTimeout = 0
	resp.Body.Filename = time.Now().Format("20060102150405") + "_placeholder.bmp"
	resp.Body.RefreshRate = "1800"
	resp.Body.UpdateFirmware = false
	resp.Body.FirmwareURL = ""
	resp.Body.ResetFirmware = false
	resp.Body.SpecialFunction = "restart_playlist"
	resp.Body.Action = ""
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
	resp.Status = http.StatusNoContent
	return resp, nil
}
