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
	"github.com/hasura/go-graphql-client"
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
	*graphql.Client
}

type RestClient struct {
	baseUrl string
	Doer
}

func NewGraphQLClient(url, token string) *GraphQLClient {
	return &GraphQLClient{
		graphql.NewClient(
			fmt.Sprintf("%s%s", url, GraphqlPath),
			NewAuthorizedHTTPClient(token)),
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
	var m struct {
		UpdateS3Source struct {
			client.UpdateS3SourceOutput
		} `graphql:"updateS3Source(input: $input)"`
	}
	err := c.Mutate(ctx, &m, map[string]interface{}{
		"input": input,
	}, graphql.OperationName("UpdateS3Source"))
	if err != nil {
		return client.UpdateS3SourceOutput{}, fmt.Errorf("GraphQL mutation failed: %v", err)
	}
	return m.UpdateS3Source.UpdateS3SourceOutput, nil
}

func (c *GraphQLClient) DeleteSource(ctx context.Context, input client.DeleteSourceInput) (client.DeleteSourceOutput, error) {
	var m struct {
		DeleteSource struct {
			client.DeleteSourceOutput
		} `graphql:"deleteSource(input: $input)"`
	}
	err := c.Mutate(ctx, &m, map[string]interface{}{
		"input": input,
	}, graphql.OperationName("DeleteSource"))
	if err != nil {
		return client.DeleteSourceOutput{}, fmt.Errorf("GraphQL mutation failed: %v", err)
	}
	return m.DeleteSource.DeleteSourceOutput, nil
}

func (c *GraphQLClient) GetS3Source(ctx context.Context, id string) (*client.S3LogIntegration, error) {
	var q struct {
		Source struct {
			S3LogIntegration client.S3LogIntegration `graphql:"... on S3LogIntegration"`
		} `graphql:"source(id: $id)"`
	}

	err := c.Query(ctx, &q, map[string]interface{}{
		"id": graphql.ID(id),
	}, graphql.OperationName("Source"))
	if err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %v", err)
	}
	return &q.Source.S3LogIntegration, nil
}

func (c *GraphQLClient) CreateS3Source(ctx context.Context, input client.CreateS3SourceInput) (client.CreateS3SourceOutput, error) {
	var m struct {
		CreateS3Source struct {
			client.CreateS3SourceOutput
		} `graphql:"createS3Source(input: $input)"`
	}
	err := c.Mutate(ctx, &m, map[string]any{
		"input": input,
	}, graphql.OperationName("CreateS3Source"))
	if err != nil {
		return client.CreateS3SourceOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return m.CreateS3Source.CreateS3SourceOutput, nil
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
	var m struct {
		CreateCloudAccount struct {
			client.CreateCloudAccountOutput
		} `graphql:"createCloudAccount(input: $input)"`
	}
	err := c.Mutate(ctx, &m, map[string]interface{}{
		"input": input,
	}, graphql.OperationName("CreateCloudAccount"))
	if err != nil {
		return client.CreateCloudAccountOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return m.CreateCloudAccount.CreateCloudAccountOutput, nil
}

func (c *GraphQLClient) UpdateCloudAccount(ctx context.Context, input client.UpdateCloudAccountInput) (client.UpdateCloudAccountOutput, error) {
	var m struct {
		UpdateCloudAccount struct {
			client.UpdateCloudAccountOutput
		} `graphql:"updateCloudAccount(input: $input)"`
	}
	err := c.Mutate(ctx, &m, map[string]interface{}{
		"input": input,
	}, graphql.OperationName("UpdateCloudAccount"))
	if err != nil {
		return client.UpdateCloudAccountOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return m.UpdateCloudAccount.UpdateCloudAccountOutput, nil
}

func (c *GraphQLClient) GetCloudAccount(ctx context.Context, id string) (*client.CloudAccount, error) {
	var q struct {
		CloudAccount client.CloudAccount `graphql:"cloudAccount(id: $id)"`
	}

	err := c.Query(ctx, &q, map[string]interface{}{
		"id": graphql.ID(id),
	}, graphql.OperationName("CloudAccount"))
	if err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %w", err)
	}
	return &q.CloudAccount, nil
}

func (c *GraphQLClient) DeleteCloudAccount(ctx context.Context, input client.DeleteCloudAccountInput) (client.DeleteCloudAccountOutput, error) {
	var m struct {
		DeleteCloudAccount struct {
			client.DeleteCloudAccountOutput
		} `graphql:"deleteCloudAccount(input: $input)"`
	}
	err := c.Mutate(ctx, &m, map[string]interface{}{
		"input": input,
	}, graphql.OperationName("DeleteCloudAccount"))
	if err != nil {
		return client.DeleteCloudAccountOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return m.DeleteCloudAccount.DeleteCloudAccountOutput, nil
}

// GraphQL Client methods for Schemas
func (c *GraphQLClient) CreateSchema(ctx context.Context, input client.CreateSchemaInput) (client.CreateSchemaOutput, error) {
	var m struct {
		CreateSchema struct {
			client.CreateSchemaOutput
		} `graphql:"createSchema(input: $input)"`
	}
	err := c.Mutate(ctx, &m, map[string]interface{}{
		"input": input,
	}, graphql.OperationName("CreateSchema"))
	if err != nil {
		return client.CreateSchemaOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return m.CreateSchema.CreateSchemaOutput, nil
}

func (c *GraphQLClient) UpdateSchema(ctx context.Context, input client.UpdateSchemaInput) (client.UpdateSchemaOutput, error) {
	var m struct {
		UpdateSchema struct {
			client.UpdateSchemaOutput
		} `graphql:"updateSchema(input: $input)"`
	}
	err := c.Mutate(ctx, &m, map[string]interface{}{
		"input": input,
	}, graphql.OperationName("UpdateSchema"))
	if err != nil {
		return client.UpdateSchemaOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return m.UpdateSchema.UpdateSchemaOutput, nil
}

func (c *GraphQLClient) GetSchema(ctx context.Context, id string) (*client.Schema, error) {
	var q struct {
		Schema client.Schema `graphql:"schema(id: $id)"`
	}

	err := c.Query(ctx, &q, map[string]interface{}{
		"id": graphql.ID(id),
	}, graphql.OperationName("Schema"))
	if err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %w", err)
	}
	return &q.Schema, nil
}

func (c *GraphQLClient) DeleteSchema(ctx context.Context, input client.DeleteSchemaInput) (client.DeleteSchemaOutput, error) {
	var m struct {
		DeleteSchema struct {
			client.DeleteSchemaOutput
		} `graphql:"deleteSchema(input: $input)"`
	}
	err := c.Mutate(ctx, &m, map[string]interface{}{
		"input": input,
	}, graphql.OperationName("DeleteSchema"))
	if err != nil {
		return client.DeleteSchemaOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return m.DeleteSchema.DeleteSchemaOutput, nil
}

// GraphQL Role methods
func (c *GraphQLClient) CreateRole(ctx context.Context, input client.CreateRoleInput) (client.CreateRoleOutput, error) {
	var m struct {
		CreateRole struct {
			client.CreateRoleOutput
		} `graphql:"createRole(input: $input)"`
	}
	err := c.Mutate(ctx, &m, map[string]interface{}{
		"input": input,
	}, graphql.OperationName("CreateRole"))
	if err != nil {
		return client.CreateRoleOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return m.CreateRole.CreateRoleOutput, nil
}

func (c *GraphQLClient) UpdateRole(ctx context.Context, input client.UpdateRoleInput) (client.UpdateRoleOutput, error) {
	var m struct {
		UpdateRole struct {
			client.UpdateRoleOutput
		} `graphql:"updateRole(input: $input)"`
	}
	err := c.Mutate(ctx, &m, map[string]interface{}{
		"input": input,
	}, graphql.OperationName("UpdateRole"))
	if err != nil {
		return client.UpdateRoleOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return m.UpdateRole.UpdateRoleOutput, nil
}

func (c *GraphQLClient) GetRoleById(ctx context.Context, id string) (*client.Role, error) {
	var q struct {
		RoleById client.Role `graphql:"roleById(id: $id)"`
	}

	err := c.Query(ctx, &q, map[string]interface{}{
		"id": graphql.ID(id),
	}, graphql.OperationName("RoleById"))
	if err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %w", err)
	}
	return &q.RoleById, nil
}

func (c *GraphQLClient) GetRoleByName(ctx context.Context, name string) (*client.Role, error) {
	var q struct {
		RoleByName client.Role `graphql:"roleByName(name: $name)"`
	}

	err := c.Query(ctx, &q, map[string]interface{}{
		"name": name,
	}, graphql.OperationName("RoleByName"))
	if err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %w", err)
	}
	return &q.RoleByName, nil
}

func (c *GraphQLClient) DeleteRole(ctx context.Context, input client.DeleteRoleInput) (client.DeleteRoleOutput, error) {
	var m struct {
		DeleteRole struct {
			client.DeleteRoleOutput
		} `graphql:"deleteRole(input: $input)"`
	}
	err := c.Mutate(ctx, &m, map[string]interface{}{
		"input": input,
	}, graphql.OperationName("DeleteRole"))
	if err != nil {
		return client.DeleteRoleOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return m.DeleteRole.DeleteRoleOutput, nil
}

// GraphQL User methods
func (c *GraphQLClient) InviteUser(ctx context.Context, input client.InviteUserInput) (client.InviteUserOutput, error) {
	var m struct {
		InviteUser struct {
			client.InviteUserOutput
		} `graphql:"inviteUser(input: $input)"`
	}
	err := c.Mutate(ctx, &m, map[string]interface{}{
		"input": input,
	}, graphql.OperationName("InviteUser"))
	if err != nil {
		return client.InviteUserOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return m.InviteUser.InviteUserOutput, nil
}

func (c *GraphQLClient) UpdateUser(ctx context.Context, input client.UpdateUserInput) (client.UpdateUserOutput, error) {
	var m struct {
		UpdateUser struct {
			client.UpdateUserOutput
		} `graphql:"updateUser(input: $input)"`
	}
	err := c.Mutate(ctx, &m, map[string]interface{}{
		"input": input,
	}, graphql.OperationName("UpdateUser"))
	if err != nil {
		return client.UpdateUserOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return m.UpdateUser.UpdateUserOutput, nil
}

func (c *GraphQLClient) GetUserById(ctx context.Context, id string) (*client.User, error) {
	var q struct {
		UserById client.User `graphql:"userById(id: $id)"`
	}

	err := c.Query(ctx, &q, map[string]interface{}{
		"id": graphql.ID(id),
	}, graphql.OperationName("UserById"))
	if err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %w", err)
	}
	return &q.UserById, nil
}

func (c *GraphQLClient) GetUserByEmail(ctx context.Context, email string) (*client.User, error) {
	var q struct {
		UserByEmail client.User `graphql:"userByEmail(email: $email)"`
	}

	err := c.Query(ctx, &q, map[string]interface{}{
		"email": email,
	}, graphql.OperationName("UserByEmail"))
	if err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %w", err)
	}
	return &q.UserByEmail, nil
}

func (c *GraphQLClient) DeleteUser(ctx context.Context, input client.DeleteUserInput) (client.DeleteUserOutput, error) {
	var m struct {
		DeleteUser struct {
			client.DeleteUserOutput
		} `graphql:"deleteUser(input: $input)"`
	}
	err := c.Mutate(ctx, &m, map[string]interface{}{
		"input": input,
	}, graphql.OperationName("DeleteUser"))
	if err != nil {
		return client.DeleteUserOutput{}, fmt.Errorf("GraphQL mutation failed: %w", err)
	}
	return m.DeleteUser.DeleteUserOutput, nil
}
