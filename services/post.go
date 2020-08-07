package services

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var POST_SERVICE_URL = os.Getenv("POST_SERVICE_URL")

type PostService interface {
	GetAll() error
}

type postService struct {
	Results []string `json:"results"`
}

func NewPostService() *postService {
	return &postService{}
}

func (post *postService) GetAll() error {
	tracer := opentracing.GlobalTracer()
	childSpan := tracer.StartSpan(
		"get all post",
	)
	defer childSpan.Finish()
	endpoint := POST_SERVICE_URL + "/allpost"
	req, _ := http.NewRequest("GET", endpoint, nil)

	ext.SpanKindRPCClient.Set(childSpan)
	ext.HTTPUrl.Set(childSpan, endpoint)
	ext.HTTPMethod.Set(childSpan, "GET")

	tracer.Inject(childSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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
