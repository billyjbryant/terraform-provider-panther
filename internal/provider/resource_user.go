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
	_ resource.Resource              = (*userResource)(nil)
	_ resource.ResourceWithConfigure = (*userResource)(nil)
	_ resource.ResourceWithImportState = (*userResource)(nil)
)

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	client client.GraphQLClient
}

type userResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Email      types.String `tfsdk:"email"`
	GivenName  types.String `tfsdk:"given_name"`
	FamilyName types.String `tfsdk:"family_name"`
	Role       types.String `tfsdk:"role"`
}

func (r *userResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Panther user account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "User identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Email address of the user",
			},
			"given_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Given name (first name) of the user",
			},
			"family_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Family name (last name) of the user",
			},
			"role": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Role name assigned to the user",
			},
		},
	}
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use inviteUser mutation as per API documentation
	// The API requires a role, so if none provided, we'll need to handle this
	var roleInput client.RoleInput
	if !data.Role.IsNull() && !data.Role.IsUnknown() && data.Role.ValueString() != "" {
		roleInput = client.RoleInput{
			Kind:  client.UserRoleInputKindName,
			Value: data.Role.ValueString(),
		}
	} else {
		// If no role specified, we can't create the user - the API requires it
		resp.Diagnostics.AddError(
			"Missing required field",
			"Role is required to create a user. Please specify a role.",
		)
		return
	}
	
	input := client.InviteUserInput{
		Email:      data.Email.ValueString(),
		GivenName:  data.GivenName.ValueString(), 
		FamilyName: data.FamilyName.ValueString(),
		Role:       roleInput,
	}

	output, err := r.client.InviteUser(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating User",
			"Could not create user, unexpected error: "+err.Error(),
		)
		return
	}

	user := output.User

	tflog.Debug(ctx, "Created User", map[string]any{
		"id": user.ID,
	})

	// Update the data with the response from the API
	data.ID = types.StringValue(user.ID)
	data.Email = types.StringValue(user.Email)
	data.GivenName = types.StringValue(user.GivenName)
	data.FamilyName = types.StringValue(user.FamilyName)
	
	// Set role from the API response if available
	if user.Role != nil {
		data.Role = types.StringValue(user.Role.Name)
	} else {
		data.Role = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.GetUserById(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading User",
			fmt.Sprintf("Could not read user with id %s, unexpected error: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	tflog.Debug(ctx, "Got User", map[string]any{
		"id": user.ID,
	})

	data.ID = types.StringValue(user.ID)
	data.Email = types.StringValue(user.Email)
	data.GivenName = types.StringValue(user.GivenName)
	data.FamilyName = types.StringValue(user.FamilyName)
	
	// Set role from the API response if available
	if user.Role != nil {
		data.Role = types.StringValue(user.Role.Name)
	} else {
		data.Role = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Role is now required, so it should always be present
	roleValue := data.Role.ValueString()
	if data.Role.IsNull() || data.Role.IsUnknown() || roleValue == "" {
		resp.Diagnostics.AddError(
			"Missing required field",
			"Role is required to update a user. Please specify a role.",
		)
		return
	}
	
	input := client.UpdateUserInput{
		ID:         data.ID.ValueString(),
		Email:      data.Email.ValueString(),
		GivenName:  data.GivenName.ValueString(),
		FamilyName: data.FamilyName.ValueString(),
		Role: client.RoleInput{
			Kind:  client.UserRoleInputKindName,  // Use role name for consistency
			Value: roleValue,
		},
	}

	output, err := r.client.UpdateUser(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating User",
			fmt.Sprintf("Could not update user with id %s, unexpected error: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	user := output.User

	tflog.Debug(ctx, "Updated User", map[string]any{
		"id": user.ID,
	})

	// Update the data with the response from the API
	data.ID = types.StringValue(user.ID)
	data.Email = types.StringValue(user.Email)
	data.GivenName = types.StringValue(user.GivenName)
	data.FamilyName = types.StringValue(user.FamilyName)
	
	// Set role from the API response if available
	if user.Role != nil {
		data.Role = types.StringValue(user.Role.Name)
	} else {
		data.Role = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := client.DeleteUserInput{
		ID: data.ID.ValueString(),
	}

	_, err := r.client.DeleteUser(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting User",
			fmt.Sprintf("Could not delete user with id %s, unexpected error: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	tflog.Debug(ctx, "Deleted User", map[string]any{
		"id": data.ID.ValueString(),
	})
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}