package project

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/parsiya/borrowedtime/shared"
)

// FileInfo is a struct created from os.FileInfo interface for serialization.
type FileInfo struct {
	Name     string `json:"name"`
	IsDir    bool   `json:"isdir"`
	Template string `json:"template"` // File template.
}

// Node represents a node in a directory tree.
type Node struct {
	FullPath   string    `json:"path"`
	Info       *FileInfo `json:"info"`
	Children   []*Node   `json:"children"`
	PostCreate []string  `json:"postcreate,omitempty"`
}

// Create creates the file or directory represented by the node and its children.
func (n *Node) Create(p Project, overwrite bool) error {
	if n == nil || n.Info == nil {
		return fmt.Errorf("project.Node.Create: missing node info")
	}
	if shared.Debug {
		log.Printf("create path %q (directory=%t, overwrite=%t)", n.FullPath, n.Info.IsDir, overwrite)
	}

	exists, err := shared.PathExists(n.FullPath)
	// Check if we have access to the path.
	if err != nil {
		return fmt.Errorf("project.Node.Create: %s", err.Error())
	}
	// Check if path exists and return an error if overwrite is not set.
	if exists && !overwrite {
		// TODO: See how we can change this to support overwriting partial parts
		// of the project.
		return fmt.Errorf("project.Node.Create: %s already exists, use overwrite", n.FullPath)
	}
	// If node is a directory, create the directory. No need to create if it already exists.
	if n.Info.IsDir && !exists {
		if err := os.MkdirAll(n.FullPath, os.ModePerm); err != nil {
			return err
		}
	}
	// If file, create the file.
	if !n.Info.IsDir {
		tmplString := ""
		if n.Info.Template != "" {
			tmplString, err = p.generateTemplate(n.Info.Template, false)
			if err != nil {
				return fmt.Errorf("project.Node.Create: render %s - %s", n.FullPath, err)
			}
		}
		flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL
		if overwrite {
			flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		}
		f, err := os.OpenFile(n.FullPath, flags, 0666)
		if err != nil {
			return err
		}
		num, writeErr := f.WriteString(tmplString)
		closeErr := f.Close()
		if writeErr != nil {
			return fmt.Errorf("project.Node.Create: error writing %s - %s", n.FullPath, writeErr)
		}
		if num != len(tmplString) {
			return fmt.Errorf("project.Node.Create: short write to %s", n.FullPath)
		}
		if closeErr != nil {
			return closeErr
		}
	}
	// Create all node's children.
	for _, child := range n.Children {
		// Calculate and populate FullPath based on the parent. Everything but
		// root should have a parent and if we are here, then we are not
		// populating root.
		// fmt.Printf("processing child.FullPath: %s\n", child.FullPath)
		// fmt.Printf("Parent's FullPath is %v\n", n.FullPath)
		// fmt.Printf("joined: %s\n", filepath.Join(n.FullPath, child.FullPath))
		// This is a good place to place logging statements.
		// E.g., creating blahblah/whatever.txt using X template.
		child.FullPath = filepath.Join(n.FullPath, child.FullPath)
		if err := child.Create(p, overwrite); err != nil {
			return err
		}
	}
	return nil
}
