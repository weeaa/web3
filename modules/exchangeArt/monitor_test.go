package exchangeArt

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestMonitorArtists(t *testing.T) {
	resp, err := doRequest("I2LwzWoHzdcibq3ngiFtumfqmJV2")
	if err != nil {
		assert.Error(t, fmt.Errorf("expected no error but got %w", err))
	}

	if resp.StatusCode != http.StatusOK {
		assert.Error(t, fmt.Errorf("bad response status: expected %d but got %d", http.StatusOK, resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		assert.Error(t, err)
	}

	var response map[string]any
	if err = json.Unmarshal(body, &response); err != nil {
		assert.Error(t, err)
	}

	if response["statusCode"].(float64) != http.StatusOK {
		assert.Error(t, fmt.Errorf("expected status 200, got %f [%s]", response["statusCode"].(float64), response["message"].(string)))
	} else {
		assert.NoError(t, nil)
	}
}
