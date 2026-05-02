package cmd

import (
	"github.com/spf13/cobra"
)

type GlobalOptions struct {
	File      string
	Config    string
	Portfolio string
	Account   string
	Offline   bool
	Refresh   bool
	NoColor   bool
}

var opts GlobalOptions

var rootCmd = &cobra.Command{
	Use:   "stkq",
	Short: "Plain-text stock portfolio tracker and stock research for engineers who invest",
	Long: `stkq keeps your portfolio in a readable text file and lets you
query holdings, estimated value, SEC filings, and structured company facts.`,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&opts.File, "file", "f", "", "portfolio file path")
	rootCmd.PersistentFlags().StringVar(&opts.Config, "config", "", "config file path")
	rootCmd.PersistentFlags().StringVar(&opts.Portfolio, "portfolio", "", "portfolio filter")
	rootCmd.PersistentFlags().StringVar(&opts.Account, "account", "", "account filter")
	rootCmd.PersistentFlags().BoolVar(&opts.Offline, "offline", false, "use only local files/cache")
	rootCmd.PersistentFlags().BoolVar(&opts.Refresh, "refresh", false, "refresh cached data")
	rootCmd.PersistentFlags().BoolVar(&opts.NoColor, "no-color", false, "disable color output")
}

func Execute() error {
	return rootCmd.Execute()
}
