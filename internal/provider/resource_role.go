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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = (*roleResource)(nil)
	_ resource.ResourceWithConfigure = (*roleResource)(nil)
)

func NewRoleResource() resource.Resource {
	return &roleResource{}
}

type roleResource struct {
	client client.GraphQLClient
}

type roleResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Permissions types.List   `tfsdk:"permissions"`
}

func (r *roleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (r *roleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Panther role for user permissions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Role identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the role",
			},
			"permissions": schema.ListAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "List of permissions assigned to the role",
			},
		},
	}
}

func (r *roleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data roleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var permissions []string
	data.Permissions.ElementsAs(ctx, &permissions, false)

	input := client.CreateRoleInput{
		Name:        data.Name.ValueString(),
		Permissions: permissions,
	}

	output, err := r.client.CreateRole(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Role",
			"Could not create role, unexpected error: "+err.Error(),
		)
		return
	}

	role := output.Role

	tflog.Debug(ctx, "Created Role", map[string]any{
		"id": role.ID,
	})

	// Update the data with the response from the API
	data.ID = types.StringValue(role.ID)  
	data.Name = types.StringValue(role.Name)
	
	// Convert permissions back to list
	if len(role.Permissions) > 0 {
		elements := make([]types.String, len(role.Permissions))
		for i, permission := range role.Permissions {
			elements[i] = types.StringValue(permission)
		}
		permissionsList, diags := types.ListValueFrom(ctx, types.StringType, elements)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		data.Permissions = permissionsList
	} else {
		data.Permissions = types.ListNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data roleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	role, err := r.client.GetRoleById(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Role",
			fmt.Sprintf("Could not read role with id %s, unexpected error: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	tflog.Debug(ctx, "Got Role", map[string]any{
		"id": role.ID,
	})

	data.ID = types.StringValue(role.ID)
	data.Name = types.StringValue(role.Name)

	// Convert permissions back to list
	if len(role.Permissions) > 0 {
		elements := make([]types.String, len(role.Permissions))
		for i, permission := range role.Permissions {
			elements[i] = types.StringValue(permission)
		}
		permissionsList, diags := types.ListValueFrom(ctx, types.StringType, elements)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		data.Permissions = permissionsList
	} else {
		data.Permissions = types.ListNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data roleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var permissions []string
	data.Permissions.ElementsAs(ctx, &permissions, false)

	input := client.UpdateRoleInput{
		ID:          data.ID.ValueString(),
		Name:        data.Name.ValueString(),
		Permissions: permissions,
	}

	output, err := r.client.UpdateRole(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Role",
			fmt.Sprintf("Could not update role with id %s, unexpected error: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	role := output.Role

	tflog.Debug(ctx, "Updated Role", map[string]any{
		"id": role.ID,
	})

	// Update the data with the response from the API
	data.ID = types.StringValue(role.ID)  
	data.Name = types.StringValue(role.Name)
	
	// Convert permissions back to list
	if len(role.Permissions) > 0 {
		elements := make([]types.String, len(role.Permissions))
		for i, permission := range role.Permissions {
			elements[i] = types.StringValue(permission)
		}
		permissionsList, diags := types.ListValueFrom(ctx, types.StringType, elements)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		data.Permissions = permissionsList
	} else {
		data.Permissions = types.ListNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data roleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := client.DeleteRoleInput{
		ID: data.ID.ValueString(),
	}

	_, err := r.client.DeleteRole(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Role",
			fmt.Sprintf("Could not delete role with id %s, unexpected error: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	tflog.Debug(ctx, "Deleted Role", map[string]any{
		"id": data.ID.ValueString(),
	})
}

func (r *roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}