package opensea

import (
	"github.com/foundVanting/opensea-stream-go/entity"
	"github.com/mitchellh/mapstructure"
	"math/big"
	"testing"
)

func TestMonitorListings(t *testing.T) {
	mockClient := &MockClient{OpenseaClient: nil} // todo update with right client
	ItemListedEvent := entity.ItemListedEvent{}

	tests := []struct {
		name string
		slug string
	}{
		{
			name: "valid",
			slug: "boredapeyachtclub",
		},
		{
			name: "invalid",
			slug: "rAnd0m_SlUg971",
		},
	}

	for _, test := range tests {
		mockClient.OnItemListed(test.slug, func(response interface{}) {
			var expectedResponse string

			if test.name == "valid" {
				expectedResponse = ``
			} else {
				expectedResponse = ``
			}

			if response != expectedResponse {
				t.Errorf("Expected response: %s, got: %s", expectedResponse, response)
			}

			if err := mapstructure.Decode(response, &ItemListedEvent); err != nil {
				t.Error(err)
			}

			//mockClient.OpenseaClient.GetFloor()

		})

		wei := new(big.Int)
		wei.SetString(ItemListedEvent.Payload.BasePrice, 10)
	}
}

type MockClient struct {
	OpenseaClient  *Client
	itemListedData []ItemListedData
}

type ItemListedData struct {
	Slug   string
	Floor  string
	Seller string
}

func (m *MockClient) OnItemListed(slug string, callback func(interface{})) {

	itemData := ItemListedData{Slug: slug}
	m.itemListedData = append(m.itemListedData, itemData)

	// Invoke the callback function
	callback("Mock response") // Simulate a response for testing
}
