package fields

import (
	"go/ast"
	"go/types"
	"reflect"
)

const (
	TagName          = "structfieldorder"
	OptionalTagValue = "optional"
)

type StructField struct {
	Name     string
	Exported bool
	Optional bool
}

type StructFields []*StructField

// NewStructFields creates a new [StructFields] from a given struct type.
// StructFields items are listed in order they appear in the struct.
func NewStructFields(strct *types.Struct) StructFields {
	sf := make(StructFields, 0, strct.NumFields())

	for i := 0; i < strct.NumFields(); i++ {
		f := strct.Field(i)

		sf = append(sf, &StructField{
			Name:     f.Name(),
			Exported: f.Exported(),
			Optional: HasOptionalTag(strct.Tag(i)),
		})
	}

	return sf
}

func HasOptionalTag(tags string) bool {
	return reflect.StructTag(tags).Get(TagName) == OptionalTagValue
}

// String returns a comma-separated list of field names.
func (sf StructFields) String() (res string) {
	for i := 0; i < len(sf); i++ {
		if res != "" {
			res += ", "
		}

		res += sf[i].Name
	}

	return res
}

func isOrderedSubset(A, B []string) bool {
	var i, j int

	for i < len(A) && j < len(B) {
		if A[i] == B[j] {
			j++
		}
		i++
	}

	return j == len(B)
}

// isInSlice checks if an element is in a given slice.
func isInSlice(fields []string, field *StructField) bool {
	for _, f := range fields {
		if f == field.Name {
			return true
		}
	}
	return false
}

// orderedIntersection finds the intersection of A and B in the order of A.
func orderedIntersection(A StructFields, B []string) StructFields {
	var intersection StructFields
	for _, a := range A {
		if isInSlice(B, a) {
			intersection = append(intersection, a)
		}
	}
	return intersection
}

// OrderedFields returns true if the fields are the correct order otherwise false.
//
//revive:disable-next-line:cyclomatic
func (sf StructFields) OrderedFields(lit *ast.CompositeLit, onlyExported bool) (bool, StructFields) {
	if len(lit.Elts) != 0 && !isNamedLiteral(lit) {
		return true, nil
	}

	structFieldNames := make([]string, len(sf))
	for i, field := range sf {
		structFieldNames[i] = field.Name
	}

	instanceFieldNames := make([]string, len(lit.Elts))
	for i := range instanceFieldNames {
		kv, ok := lit.Elts[i].(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		k, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		instanceFieldNames[i] = k.Name
	}

	if isOrderedSubset(structFieldNames, instanceFieldNames) {
		return true, nil
	} else {
		return false, orderedIntersection(sf, instanceFieldNames)
	}
}

func (sf StructFields) existenceMap() map[string]bool {
	m := make(map[string]bool, len(sf))

	for i := 0; i < len(sf); i++ {
		m[sf[i].Name] = false
	}

	return m
}

// isNamedLiteral returns true if the given literal is unnamed.
//
// The logic is basing on the principle that literal is named or unnamed,
// therefore is literal's first element is a [ast.KeyValueExpr], it is named.
//
// Method will panic if the given literal is empty.
func isNamedLiteral(lit *ast.CompositeLit) bool {
	if _, ok := lit.Elts[0].(*ast.KeyValueExpr); !ok {
		return false
	}

	return true
}
