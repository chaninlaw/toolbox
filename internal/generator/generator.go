package generator

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/chaninlaw/toolbox/pkgs/logs"
	"github.com/chaninlaw/toolbox/pkgs/utils"
)

type Options struct {
	ProjectName string
	LiveReload  bool
}

// Generate creates a new project structure
func Generate(options Options) error {
	dir := filepath.Join(options.ProjectName)

	// Generating folder...
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		logs.Error("failed to create directory: %v", err)
		return err
	}

	// Creating internal directory...
	err = os.MkdirAll(filepath.Join(dir, "internal"), os.ModePerm)
	if err != nil {
		logs.Error("failed to create internal directory: %v", err)
		return err
	}
	// Creating pkg directory...
	err = os.MkdirAll(filepath.Join(dir, "pkg"), os.ModePerm)
	if err != nil {
		logs.Error("failed to create pkg directory: %v", err)
		return err
	}

	// Creating main.go file...
	filePath := filepath.Join(dir, "main.go")
	tmplPath := utils.AbsolutePath("boilerplate.go.tmpl")
	if err := createFileAndParseTemplate(filePath, tmplPath, options); err != nil {
		logs.Error("failed to create main.go file: %v", err)
		return err
	}

	// Creating README.md file...
	filePath = filepath.Join(dir, "README.md")
	tmplPath = utils.AbsolutePath("readme.go.tmpl")
	if err := createFileAndParseTemplate(filePath, tmplPath, options); err != nil {
		logs.Error("failed to create README.md file: %v", err)
		return err
	}

	// Creating Makefile...
	filePath = filepath.Join(dir, "Makefile")
	tmplPath = utils.AbsolutePath("makefile.go.tmpl")
	if err := createFileAndParseTemplate(filePath, tmplPath, options); err != nil {
		logs.Error("failed to create Makefile: %v", err)
		return err
	}

	// Executing go mod init...
	if err := utils.ExecCommandInDir(dir, "go", "mod", "init", options.ProjectName); err != nil {
		logs.Error("failed to execute go mod init: %v", err)
		return err
	}

	// Initializing git...
	if err := utils.ExecCommandInDir(dir, "git", "init"); err != nil {
		logs.Error("failed to initialize git: %v", err)
		return err
	}

	// Creating gitignore file...
	filePath = filepath.Join(dir, ".gitignore")
	tmplPath = utils.AbsolutePath("gitignore.go.tmpl")
	if err := createFileAndParseTemplate(filePath, tmplPath, options); err != nil {
		logs.Error("failed to create .gitignore file: %v", err)
		return err
	}

	if options.LiveReload {
		// Checking air -v command...
		if err := utils.ExecCommandInDir(dir, "air", "-v"); err != nil {
			logs.Warn("air is not installed, please install it first")
			logs.Warn("you can install it by running: go install github.com/cosmtrek/air@latest")
		}

		// Creating air.toml file...
		filePath = filepath.Join(dir, ".air.toml")
		tmplPath = utils.AbsolutePath("air.go.tmpl")
		if err := createFileAndParseTemplate(filePath, tmplPath, options); err != nil {
			logs.Error("failed to create .air.toml file: %v", err)
			return err
		}
	}

	return nil
}

func createFileAndParseTemplate(filePath, tmplPath string, options Options) error {
	file, err := os.Create(filePath)
	if err != nil {
		logs.Error("failed to create file: %v", err)
		return err
	}
	defer file.Close()

	// Parsing and executin
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		logs.Error("failed to parse template: %v", err)
		return err
	}
	// Writing template to file
	return tmpl.Execute(file, options)
}
