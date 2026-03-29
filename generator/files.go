package generator

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func (g *Generator) generateFile(projectPath, relativePath, tmplContent string, data templateData) error {
	tmpl, err := template.New("template").Funcs(template.FuncMap{
		"ToClassName": toClassName,
	}).Parse(tmplContent)
	if err != nil {
		return fmt.Errorf("failed to parse template for %s: %w", relativePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template for %s: %w", relativePath, err)
	}

	fullPath := filepath.Join(projectPath, relativePath)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := os.WriteFile(fullPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", fullPath, err)
	}

	fmt.Printf("  [OK] Generated: %s\n", relativePath)
	return nil
}

func (g *Generator) generateWrapperScript(projectPath, filename, content string, perm os.FileMode) error {
	fullPath := filepath.Join(projectPath, filename)
	if err := os.WriteFile(fullPath, []byte(content), perm); err != nil {
		return fmt.Errorf("failed to write file %s: %w", fullPath, err)
	}
	fmt.Printf("  [OK] Generated: %s\n", filename)
	return nil
}

func (g *Generator) downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func toClassName(modID string) string {
	parts := strings.Split(modID, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}
