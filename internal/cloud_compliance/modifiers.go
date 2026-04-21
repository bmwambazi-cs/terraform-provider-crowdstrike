package cloudcompliance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// cloneWarningModifier emits a warning during plan when a framework will be
// created by cloning a parent framework.
type cloneWarningModifier struct{}

func (m cloneWarningModifier) Description(_ context.Context) string {
	return "Warns the user when a framework will be created by cloning a parent framework."
}

func (m cloneWarningModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m cloneWarningModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Only warn on create when parent_framework_id is set.
	if !req.StateValue.IsNull() || req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	resp.Diagnostics.AddWarning(
		"Cloning compliance framework",
		fmt.Sprintf(
			"This framework will be created by cloning the parent framework %q. Sections and controls will be inherited from the parent. If sections are also specified in config, they will be merged on top of the cloned sections.",
			req.PlanValue.ValueString(),
		),
	)
}

func cloneWarning() planmodifier.String {
	return cloneWarningModifier{}
}

// sectionsDefaultNullModifier sets sections to null during plan when neither
// sections nor parent_framework_id are in the config. This prevents the
// Computed attribute from showing "(known after apply)" on non-clone creates.
type sectionsDefaultNullModifier struct{}

func (m sectionsDefaultNullModifier) Description(_ context.Context) string {
	return "Sets sections to null when not configured and not cloning."
}

func (m sectionsDefaultNullModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m sectionsDefaultNullModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	// Only act when sections is unknown (Computed kicking in).
	if !req.PlanValue.IsUnknown() {
		return
	}

	// Check if parent_framework_id is set in the plan.
	var parentID types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("parent_framework_id"), &parentID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If cloning, leave sections as unknown — the server will populate them.
	if !parentID.IsNull() && !parentID.IsUnknown() {
		return
	}

	// Not cloning and no sections configured — plan as null.
	resp.PlanValue = types.MapNull(req.PlanValue.ElementType(ctx))
}

func sectionsDefaultNull() planmodifier.Map {
	return sectionsDefaultNullModifier{}
}
