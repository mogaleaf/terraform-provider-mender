package terraform

import (
	"context"
	"crypto/md5"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceArtifact struct{}

func (r resourceArtifact) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"source_file": {
				Type:     types.StringType,
				Required: true,
			},
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"md5": {
				Type:     types.StringType,
				Computed: true,
			},
		},
	}, nil
}

func (r resourceArtifact) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceOrder{
		p: *(p.(*provider)),
	}, nil
}

type resourceOrder struct {
	p provider
}

// Create a new resource
func (r resourceOrder) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply",
		)
		return
	}
	var artifact Artifact
	diags := req.Plan.Get(ctx, &artifact)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if artifact.SourceFile.Null {
		resp.Diagnostics.AddError(
			"source file is mandatory",
			"you need to provide the path of the artifact",
		)
		return
	}

	if _, err := os.Stat(artifact.SourceFile.Value); err != nil {
		resp.Diagnostics.AddError(
			"source file does not exist",
			"the file does not exist",
		)
		return
	}

	dat, err := os.ReadFile(artifact.SourceFile.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"can't read file",
			"can't read file",
		)
		return
	}
	id, err := r.p.client.UploadArtifact(dat)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating artifact",
			"Could not create artifact, unexpected error: "+err.Error(),
		)
		return
	}
	md5 := fmt.Sprintf("%x", md5.Sum(dat))
	// Generate resource state struct
	var result = Artifact{
		ID:         types.String{Value: id},
		SourceFile: artifact.SourceFile,
		Md5:        types.String{Value: string(md5)},
	}
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read resource information
func (r resourceOrder) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Set state
	var artifact Artifact
	diags := req.State.Get(ctx, &artifact)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if _, err := os.Stat(artifact.SourceFile.Value); err != nil {
		resp.Diagnostics.AddError(
			"source file does not exist",
			"the file does not exist",
		)
		return
	}

	dat, err := os.ReadFile(artifact.SourceFile.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"can't read file",
			"can't read file",
		)
		return
	}
	md5 := fmt.Sprintf("%x", md5.Sum(dat))
	artifact.Md5 = types.String{Value: md5}

	diags = resp.State.Set(ctx, &artifact)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceOrder) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
}

// Delete resource
func (r resourceOrder) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
}
