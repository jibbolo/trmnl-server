package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	_, api := humatest.New(t)

	addRoutes(api)

	macAddress := "00:11:22:33:44:55"

	t.Run("setup device", func(t *testing.T) {
		resp := api.Get("/api/setup/", &SetupRequest{
			ID: macAddress,
		})
		require.Equal(t, 200, resp.Result().StatusCode)
		require.Equal(t, 200, resp.Result().StatusCode)
		output, err := convertTestResponse[*SetupResponse](resp.Body)
		require.NoError(t, err)
		assert.Equal(t, output.Body.Status, resp.Result().StatusCode, "Body status")
		assert.NotEmpty(t, output.Body.APIKey, "Api Key")
		assert.NotEmpty(t, output.Body.FriendlyID, "Friendly ID")
		assert.Contains(t, output.Body.ImageURL, "http://", "Image url")
		assert.NotEmpty(t, output.Body.Message, "message")
		assert.NotEmpty(t, output.Body.Filename, "empty_state")
	})
}

func TestDisplay(t *testing.T) {
	_, api := humatest.New(t)

	addRoutes(api)

	macAddress := "00:11:22:33:44:55"
	apiKey := "sk1234567890"

	t.Run("display", func(t *testing.T) {
		resp := api.Get("/api/display", &DisplayRequest{
			ID:             macAddress,
			AccessToken:    apiKey,
			RefreshRate:    1800,
			BatteryVoltage: 3.6,
			FWVersion:      "1.5.5",
			RSSI:           -69,
		})
		require.Equal(t, 200, resp.Result().StatusCode)
		output, err := convertTestResponse[*DisplayResponse](resp.Body)
		require.NoError(t, err)
		assert.Equal(t, output.Body.Status, 0, "Body status")
		assert.NotEmpty(t, output.Body.Filename, "2025-06-08T00:00:00")
		assert.NotEmpty(t, output.Body.RefreshRate, "1800")
		assert.Contains(t, output.Body.ImageURL, "http://", "Image url")
		assert.Empty(t, output.Body.FirmwareURL)
		assert.False(t, output.Body.ResetFirmware)
		assert.False(t, output.Body.UpdateFirmware)
	})
}

func convertTestResponse[T any](body *bytes.Buffer) (T, error) {
	payload := fmt.Sprintf(`{"Body": %s}`, body.String())
	var output T
	return output, json.Unmarshal([]byte(payload), &output)

}
