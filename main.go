package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vinhut/feed-service/services"

	"encoding/json"
	"fmt"
)

var SERVICE_NAME = "feed-service"

func setupRouter(authservice services.AuthService, postservice *services.PostService) *gin.Engine {

	router := gin.Default()

	router.GET(SERVICE_NAME+"/ping", func(c *gin.Context) {
		c.String(200, "OK")
	})

	router.GET(SERVICE_NAME+"/feed", func(c *gin.Context) {
		value, err := c.Cookie("token")
		if err != nil {
			panic("failed get token")
		}
		user_data, auth_error := authservice.Check(SERVICE_NAME, value)
		if auth_error != nil {
			panic(auth_error)
		}
		var raw struct {
			Uid     string
			Email   string
			Role    string
			Created string
		}
		if err := json.Unmarshal([]byte(user_data), &raw); err != nil {
			fmt.Println(err)
			panic(err)
		}

		if raw.Email == "" {
			c.String(403, "")
		}

		post_err := postservice.GetAll()
		if post_err != nil {
			fmt.Println(post_err)
			panic(post_err)
		}
		out, json_err := json.Marshal(postservice)
		if json_err != nil {
			fmt.Println(json_err)
			panic(json_err)
		}

		c.String(200, string(out))

	})

	return router
}

func main() {

	authservice := services.NewUserAuthService()
	postservice := services.NewPostService()
	router := setupRouter(authservice, postservice)
	router.Run(":8080")

}
