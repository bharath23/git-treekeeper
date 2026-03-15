package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/spf13/cobra"
)

var version = "dev"

var templateOnce sync.Once

const helpTemplate = `{{- with .Long}}{{- .}}

{{end -}}Usage:
  {{displayUseLine .}}

{{- if .HasAvailableSubCommands}}

Commands:
{{- range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{- end}}
{{- if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{- end}}
{{- if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{- end}}
`

func displayRootName() string {
	base := filepath.Base(os.Args[0])
	if strings.HasPrefix(base, "git-") {
		suffix := strings.TrimPrefix(base, "git-")
		if suffix == "" {
			return "git"
		}
		return "git " + suffix
	}
	return base
}

func displayCommandPath(cmd *cobra.Command) string {
	path := cmd.CommandPath()
	if path == "" {
		return displayRootName()
	}
	parts := strings.SplitN(path, " ", 2)
	if len(parts) == 1 {
		return displayRootName()
	}
	return displayRootName() + " " + parts[1]
}

func displayUseLine(cmd *cobra.Command) string {
	use := strings.Fields(cmd.Use)
	if len(use) <= 1 {
		return displayCommandPath(cmd)
	}
	return displayCommandPath(cmd) + " " + strings.Join(use[1:], " ")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func NewRootCmd() *cobra.Command {
	var quiet bool
	var verbose bool

	templateOnce.Do(func() {
		cobra.AddTemplateFunc("displayUseLine", displayUseLine)
	})

	root := &cobra.Command{
		Use:   "git-tk <command>",
		Short: "Manage Git branches and worktrees efficiently.",
		Long: `Git TreeKeeper

Manage Git branches and worktrees efficiently.

Git TreeKeeper automatically creates and manages worktrees when
creating or switching branches, making it easier to work with
large repositories and multiple branches simultaneously.`,
		DisableSuggestions: true,
		SilenceErrors:      true,
		SilenceUsage:       true,
		Version:            version,
	}

	root.CompletionOptions.DisableDefaultCmd = true
	root.AddCommand(NewBranchCmd())
	root.AddCommand(NewCheckoutCmd())
	root.AddCommand(NewCloneCmd())
	root.AddCommand(NewListCmd())
	root.AddCommand(NewDoctorCmd())
	root.PersistentFlags().BoolVar(&quiet, "quiet", false, "Suppress informational output")
	root.PersistentFlags().BoolVar(&verbose, "verbose", false, "Verbose output")
	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		treekeeper.SetOutputMode(quiet, verbose)
		return nil
	}
	root.SetHelpTemplate(helpTemplate)
	root.SetUsageTemplate(helpTemplate)
	return root
}

var RootCmd = NewRootCmd()
