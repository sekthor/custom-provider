// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/client"
)

var _ provider.Provider = &CustomProvider{}
var _ provider.ProviderWithFunctions = &CustomProvider{}
var _ provider.ProviderWithEphemeralResources = &CustomProvider{}
var _ provider.ProviderWithActions = &CustomProvider{}

type CustomProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type CustomProviderModel struct {
	Directory types.String `tfsdk:"directory"`
}

func (p *CustomProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "custom"
	resp.Version = p.version
}

func (p *CustomProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"directory": schema.StringAttribute{
				MarkdownDescription: "the data directory for persistence",
				Required:            true,
			},
		},
	}
}

func (p *CustomProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CustomProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Directory.IsNull() {
		resp.Diagnostics.AddError("could not create filesystem client", "directory must not be null")
		return
	}

	dir := data.Directory.ValueString()
	client, err := client.NewFilesystemClient(dir)

	if err != nil {
		resp.Diagnostics.AddError("could not create filesystem client", err.Error())
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *CustomProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewResourceResource,
	}
}

func (p *CustomProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *CustomProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewResourceDataSource,
	}
}

func (p *CustomProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *CustomProvider) Actions(ctx context.Context) []func() action.Action {
	return []func() action.Action{
		//NewExampleAction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CustomProvider{
			version: version,
		}
	}
}
