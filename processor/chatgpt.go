package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/xatta-trone/words-combinator/database"
	"github.com/xatta-trone/words-combinator/model"
	"github.com/xatta-trone/words-combinator/utils"
)

func GetChatGpt() {

	words := []model.WordGetStruct{}
	database.Gdb.Select(&words, "SELECT id, word from wordlist where is_parsed_gpt=0")

	for _, word := range words {
		buildChatGpt(word.ID, word.Word)
	}

}

func buildChatGpt(wordId int64, word string) {
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Ask to chat gpt: ")
	// // text, err := reader.ReadString('\n')
	// if err != nil {
	// 	fmt.Println("Error reading input:", err)
	// 	return
	// }

	fmt.Printf("Sending question to openaiAPI... with word %s \n", word)

	c := openai.NewClient("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")

	resp, err := c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: fmt.Sprintf("give me definition, example and synonyms of %s", word),
				},
			},
		},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Choices[0].Message.Content)

	content := resp.Choices[0].Message.Content

	Def := strings.Index(content, "Definition:")
	Example := strings.Index(content, "Example:")
	Syn := strings.Index(content, "Synonyms:")

	ExampleList := strings.Split(resp.Choices[0].Message.Content[Def:Example], ":")
	DefList := strings.Split(resp.Choices[0].Message.Content[Example:Syn], ":")
	SynList := strings.Split(resp.Choices[0].Message.Content[Syn:], ":")

	// strings.TrimSpace(ExampleList[len(ExampleList)-1])

	var chatResponse model.ChatResp

	chatResponse.Definition = strings.TrimSpace(DefList[len(DefList)-1])
	chatResponse.Example = strings.TrimSpace(ExampleList[len(ExampleList)-1])

	// synonyms
	lists := strings.Split(strings.TrimSpace(SynList[len(SynList)-1]), ",")

	s := []string{}

	reg, err := regexp.Compile(`[^a-zA-Z\s]+`)
	if err != nil {
		log.Fatal(err)
	}

	for _, a := range lists {
		s = append(s, reg.ReplaceAllString(strings.Replace(a, ".", "", -1), ""))
	}

	chatResponse.Synonyms = s

	data, _ := json.Marshal(chatResponse)

	_, err = database.Gdb.Exec("Update wordlist set gpt=?,is_parsed_gpt=1 where id = ? ", data, wordId)

	if err != nil {
		fmt.Println(err)
	}

	// fmt.Printf("Inserted %v - %s from google \n", word.ID, word.Word)
	str := fmt.Sprintf("Inserted %v - %s from gpt \n", wordId, word)
	utils.PrintG(str)

	//  resp.Choices[0].Message.Content[Example:Syn], resp.Choices[0].Message.Content[Syn:]

	fmt.Println(chatResponse)
}
