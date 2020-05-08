package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var POST_SERVICE_URL = os.Getenv("POST_SERVICE_URL")

type PostService struct {
	Results []string `json:"results"`
}

func NewPostService() *PostService {
	return &PostService{}
}

func (post *PostService) GetAll() error {
	resp, err := http.Get(POST_SERVICE_URL + "/allpost")
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("body = ", string(body))
		json_err := json.Unmarshal(body, &post)
		fmt.Println("post = ", post)
		if json_err != nil {
			return json_err
		}
		return nil
	} else {
		panic(resp.StatusCode)
	}

}
