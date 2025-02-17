package parse

import (
	"fmt"
	"strings"

	"github.com/open-policy-agent/opa/ast"

	rio "github.com/styrainc/regal/internal/io"
)

// ParserOptions provides the parse options necessary to include location data in AST results.
func ParserOptions() ast.ParserOptions {
	return ast.ParserOptions{
		ProcessAnnotation: true,
		JSONOptions: &ast.JSONOptions{
			MarshalOptions: ast.JSONMarshalOptions{
				IncludeLocation: ast.NodeToggle{
					Term:           true,
					Package:        true,
					Comment:        true,
					Import:         true,
					Rule:           true,
					Head:           true,
					Expr:           true,
					SomeDecl:       true,
					Every:          true,
					With:           true,
					Annotations:    true,
					AnnotationsRef: true,
				},
			},
		},
	}
}

// MustParseModule works like ast.MustParseModule but with the Regal parser options applied.
func MustParseModule(policy string) *ast.Module {
	return ast.MustParseModuleWithOpts(policy, ParserOptions())
}

// EnhanceAST TODO rename with https://github.com/StyraInc/regal/issues/86.
func EnhanceAST(name string, content string, module *ast.Module) (map[string]any, error) {
	var enhancedAst map[string]any

	if err := rio.JSONRoundTrip(module, &enhancedAst); err != nil {
		return nil, fmt.Errorf("JSON rountrip failed for module: %w", err)
	}

	enhancedAst["regal"] = map[string]any{
		"file": map[string]any{
			"name":  name,
			"lines": strings.Split(content, "\n"),
		},
	}

	return enhancedAst, nil
}
