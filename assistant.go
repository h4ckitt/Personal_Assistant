package main

import (
	"encoding/json"
	"fmt"
	"github.com/yoruba-codigy/goTelegram/v2"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	//"strings"
)

var (
	b          goTelegram.Bot
	qualityReg *regexp.Regexp
	urlReg     *regexp.Regexp
	choice     map[int][]string
	client     *http.Client
)

func main() {
	var err error

	qualityReg, _ = regexp.Compile("[0-9]+x[0-9]+")
	urlReg = regexp.MustCompile(`https:\/\/twitter.com\/\w+\/status\/([0-9]+)\?s=[0-9]+`)

	choice = make(map[int][]string)

	client = &http.Client{}

	b, err = goTelegram.NewBot("911112005:AAFSM5o7Cb1dmvETAqwwk496tp8RPSxpdjQ")

	if err != nil {
		log.Println("Couldn't Create Bot Successfully, Check Logs For More Details")
		log.Panic(err)
	}

	b.SetHandler(handler)

	fmt.Println("Starting Server")

	err = http.ListenAndServe(":8080", http.HandlerFunc(b.UpdateHandler))
}

func handler(update goTelegram.Update) {
	switch update.Type {
	case "text":
		processText(update)

	case "callback":
		defer b.AnswerCallback(update.CallbackQuery.ID)
		if _, ok := choice[update.CallbackQuery.From.ID]; ok {
			processCallback(update)

		}
	}
}

func processText(update goTelegram.Update) {
	defer b.DeleteKeyboard()
	var data data
	i := 0

	//text := strings.Fields(update.Message.Text)[0]

	id := urlReg.FindStringSubmatch(update.Message.Text)

	fmt.Println(id)

	if len(id) < 2 {
		err := b.SendMessage("Not A Valid Twitter Status URL", update.Message.Chat)

		if err != nil {
			log.Println(err)
		}

		return
	}

	req, _ := http.NewRequest("GET", "https://api.twitter.com/1.1/statuses/show.json?id="+id[1], nil)
	req.Header.Set("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAAE3eQAEAAAAAn26iOFtkEaizxhUXNyi32mR8K6Q%3DUkTMU4C1ARdMt9M6xTH7h9sz50QUuO4o0n8X2nlVcuJLNHnLtV")
	res, _ := client.Do(req)

	err := json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		log.Println("There Was An Error Marshaling Response")
		return
	}

	if data.ExtendedEntities.Media[0].Type != "video" {
		err := b.SendMessage("No Video Found In Provided URL", update.Message.Chat)

		if err != nil {
			log.Println(err)
		}

		return
	}

	for _, elem := range data.ExtendedEntities.Media[0].VideoInfo.Variants {
		if elem.Bitrate != 0 {
			resolution := qualityReg.FindString(elem.URL)
			choice[update.Message.From.ID] = append(choice[update.Message.From.ID], elem.URL)
			b.AddButton(resolution, strconv.Itoa(i))
			i++
		}
	}

	b.MakeKeyboard(1)
	fmt.Println(update.Message.Chat)
	err = b.SendMessage("Choose Quality", update.Message.Chat)

	if err != nil {
		log.Println(err)
	}

}

func processCallback(update goTelegram.Update) {
	index, _ := strconv.Atoi(update.CallbackQuery.Data)
	link := choice[update.CallbackQuery.From.ID][index]

	fmt.Println(link)

	resp, err := http.Get(link)

	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	out, err := os.Create("TwitterVideo.mp4")

	if err != nil {
		log.Println(err)
		return
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		log.Println(err)
	}

}
