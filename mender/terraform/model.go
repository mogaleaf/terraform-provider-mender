package terraform

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Artifact struct {
	SourceFile types.String `tfsdk:"source_file"`
	ID         types.String `tfsdk:"id"`
	Md5        types.String `tfsdk:"md5"`
}
