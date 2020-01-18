package crud

import (
	"fmt"
	"github.com/elvsn/scim.go/crud/expr"
	"github.com/elvsn/scim.go/prop"
	"github.com/elvsn/scim.go/spec"
)

// Add value to SCIM resource at the given SCIM path. If SCIM path is empty, value will be added
// to the root of the resource. The supplied value must be compatible with the target property attribute,
// otherwise error will be returned.
func Add(resource *prop.Resource, path string, value interface{}) error {
	if len(path) == 0 {
		return resource.Navigator().Add(value)
	}

	head, err := expr.CompilePath(path)
	if err != nil {
		return err
	}

	return defaultTraverse(resource.RootProperty(), skipMainSchemaNamespace(resource, head), func(nav *prop.Navigator) error {
		return nav.Add(value)
	})
}

// Replace value in SCIM resource at the given SCIM path. If SCIM path is empty, the root of the resource
// will be replaced. The supplied value must be compatible with the target property attribute, otherwise
// error will be returned.
func Replace(resource *prop.Resource, path string, value interface{}) error {
	if len(path) == 0 {
		return resource.Navigator().Replace(value)
	}

	head, err := expr.CompilePath(path)
	if err != nil {
		return err
	}

	return defaultTraverse(resource.RootProperty(), skipMainSchemaNamespace(resource, head), func(nav *prop.Navigator) error {
		return nav.Replace(value)
	})
}

// Delete value from the SCIM resource at the specified SCIM path. The path cannot be empty.
func Delete(resource *prop.Resource, path string) error {
	if len(path) == 0 {
		return fmt.Errorf("%w: path must be specified for delete operation", spec.ErrInvalidPath)
	}

	head, err := expr.CompilePath(path)
	if err != nil {
		return err
	}

	return defaultTraverse(resource.RootProperty(), skipMainSchemaNamespace(resource, head), func(nav *prop.Navigator) error {
		return nav.Delete()
	})
}

func skipMainSchemaNamespace(resource *prop.Resource, query *expr.Expression) *expr.Expression {
	if query == nil {
		return nil
	}

	if query.IsPath() && query.Token() == resource.ResourceType().Schema().ID() {
		return query.Next()
	}

	return query
}
