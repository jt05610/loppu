/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// watchdogCmd represents the watchdog command
var watchdogCmd = &cobra.Command{
	Use:   "watchdog",
	Short: "watchdog watches values stream data",
	Long:  `watchdog currently watches for values to cross a value`,
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run watches values stream data",
	Long:  `run currently watches for values to cross a value`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize the watchdog",
	Long:  `run init to generate the default watchdog structure`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(watchdogCmd)
	watchdogCmd.AddCommand(runCmd)
	watchdogCmd.AddCommand(initCmd)

}
