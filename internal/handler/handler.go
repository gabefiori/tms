package handler

import (
	"bytes"
	"io"
	"os"

	"github.com/gabefiori/tms/internal/config"
	"github.com/gabefiori/tms/internal/handler/tmux"
	"github.com/gabefiori/tms/internal/selector"
	"github.com/gabefiori/tms/internal/targets"
	hd "github.com/mitchellh/go-homedir"
)

func Run(cfg *config.Config) error {
	tg, err := targets.Collect(cfg.Targets)

	if err != nil {
		return err
	}

	if cfg.List {
		if cfg.Filter != "" {
			tg.Filter(cfg.Filter)
		}

		return tg.Print(os.Stdout)
	}

	if cfg.Filter != "" {
		// In this case, we delegate the responsibility of filtering to the selector.
		// This way, we avoid losing any targets.
		cfg.Selector = append(cfg.Selector, "--query="+cfg.Filter)
	}

	selected, err := selector.Run(tg.Buf, cfg.Selector)

	if err != nil {
		return err
	}

	if selected == "" {
		return nil
	}

	if cfg.OutputTarget {
		return toStdout(selected)
	}

	return tmux.Run(selected)
}

func RunSingle(target string) error {
	if err := targets.FindSingle(target); err != nil {
		return err
	}

	return tmux.Run(target)
}

func toStdout(target string) error {
	expanded, err := hd.Expand(target)

	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, bytes.NewBufferString(expanded))

	return err
}
