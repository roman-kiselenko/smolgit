package route

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"smolgit/pkg/model"

	"github.com/gin-gonic/gin"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var ErrLimitReached = fmt.Errorf("limit reached")

// Commits returns the last 50 commits for a given repository.
// GET /api/repos/commits/:user/:path
func (r *Route) Commits(c *gin.Context) {
	user, repoPath := c.Param("user"), c.Param("path")
	fullPath := "/" + user + "/" + repoPath
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
	commitIter, err := repo.Log(&gogit.LogOptions{From: headRef.Hash()})
	if err != nil {
		slog.Error("cant get commit iterator", "err", err)
		c.JSON(http.StatusNotFound, gin.H{"title": "cant get commit iterator"})
		return
	}
	type commitItem struct {
		Hash    string `json:"hash"`
		Message string `json:"message"`
		Author  string `json:"author"`
		Email   string `json:"email"`
		Date    string `json:"date"`
	}
	commits := []commitItem{}
	count := 0
	if err := commitIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, commitItem{
			Hash:    c.Hash.String()[:7],
			Message: c.Message,
			Author:  c.Author.Name,
			Email:   c.Author.Email,
			Date:    c.Author.When.Format(time.RFC3339),
		})
		count++
		if count >= 10 {
			return ErrLimitReached // Break the iteration
		}
		return nil
	}); err != nil && errors.Is(err, ErrLimitReached) {
		c.JSON(http.StatusOK, gin.H{
			"title":   "Commits",
			"repo":    model.Repository{Path: repoPath, User: &model.User{Name: user}},
			"commits": commits,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"title":   "Commits",
		"repo":    model.Repository{Path: repoPath, User: &model.User{Name: user}},
		"commits": commits,
	})
}
