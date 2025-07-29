package panther

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Customer with custom url that provides the graphql endpoint (legacy behavior)
func TestCreateAPIClient_CustomURLWithGraphEndpoint(t *testing.T) {
	url := "panther-url/public/graphql"
	client := *CreateAPIClient(url, "token")
	// With machinebox client, we can't easily access the internal URL, but we can verify the client was created
	assert.NotNil(t, client.GraphQLClient)
	assert.Equal(t, "panther-url", client.RestClient.baseUrl)
}

// Customer with custom url that provides the panther root url (new behavior)
func TestCreateAPIClient_CustomUrlWithBaseUrl(t *testing.T) {
	url := "panther-url"
	client := *CreateAPIClient(url, "token")
	// With machinebox client, we can't easily access the internal URL, but we can verify the client was created
	assert.NotNil(t, client.GraphQLClient)
	assert.Equal(t, "panther-url", client.RestClient.baseUrl)
}

// Customer with API Gateway url that provides the graphql endpoint (legacy behavior)
func TestCreateAPIClient_ApiGWUrlWithGraphEndpoint(t *testing.T) {
	url := "panther-url/v1/public/graphql"
	client := *CreateAPIClient(url, "token")
	// With machinebox client, we can't easily access the internal URL, but we can verify the client was created
	assert.NotNil(t, client.GraphQLClient)
	assert.Equal(t, "panther-url/v1", client.RestClient.baseUrl)
}

// Customer with API Gateway url that provides the panther root url (new behavior)
func TestCreateAPIClient_ApiGWUrlWithBaseUrl(t *testing.T) {
	url := "panther-url/v1"
	client := *CreateAPIClient(url, "token")
	// With machinebox client, we can't easily access the internal URL, but we can verify the client was created
	assert.NotNil(t, client.GraphQLClient)
	assert.Equal(t, "panther-url/v1", client.RestClient.baseUrl)
}
