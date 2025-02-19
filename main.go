package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"
)

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func generateTag(fieldName, newTag string, existingTag string) string {
	if existingTag != "" {
		existingTag = strings.Trim(existingTag, "`")
	}
	fieldName = toSnakeCase(fieldName)
	if len(existingTag) > 0 {
		return fmt.Sprintf("`%s %s:\"%s\"`", existingTag, newTag, fieldName)
	} else {
		return fmt.Sprintf("`%s:\"%s\"`", newTag, fieldName)
	}
}

func addTags(a *ast.File, tagname string) error {
	ast.Inspect(a, func(n ast.Node) bool {
		// continue only if its starts with type declaration
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		// check type declaration is for a struct
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}
		if structType.Fields == nil {
			return true
		}
		for _, field := range structType.Fields.List {
			if len(field.Names) == 0 {
				continue
			}
			var existingTag string
			if field.Tag != nil {
				existingTag = field.Tag.Value
			}
			newTag := generateTag(field.Names[0].Name, tagname, existingTag)
			field.Tag = &ast.BasicLit{
				Kind:  token.STRING,
				Value: newTag,
			}
		}
		return true
	})
	return nil
}

func Parse(content string, tagName string) error {
	fst := token.NewFileSet()
	a, err := parser.ParseFile(fst, "stdin.go", content, parser.ParseComments)
	if err != nil {
		return err
	}
	addTags(a, tagName)
	buf := bytes.NewBuffer(nil)
	err = format.Node(buf, fst, a)
	if err != nil {
		return err
	}
	out := buf.String()
	idx := strings.Index(out, "type")
	if idx != -1 {
		out = out[idx:]
	}
	fmt.Print(out)
	return nil
}

func main() {
	tagName := "missingTagArg"
	if len(os.Args) == 2 {
		tagName = os.Args[1]
	}
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	isPipe := fi.Mode()&os.ModeNamedPipe != 0
	if !isPipe {
		fmt.Fprintf(os.Stderr, "%s has to be run with stdin in pipe mode\n", os.Args[0])
		return
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	content := new(strings.Builder)
	_, err = content.WriteString("package stdin\n")
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	_, err = content.WriteString(string(data))
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	err = Parse(content.String(), tagName)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
}
