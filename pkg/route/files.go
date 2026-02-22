package route

import (
	"log/slog"
	"net/http"
	"strings"

	"smolgit/pkg/model"

	"github.com/gin-gonic/gin"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// type Node struct {
// 	Name     string           `json:"filename"`
// 	IsDir    bool             `json:"is_dir"`
// 	Children map[string]*Node `json:"children"`
// 	Path     string           `json:"path"`
// }

type Node struct {
	Name  string `json:"name,omitempty"`
	Items []Node `json:"items,omitempty"`
}

// buildFileTree constructs the hierarchical tree from a list of file paths
func buildFileTree(paths []string) *Node {
	root := &Node{
		Name:  "",
		Items: []Node{},
	}

	for _, path := range paths {
		isDir := strings.HasSuffix(path, "/")
		path = strings.TrimSuffix(path, "/") // Remove trailing slash

		parts := strings.Split(path, "/")
		current := root

		// Traverse through the directory structure
		for i := 0; i < len(parts)-1; i++ {
			part := parts[i]
			found := false

			// Check if a child with this name already exists
			for j := range current.Items {
				if current.Items[j].Name == part {
					current = &current.Items[j]
					found = true
					break
				}
			}

			if !found {
				// Create a new directory node
				newNode := &Node{
					Name:  part,
					Items: []Node{},
				}
				current.Items = append(current.Items, *newNode)
				current = newNode
			}
		}

		// Handle the final part (file or directory)
		lastPart := parts[len(parts)-1]
		if isDir {
			// Add a new directory node
			newNode := &Node{
				Name:  lastPart,
				Items: []Node{},
			}
			current.Items = append(current.Items, *newNode)
		} else {
			// Add a file node
			current.Items = append(current.Items, Node{Name: lastPart})
		}
	}

	return root
}

// collectFlatNodes returns a flat list of nodes, with each directory containing its items
func collectFlatNodes(node *Node) []Node {
	var result []Node

	// Function to recursively traverse the tree
	var traverse func(*Node)
	traverse = func(n *Node) {
		// Add current node to the result
		result = append(result, Node{Name: n.Name, Items: n.Items})

		// If this node has children, recursively process them
		for _, child := range n.Items {
			traverse(&child)
		}
	}

	traverse(node)
	return result
}

func (r *Route) Files(c *gin.Context) {
	user, repoPath := c.Param("user"), c.Param("path")
	fullPath := "/" + user + "/" + repoPath
	// TODO more secure way to open repo
	// check user permissions
	slog.Debug("open repository", "fullpath", r.cfg.GitPath+fullPath)
	repo, err := gogit.PlainOpen(r.cfg.GitPath + fullPath)
	if err != nil {
		slog.Error("cant find repository", "err", err)
		c.JSON(http.StatusNotFound, gin.H{"title": "cant find any repository"})
		return
	}
	headRef, err := repo.Head()
	if err != nil {
		slog.Error("cant get repo HEAD", "err", err)
		c.JSON(http.StatusNotFound, gin.H{"title": "cant get repo HEAD"})
		return
	}
	lastCommit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		slog.Error("cant get commit object", "err", err)
		c.JSON(http.StatusNotFound, gin.H{"title": "cant get commit object"})
		return
	}
	files, err := lastCommit.Files()
	if err != nil {
		slog.Error("cant get commit files", "err", err)
		c.JSON(http.StatusNotFound, gin.H{"title": "cant get commit files"})
		return
	}

	// Build file tree
	paths := []string{}
	if err := files.ForEach(func(f *object.File) error {
		paths = append(paths, f.Name)
		return nil
	}); err != nil {
		slog.Error("repo files", "err", err)
		c.JSON(http.StatusNotFound, gin.H{"title": repoPath + " not found"})
		return
	}
	root := buildFileTree(paths)

	c.JSON(http.StatusOK, gin.H{
		"hash":    lastCommit.Hash.String(),
		"message": lastCommit.Message,
		"author":  lastCommit.Author.Name,
		"email":   lastCommit.Author.Email,
		"date":    lastCommit.Author.When,
		"title":   "Files",
		"repo":    model.Repository{Path: repoPath, User: &model.User{Name: user}},
		"files":   root.Items,
	})
}
