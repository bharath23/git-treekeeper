package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "git tk <command>",
	Short: "Manage Git branches and worktrees",
	Long: `Git TreeKeeper

Manage Git branches and worktrees efficiently.

Git TreeKeeper automatically creates and manages worktrees when
creating or switching branches, making it easier to work with
large repositories and multiple branches simultaneously.`,
	DisableSuggestions: true,
	SilenceErrors:      true,
	SilenceUsage:       true,
	Version:            "0.1.0",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.AddCommand(NewBranchCmd())
	RootCmd.AddCommand(NewCheckoutCmd())
	RootCmd.AddCommand(NewCloneCmd())
	RootCmd.SetHelpTemplate(`{{with .Long}}{{.}}{{end}}

Usage:
  {{.UseLine}}

{{if .HasAvailableSubCommands}}Commands:
{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}

{{if .HasAvailableLocalFlags}}Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

{{if .HasAvailableInheritedFlags}}Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
`)
}
