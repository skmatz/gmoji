package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/integralist/go-findroot/find"
)

const (
	dirName  string = "gmoji"
	fileName string = "gmojis.json"
)

// CLI represents this application itself.
type CLI struct {
	// CheckGit checks if a git repository exists.
	CheckGit bool
	// Clipboard doesn't commit, but copies the commit message to the clipboard.
	Clipboard bool
	// HookPath is the path to the file to write the commit message to.
	HookPath string
}

// NewCLI returns a new CLI.
func NewCLI(checkGit bool, clipboard bool, hook string) (*CLI, error) {
	c := &CLI{CheckGit: checkGit, Clipboard: clipboard, HookPath: hook}
	if c.CheckGit && hook == "" {
		p, err := c.hookPath()
		if err != nil {
			return nil, err
		}
		c.HookPath = p
	}
	return c, nil
}

// GitRoot returns the path to the git root directory.
func (c CLI) GitRoot() (string, error) {
	root, err := find.Repo()
	if err != nil {
		return "", fmt.Errorf("not a git repository (or any of the parent directories)")
	}

	return root.Path, nil
}

// ListPath returns the path to the gitmojis.json file.
func (c CLI) ListPath() (string, error) {
	cfg, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(cfg, dirName, fileName), nil
}

// Wget downloads a file from the given URL.
func (c CLI) Wget(src, dst string) (err error) {
	response, err := http.Get(src)
	if err != nil {
		return err
	}
	defer func() {
		cerr := response.Body.Close()
		if err == nil {
			err = cerr
		}
	}()

	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return err
}

func (c CLI) commitMessage(code, title, message string) string {
	if message != "" {
		return fmt.Sprintf("%s %s\n\n%s\n", code, title, message)
	}
	return fmt.Sprintf("%s %s\n", code, title)
}
