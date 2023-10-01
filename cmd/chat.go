/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start ChatGPT",
	Long: `Start ChatGPT on ChatGPTAPI. This program is a CLI tool to start ChatGPT on ChatGPTAPI. 
	Chat content is stored database.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("chat called")
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
