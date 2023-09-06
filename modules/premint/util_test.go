package premint

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/pkg/tls"
	"net/http"
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name       string
		publicKey  string
		privateKey string
		proxy      string
	}{
		{
			name:       "valid login response – with proxy",
			publicKey:  os.Getenv("PREMINT_PUB_KEY"),
			privateKey: os.Getenv("PREMINT_PRIV_KEY"),
			proxy:      tls.TestProxy,
		},
		{
			name:       "valid login response – without proxy",
			publicKey:  os.Getenv("PREMINT_PUB_KEY"),
			privateKey: os.Getenv("PREMINT_PRIV_KEY"),
			proxy:      "",
		},
		{
			name:       "invalid login response – random public & private key",
			publicKey:  "0xRandomValue",
			privateKey: "0x123myPrivateKey0x",
			proxy:      "",
		},
	}

	for _, test := range tests {

		profile := NewProfile(test.publicKey, test.privateKey, test.proxy, 5000)
		err := profile.login()
		if err != nil {
			assert.Error(t, err)
		}
		assert.NoError(t, err)

		profile.Monitor(discord.NewClient(discord.UnisatSettings{},
			discord.PremintSettings{
				DiscordWebhook: os.Getenv("PREMINT_WEBHOOK"),
				Verbose:        false,
			},
			discord.EtherscanSettings{},
			discord.OpenseaSettings{},
			discord.ExchangeartSettings{},
			discord.LaunchmynftSettings{},
			"", "", 1), []RaffleType{Daily})
	}

}

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestProfile_Login_Success(t *testing.T) {
	// Create a new instance of the Profile with your mocked HTTP client
	mockClient := &MockHTTPClient{}

	profile := &Profile{
		publicAddress: "your_public_address",
		privateKey:    "your_private_key",
	}

	// Set up expectations for the first request (GET)
	getRequest := &http.Request{
		// Set your expected request details here
	}
	mockClient.On("Do", getRequest).Return(&http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		// Set your expected response details here
	}, nil)

	// Set up expectations for the second request (POST)
	postRequest := &http.Request{
		// Set your expected request details here
	}
	mockClient.On("Do", postRequest).Return(&http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		// Set your expected response details here
	}, nil)

	// Set up expectations for the third request (POST)
	postRequest2 := &http.Request{
		// Set your expected request details here
	}
	mockClient.On("Do", postRequest2).Return(&http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		// Set your expected response details here
	}, nil)

	// Call the login function
	err := profile.login()

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the expectations for the HTTP requests were met
	mockClient.AssertExpectations(t)
}

func TestProfile_Login_Error(t *testing.T) {
	// Create a new instance of the Profile with your mocked HTTP client
	profile := &Profile{
		Client:        &MockHTTPClient{},
		publicAddress: "your_public_address",
		privateKey:    "your_private_key",
	}

	// Set up expectations for the first request (GET) to return an error
	mockHTTP := profile.Client.(*MockHTTPClient)
	getRequest := &http.Request{
		// Set your expected request details here
	}
	mockHTTP.On("Do", getRequest).Return(nil, errors.New("HTTP request failed"))

	// Call the login function
	err := profile.login()

	// Assert that there was an error
	assert.Error(t, err)

	// Assert that the expectations for the HTTP requests were met
	mockHTTP.AssertExpectations(t)
}
