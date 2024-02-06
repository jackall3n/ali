package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name: "create",
		Action: func(ctx context.Context, command *cli.Command) error {
			name := command.Args().Get(0)
			cmd := command.Args().Get(1)

			fmt.Println("create", name, cmd)

			file, err := getConfigFile()

			if err != nil {
				return err
			}

			defer file.Close()

			if _, err := io.WriteString(file, fmt.Sprintf("alias %s=\"%s\"\n", name, cmd)); err != nil {
				return err
			}

			return nil
		},

		Commands: []*cli.Command{
			{
				Name: "init",
				Action: func(ctx context.Context, command *cli.Command) error {
					fmt.Println("init")

					file, err := getConfigFile()

					if err != nil {
						return err
					}

					defer file.Close()

					fmt.Println("add ")

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func getConfigFile() (*os.File, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".alirc")

	file, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, err
	}

	exists, err := fileExists(configPath)

	if err != nil {
		return nil, err
	}

	if !exists {
		_, err := io.WriteString(file, "")

		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
