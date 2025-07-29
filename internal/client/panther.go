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

package client

import (
	"context"
)

type GraphQLClient interface {
	CreateS3Source(ctx context.Context, input CreateS3SourceInput) (CreateS3SourceOutput, error)
	UpdateS3Source(ctx context.Context, input UpdateS3SourceInput) (UpdateS3SourceOutput, error)
	GetS3Source(ctx context.Context, id string) (*S3LogIntegration, error)
	DeleteSource(ctx context.Context, input DeleteSourceInput) (DeleteSourceOutput, error)
	
	// Cloud Account management
	CreateCloudAccount(ctx context.Context, input CreateCloudAccountInput) (CreateCloudAccountOutput, error)  
	UpdateCloudAccount(ctx context.Context, input UpdateCloudAccountInput) (UpdateCloudAccountOutput, error)
	GetCloudAccount(ctx context.Context, id string) (*CloudAccount, error)
	DeleteCloudAccount(ctx context.Context, input DeleteCloudAccountInput) (DeleteCloudAccountOutput, error)
	
	// Schema management
	CreateSchema(ctx context.Context, input CreateSchemaInput) (CreateSchemaOutput, error)
	UpdateSchema(ctx context.Context, input UpdateSchemaInput) (UpdateSchemaOutput, error)
	GetSchema(ctx context.Context, id string) (*Schema, error)
	DeleteSchema(ctx context.Context, input DeleteSchemaInput) (DeleteSchemaOutput, error)
	
	// Role management (GraphQL)
	CreateRole(ctx context.Context, input CreateRoleInput) (CreateRoleOutput, error)
	UpdateRole(ctx context.Context, input UpdateRoleInput) (UpdateRoleOutput, error)
	GetRoleById(ctx context.Context, id string) (*Role, error)
	GetRoleByName(ctx context.Context, name string) (*Role, error)
	DeleteRole(ctx context.Context, input DeleteRoleInput) (DeleteRoleOutput, error)
	
	// User management (GraphQL)
	InviteUser(ctx context.Context, input InviteUserInput) (InviteUserOutput, error)
	UpdateUser(ctx context.Context, input UpdateUserInput) (UpdateUserOutput, error)
	GetUserById(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error) 
	DeleteUser(ctx context.Context, input DeleteUserInput) (DeleteUserOutput, error)
}

type RestClient interface {
	CreateHttpSource(ctx context.Context, input CreateHttpSourceInput) (HttpSource, error)
	UpdateHttpSource(ctx context.Context, input UpdateHttpSourceInput) (HttpSource, error)
	GetHttpSource(ctx context.Context, id string) (HttpSource, error)
	DeleteHttpSource(ctx context.Context, id string) error
	
	// Rules management
	CreateRule(ctx context.Context, input CreateRuleInput) (Rule, error)
	UpdateRule(ctx context.Context, input UpdateRuleInput) (Rule, error)
	GetRule(ctx context.Context, id string) (Rule, error)
	DeleteRule(ctx context.Context, id string) error
	
	
	// Data Models management
	CreateDataModel(ctx context.Context, input CreateDataModelInput) (DataModel, error)
	UpdateDataModel(ctx context.Context, input UpdateDataModelInput) (DataModel, error)
	GetDataModel(ctx context.Context, id string) (DataModel, error)
	DeleteDataModel(ctx context.Context, id string) error
	
	// Globals management
	CreateGlobal(ctx context.Context, input CreateGlobalInput) (Global, error)
	UpdateGlobal(ctx context.Context, input UpdateGlobalInput) (Global, error)
	GetGlobal(ctx context.Context, id string) (Global, error)
	DeleteGlobal(ctx context.Context, id string) error
}

// CreateS3SourceInput Input for the createS3LogSource mutation
type CreateS3SourceInput struct {
	AwsAccountID               string                  `json:"awsAccountId"`
	KmsKey                     string                  `json:"kmsKey"`
	Label                      string                  `json:"label"`
	LogProcessingRole          string                  `json:"logProcessingRole"`
	LogStreamType              string                  `json:"logStreamType"`
	ManagedBucketNotifications bool                    `json:"managedBucketNotifications"`
	S3Bucket                   string                  `json:"s3Bucket"`
	S3PrefixLogTypes           []S3PrefixLogTypesInput `json:"s3PrefixLogTypes"`
}

// CreateS3SourceOutput output for the createS3LogSource mutation
type CreateS3SourceOutput struct {
	LogSource *S3LogIntegration `graphql:"logSource"`
}

// UpdateS3SourceInput input for the updateS3Source mutation
type UpdateS3SourceInput struct {
	ID                         string                  `json:"id"`
	KmsKey                     string                  `json:"kmsKey"`
	Label                      string                  `json:"label"`
	LogProcessingRole          string                  `json:"logProcessingRole"`
	LogStreamType              string                  `json:"logStreamType"`
	ManagedBucketNotifications bool                    `json:"managedBucketNotifications"`
	S3PrefixLogTypes           []S3PrefixLogTypesInput `json:"s3PrefixLogTypes"`
}

// UpdateS3SourceOutput output for the updateS3LogSource mutation
type UpdateS3SourceOutput struct {
	LogSource *S3LogIntegration `graphql:"logSource"`
}

// DeleteSourceInput input for the deleteSource mutation
type DeleteSourceInput struct {
	ID string `json:"id"`
}

// DeleteSourceOutput output for the deleteSource mutation
type DeleteSourceOutput struct {
	ID string `json:"id"`
}

// S3LogIntegration Represents an S3 Log Source Integration
type S3LogIntegration struct {
	// The ID of the AWS Account where the S3 Bucket is located
	AwsAccountID string `graphql:"awsAccountId"`
	// The ID of the Log Source integration
	IntegrationID string `graphql:"integrationId"`
	// The name of the Log Source integration
	IntegrationLabel string `graphql:"integrationLabel"`
	// The type of Log Source integration
	IntegrationType string `graphql:"integrationType"`
	// True if the Log Source can be modified
	IsEditable bool `graphql:"isEditable"`
	// KMS key used to access the S3 Bucket
	KmsKey string `graphql:"kmsKey"`
	// The AWS Role used to access the S3 Bucket
	LogProcessingRole *string `graphql:"logProcessingRole"`
	// The format of the log files being ingested
	LogStreamType *string `graphql:"logStreamType"`
	// True if bucket notifications are being managed by Panther
	ManagedBucketNotifications bool `json:"managedBucketNotifications"`
	// The S3 Bucket name being ingested
	S3Bucket string `graphql:"s3Bucket"`
	// The prefix on the S3 Bucket name being ingested
	S3Prefix *string `graphql:"s3Prefix"`
	// Used to map prefixes to log types
	S3PrefixLogTypes []S3PrefixLogTypes `graphql:"s3PrefixLogTypes"`
}

// S3PrefixLogTypesInput Mapping of S3 prefixes to log types
type S3PrefixLogTypesInput struct {
	// S3 Prefixes to exclude
	ExcludedPrefixes []string `json:"excludedPrefixes"`
	// Log types to map to prefix
	LogTypes []string `json:"logTypes"`
	// S3 Prefix to map to log types
	Prefix string `json:"prefix"`
}

type S3PrefixLogTypes struct {
	// S3 Prefixes to exclude
	ExcludedPrefixes []string `graphql:"excludedPrefixes"`
	// Log types to map to prefix
	LogTypes []string `graphql:"logTypes"`
	// S3 Prefix to map to log types
	Prefix string `graphql:"prefix"`
}

type HttpSource struct {
	IntegrationId string
	HttpSourceModifiableAttributes
}

// LogStreamTypeOptions contains options specific to the log stream type
type LogStreamTypeOptions struct {
	JsonArrayEnvelopeField string
}

// HttpSourceModifiableAttributes attributes that can be modified on an http log source
type HttpSourceModifiableAttributes struct {
	IntegrationLabel     string
	LogStreamType        string
	LogTypes             []string
	LogStreamTypeOptions *LogStreamTypeOptions
	AuthHmacAlg          string
	AuthHeaderKey        string
	AuthPassword         string
	AuthSecretValue      string
	AuthMethod           string
	AuthUsername         string
	AuthBearerToken      string
}

// CreateHttpSourceInput Input for creating an http log source
type CreateHttpSourceInput struct {
	HttpSourceModifiableAttributes
}

// UpdateHttpSourceInput input for updating an http log source
type UpdateHttpSourceInput struct {
	IntegrationId string
	HttpSourceModifiableAttributes
}

type HttpErrorResponse struct {
	Message string
}

// Rule types
type Rule struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	RuleModifiableAttributes
}

type RuleModifiableAttributes struct {
	DisplayName         string   `json:"displayName"`
	Body                string   `json:"body"`
	Description         string   `json:"description,omitempty"`
	Severity            string   `json:"severity,omitempty"`
	LogTypes            []string `json:"logTypes,omitempty"`
	Tags                []string `json:"tags,omitempty"`
	References          []string `json:"references,omitempty"`
	Runbook             string   `json:"runbook,omitempty"`
	DedupPeriodMinutes  int      `json:"dedupPeriodMinutes,omitempty"`
	Enabled             bool     `json:"enabled,omitempty"`
}

type CreateRuleInput struct {
	ID string `json:"id"`
	RuleModifiableAttributes
}

type UpdateRuleInput struct {
	ID string `json:"id"`
	RuleModifiableAttributes
}



// Data Model types
type DataModel struct {
	ID       string `json:"id"`
	DataModelModifiableAttributes
}

type DataModelModifiableAttributes struct {
	LogType     string            `json:"logType"`
	Description string            `json:"description"`
	Enabled     bool              `json:"enabled"`
	Mappings    []DataModelMapping `json:"mappings"`
	Body        string            `json:"body"`
}

type DataModelMapping struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type CreateDataModelInput struct {
	DataModelModifiableAttributes
}

type UpdateDataModelInput struct {
	ID string `json:"id"`
	DataModelModifiableAttributes
}

// Global types
type Global struct {
	ID       string `json:"id"`
	GlobalModifiableAttributes
}

type GlobalModifiableAttributes struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	Tags        []string `json:"tags"`
}

type CreateGlobalInput struct {
	GlobalModifiableAttributes
}

type UpdateGlobalInput struct {
	ID string `json:"id"`
	GlobalModifiableAttributes
}

// GraphQL Cloud Account types
type AWSScanConfig struct {
	AuditRole string `graphql:"auditRole"`
}

type CloudAccount struct {
	ID                        string         `graphql:"id"`
	AWSAccountID              string         `graphql:"awsAccountId"`
	Label                     string         `graphql:"label"`
	AWSStackName              string         `graphql:"awsStackName"`
	AWSScanConfig             AWSScanConfig  `graphql:"awsScanConfig"`
	AWSRegionIgnoreList       []string       `graphql:"awsRegionIgnoreList"`
	ResourceRegexIgnoreList   []string       `graphql:"resourceRegexIgnoreList"`
	ResourceTypeIgnoreList    []string       `graphql:"resourceTypeIgnoreList"`
	IsEditable                bool           `graphql:"isEditable"`
}

type AWSScanConfigInput struct {
	AuditRole string `json:"auditRole"`
}

type CreateCloudAccountInput struct {
	AWSAccountID              string             `json:"awsAccountId"`
	Label                     string             `json:"label"`
	AWSScanConfig             AWSScanConfigInput `json:"awsScanConfig"`
	AWSRegionIgnoreList       []string           `json:"awsRegionIgnoreList,omitempty"`
	ResourceRegexIgnoreList   []string           `json:"resourceRegexIgnoreList,omitempty"`
	ResourceTypeIgnoreList    []string           `json:"resourceTypeIgnoreList,omitempty"`
}

type CreateCloudAccountOutput struct {
	CloudAccount *CloudAccount `graphql:"cloudAccount"`
}

type UpdateCloudAccountInput struct {
	ID                        string             `json:"id"`
	Label                     string             `json:"label"`
	AWSScanConfig             AWSScanConfigInput `json:"awsScanConfig"`
	AWSRegionIgnoreList       []string           `json:"awsRegionIgnoreList,omitempty"`
	ResourceRegexIgnoreList   []string           `json:"resourceRegexIgnoreList,omitempty"`
	ResourceTypeIgnoreList    []string           `json:"resourceTypeIgnoreList,omitempty"`
}

type UpdateCloudAccountOutput struct {
	CloudAccount *CloudAccount `graphql:"cloudAccount"`
}

type DeleteCloudAccountInput struct {
	ID string `json:"id"`
}

type DeleteCloudAccountOutput struct {
	ID string `json:"id"`
}

// GraphQL Schema types
type Schema struct {
	ID                      string   `graphql:"id"`
	Name                    string   `graphql:"name"`
	Description             string   `graphql:"description"`
	Spec                    string   `graphql:"spec"`
	Version                 int      `graphql:"version"`
	LogTypes                []string `graphql:"logTypes"`
	IsFieldDiscoveryEnabled bool     `graphql:"isFieldDiscoveryEnabled"`
}

type CreateSchemaInput struct {
	Name                    string   `json:"name"`
	Description             string   `json:"description"`
	Spec                    string   `json:"spec"`
	LogTypes                []string `json:"logTypes"`
	IsFieldDiscoveryEnabled bool     `json:"isFieldDiscoveryEnabled"`
}

type CreateSchemaOutput struct {
	Schema *Schema `graphql:"schema"`
}

type UpdateSchemaInput struct {
	ID                      string   `json:"id"`
	Name                    string   `json:"name"`
	Description             string   `json:"description"`
	Spec                    string   `json:"spec"`
	LogTypes                []string `json:"logTypes"`
	IsFieldDiscoveryEnabled bool     `json:"isFieldDiscoveryEnabled"`
}

type UpdateSchemaOutput struct {
	Schema *Schema `graphql:"schema"`
}

type DeleteSchemaInput struct {
	ID string `json:"id"`
}

type DeleteSchemaOutput struct {
	ID string `json:"id"`
}

// GraphQL Role types (matching actual API schema)
type Role struct {
	ID          string   `graphql:"id"`
	Name        string   `graphql:"name"`
	Permissions []string `graphql:"permissions"`
}

type CreateRoleInput struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type CreateRoleOutput struct {
	Role *Role `graphql:"role"`
}

type UpdateRoleInput struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type UpdateRoleOutput struct {
	Role *Role `graphql:"role"`
}

type DeleteRoleInput struct {
	ID string `json:"id"`
}

type DeleteRoleOutput struct {
	ID string `json:"id"`
}

// GraphQL User types (matching actual API schema)
type User struct {
	ID         string `graphql:"id"`
	Email      string `graphql:"email"`
	GivenName  string `graphql:"givenName"`
	FamilyName string `graphql:"familyName"`
	Status     string `graphql:"status"`
	CreatedAt  string `graphql:"createdAt"`
	Role       *Role  `graphql:"role"`
}

type UserRoleInputKind string

const (
	UserRoleInputKindID   UserRoleInputKind = "ID"
	UserRoleInputKindName UserRoleInputKind = "NAME"
)

type RoleInput struct {
	Kind  UserRoleInputKind `json:"kind"`
	Value string            `json:"value"`
}

type InviteUserInput struct {
	Email      string    `json:"email"`
	GivenName  string    `json:"givenName"`
	FamilyName string    `json:"familyName"`
	Role       RoleInput `json:"role"`
}

type InviteUserOutput struct {
	User *User `graphql:"user"`
}

type UpdateUserInput struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	GivenName  string    `json:"givenName"`  
	FamilyName string    `json:"familyName"`
	Role       RoleInput `json:"role"`
}

type UpdateUserOutput struct {
	User *User `graphql:"user"`
}

type DeleteUserInput struct {
	ID string `json:"id"`
}

type DeleteUserOutput struct {
	ID string `json:"id"`
}
