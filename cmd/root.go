/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mjehanno/new/internal/pipe"
	"github.com/mjehanno/new/internal/tui"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "new",
	Short: "new is a custom wrapper around cookiecutter to start new projects",
	Long: `new is a custom wrapper around cookiecutter made to bootstrap new project. 
CookieCutter gives you the possibility to create template of projects. 
new will give you the possibility to fetch directly template you've created and pushed on Github.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		errs := make([]error, 0)
		var interactiveMode bool
		if len(args) > 1 {
			err = errors.New("too much args provided, you only need to pass one")
			errs = append(errs, err)
		}

		interactiveMode, err = cmd.Flags().GetBool("interactive")
		if err != nil {
			err = fmt.Errorf("error while parsing flags : %s", err.Error())
			errs = append(errs, err)
		}

		if interactiveMode && len(args) > 0 {
			err = errors.New("interactive mode is not compatible with args")
			errs = append(errs, err)
		}

		if len(errs) > 0 {
			for _, e := range errs {
				log.Println(e.Error())
			}
			os.Exit(1)
		}

		if interactiveMode {

			p := tea.NewProgram(tea.Model(tui.InitialModel()), tea.WithAltScreen())
			if err := p.Start(); err != nil {
				log.Println(fmt.Errorf("alas, there's been an error: %w", err).Error())
				os.Exit(1)
			}

			if url := <-pipe.Chan; url != "" {
				cmd := exec.Command("cookiecutter", url)
				cmd.Stderr = os.Stderr
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout

				err := cmd.Run()
				if err != nil {
					log.Println(fmt.Errorf("failed to run command : %w", err).Error())
				}
			}

		} else {
			cmd := exec.Command("cookiecutter", args[0])
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			err := cmd.Run()
			if err != nil {
				log.Println(fmt.Errorf("error while running cookiecutter : %w", err).Error())
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.new.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("interactive", "i", false, "Enable interactive mode, it will fetch your github cookiecutter template.")
}
