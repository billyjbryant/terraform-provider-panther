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

package provider

import (
	"context"
	"fmt"
	"terraform-provider-panther/internal/client"
	"terraform-provider-panther/internal/client/panther"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = (*schemaResource)(nil)
	_ resource.ResourceWithConfigure = (*schemaResource)(nil)
	_ resource.ResourceWithImportState = (*schemaResource)(nil)
)

func NewSchemaResource() resource.Resource {
	return &schemaResource{}
}

type schemaResource struct {
	client client.GraphQLClient
}

type schemaResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	Spec                    types.String `tfsdk:"spec"`
	Version                 types.Int64  `tfsdk:"version"`
	LogTypes                types.List   `tfsdk:"log_types"`
	IsFieldDiscoveryEnabled types.Bool   `tfsdk:"is_field_discovery_enabled"`
}

func (r *schemaResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

func (r *schemaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Panther log schema for custom log types.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Schema identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the schema (e.g., 'Custom.MyLog')",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Description of the schema",
			},
			"spec": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "YAML specification of the schema fields",
			},
			"version": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Schema version number",
			},
			"log_types": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "List of log types associated with this schema",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"is_field_discovery_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether field discovery is enabled for this schema",
			},
		},
	}
}

func (r *schemaResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*panther.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *panther.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c.GraphQLClient
}

func (r *schemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schemaResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	logTypesSlice := make([]string, 0, len(data.LogTypes.Elements()))
	resp.Diagnostics.Append(data.LogTypes.ElementsAs(ctx, &logTypesSlice, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := client.CreateSchemaInput{
		Name:                    data.Name.ValueString(),
		Description:             data.Description.ValueString(),
		Spec:                    data.Spec.ValueString(),
		LogTypes:                logTypesSlice,
		IsFieldDiscoveryEnabled: data.IsFieldDiscoveryEnabled.ValueBool(),
	}

	output, err := r.client.CreateSchema(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create schema, got error: %s", err))
		return
	}

	if output.Schema == nil {
		resp.Diagnostics.AddError("Client Error", "Create schema response was nil")
		return
	}

	data.ID = types.StringValue(output.Schema.ID)
	data.Name = types.StringValue(output.Schema.Name)
	data.Description = types.StringValue(output.Schema.Description)
	data.Spec = types.StringValue(output.Schema.Spec)
	data.Version = types.Int64Value(int64(output.Schema.Version))
	data.IsFieldDiscoveryEnabled = types.BoolValue(output.Schema.IsFieldDiscoveryEnabled)

	logTypesList, diags := types.ListValueFrom(ctx, types.StringType, output.Schema.LogTypes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.LogTypes = logTypesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *schemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schemaResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	schema, err := r.client.GetSchema(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read schema, got error: %s", err))
		return
	}

	if schema == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.ID = types.StringValue(schema.ID)
	data.Name = types.StringValue(schema.Name)
	data.Description = types.StringValue(schema.Description)
	data.Spec = types.StringValue(schema.Spec)
	data.Version = types.Int64Value(int64(schema.Version))
	data.IsFieldDiscoveryEnabled = types.BoolValue(schema.IsFieldDiscoveryEnabled)

	logTypesList, diags := types.ListValueFrom(ctx, types.StringType, schema.LogTypes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.LogTypes = logTypesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *schemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schemaResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	logTypesSlice := make([]string, 0, len(data.LogTypes.Elements()))
	resp.Diagnostics.Append(data.LogTypes.ElementsAs(ctx, &logTypesSlice, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := client.UpdateSchemaInput{
		ID:                      data.ID.ValueString(),
		Name:                    data.Name.ValueString(),
		Description:             data.Description.ValueString(),
		Spec:                    data.Spec.ValueString(),
		LogTypes:                logTypesSlice,
		IsFieldDiscoveryEnabled: data.IsFieldDiscoveryEnabled.ValueBool(),
	}

	output, err := r.client.UpdateSchema(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update schema, got error: %s", err))
		return
	}

	if output.Schema == nil {
		resp.Diagnostics.AddError("Client Error", "Update schema response was nil")
		return
	}

	data.ID = types.StringValue(output.Schema.ID)
	data.Name = types.StringValue(output.Schema.Name)
	data.Description = types.StringValue(output.Schema.Description)
	data.Spec = types.StringValue(output.Schema.Spec)
	data.Version = types.Int64Value(int64(output.Schema.Version))
	data.IsFieldDiscoveryEnabled = types.BoolValue(output.Schema.IsFieldDiscoveryEnabled)

	logTypesList, diags := types.ListValueFrom(ctx, types.StringType, output.Schema.LogTypes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.LogTypes = logTypesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *schemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemaResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := client.DeleteSchemaInput{
		ID: data.ID.ValueString(),
	}

	_, err := r.client.DeleteSchema(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete schema, got error: %s", err))
		return
	}
}

func (r *schemaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}