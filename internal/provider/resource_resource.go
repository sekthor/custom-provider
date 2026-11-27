package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/client"
)

var _ resource.ResourceWithConfigure = &ResourceResource{}
var _ resource.ResourceWithImportState = &ResourceResource{}

func NewResourceResource() resource.Resource {
	return &ResourceResource{}
}

type ResourceResource struct {
	client client.Client
}

type ResourceResourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (r *ResourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource"
}

func (r *ResourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "resource resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "id of the resource",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "name of the resource",
				Optional:            true,
			},
		},
	}
}

func (r *ResourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ResourceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resource := client.Resource{
		Id:   data.Id.ValueString(),
		Name: data.Name.ValueString(),
	}

	resource, err := r.client.Create(ctx, resource)
	if err != nil {
		resp.Diagnostics.AddError("could not create resource", err.Error())
		return
	}

	data.Id = types.StringValue(resource.Id)
	data.Name = types.StringValue(resource.Name)

	tflog.Trace(ctx, "created a resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ResourceDataSourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := r.client.GetResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	data.Id = types.StringValue(resource.Id)
	data.Name = types.StringValue(resource.Name)

	tflog.Trace(ctx, "read custom_resource resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ResourceDataSourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resource := client.Resource{
		Id:   data.Id.ValueString(),
		Name: data.Name.ValueString(),
	}

	r.client.Update(ctx, resource)

	data.Id = types.StringValue(resource.Id)
	data.Name = types.StringValue(resource.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ResourceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("could not delete resource", err.Error())
	}
}

func (r *ResourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
