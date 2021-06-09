package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type GetUnsplashPhotosInput struct {
	Query   string `json:"query"`
	Page    int    `json:"page"`
	PerPage int    `json:"perPage"`
}

type GetUnsplashPhotosOutput struct {
	MediaId           string `json:"mediaId"`
	Url               string `json:"url"`
	Kind              string `json:"kind"`
	MimeType          string `json:"mimeType"`
	Name              string `json:"name"`
	Size              int    `json:"size"`
	AuthorName        string `json:"authorName"`
	AuthorProfileLink string `json:"authorProfileLink"`
	DownloadLocation  string `json:"downloadLocation"`
	ThumbUrl          string `json:"thumbUrl"`
}

func main() {
	app := fiber.New()

	app.Post("/searchImages", func(c *fiber.Ctx) error {
		godotenv.Load()
		getUnsplashPhotosInput := new(GetUnsplashPhotosInput)

		if err := c.BodyParser(getUnsplashPhotosInput); err != nil {
			return err
		}

		client := resty.New()
		res, _ := client.R().SetQueryParams(map[string]string{
			"query":    getUnsplashPhotosInput.Query,
			"page":     strconv.Itoa(getUnsplashPhotosInput.Page),
			"per_page": strconv.Itoa(getUnsplashPhotosInput.PerPage),
		}).
			SetHeader("Accept", "application/json").
			SetHeader("Authorization", fmt.Sprintf("Client-ID %s", os.Getenv("UNSPLASH_ACCESS_TOKEN"))).
			Get("https://api.unsplash.com/search/photos")

		var result map[string]interface{}
		json.Unmarshal(res.Body(), &result)
		photosArray := result["results"].([]interface{})

		var output []GetUnsplashPhotosOutput

		for _, photo := range photosArray {
			p := photo.(map[string]interface{})

			photoRes := GetUnsplashPhotosOutput{
				MediaId:           p["id"].(string),
				Url:               p["urls"].(map[string]interface{})["regular"].(string),
				Kind:              "image",
				MimeType:          "image/jpeg",
				Name:              p["alt_description"].(string),
				Size:              0,
				AuthorName:        p["user"].(map[string]interface{})["name"].(string),
				AuthorProfileLink: fmt.Sprintf("%s?utm_source=lyearn&utm_medium=referral", p["user"].(map[string]interface{})["links"].(map[string]interface{})["html"]),
				DownloadLocation:  p["links"].(map[string]interface{})["download_location"].(string),
				ThumbUrl:          p["urls"].(map[string]interface{})["thumb"].(string),
			}

			output = append(output, photoRes)
		}

		return c.JSON(output)
	})

	app.Listen(":3021")
}
