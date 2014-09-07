package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/MariaTerzieva/gotumblr"
	"github.com/kennygrant/sanitize"
)

func main() {

	blogname := flag.String("blogname", "", "the blogname like example.tumblr.com")
	flag.Parse()

	client := gotumblr.NewTumblrRestClient(
		"consumerKey",
		"consumerSecret",
		"oauthToken",
		"oauthSecret",
		"callbackUrl",
		"http://api.tumblr.com")

	response := client.Posts(
		*blogname,
		"text",
		map[string]string{})

	fmt.Println("total text posts:", response.Total_posts)

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
				*blogname,
				"text",
				map[string]string{
					"offset": strconv.FormatInt(i, 10)})
		} else {
			break
		}
	}

	for _, post := range alltheposts {

		quote := strings.Trim(sanitize.HTML(post.Body), "\n")

		if err := client.CreateQuote(
			*blogname,
			map[string]string{
				"quote": quote,
				"date":  post.Date,
			},
		); err == nil {

			if err := client.DeletePost(
				*blogname,
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
