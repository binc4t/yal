// this is a trial for go-openai

package test1

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const proxyURL = "http://127.0.0.1:1087"

var Proxy *url.URL
var Client *openai.Client
var Token = os.Getenv("OPENAI_API_KEY")

func init() {
	var err error
	if Proxy, err = url.Parse(proxyURL); err != nil {
		log.Fatal(err)
	}

	config := openai.DefaultConfig(Token)
	config.HTTPClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(Proxy),
		},
	}
	Client = openai.NewClientWithConfig(config)
}

func dealErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func text() {
	resp, err := Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}

func speech() {
	b, err := os.ReadFile("/Users/yyymagic/proj/ai/a.txt")
	if err != nil {
		log.Fatal(err)
	}
	req := openai.CreateSpeechRequest{
		Model: openai.TTSModel1,
		Input: string(b),
		Voice: openai.VoiceNova,
	}

	resp, err := Client.CreateSpeech(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("out.mp3")
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(f, resp)
	if err != nil {
		log.Fatal(err)
	}
}

func genImage() {
	f, err := os.Open("/Users/yyymagic/proj/ai/aa.png")
	dealErr(err)
	req := openai.ImageEditRequest{
		Image:          f,
		Prompt:         "change the background to Disney Castle",
		Model:          openai.CreateImageModelDallE2,
		N:              1,
		Size:           openai.CreateImageSize1024x1024,
		ResponseFormat: openai.CreateImageResponseFormatURL,
	}
	resp, err := Client.CreateEditImage(context.Background(), req)
	dealErr(err)
	for _, v := range resp.Data {
		fmt.Println(v.URL)
	}
}

func assistant() {
	name := "yyy"
	instruction := "You are a personal math tutor. Write and run code to answer math questions."
	req := openai.AssistantRequest{
		Model:        openai.GPT3Dot5Turbo1106,
		Name:         &name,
		Instructions: &instruction,
		Tools: []openai.AssistantTool{{
			openai.AssistantToolTypeCodeInterpreter,
			nil,
		}},
	}
	assist, err := Client.CreateAssistant(context.Background(), req)
	dealErr(err)

	tr := openai.ThreadRequest{}
	thread, err := Client.CreateThread(context.Background(), tr)
	dealErr(err)

	mr := openai.MessageRequest{
		Role:    string(openai.ThreadMessageRoleUser),
		Content: "I need to solve the equation `3x + 11 = 14`. Can you help me?",
	}
	_, err = Client.CreateMessage(context.Background(), thread.ID, mr)
	dealErr(err)

	rr := openai.RunRequest{
		AssistantID: assist.ID,
	}
	r, err := Client.CreateRun(context.Background(), thread.ID, rr)
	dealErr(err)

	for {
		cr, err := Client.RetrieveRun(context.Background(), r.ThreadID, r.ID)
		dealErr(err)
		if cr.Status == openai.RunStatusCompleted {
			break
		}
		fmt.Println("status is: ", cr.Status)
		time.Sleep(time.Second)
	}

	ms, err := Client.ListMessage(context.Background(), thread.ID, nil, nil, nil, nil)
	dealErr(err)

	for _, cur := range ms.Messages {
		fmt.Println("-------------------------------------")
		fmt.Printf("%+v\n", cur)
		for _, c := range cur.Content {
			if c.Text != nil {
				fmt.Printf("%s, %s\n", c.Type, c.Text.Value)
			}
		}

	}
}

func do() {
	assistant()
}
