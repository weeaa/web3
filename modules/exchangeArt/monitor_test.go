package exchangeArt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"testing"
)

func TestMonitorArtists(t *testing.T) {

	resp, err := doRequest("I2LwzWoHzdcibq3ngiFtumfqmJV2")
	if err != nil {
		assert.Error(t, fmt.Errorf("expected no error but got %w", err))
	}

	if resp.StatusCode != 200 {
		assert.Error(t, fmt.Errorf("expected %d but got %d", http.StatusOK, resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	log.Println(string(body))
}
