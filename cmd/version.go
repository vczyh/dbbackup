package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime/debug"
)

var (
	version string
	commit  = func() string {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					return setting.Value
				}
			}
		}
		return ""
	}()
)

// represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Build version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Git Commit: %s\n", commit)
		fmt.Printf("Github: https://github.com/vczyh/dbbackup\n")
	},
}
