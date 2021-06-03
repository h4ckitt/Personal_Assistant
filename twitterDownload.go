package main

import (
	"encoding/json"
	"fmt"
	"github.com/yoruba-codigy/goTelegram"
	"log"
	"net/http"
	"strconv"
)

func getVideo(update goTelegram.Update) {
	var data data
	i := 0

	id := urlReg.FindStringSubmatch(update.Message.Text)

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

	choice[update.Message.From.ID] = &replyData{Data: data}
	for _, elem := range data.ExtendedEntities.Media[0].VideoInfo.Variants {
		if elem.Bitrate != 0 {
			resolution := qualityReg.FindString(elem.URL)
			choice[update.Message.From.ID].Qualities = append(choice[update.Message.From.ID].Qualities, elem.URL)
			b.AddButton("▶️ "+resolution, "tw-"+strconv.Itoa(i))
			i++
		}
	}

	b.MakeKeyboard(1)

	err = b.ReplyMessage("⬇️ Choose A Quality", update.Message)

	if err != nil {
		log.Println(err)
	}
}

func sendVideo(update goTelegram.Update) {
	data := choice[update.CallbackQuery.From.ID].Data
	index, _ := strconv.Atoi(update.CallbackQuery.Data)
	link := choice[update.CallbackQuery.From.ID].Qualities[index]
	defer delete(choice, update.CallbackQuery.From.ID)

	text := fmt.Sprintf("👤Name: %s\n❤️Likes: %s\n🔁Retweets: %s\n", data.User.Name, strconv.Itoa(data.FavoriteCount), strconv.Itoa(data.RetweetCount))

	b.DeleteMessage(update.CallbackQuery.Message)

	b.SendVideo(link, text, update.CallbackQuery.Message.Chat)
}
