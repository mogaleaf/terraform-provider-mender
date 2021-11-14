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
type md5Modifier struct{}
type idModifier struct{}

func (r resourceArtifact) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"source_file": {
				Type:     types.StringType,
				Required: true,
			},
			"description": {
				Type:     types.StringType,
				Optional: true,
			},
			"id": {
				Type:     types.StringType,
				Computed: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					&idModifier{},
				},
			},
			"md5": {
				Type:     types.StringType,
				Computed: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					&md5Modifier{},
				},
			},
		},
	}, nil
}

func (m *md5Modifier) Description(context.Context) string {
	return ""
}

func (m *md5Modifier) MarkdownDescription(context.Context) string {
	return ""
}

func (m *md5Modifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	artifact, diags, md5, done := modifiedFile(ctx, req, resp)
	if done {
		return
	}
	if artifact.Md5.Value != md5 {
		resp.AttributePlan = types.String{Value: md5}
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (m *idModifier) Description(context.Context) string {
	return ""
}

func (m *idModifier) MarkdownDescription(context.Context) string {
	return ""
}

func (m *idModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	artifact, diags, md5, done := modifiedFile(ctx, req, resp)
	if done {
		return
	}
	if artifact.Md5.Value != md5 {
		resp.AttributePlan = types.String{Value: "", Unknown: true}
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func modifiedFile(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) (Artifact, diag.Diagnostics, string, bool) {
	// Set state
	var artifact Artifact
	diags := req.Plan.Get(ctx, &artifact)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return Artifact{}, nil, "", true
	}
	if _, err := os.Stat(artifact.SourceFile.Value); err != nil {
		resp.Diagnostics.AddError(
			"source file does not exist",
			"the file does not exist",
		)
		return Artifact{}, nil, "", true
	}

	dat, err := os.ReadFile(artifact.SourceFile.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"can't read file",
			"can't read file",
		)
		return Artifact{}, nil, "", true
	}
	md5 := fmt.Sprintf("%x", md5.Sum(dat))
	return artifact, diags, md5, false
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
	r.upload(ctx, &req.Plan, &resp.Diagnostics, &resp.State)
}

// Read resource information
func (r resourceOrder) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Set state
	var artifact Artifact
	diags := req.State.Get(ctx, &artifact)
	//TODO
	diags = resp.State.Set(ctx, &artifact)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceOrder) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	r.upload(ctx, &req.Plan, &resp.Diagnostics, &resp.State)
}

func (r resourceOrder) upload(ctx context.Context, plan *tfsdk.Plan, diag *diag.Diagnostics, state *tfsdk.State) {
	if !r.p.configured {
		diag.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply",
		)
		return
	}
	var artifact Artifact
	diags := plan.Get(ctx, &artifact)
	diag.Append(diags...)
	if diag.HasError() {
		return
	}

	if artifact.SourceFile.Null {
		diag.AddError(
			"source file is mandatory",
			"you need to provide the path of the artifact",
		)
		return
	}

	if _, err := os.Stat(artifact.SourceFile.Value); err != nil {
		diag.AddError(
			"source file does not exist",
			"the file does not exist",
		)
		return
	}

	dat, err := os.ReadFile(artifact.SourceFile.Value)
	if err != nil {
		diag.AddError(
			"can't read file",
			"can't read file",
		)
		return
	}

	var description string
	if artifact.Description.Value != "" {
		description = artifact.Description.Value
	}
	id, err := r.p.client.UploadArtifact(dat, description)
	if err != nil {
		diag.AddError(
			"Error creating artifact",
			"Could not create artifact, unexpected error: "+err.Error(),
		)
		return
	}
	md5 := fmt.Sprintf("%x", md5.Sum(dat))
	// Generate resource state struct
	var result = Artifact{
		ID:          types.String{Value: id},
		SourceFile:  artifact.SourceFile,
		Md5:         types.String{Value: md5},
		Description: types.String{Value: description},
	}
	diags = state.Set(ctx, result)
	diag.Append(diags...)
	if diag.HasError() {
		return
	}
}

// Delete resource
func (r resourceOrder) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
}
