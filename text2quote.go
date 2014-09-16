package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/MariaTerzieva/gotumblr"
	"github.com/kennygrant/sanitize"
)

func readConfig(filename string) (map[string]string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var c map[string]string
	json.Unmarshal(content, &c)

	return c, nil
}

func writeBackup(filename string, data interface{}) error {

	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, d, 0666)
}

func main() {

	conf, err := readConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	blogname := conf["blogname"]

	client := gotumblr.NewTumblrRestClient(
		conf["consumer_key"],
		conf["consumer_secret"],
		conf["token"],
		conf["token_secret"],
		"http://localhost",
		"http://api.tumblr.com")

	response := client.Posts(
		blogname,
		"text",
		map[string]string{})

	fmt.Println("total text posts:", response.Total_posts)

	if response.Total_posts > 0 {

		var alltheposts []gotumblr.TextPost
		var i int64
		// fetch all the posts
		for {
			fmt.Printf("fetched posts from %d to %d\n", i, i+19)

			for _, jsonpost := range response.Posts {
				post := gotumblr.TextPost{}
				json.Unmarshal(jsonpost, &post)
				alltheposts = append(alltheposts, post)
			}

			if i < response.Total_posts-20 {
				i += 20
				response = client.Posts(
					blogname,
					"text",
					map[string]string{
						"offset": strconv.FormatInt(i, 10)})
			} else {
				break
			}
		}

		if alltheposts != nil {

			if err := writeBackup("./backup.json", alltheposts); err != nil {
				log.Fatal(err)
			}

			for _, post := range alltheposts {

				quote := strings.Trim(sanitize.HTML(post.Body), "\n")

				if err := client.CreateQuote(
					blogname,
					map[string]string{
						"quote": quote,
						"date":  post.Date,
					},
				); err == nil {

					if err := client.DeletePost(
						blogname,
						strconv.FormatInt(post.Id, 10),
					); err != nil {
						log.Fatal(err)
					}

				} else {
					log.Fatal(err)
				}

				fmt.Println("created new quote:", quote)

			}
		}
	}

}
