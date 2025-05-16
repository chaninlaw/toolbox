package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/chaninlaw/toolbox/pkgs/utils"
)

type Options struct {
	ProjectName string
	LiveReload  bool
}

// Generate creates a new project structure
func Generate(options Options) error {
	// Support custom path, but use last element as project name
	dir := filepath.Clean(options.ProjectName)
	projectName := filepath.Base(dir)

	// Generating folder...
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("failed to create directory: %v\n", err)
		return err
	}

	// Creating internal directory...
	err = os.MkdirAll(filepath.Join(dir, "internal"), os.ModePerm)
	if err != nil {
		fmt.Printf("failed to create internal directory: %v\n", err)
		return err
	}
	// Creating pkgs directory...
	err = os.MkdirAll(filepath.Join(dir, "pkgs"), os.ModePerm)
	if err != nil {
		fmt.Printf("failed to create pkg directory: %v\n", err)
		return err
	}

	// Creating main.go file...
	filePath := filepath.Join(dir, "main.go")
	tmplPath := utils.AbsolutePath("boilerplate.go.tmpl")
	options.ProjectName = projectName // Use only the last element for templates and go mod
	if err := createFileAndParseTemplate(filePath, tmplPath, options); err != nil {
		fmt.Printf("failed to create main.go file: %v\n", err)
		return err
	}

	// Creating README.md file...
	filePath = filepath.Join(dir, "README.md")
	tmplPath = utils.AbsolutePath("readme.go.tmpl")
	if err := createFileAndParseTemplate(filePath, tmplPath, options); err != nil {
		fmt.Printf("failed to create README.md file: %v\n", err)
		return err
	}

	// Creating Makefile...
	filePath = filepath.Join(dir, "Makefile")
	tmplPath = utils.AbsolutePath("makefile.go.tmpl")
	if err := createFileAndParseTemplate(filePath, tmplPath, options); err != nil {
		fmt.Printf("failed to create Makefile: %v\n", err)
		return err
	}

	// Executing go mod init...
	if err := utils.ExecCommandInDir(dir, "go", "mod", "init", projectName); err != nil {
		fmt.Printf("failed to execute go mod init: %v\n", err)
		return err
	}

	// Initializing git...
	if err := utils.ExecCommandInDir(dir, "git", "init"); err != nil {
		fmt.Printf("failed to initialize git: %v\n", err)
		return err
	}

	// Creating gitignore file...
	filePath = filepath.Join(dir, ".gitignore")
	tmplPath = utils.AbsolutePath("gitignore.go.tmpl")
	if err := createFileAndParseTemplate(filePath, tmplPath, options); err != nil {
		fmt.Printf("failed to create .gitignore file: %v\n", err)
		return err
	}

	if options.LiveReload {
		// Checking air -v command...
		if err := utils.ExecCommandInDir(dir, "air", "-v"); err != nil {
			fmt.Println("air is not installed, please install it first")
			fmt.Println("you can install it by running: go install github.com/air-verse/air@latest")
		}

		// Creating air.toml file...
		filePath = filepath.Join(dir, ".air.toml")
		tmplPath = utils.AbsolutePath("air.go.tmpl")
		if err := createFileAndParseTemplate(filePath, tmplPath, options); err != nil {
			fmt.Printf("failed to create .air.toml file: %v", err)
			return err
		}
	}

	return nil
}

func createFileAndParseTemplate(filePath, tmplPath string, options Options) error {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("failed to create file: %v\n", err)
		return err
	}
	defer file.Close()

	// Parsing and executin
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		fmt.Printf("failed to parse template: %v\n", err)
		return err
	}
	// Writing template to file
	return tmpl.Execute(file, options)
}
