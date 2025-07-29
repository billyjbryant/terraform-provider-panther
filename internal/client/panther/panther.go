/*
Copyright 2023 Panther Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package panther

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/machinebox/graphql"
	"io"
	"net/http"
	"strings"
	"terraform-provider-panther/internal/client"
)

const GraphqlPath = "/public/graphql"
const RestBasePath = "/v1"
const RestHttpSourcePath = "/log-sources/http"

var _ client.GraphQLClient = (*GraphQLClient)(nil)

var _ client.RestClient = (*RestClient)(nil)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type APIClient struct {
	*GraphQLClient
	*RestClient
}

type GraphQLClient struct {
	client *graphql.Client
	token  string
}

type RestClient struct {
	baseUrl string
	Doer
}

func NewGraphQLClient(url, token string) *GraphQLClient {
	return &GraphQLClient{
		client: graphql.NewClient(fmt.Sprintf("%s%s", url, GraphqlPath)),
		token:  token,
	}
}

func NewRestClient(url, token string) *RestClient {
	// Detect if this is an API Gateway URL (contains /v1) vs direct Panther URL
	var baseUrl string
	if strings.Contains(url, RestBasePath) {
		// API Gateway URL - use /v1 prefix
		if strings.HasSuffix(url, RestBasePath) {
			baseUrl = url
		} else {
			baseUrl = fmt.Sprintf("%s%s", url, RestBasePath)
		}
	} else {
		// Direct Panther URL - no /v1 prefix for HTTP sources, but /v1 for other REST endpoints
		baseUrl = url
	}
	
	return &RestClient{
		baseUrl: baseUrl,
		Doer:    NewAuthorizedHTTPClient(token),
	}
}

func NewAPIClient(graphClient *GraphQLClient, restClient *RestClient) *APIClient {
	return &APIClient{
		graphClient,
		restClient,
	}
}

func CreateAPIClient(url, token string) *APIClient {
	// url in previous versions was provided including graphql endpoint,
	// we strip it here to keep it backwards compatible
	pantherUrl := strings.TrimSuffix(url, GraphqlPath)
	graphClient := NewGraphQLClient(pantherUrl, token)
	restClient := NewRestClient(pantherUrl, token)

	return NewAPIClient(graphClient, restClient)
}

func (c *RestClient) CreateHttpSource(ctx context.Context, input client.CreateHttpSourceInput) (client.HttpSource, error) {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return client.HttpSource{}, fmt.Errorf("error marshaling data: %w", err)
	}
	// Construct HTTP source URL based on whether this is API Gateway or direct Panther URL
	var url string
	if strings.Contains(c.baseUrl, RestBasePath) {
		// API Gateway URL - append to /v1 base
		url = fmt.Sprintf("%s%s", c.baseUrl, RestHttpSourcePath)
	} else {
		// Direct Panther URL - append directly
		url = fmt.Sprintf("%s%s", c.baseUrl, RestHttpSourcePath)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return client.HttpSource{}, fmt.Errorf("failed to create http request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return client.HttpSource{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return client.HttpSource{}, fmt.Errorf("failed to make request, status: %d, message: %s", resp.StatusCode, getErrorResponseMsg(resp))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return client.HttpSource{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var response client.HttpSource
	if err = json.Unmarshal(body, &response); err != nil {
		return client.HttpSource{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return response, nil
}

func (c *RestClient) UpdateHttpSource(ctx context.Context, input client.UpdateHttpSourceInput) (client.HttpSource, error) {
	// Construct HTTP source URL based on whether this is API Gateway or direct Panther URL
	var baseURL string
	if strings.Contains(c.baseUrl, RestBasePath) {
		// API Gateway URL - append to /v1 base
		baseURL = fmt.Sprintf("%s%s", c.baseUrl, RestHttpSourcePath)
	} else {
		// Direct Panther URL - append directly
		baseURL = fmt.Sprintf("%s%s", c.baseUrl, RestHttpSourcePath)
	}
	reqURL := fmt.Sprintf("%s/%s", baseURL, input.IntegrationId)
	jsonData, err := json.Marshal(input)
	if err != nil {
		return client.HttpSource{}, fmt.Errorf("error marshaling data: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, reqURL, bytes.NewReader(jsonData))
	if err != nil {
		return client.HttpSource{}, fmt.Errorf("failed to create http request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return client.HttpSource{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.HttpSource{}, fmt.Errorf("failed to make request, status: %d, message: %s", resp.StatusCode, getErrorResponseMsg(resp))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return client.HttpSource{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var response client.HttpSource
	if err = json.Unmarshal(body, &response); err != nil {
		return client.HttpSource{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return response, nil
}

func (c *RestClient) GetHttpSource(ctx context.Context, id string) (client.HttpSource, error) {
	// Construct HTTP source URL based on whether this is API Gateway or direct Panther URL
	var baseURL string
	if strings.Contains(c.baseUrl, RestBasePath) {
		// API Gateway URL - append to /v1 base
		baseURL = fmt.Sprintf("%s%s", c.baseUrl, RestHttpSourcePath)
	} else {
		// Direct Panther URL - append directly
		baseURL = fmt.Sprintf("%s%s", c.baseUrl, RestHttpSourcePath)
	}
	reqURL := fmt.Sprintf("%s/%s", baseURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return client.HttpSource{}, fmt.Errorf("failed to create http request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return client.HttpSource{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.HttpSource{}, fmt.Errorf("failed to make request, status: %d, message: %s", resp.StatusCode, getErrorResponseMsg(resp))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return client.HttpSource{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var response client.HttpSource
	if err = json.Unmarshal(body, &response); err != nil {
		return client.HttpSource{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return response, nil
}

func (c *RestClient) DeleteHttpSource(ctx context.Context, id string) error {
	// Construct HTTP source URL based on whether this is API Gateway or direct Panther URL
	var baseURL string
	if strings.Contains(c.baseUrl, RestBasePath) {
		// API Gateway URL - append to /v1 base
		baseURL = fmt.Sprintf("%s%s", c.baseUrl, RestHttpSourcePath)
	} else {
		// Direct Panther URL - append directly
		baseURL = fmt.Sprintf("%s%s", c.baseUrl, RestHttpSourcePath)
	}
	reqURL := fmt.Sprintf("%s/%s", baseURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to make request, status: %d, message: %s", resp.StatusCode, getErrorResponseMsg(resp))
	}

	return nil
}

func (c *GraphQLClient) UpdateS3Source(ctx context.Context, input client.UpdateS3SourceInput) (client.UpdateS3SourceOutput, error) {
	req := graphql.NewRequest(`
		mutation($input: UpdateS3SourceInput!) {
			updateS3Source(input: $input) {
				logSource {
					integrationId
					integrationLabel
					s3Bucket
					s3Prefix
					s3PrefixLogTypes {
						prefix
						logTypes
					}
				}
			}
		}
	`)

	req.Var("input", input)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		UpdateS3Source client.UpdateS3SourceOutput `json:"updateS3Source"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return client.UpdateS3SourceOutput{}, fmt.Errorf("GraphQL mutation failed: %v", err)
	}
	return response.UpdateS3Source, nil
}

func (c *GraphQLClient) DeleteSource(ctx context.Context, input client.DeleteSourceInput) (client.DeleteSourceOutput, error) {
	req := graphql.NewRequest(`
		mutation($input: DeleteSourceInput!) {
			deleteSource(input: $input) {
				__typename
			}
		}
	`)

	req.Var("input", input)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		DeleteSource struct {
			Typename string `json:"__typename"`
		} `json:"deleteSource"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return client.DeleteSourceOutput{}, fmt.Errorf("GraphQL mutation failed: %v", err)
	}
	return client.DeleteSourceOutput{}, nil
}

func (c *GraphQLClient) GetS3Source(ctx context.Context, id string) (*client.S3LogIntegration, error) {
	req := graphql.NewRequest(`
		query($id: ID!) {
			source(id: $id) {
				... on S3LogIntegration {
					integrationId
					integrationLabel
					s3Bucket
					s3Prefix
					s3PrefixLogTypes {
						prefix
						logTypes
						excludedPrefixes
					}
					awsAccountId
					kmsKey
					logProcessingRole
					logStreamType
					managedBucketNotifications
				}
			}
		}
	`)

	req.Var("id", id)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		Source client.S3LogIntegration `json:"source"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %v", err)
	}
	return &response.Source, nil
}

func (c *GraphQLClient) CreateS3Source(ctx context.Context, input client.CreateS3SourceInput) (client.CreateS3SourceOutput, error) {
	req := graphql.NewRequest(`
		mutation($input: CreateS3SourceInput!) {
			createS3Source(input: $input) {
				logSource {
					integrationId
					integrationLabel
					s3Bucket
					s3Prefix
					s3PrefixLogTypes {
						prefix
						logTypes
					}
					awsAccountId
				}
			}
		}
	`)

	req.Var("input", input)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		CreateS3Source client.CreateS3SourceOutput `json:"createS3Source"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return client.CreateS3SourceOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return response.CreateS3Source, nil
}

func getErrorResponseMsg(resp *http.Response) string {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("failed to read response body: %s", err.Error())
	}

	var errResponse client.HttpErrorResponse
	if err = json.Unmarshal(body, &errResponse); err != nil {
		return fmt.Sprintf("failed to unmarshal response body to get error response: %s", err.Error())
	}

	return errResponse.Message
}

// Generic REST helper methods
func (c *RestClient) doRequest(ctx context.Context, method, path string, input interface{}, expectedStatus int) ([]byte, error) {
	var body io.Reader
	if input != nil {
		jsonData, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("error marshaling data: %w", err)
		}
		body = bytes.NewReader(jsonData)
	}

	url := fmt.Sprintf("%s%s", c.baseUrl, path)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	if input != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		return nil, fmt.Errorf("failed to make request, status: %d, message: %s", resp.StatusCode, getErrorResponseMsg(resp))
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return responseBody, nil
}

// Rules methods
func (c *RestClient) CreateRule(ctx context.Context, input client.CreateRuleInput) (client.Rule, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/rules", input, http.StatusOK)
	if err != nil {
		return client.Rule{}, err
	}

	var response client.Rule
	if err = json.Unmarshal(body, &response); err != nil {
		return client.Rule{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return response, nil
}

func (c *RestClient) UpdateRule(ctx context.Context, input client.UpdateRuleInput) (client.Rule, error) {
	path := fmt.Sprintf("/rules/%s", input.ID)
	body, err := c.doRequest(ctx, http.MethodPut, path, input, http.StatusOK)
	if err != nil {
		return client.Rule{}, err
	}

	var response client.Rule
	if err = json.Unmarshal(body, &response); err != nil {
		return client.Rule{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return response, nil
}

func (c *RestClient) GetRule(ctx context.Context, id string) (client.Rule, error) {
	path := fmt.Sprintf("/rules/%s", id)
	body, err := c.doRequest(ctx, http.MethodGet, path, nil, http.StatusOK)
	if err != nil {
		return client.Rule{}, err
	}

	var response client.Rule
	if err = json.Unmarshal(body, &response); err != nil {
		return client.Rule{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return response, nil
}

func (c *RestClient) DeleteRule(ctx context.Context, id string) error {
	path := fmt.Sprintf("/rules/%s", id)
	_, err := c.doRequest(ctx, http.MethodDelete, path, nil, http.StatusNoContent)
	return err
}



// Data Models methods
func (c *RestClient) CreateDataModel(ctx context.Context, input client.CreateDataModelInput) (client.DataModel, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/data-models", input, http.StatusCreated)
	if err != nil {
		return client.DataModel{}, err
	}

	var response client.DataModel
	if err = json.Unmarshal(body, &response); err != nil {
		return client.DataModel{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return response, nil
}

func (c *RestClient) UpdateDataModel(ctx context.Context, input client.UpdateDataModelInput) (client.DataModel, error) {
	path := fmt.Sprintf("/data-models/%s", input.ID)
	body, err := c.doRequest(ctx, http.MethodPut, path, input, http.StatusOK)
	if err != nil {
		return client.DataModel{}, err
	}

	var response client.DataModel
	if err = json.Unmarshal(body, &response); err != nil {
		return client.DataModel{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return response, nil
}

func (c *RestClient) GetDataModel(ctx context.Context, id string) (client.DataModel, error) {
	path := fmt.Sprintf("/data-models/%s", id)
	body, err := c.doRequest(ctx, http.MethodGet, path, nil, http.StatusOK)
	if err != nil {
		return client.DataModel{}, err
	}

	var response client.DataModel
	if err = json.Unmarshal(body, &response); err != nil {
		return client.DataModel{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return response, nil
}

func (c *RestClient) DeleteDataModel(ctx context.Context, id string) error {
	path := fmt.Sprintf("/data-models/%s", id)
	_, err := c.doRequest(ctx, http.MethodDelete, path, nil, http.StatusNoContent)
	return err
}

// Globals methods
func (c *RestClient) CreateGlobal(ctx context.Context, input client.CreateGlobalInput) (client.Global, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/globals", input, http.StatusCreated)
	if err != nil {
		return client.Global{}, err
	}

	var response client.Global
	if err = json.Unmarshal(body, &response); err != nil {
		return client.Global{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return response, nil
}

func (c *RestClient) UpdateGlobal(ctx context.Context, input client.UpdateGlobalInput) (client.Global, error) {
	path := fmt.Sprintf("/globals/%s", input.ID)
	body, err := c.doRequest(ctx, http.MethodPut, path, input, http.StatusOK)
	if err != nil {
		return client.Global{}, err
	}

	var response client.Global
	if err = json.Unmarshal(body, &response); err != nil {
		return client.Global{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return response, nil
}

func (c *RestClient) GetGlobal(ctx context.Context, id string) (client.Global, error) {
	path := fmt.Sprintf("/globals/%s", id)
	body, err := c.doRequest(ctx, http.MethodGet, path, nil, http.StatusOK)
	if err != nil {
		return client.Global{}, err
	}

	var response client.Global
	if err = json.Unmarshal(body, &response); err != nil {
		return client.Global{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return response, nil
}

func (c *RestClient) DeleteGlobal(ctx context.Context, id string) error {
	path := fmt.Sprintf("/globals/%s", id)
	_, err := c.doRequest(ctx, http.MethodDelete, path, nil, http.StatusNoContent)
	return err
}

// GraphQL Client methods for Cloud Accounts
func (c *GraphQLClient) CreateCloudAccount(ctx context.Context, input client.CreateCloudAccountInput) (client.CreateCloudAccountOutput, error) {
	req := graphql.NewRequest(`
		mutation($input: CreateCloudAccountInput!) {
			createCloudAccount(input: $input) {
				cloudAccount {
					id
					awsAccountId
					label
					awsStackName
					awsScanConfig {
						auditRole
					}
					awsRegionIgnoreList
					resourceTypeIgnoreList
					resourceRegexIgnoreList
					isEditable
				}
			}
		}
	`)

	req.Var("input", input)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		CreateCloudAccount client.CreateCloudAccountOutput `json:"createCloudAccount"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return client.CreateCloudAccountOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return response.CreateCloudAccount, nil
}

func (c *GraphQLClient) UpdateCloudAccount(ctx context.Context, input client.UpdateCloudAccountInput) (client.UpdateCloudAccountOutput, error) {
	req := graphql.NewRequest(`
		mutation($input: UpdateCloudAccountInput!) {
			updateCloudAccount(input: $input) {
				cloudAccount {
					id
					awsAccountId
					label
					awsStackName
					awsScanConfig {
						auditRole
					}
					awsRegionIgnoreList
					resourceTypeIgnoreList
					resourceRegexIgnoreList
					isEditable
				}
			}
		}
	`)

	req.Var("input", input)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		UpdateCloudAccount client.UpdateCloudAccountOutput `json:"updateCloudAccount"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return client.UpdateCloudAccountOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return response.UpdateCloudAccount, nil
}

func (c *GraphQLClient) GetCloudAccount(ctx context.Context, id string) (*client.CloudAccount, error) {
	req := graphql.NewRequest(`
		query($id: ID!) {
			cloudAccount(id: $id) {
				id
				awsAccountId
				label
				awsStackName
				awsScanConfig {
					auditRole
				}
				awsRegionIgnoreList
				resourceTypeIgnoreList
				resourceRegexIgnoreList
				isEditable
			}
		}
	`)

	req.Var("id", id)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		CloudAccount client.CloudAccount `json:"cloudAccount"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %w", err)
	}
	return &response.CloudAccount, nil
}

func (c *GraphQLClient) DeleteCloudAccount(ctx context.Context, input client.DeleteCloudAccountInput) (client.DeleteCloudAccountOutput, error) {
	req := graphql.NewRequest(`
		mutation($input: DeleteCloudAccountInput!) {
			deleteCloudAccount(input: $input) {
				id
			}
		}
	`)

	req.Var("input", input)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		DeleteCloudAccount client.DeleteCloudAccountOutput `json:"deleteCloudAccount"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return client.DeleteCloudAccountOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return response.DeleteCloudAccount, nil
}

// GraphQL Client methods for Schemas (using machinebox client)
func (c *GraphQLClient) CreateSchema(ctx context.Context, input client.CreateSchemaInput) (client.CreateSchemaOutput, error) {
	req := graphql.NewRequest(`
		mutation($input: CreateOrUpdateSchemaInput!) {
			createOrUpdateSchema(input: $input) {
				schema {
					name
					description
					spec
					version
					revision
					isFieldDiscoveryEnabled
					createdAt
					updatedAt
				}
			}
		}
	`)

	req.Var("input", input)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		CreateOrUpdateSchema client.CreateSchemaOutput `json:"createOrUpdateSchema"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return client.CreateSchemaOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}

	return response.CreateOrUpdateSchema, nil
}

func (c *GraphQLClient) UpdateSchema(ctx context.Context, input client.UpdateSchemaInput) (client.UpdateSchemaOutput, error) {
	// Convert UpdateSchemaInput to CreateSchemaInput since they use the same mutation
	createInput := client.CreateSchemaInput{
		Name:                    input.Name,
		Description:             input.Description,
		Spec:                    input.Spec,
		IsFieldDiscoveryEnabled: input.IsFieldDiscoveryEnabled,
		Revision:                input.Revision,
	}

	req := graphql.NewRequest(`
		mutation($input: CreateOrUpdateSchemaInput!) {
			createOrUpdateSchema(input: $input) {
				schema {
					name
					description
					spec
					version
					revision
					isFieldDiscoveryEnabled
					createdAt
					updatedAt
				}
			}
		}
	`)

	req.Var("input", createInput)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		CreateOrUpdateSchema client.UpdateSchemaOutput `json:"createOrUpdateSchema"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return client.UpdateSchemaOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}

	return response.CreateOrUpdateSchema, nil
}

func (c *GraphQLClient) GetSchema(ctx context.Context, name string) (*client.Schema, error) {
	req := graphql.NewRequest(`
		query($input: SchemasInput!) {
			schemas(input: $input) {
				edges {
					node {
						name
						description
						spec
						version
						revision
						isFieldDiscoveryEnabled
						createdAt
						updatedAt
					}
				}
			}
		}
	`)

	input := map[string]interface{}{
		"cursor": "",
	}

	req.Var("input", input)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		Schemas struct {
			Edges []struct {
				Node client.Schema `json:"node"`
			} `json:"edges"`
		} `json:"schemas"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %w", err)
	}

	// Find schema by name
	for _, edge := range response.Schemas.Edges {
		if edge.Node.Name == name {
			return &edge.Node, nil
		}
	}

	return nil, nil // Schema not found
}

func (c *GraphQLClient) DeleteSchema(ctx context.Context, input client.DeleteSchemaInput) (client.DeleteSchemaOutput, error) {
	req := graphql.NewRequest(`
		mutation($input: UpdateSchemaStatusInput!) {
			updateSchemaStatus(input: $input) {
				schema {
					name
				}
			}
		}
	`)

	statusInput := map[string]interface{}{
		"name":       input.Name,
		"isArchived": true,
	}

	req.Var("input", statusInput)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		UpdateSchemaStatus struct {
			Schema client.Schema `json:"schema"`
		} `json:"updateSchemaStatus"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return client.DeleteSchemaOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}

	return client.DeleteSchemaOutput{Name: response.UpdateSchemaStatus.Schema.Name}, nil
}

// GraphQL Role methods
func (c *GraphQLClient) CreateRole(ctx context.Context, input client.CreateRoleInput) (client.CreateRoleOutput, error) {
	req := graphql.NewRequest(`
		mutation($input: CreateRoleInput!) {
			createRole(input: $input) {
				role {
					id
					name
					permissions
				}
			}
		}
	`)

	req.Var("input", input)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		CreateRole client.CreateRoleOutput `json:"createRole"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return client.CreateRoleOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return response.CreateRole, nil
}

func (c *GraphQLClient) UpdateRole(ctx context.Context, input client.UpdateRoleInput) (client.UpdateRoleOutput, error) {
	req := graphql.NewRequest(`
		mutation($input: UpdateRoleInput!) {
			updateRole(input: $input) {
				role {
					id
					name
					permissions
				}
			}
		}
	`)

	req.Var("input", input)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		UpdateRole client.UpdateRoleOutput `json:"updateRole"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return client.UpdateRoleOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return response.UpdateRole, nil
}

func (c *GraphQLClient) GetRoleById(ctx context.Context, id string) (*client.Role, error) {
	req := graphql.NewRequest(`
		query($id: ID!) {
			roleById(id: $id) {
				id
				name
				permissions
			}
		}
	`)

	req.Var("id", id)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		RoleById client.Role `json:"roleById"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %w", err)
	}
	return &response.RoleById, nil
}

func (c *GraphQLClient) GetRoleByName(ctx context.Context, name string) (*client.Role, error) {
	req := graphql.NewRequest(`
		query($name: String!) {
			roleByName(name: $name) {
				id
				name
				permissions
			}
		}
	`)

	req.Var("name", name)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		RoleByName client.Role `json:"roleByName"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %w", err)
	}
	return &response.RoleByName, nil
}

func (c *GraphQLClient) DeleteRole(ctx context.Context, input client.DeleteRoleInput) (client.DeleteRoleOutput, error) {
	req := graphql.NewRequest(`
		mutation($input: DeleteRoleInput!) {
			deleteRole(input: $input) {
				__typename
			}
		}
	`)

	req.Var("input", input)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		DeleteRole struct {
			Typename string `json:"__typename"`
		} `json:"deleteRole"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return client.DeleteRoleOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return client.DeleteRoleOutput{}, nil
}

// GraphQL User methods
func (c *GraphQLClient) InviteUser(ctx context.Context, input client.InviteUserInput) (client.InviteUserOutput, error) {
	req := graphql.NewRequest(`
		mutation($input: InviteUserInput!) {
			inviteUser(input: $input) {
				user {
					id
					email
					familyName
					givenName
					status
					role {
						id
						name
					}
				}
			}
		}
	`)

	req.Var("input", input)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		InviteUser client.InviteUserOutput `json:"inviteUser"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return client.InviteUserOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return response.InviteUser, nil
}

func (c *GraphQLClient) UpdateUser(ctx context.Context, input client.UpdateUserInput) (client.UpdateUserOutput, error) {
	req := graphql.NewRequest(`
		mutation($input: UpdateUserInput!) {
			updateUser(input: $input) {
				user {
					id
					email
					familyName
					givenName
					status
					role {
						id
						name
					}
				}
			}
		}
	`)

	req.Var("input", input)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		UpdateUser client.UpdateUserOutput `json:"updateUser"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return client.UpdateUserOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return response.UpdateUser, nil
}

func (c *GraphQLClient) GetUserById(ctx context.Context, id string) (*client.User, error) {
	req := graphql.NewRequest(`
		query($id: ID!) {
			userById(id: $id) {
				id
				email
				familyName
				givenName
				status
				role {
					id
					name
				}
			}
		}
	`)

	req.Var("id", id)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		UserById client.User `json:"userById"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %w", err)
	}
	return &response.UserById, nil
}

func (c *GraphQLClient) GetUserByEmail(ctx context.Context, email string) (*client.User, error) {
	req := graphql.NewRequest(`
		query($email: String!) {
			userByEmail(email: $email) {
				id
				email
				familyName
				givenName
				status
				role {
					id
					name
				}
			}
		}
	`)

	req.Var("email", email)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		UserByEmail client.User `json:"userByEmail"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %w", err)
	}
	return &response.UserByEmail, nil
}

func (c *GraphQLClient) DeleteUser(ctx context.Context, input client.DeleteUserInput) (client.DeleteUserOutput, error) {
	req := graphql.NewRequest(`
		mutation($input: DeleteUserInput!) {
			deleteUser(input: $input) {
				__typename
			}
		}
	`)

	req.Var("input", input)
	req.Header.Set("X-API-Key", c.token)

	var response struct {
		DeleteUser struct {
			Typename string `json:"__typename"`
		} `json:"deleteUser"`
	}

	err := c.client.Run(ctx, req, &response)
	if err != nil {
		return client.DeleteUserOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return client.DeleteUserOutput{}, nil
}
