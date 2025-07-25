// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package resource_rule

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func RuleResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"body": schema.StringAttribute{
				Required:            true,
				Description:         "The python body of the rule",
				MarkdownDescription: "The python body of the rule",
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"created_by": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"type": schema.StringAttribute{
						Computed: true,
					},
				},
				CustomType: CreatedByType{
					ObjectType: types.ObjectType{
						AttrTypes: CreatedByValue{}.AttributeTypes(ctx),
					},
				},
				Computed:            true,
				Description:         "The actor who created the rule",
				MarkdownDescription: "The actor who created the rule",
			},
			"created_by_external": schema.StringAttribute{
				Computed:            true,
				Description:         "The text of the user-provided CreatedBy field when uploaded via CI/CD",
				MarkdownDescription: "The text of the user-provided CreatedBy field when uploaded via CI/CD",
			},
			"dedup_period_minutes": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The amount of time in minutes for grouping alerts",
				MarkdownDescription: "The amount of time in minutes for grouping alerts",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				Default: int64default.StaticInt64(60),
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The description of the rule",
				MarkdownDescription: "The description of the rule",
			},
			"display_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The display name of the rule",
				MarkdownDescription: "The display name of the rule",
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Determines whether or not the rule is active",
				MarkdownDescription: "Determines whether or not the rule is active",
			},
			"inline_filters": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The filter for the rule represented in YAML",
				MarkdownDescription: "The filter for the rule represented in YAML",
			},
			"last_modified": schema.StringAttribute{
				Computed: true,
			},
			"log_types": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "log types",
				MarkdownDescription: "log types",
			},
			"managed": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Determines if the rule is managed by panther",
				MarkdownDescription: "Determines if the rule is managed by panther",
			},
			"output_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "Destination IDs that override default alert routing based on severity",
				MarkdownDescription: "Destination IDs that override default alert routing based on severity",
			},
			"reports": schema.MapAttribute{
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
				Optional:            true,
				Computed:            true,
				Description:         "reports",
				MarkdownDescription: "reports",
			},
			"runbook": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "How to handle the generated alert",
				MarkdownDescription: "How to handle the generated alert",
			},
			"severity": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"INFO",
						"LOW",
						"MEDIUM",
						"HIGH",
						"CRITICAL",
					),
				},
			},
			"summary_attributes": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "A list of fields in the event to create top 5 summaries for",
				MarkdownDescription: "A list of fields in the event to create top 5 summaries for",
			},
			"tags": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "The tags for the rule",
				MarkdownDescription: "The tags for the rule",
			},
			"tests": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"expected_result": schema.BoolAttribute{
							Required:            true,
							Description:         "The expected result",
							MarkdownDescription: "The expected result",
						},
						"mocks": schema.ListAttribute{
							ElementType: types.MapType{
								ElemType: types.StringType,
							},
							Optional:            true,
							Computed:            true,
							Description:         "mocks",
							MarkdownDescription: "mocks",
						},
						"name": schema.StringAttribute{
							Required:            true,
							Description:         "name",
							MarkdownDescription: "name",
						},
						"resource": schema.StringAttribute{
							Required:            true,
							Description:         "resource",
							MarkdownDescription: "resource",
						},
					},
					CustomType: TestsType{
						ObjectType: types.ObjectType{
							AttrTypes: TestsValue{}.AttributeTypes(ctx),
						},
					},
				},
				Optional:            true,
				Computed:            true,
				Description:         "Unit tests for the Rule. Best practice is to include a positive and negative case",
				MarkdownDescription: "Unit tests for the Rule. Best practice is to include a positive and negative case",
			},
			"threshold": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "the number of events that must match before an alert is triggered",
				MarkdownDescription: "the number of events that must match before an alert is triggered",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				Default: int64default.StaticInt64(1),
			},
		},
	}
}

type RuleModel struct {
	Id                 types.String   `tfsdk:"id"`
	Body               types.String   `tfsdk:"body"`
	CreatedAt          types.String   `tfsdk:"created_at"`
	CreatedBy          CreatedByValue `tfsdk:"created_by"`
	CreatedByExternal  types.String   `tfsdk:"created_by_external"`
	DedupPeriodMinutes types.Int64    `tfsdk:"dedup_period_minutes"`
	Description        types.String   `tfsdk:"description"`
	DisplayName        types.String   `tfsdk:"display_name"`
	Enabled            types.Bool     `tfsdk:"enabled"`
	InlineFilters      types.String   `tfsdk:"inline_filters"`
	LastModified       types.String   `tfsdk:"last_modified"`
	LogTypes           types.List     `tfsdk:"log_types"`
	Managed            types.Bool     `tfsdk:"managed"`
	OutputIds          types.List     `tfsdk:"output_ids"`
	Reports            types.Map      `tfsdk:"reports"`
	Runbook            types.String   `tfsdk:"runbook"`
	Severity           types.String   `tfsdk:"severity"`
	SummaryAttributes  types.List     `tfsdk:"summary_attributes"`
	Tags               types.List     `tfsdk:"tags"`
	Tests              types.List     `tfsdk:"tests"`
	Threshold          types.Int64    `tfsdk:"threshold"`
}

var _ basetypes.ObjectTypable = CreatedByType{}

type CreatedByType struct {
	basetypes.ObjectType
}

func (t CreatedByType) Equal(o attr.Type) bool {
	other, ok := o.(CreatedByType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t CreatedByType) String() string {
	return "CreatedByType"
}

func (t CreatedByType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return nil, diags
	}

	idVal, ok := idAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be basetypes.StringValue, was: %T`, idAttribute))
	}

	typeAttribute, ok := attributes["type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`type is missing from object`)

		return nil, diags
	}

	typeVal, ok := typeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`type expected to be basetypes.StringValue, was: %T`, typeAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return CreatedByValue{
		Id:            idVal,
		CreatedByType: typeVal,
		state:         attr.ValueStateKnown,
	}, diags
}

func NewCreatedByValueNull() CreatedByValue {
	return CreatedByValue{
		state: attr.ValueStateNull,
	}
}

func NewCreatedByValueUnknown() CreatedByValue {
	return CreatedByValue{
		state: attr.ValueStateUnknown,
	}
}

func NewCreatedByValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CreatedByValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing CreatedByValue Attribute Value",
				"While creating a CreatedByValue value, a missing attribute value was detected. "+
					"A CreatedByValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CreatedByValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid CreatedByValue Attribute Type",
				"While creating a CreatedByValue value, an invalid attribute value was detected. "+
					"A CreatedByValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CreatedByValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("CreatedByValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra CreatedByValue Attribute Value",
				"While creating a CreatedByValue value, an extra attribute value was detected. "+
					"A CreatedByValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra CreatedByValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewCreatedByValueUnknown(), diags
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return NewCreatedByValueUnknown(), diags
	}

	idVal, ok := idAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be basetypes.StringValue, was: %T`, idAttribute))
	}

	typeAttribute, ok := attributes["type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`type is missing from object`)

		return NewCreatedByValueUnknown(), diags
	}

	typeVal, ok := typeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`type expected to be basetypes.StringValue, was: %T`, typeAttribute))
	}

	if diags.HasError() {
		return NewCreatedByValueUnknown(), diags
	}

	return CreatedByValue{
		Id:            idVal,
		CreatedByType: typeVal,
		state:         attr.ValueStateKnown,
	}, diags
}

func NewCreatedByValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) CreatedByValue {
	object, diags := NewCreatedByValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewCreatedByValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t CreatedByType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewCreatedByValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewCreatedByValueUnknown(), nil
	}

	if in.IsNull() {
		return NewCreatedByValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewCreatedByValueMust(CreatedByValue{}.AttributeTypes(ctx), attributes), nil
}

func (t CreatedByType) ValueType(ctx context.Context) attr.Value {
	return CreatedByValue{}
}

var _ basetypes.ObjectValuable = CreatedByValue{}

type CreatedByValue struct {
	Id            basetypes.StringValue `tfsdk:"id"`
	CreatedByType basetypes.StringValue `tfsdk:"type"`
	state         attr.ValueState
}

func (v CreatedByValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["id"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["type"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.Id.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["id"] = val

		val, err = v.CreatedByType.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["type"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v CreatedByValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CreatedByValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CreatedByValue) String() string {
	return "CreatedByValue"
}

func (v CreatedByValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"id":   basetypes.StringType{},
		"type": basetypes.StringType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"id":   v.Id,
			"type": v.CreatedByType,
		})

	return objVal, diags
}

func (v CreatedByValue) Equal(o attr.Value) bool {
	other, ok := o.(CreatedByValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Id.Equal(other.Id) {
		return false
	}

	if !v.CreatedByType.Equal(other.CreatedByType) {
		return false
	}

	return true
}

func (v CreatedByValue) Type(ctx context.Context) attr.Type {
	return CreatedByType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v CreatedByValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":   basetypes.StringType{},
		"type": basetypes.StringType{},
	}
}

var _ basetypes.ObjectTypable = TestsType{}

type TestsType struct {
	basetypes.ObjectType
}

func (t TestsType) Equal(o attr.Type) bool {
	other, ok := o.(TestsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t TestsType) String() string {
	return "TestsType"
}

func (t TestsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	expectedResultAttribute, ok := attributes["expected_result"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`expected_result is missing from object`)

		return nil, diags
	}

	expectedResultVal, ok := expectedResultAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`expected_result expected to be basetypes.BoolValue, was: %T`, expectedResultAttribute))
	}

	mocksAttribute, ok := attributes["mocks"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`mocks is missing from object`)

		return nil, diags
	}

	mocksVal, ok := mocksAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`mocks expected to be basetypes.ListValue, was: %T`, mocksAttribute))
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return nil, diags
	}

	nameVal, ok := nameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be basetypes.StringValue, was: %T`, nameAttribute))
	}

	resourceAttribute, ok := attributes["resource"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`resource is missing from object`)

		return nil, diags
	}

	resourceVal, ok := resourceAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`resource expected to be basetypes.StringValue, was: %T`, resourceAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return TestsValue{
		ExpectedResult: expectedResultVal,
		Mocks:          mocksVal,
		Name:           nameVal,
		Resource:       resourceVal,
		state:          attr.ValueStateKnown,
	}, diags
}

func NewTestsValueNull() TestsValue {
	return TestsValue{
		state: attr.ValueStateNull,
	}
}

func NewTestsValueUnknown() TestsValue {
	return TestsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewTestsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (TestsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing TestsValue Attribute Value",
				"While creating a TestsValue value, a missing attribute value was detected. "+
					"A TestsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("TestsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid TestsValue Attribute Type",
				"While creating a TestsValue value, an invalid attribute value was detected. "+
					"A TestsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("TestsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("TestsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra TestsValue Attribute Value",
				"While creating a TestsValue value, an extra attribute value was detected. "+
					"A TestsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra TestsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewTestsValueUnknown(), diags
	}

	expectedResultAttribute, ok := attributes["expected_result"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`expected_result is missing from object`)

		return NewTestsValueUnknown(), diags
	}

	expectedResultVal, ok := expectedResultAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`expected_result expected to be basetypes.BoolValue, was: %T`, expectedResultAttribute))
	}

	mocksAttribute, ok := attributes["mocks"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`mocks is missing from object`)

		return NewTestsValueUnknown(), diags
	}

	mocksVal, ok := mocksAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`mocks expected to be basetypes.ListValue, was: %T`, mocksAttribute))
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return NewTestsValueUnknown(), diags
	}

	nameVal, ok := nameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be basetypes.StringValue, was: %T`, nameAttribute))
	}

	resourceAttribute, ok := attributes["resource"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`resource is missing from object`)

		return NewTestsValueUnknown(), diags
	}

	resourceVal, ok := resourceAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`resource expected to be basetypes.StringValue, was: %T`, resourceAttribute))
	}

	if diags.HasError() {
		return NewTestsValueUnknown(), diags
	}

	return TestsValue{
		ExpectedResult: expectedResultVal,
		Mocks:          mocksVal,
		Name:           nameVal,
		Resource:       resourceVal,
		state:          attr.ValueStateKnown,
	}, diags
}

func NewTestsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) TestsValue {
	object, diags := NewTestsValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewTestsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t TestsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewTestsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewTestsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewTestsValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewTestsValueMust(TestsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t TestsType) ValueType(ctx context.Context) attr.Value {
	return TestsValue{}
}

var _ basetypes.ObjectValuable = TestsValue{}

type TestsValue struct {
	ExpectedResult basetypes.BoolValue   `tfsdk:"expected_result"`
	Mocks          basetypes.ListValue   `tfsdk:"mocks"`
	Name           basetypes.StringValue `tfsdk:"name"`
	Resource       basetypes.StringValue `tfsdk:"resource"`
	state          attr.ValueState
}

func (v TestsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 4)

	var val tftypes.Value
	var err error

	attrTypes["expected_result"] = basetypes.BoolType{}.TerraformType(ctx)
	attrTypes["mocks"] = basetypes.ListType{
		ElemType: types.MapType{
			ElemType: types.StringType,
		},
	}.TerraformType(ctx)
	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["resource"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 4)

		val, err = v.ExpectedResult.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["expected_result"] = val

		val, err = v.Mocks.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["mocks"] = val

		val, err = v.Name.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["name"] = val

		val, err = v.Resource.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["resource"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v TestsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v TestsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v TestsValue) String() string {
	return "TestsValue"
}

func (v TestsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	var mocksVal basetypes.ListValue
	switch {
	case v.Mocks.IsUnknown():
		mocksVal = types.ListUnknown(types.MapType{
			ElemType: types.StringType,
		})
	case v.Mocks.IsNull():
		mocksVal = types.ListNull(types.MapType{
			ElemType: types.StringType,
		})
	default:
		var d diag.Diagnostics
		mocksVal, d = types.ListValue(types.MapType{
			ElemType: types.StringType,
		}, v.Mocks.Elements())
		diags.Append(d...)
	}

	if diags.HasError() {
		return types.ObjectUnknown(map[string]attr.Type{
			"expected_result": basetypes.BoolType{},
			"mocks": basetypes.ListType{
				ElemType: types.MapType{
					ElemType: types.StringType,
				},
			},
			"name":     basetypes.StringType{},
			"resource": basetypes.StringType{},
		}), diags
	}

	attributeTypes := map[string]attr.Type{
		"expected_result": basetypes.BoolType{},
		"mocks": basetypes.ListType{
			ElemType: types.MapType{
				ElemType: types.StringType,
			},
		},
		"name":     basetypes.StringType{},
		"resource": basetypes.StringType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"expected_result": v.ExpectedResult,
			"mocks":           mocksVal,
			"name":            v.Name,
			"resource":        v.Resource,
		})

	return objVal, diags
}

func (v TestsValue) Equal(o attr.Value) bool {
	other, ok := o.(TestsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.ExpectedResult.Equal(other.ExpectedResult) {
		return false
	}

	if !v.Mocks.Equal(other.Mocks) {
		return false
	}

	if !v.Name.Equal(other.Name) {
		return false
	}

	if !v.Resource.Equal(other.Resource) {
		return false
	}

	return true
}

func (v TestsValue) Type(ctx context.Context) attr.Type {
	return TestsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v TestsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"expected_result": basetypes.BoolType{},
		"mocks": basetypes.ListType{
			ElemType: types.MapType{
				ElemType: types.StringType,
			},
		},
		"name":     basetypes.StringType{},
		"resource": basetypes.StringType{},
	}
}
