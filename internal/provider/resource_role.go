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
	"sort"
	"strings"
	"terraform-provider-panther/internal/client"
	"terraform-provider-panther/internal/client/panther"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = (*roleResource)(nil)
	_ resource.ResourceWithConfigure = (*roleResource)(nil)
	_ resource.ResourceWithImportState = (*roleResource)(nil)
)

// NewRoleResource creates a new role resource.
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

// Metadata returns the resource type name.
func (r *roleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Schema defines the schema for the role resource.
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
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure configures the resource with the provider client.
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

// filterRedundantReadPermissions removes *Read permissions when corresponding *Modify permissions exist.
// The Panther API automatically grants Read access with Modify permissions.
func filterRedundantReadPermissions(permissions []string) []string {
	permissionSet := make(map[string]bool)
	for _, perm := range permissions {
		permissionSet[perm] = true
	}
	
	var filtered []string
	for _, perm := range permissions {
		if strings.HasSuffix(perm, "Read") {
			baseName := strings.TrimSuffix(perm, "Read")
			modifyPerm := baseName + "Modify"
			if permissionSet[modifyPerm] {
				continue
			}
		}
		filtered = append(filtered, perm)
	}
	
	return filtered
}

// restoreOriginalPermissions adds back *Read permissions from the original configuration
// that were filtered out by the API due to corresponding *Modify permissions.
func restoreOriginalPermissions(apiPermissions []string, originalPermissions []string) []string {
	apiSet := make(map[string]bool)
	for _, perm := range apiPermissions {
		apiSet[perm] = true
	}
	
	originalSet := make(map[string]bool)
	for _, perm := range originalPermissions {
		originalSet[perm] = true
	}
	
	result := make([]string, len(apiPermissions))
	copy(result, apiPermissions)
	
	for _, perm := range originalPermissions {
		if strings.HasSuffix(perm, "Read") && !apiSet[perm] {
			baseName := strings.TrimSuffix(perm, "Read")
			modifyPerm := baseName + "Modify"
			if apiSet[modifyPerm] && originalSet[perm] {
				result = append(result, perm)
			}
		}
	}
	
	return result
}

// Create creates a new role resource.
func (r *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data roleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var permissions []string
	data.Permissions.ElementsAs(ctx, &permissions, false)
	
	filteredPermissions := filterRedundantReadPermissions(permissions)

	input := client.CreateRoleInput{
		Name:        data.Name.ValueString(),
		Permissions: filteredPermissions,
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

	data.ID = types.StringValue(role.ID)
	data.Name = types.StringValue(role.Name)
	
	// Store original configuration permissions in state to maintain consistency
	if len(permissions) > 0 {
		sortedPermissions := make([]string, len(permissions))
		copy(sortedPermissions, permissions)
		sort.Strings(sortedPermissions)
		
		elements := make([]types.String, len(sortedPermissions))
		for i, permission := range sortedPermissions {
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

// Read reads the role resource.
func (r *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data roleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract current permissions from state to preserve order
	var currentPermissions []string
	if !data.Permissions.IsNull() && !data.Permissions.IsUnknown() {
		data.Permissions.ElementsAs(ctx, &currentPermissions, false)
	}

	role, err := r.client.GetRoleById(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Role",
			fmt.Sprintf("Could not read role with id %s, unexpected error: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}


	data.ID = types.StringValue(role.ID)
	data.Name = types.StringValue(role.Name)

	// Restore original permissions from configuration to maintain state consistency
	restoredPermissions := restoreOriginalPermissions(role.Permissions, currentPermissions)
	
	if len(restoredPermissions) > 0 {
		sort.Strings(restoredPermissions)
		
		elements := make([]types.String, len(restoredPermissions))
		for i, permission := range restoredPermissions {
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

// Update updates the role resource.
func (r *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data roleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var permissions []string
	data.Permissions.ElementsAs(ctx, &permissions, false)
	
	filteredPermissions := filterRedundantReadPermissions(permissions)

	input := client.UpdateRoleInput{
		ID:          data.ID.ValueString(),
		Name:        data.Name.ValueString(),
		Permissions: filteredPermissions,
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

	data.ID = types.StringValue(role.ID)
	data.Name = types.StringValue(role.Name)
	
	// Store original configuration permissions in state to maintain consistency
	if len(permissions) > 0 {
		sortedPermissions := make([]string, len(permissions))
		copy(sortedPermissions, permissions)
		sort.Strings(sortedPermissions)
		
		elements := make([]types.String, len(sortedPermissions))
		for i, permission := range sortedPermissions {
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

// Delete deletes the role resource.
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

}

// ImportState imports an existing role resource.
func (r *roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
