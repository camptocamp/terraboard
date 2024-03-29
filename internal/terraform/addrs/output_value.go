package addrs

import (
	"fmt"

	"github.com/camptocamp/terraboard/internal/terraform/tfdiags"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// OutputValue is the address of an output value, in the context of the module
// that is defining it.
//
// This is related to but separate from ModuleCallOutput, which represents
// a module output from the perspective of its parent module. Since output
// values cannot be represented from the module where they are defined,
// OutputValue is not Referenceable, while ModuleCallOutput is.
type OutputValue struct {
	Name string
}

func (v OutputValue) String() string {
	return "output." + v.Name
}

// Absolute converts the receiver into an absolute address within the given
// module instance.
func (v OutputValue) Absolute(m ModuleInstance) AbsOutputValue {
	return AbsOutputValue{
		Module:      m,
		OutputValue: v,
	}
}

// AbsOutputValue is the absolute address of an output value within a module instance.
//
// This represents an output globally within the namespace of a particular
// configuration. It is related to but separate from ModuleCallOutput, which
// represents a module output from the perspective of its parent module.
type AbsOutputValue struct {
	checkable
	Module      ModuleInstance
	OutputValue OutputValue
}

// OutputValue returns the absolute address of an output value of the given
// name within the receiving module instance.
func (m ModuleInstance) OutputValue(name string) AbsOutputValue {
	return AbsOutputValue{
		Module: m,
		OutputValue: OutputValue{
			Name: name,
		},
	}
}

func (v AbsOutputValue) Check(t CheckType, i int) Check {
	return Check{
		Container: v,
		Type:      t,
		Index:     i,
	}
}

func (v AbsOutputValue) String() string {
	if v.Module.IsRoot() {
		return v.OutputValue.String()
	}
	return fmt.Sprintf("%s.%s", v.Module.String(), v.OutputValue.String())
}

func (v AbsOutputValue) Equal(o AbsOutputValue) bool {
	return v.OutputValue == o.OutputValue && v.Module.Equal(o.Module)
}

func ParseAbsOutputValue(traversal hcl.Traversal) (AbsOutputValue, tfdiags.Diagnostics) {
	path, remain, diags := parseModuleInstancePrefix(traversal)
	if diags.HasErrors() {
		return AbsOutputValue{}, diags
	}

	if len(remain) != 2 {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid address",
			Detail:   "An output name is required.",
			Subject:  traversal.SourceRange().Ptr(),
		})
		return AbsOutputValue{}, diags
	}

	if remain.RootName() != "output" {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid address",
			Detail:   "Output address must start with \"output.\".",
			Subject:  remain[0].SourceRange().Ptr(),
		})
		return AbsOutputValue{}, diags
	}

	var name string
	switch tt := remain[1].(type) {
	case hcl.TraverseAttr:
		name = tt.Name
	default:
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid address",
			Detail:   "An output name is required.",
			Subject:  remain[1].SourceRange().Ptr(),
		})
		return AbsOutputValue{}, diags
	}

	return AbsOutputValue{
		Module: path,
		OutputValue: OutputValue{
			Name: name,
		},
	}, diags
}

func ParseAbsOutputValueStr(str string) (AbsOutputValue, tfdiags.Diagnostics) {
	var diags tfdiags.Diagnostics

	traversal, parseDiags := hclsyntax.ParseTraversalAbs([]byte(str), "", hcl.Pos{Line: 1, Column: 1})
	diags = diags.Append(parseDiags)
	if parseDiags.HasErrors() {
		return AbsOutputValue{}, diags
	}

	addr, addrDiags := ParseAbsOutputValue(traversal)
	diags = diags.Append(addrDiags)
	return addr, diags
}

// ModuleCallOutput converts an AbsModuleOutput into a ModuleCallOutput,
// returning also the module instance that the ModuleCallOutput is relative
// to.
//
// The root module does not have a call, and so this method cannot be used
// with outputs in the root module, and will panic in that case.
func (v AbsOutputValue) ModuleCallOutput() (ModuleInstance, ModuleCallInstanceOutput) {
	if v.Module.IsRoot() {
		panic("ReferenceFromCall used with root module output")
	}

	caller, call := v.Module.CallInstance()
	return caller, ModuleCallInstanceOutput{
		Call: call,
		Name: v.OutputValue.Name,
	}
}
