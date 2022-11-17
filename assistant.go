package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/yoruba-codigy/goTelegram"
)

var (
	b          goTelegram.Bot
	qualityReg *regexp.Regexp
	urlReg     *regexp.Regexp
	choice     map[int]*replyData
	client     *http.Client
)

func main() {
	var err error

	qualityReg, _ = regexp.Compile("[0-9]+x[0-9]+")
	urlReg = regexp.MustCompile(`https:\/\/twitter.com\/\w+\/status\/([0-9]+)(?:\?[a-z]+=[0-9A-Za-z_]+(?:\&[a-z]+=[0-9A-Za-z_-]+)*)*`)

	choice = make(map[int]*replyData)

	client = &http.Client{}

	b, err = goTelegram.NewBot("911112005:AAFSM5o7Cb1dmvETAqwwk496tp8RPSxpdjQ")

	if err != nil {
		log.Println("Couldn't Create Bot Successfully, Check Logs For More Details")
		log.Fatal(err)
	}

	b.SetHandler(handler)

	fmt.Println("Starting Server")

	err = http.ListenAndServe(":"+os.Getenv("PORT"), http.HandlerFunc(b.UpdateHandler))
}

func handler(update goTelegram.Update) {
	switch update.Type {
	case "text":
		processText(update)

	case "callback":
		defer b.AnswerCallback(update.CallbackQuery.ID)
		processCallback(update)
	}
}

func processText(update goTelegram.Update) {
	defer b.DeleteKeyboard()

	//texts := strings.Fields(update.Message.Text)

	if urlReg.MatchString(update.Message.Text) {
		getVideo(update)
		return
	}

	switch update.Command {
	case "/start":
		b.SendMessage("Hello "+update.Message.From.Firstname, update.Message.Chat)
	}
}

func processCallback(update goTelegram.Update) {

	if strings.HasPrefix(update.CallbackQuery.Data, "tw") {
		if _, ok := choice[update.CallbackQuery.From.ID]; ok {
			sendVideo(update)
		}
	}
}
