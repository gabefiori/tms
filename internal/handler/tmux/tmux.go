package tmux

import (
	"path/filepath"
	"strings"

	"github.com/gabefiori/gotmux"
)

func Run(target string) error {
	sessionName := strings.TrimPrefix(filepath.Base(target), ".")

	if gotmux.HasSession(sessionName) {
		if err := gotmux.AttachOrSwitchTo(sessionName); err != nil {
			return err
		}

		return nil
	}

	session, err := gotmux.NewSession(&gotmux.SessionConfig{
		Name: sessionName,
		Dir:  target,
	})

	if err != nil {
		return err
	}

	if err := session.AttachOrSwitch(); err != nil {
		return err
	}

	return nil
}
