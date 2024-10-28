package cli

import (
	"os"
	"runtime/debug"

	"github.com/gabefiori/tms/internal/config"
	"github.com/gabefiori/tms/internal/handler"
	"github.com/urfave/cli/v2"
)

func Run() error {
	var (
		path         string
		filter       string
		target       string
		list         bool
		outputTarget bool
	)

	app := &cli.App{
		Name:    "tms",
		Usage:   "Tmux Sessionizer is a tool for navigating through folders and projects as tmux sessions.",
		Version: getVersion(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Load configuration from `file`",
				Value:       "~/.config/tms/config.json",
				Destination: &path,
			},
			&cli.StringFlag{
				Name:        "filter",
				Aliases:     []string{"f"},
				Usage:       "Specify a filter to narrow down the results displayed in the selector",
				Value:       "",
				Destination: &filter,
			},
			&cli.StringFlag{
				Name:        "target",
				Aliases:     []string{"t"},
				Usage:       "Specify a target (e.g., path) to switch or attach to",
				Value:       "",
				Destination: &target,
			},
			&cli.BoolFlag{
				Name:        "list",
				Aliases:     []string{"l"},
				Usage:       "List of all discovered targets",
				Value:       false,
				Destination: &list,
			},
			&cli.BoolFlag{
				Name:        "output-target",
				Aliases:     []string{"ot"},
				Usage:       "Output the selected target directory for use in command substitution (e.g., cd \"$(tms -ot)\").",
				Value:       false,
				Destination: &outputTarget,
			},
		},
		Action: func(ctx *cli.Context) error {
			if target != "" {
				return handler.RunSingle(target)
			}

			cliCfg := config.Cli{
				Path:         path,
				Filter:       filter,
				List:         list,
				OutputTarget: outputTarget,
			}

			cfg, err := config.Load(cliCfg)

			if err != nil {
				return err
			}

			return handler.Run(cfg)
		},
	}

	if err := app.Run(os.Args); err != nil {
		return err
	}

	return nil
}

func getVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	}

	return "unknown"
}
