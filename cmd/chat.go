/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/shoet/gpt-chat/clocker"
	"github.com/shoet/gpt-chat/service"
	"github.com/shoet/gpt-chat/store"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start ChatGPT",
	Long: `Start ChatGPT on ChatGPTAPI. This program is a CLI tool to start ChatGPT on ChatGPTAPI. 
	Chat content is stored database.`,
	Run: func(cmd *cobra.Command, args []string) {
		msg, err := cmd.Flags().GetString("message")
		if err != nil {
			fmt.Println("failed to get message")
			os.Exit(1)
		}

		gpt := buildChatGPTService()
		repo := buildRepository()
		db, closer := buildDB()
		defer closer()

		storage, err := service.NewStorageRDB(db, repo)
		chatSrv, err := service.NewChatService(gpt, storage)
		if err := chatSrv.Chat("", msg); err != nil {
			fmt.Println("failed to chat")
			os.Exit(1)
		}

	},
}

func buildChatGPTService() *service.ChatGPTService {
	apikey, exist := os.LookupEnv("CHATGPT_API_SECRET")
	if !exist {
		fmt.Println("OPENAI_API_KEY is not set")
		os.Exit(1)
	}
	gpt := service.NewChatGPTService(apikey, &http.Client{})
	return gpt

}

func buildRepository() *store.Repository {
	repo, err := store.NewRepository(&clocker.RealClocker{})
	if err != nil {
		fmt.Println("failed to create repository")
		os.Exit(1)
	}
	return repo

}

func buildDB() (*sqlx.DB, func() error) {
	c := mysql.Config{
		DBName:               "gpt",
		User:                 "gpt",
		Passwd:               "gpt",
		Addr:                 "localhost:33306",
		Net:                  "tcp",
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	db, closer, err := store.NewRDB(c.FormatDSN())
	if err != nil {
		fmt.Println("failed to connect database")
		os.Exit(1)
	}
	return db, closer
}

func init() {
	chatCmd.Flags().StringP("message", "m", "", "chat message")
	rootCmd.AddCommand(chatCmd)
}
