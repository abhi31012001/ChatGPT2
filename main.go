package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type NullWriter struct{}

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func main() {
	log.SetOutput(new(NullWriter))
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	api := viper.GetString("API_KEY")
	if api == "" {
		log.Fatalln("Api Key is not present")
	}
	ctx := context.Background()
	client := gpt3.NewClient(api)
	rootcmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with chatGPT in console.",
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)
			quit := false
			for !quit {
				fmt.Print("Write quit to exit or else enjoy: ")
				if !scanner.Scan() {
					break
				}
				question := scanner.Text()
				switch question {
				case "quit":
					quit = true
				default:
					GetResponse(client, ctx, question)
				}
			}
		},
	}
	rootcmd.Execute()
}

func GetResponse(client gpt3.Client, ctx context.Context, question string) {
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			question,
		},
		MaxTokens:   gpt3.IntPtr(300),
		Temperature: gpt3.Float32Ptr(0),
	}, func(cr *gpt3.CompletionResponse) {
		fmt.Print(cr.Choices[0].Text)
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print("\n")

}
