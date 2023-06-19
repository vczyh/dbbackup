/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dbbackup",
	Short: "A database backup tool",
}

func Execute() {
	fmt.Println()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(mysqlCmd)
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
