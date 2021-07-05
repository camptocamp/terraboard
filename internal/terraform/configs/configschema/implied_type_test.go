package configschema

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestBlockImpliedType(t *testing.T) {
	tests := map[string]struct {
		Schema *Block
		Want   cty.Type
	}{
		"nil": {
			nil,
			cty.EmptyObject,
		},
		"empty": {
			&Block{},
			cty.EmptyObject,
		},
		"attributes": {
			&Block{
				Attributes: map[string]*Attribute{
					"optional": {
						Type:     cty.String,
						Optional: true,
					},
					"required": {
						Type:     cty.Number,
						Required: true,
					},
					"computed": {
						Type:     cty.List(cty.Bool),
						Computed: true,
					},
					"optional_computed": {
						Type:     cty.Map(cty.Bool),
						Optional: true,
					},
				},
			},
			cty.Object(map[string]cty.Type{
				"optional":          cty.String,
				"required":          cty.Number,
				"computed":          cty.List(cty.Bool),
				"optional_computed": cty.Map(cty.Bool),
			}),
		},
		"blocks": {
			&Block{
				BlockTypes: map[string]*NestedBlock{
					"single": &NestedBlock{
						Nesting: NestingSingle,
						Block: Block{
							Attributes: map[string]*Attribute{
								"foo": {
									Type:     cty.DynamicPseudoType,
									Required: true,
								},
							},
						},
					},
					"list": &NestedBlock{
						Nesting: NestingList,
					},
					"set": &NestedBlock{
						Nesting: NestingSet,
					},
					"map": &NestedBlock{
						Nesting: NestingMap,
					},
				},
			},
			cty.Object(map[string]cty.Type{
				"single": cty.Object(map[string]cty.Type{
					"foo": cty.DynamicPseudoType,
				}),
				"list": cty.List(cty.EmptyObject),
				"set":  cty.Set(cty.EmptyObject),
				"map":  cty.Map(cty.EmptyObject),
			}),
		},
		"deep block nesting": {
			&Block{
				BlockTypes: map[string]*NestedBlock{
					"single": &NestedBlock{
						Nesting: NestingSingle,
						Block: Block{
							BlockTypes: map[string]*NestedBlock{
								"list": &NestedBlock{
									Nesting: NestingList,
									Block: Block{
										BlockTypes: map[string]*NestedBlock{
											"set": &NestedBlock{
												Nesting: NestingSet,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			cty.Object(map[string]cty.Type{
				"single": cty.Object(map[string]cty.Type{
					"list": cty.List(cty.Object(map[string]cty.Type{
						"set": cty.Set(cty.EmptyObject),
					})),
				}),
			}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := test.Schema.ImpliedType()
			if !got.Equals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestObjectImpliedType(t *testing.T) {
	tests := map[string]struct {
		Schema *Object
		Want   cty.Type
	}{
		"nil": {
			nil,
			cty.EmptyObject,
		},
		"attributes": {
			&Object{
				Nesting: NestingSingle,
				Attributes: map[string]*Attribute{
					"optional": {
						Type:     cty.String,
						Optional: true,
					},
					"required": {
						Type:     cty.Number,
						Required: true,
					},
					"computed": {
						Type:     cty.List(cty.Bool),
						Computed: true,
					},
					"optional_computed": {
						Type:     cty.Map(cty.Bool),
						Optional: true,
					},
				},
			},
			cty.ObjectWithOptionalAttrs(
				map[string]cty.Type{
					"optional":          cty.String,
					"required":          cty.Number,
					"computed":          cty.List(cty.Bool),
					"optional_computed": cty.Map(cty.Bool),
				},
				[]string{"optional", "optional_computed"},
			),
		},
		"nested attributes": {
			&Object{
				Nesting: NestingSingle,
				Attributes: map[string]*Attribute{
					"nested_type": {
						NestedType: &Object{
							Nesting: NestingSingle,
							Attributes: map[string]*Attribute{
								"optional": {
									Type:     cty.String,
									Optional: true,
								},
								"required": {
									Type:     cty.Number,
									Required: true,
								},
								"computed": {
									Type:     cty.List(cty.Bool),
									Computed: true,
								},
								"optional_computed": {
									Type:     cty.Map(cty.Bool),
									Optional: true,
								},
							},
						},
						Optional: true,
					},
				},
			},
			cty.ObjectWithOptionalAttrs(map[string]cty.Type{
				"nested_type": cty.ObjectWithOptionalAttrs(map[string]cty.Type{
					"optional":          cty.String,
					"required":          cty.Number,
					"computed":          cty.List(cty.Bool),
					"optional_computed": cty.Map(cty.Bool),
				}, []string{"optional", "optional_computed"}),
			}, []string{"nested_type"}),
		},
		"NestingList": {
			&Object{
				Nesting: NestingList,
				Attributes: map[string]*Attribute{
					"foo": {Type: cty.String, Optional: true},
				},
			},
			cty.List(cty.ObjectWithOptionalAttrs(map[string]cty.Type{"foo": cty.String}, []string{"foo"})),
		},
		"NestingMap": {
			&Object{
				Nesting: NestingMap,
				Attributes: map[string]*Attribute{
					"foo": {Type: cty.String},
				},
			},
			cty.Map(cty.Object(map[string]cty.Type{"foo": cty.String})),
		},
		"NestingSet": {
			&Object{
				Nesting: NestingSet,
				Attributes: map[string]*Attribute{
					"foo": {Type: cty.String},
				},
			},
			cty.Set(cty.Object(map[string]cty.Type{"foo": cty.String})),
		},
		"deeply nested NestingList": {
			&Object{
				Nesting: NestingList,
				Attributes: map[string]*Attribute{
					"foo": {
						NestedType: &Object{
							Nesting: NestingList,
							Attributes: map[string]*Attribute{
								"bar": {Type: cty.String},
							},
						},
					},
				},
			},
			cty.List(cty.Object(map[string]cty.Type{"foo": cty.List(cty.Object(map[string]cty.Type{"bar": cty.String}))})),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := test.Schema.ImpliedType()
			if !got.Equals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestObjectContainsSensitive(t *testing.T) {
	tests := map[string]struct {
		Schema *Object
		Want   bool
	}{
		"object contains sensitive": {
			&Object{
				Attributes: map[string]*Attribute{
					"sensitive": {Sensitive: true},
				},
			},
			true,
		},
		"no sensitive attrs": {
			&Object{
				Attributes: map[string]*Attribute{
					"insensitive": {},
				},
			},
			false,
		},
		"nested object contains sensitive": {
			&Object{
				Attributes: map[string]*Attribute{
					"nested": {
						NestedType: &Object{
							Nesting: NestingSingle,
							Attributes: map[string]*Attribute{
								"sensitive": {Sensitive: true},
							},
						},
					},
				},
			},
			true,
		},
		"nested obj, no sensitive attrs": {
			&Object{
				Attributes: map[string]*Attribute{
					"nested": {
						NestedType: &Object{
							Nesting: NestingSingle,
							Attributes: map[string]*Attribute{
								"public": {},
							},
						},
					},
				},
			},
			false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := test.Schema.ContainsSensitive()
			if got != test.Want {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}

}
