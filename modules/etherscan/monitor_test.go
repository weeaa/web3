package etherscan

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestMonitorVerifiedContracts(t *testing.T) {
	resp, err := doRequest()
	if err != nil {
		assert.Error(t, fmt.Errorf("expected no error but got %w", err))
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		assert.Error(t, fmt.Errorf("expected %d but got %d", http.StatusOK, resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		assert.Error(t, fmt.Errorf("expected no error but got %w", err))
	}

	contract := ParseHTML(goquery.NewDocumentFromReader(strings.NewReader(string(body))))
	if contract.Name == "" || contract.Address == "" {
		assert.Error(t, fmt.Errorf("expected non empty contract name & address but got [%s – %s] ", contract.Name, contract.Address))
	}

	assert.NoError(t, nil)
}
