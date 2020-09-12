package main

import (
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	transport "github.com/uber/jaeger-client-go/transport/zipkin"
	"github.com/uber/jaeger-client-go/zipkin"
	"github.com/uber/jaeger-lib/metrics"
	"github.com/vinhut/feed-service/services"

	"encoding/json"
	"os"
)

var SERVICE_NAME = "feed-service"

func setupRouter(authservice services.AuthService, postservice services.PostService) *gin.Engine {

	var JAEGER_COLLECTOR_ENDPOINT = os.Getenv("JAEGER_COLLECTOR_ENDPOINT")
	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	trsport, _ := transport.NewHTTPTransport(
		JAEGER_COLLECTOR_ENDPOINT,
		transport.HTTPLogger(jaeger.StdLogger),
	)
	cfg := jaegercfg.Configuration{
		ServiceName: "feed-service",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: JAEGER_COLLECTOR_ENDPOINT,
		},
	}
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory
	cfg.InitGlobalTracer(
		"feed-service",
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
		jaegercfg.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		jaegercfg.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
		jaegercfg.ZipkinSharedRPCSpan(true),
		jaegercfg.Reporter(jaeger.NewRemoteReporter(trsport)),
	)
	tracer := opentracing.GlobalTracer()

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "OK")
	})

	router.GET(SERVICE_NAME+"/feed", func(c *gin.Context) {
		span := tracer.StartSpan("get feed")

		feed_range, query_exist := c.GetQuery("range")
		if query_exist == false {
			feed_range = "8"
		}
		value, token_err := c.Cookie("token")
		if token_err != nil {
			span.Finish()
			c.AbortWithStatusJSON(401, gin.H{"reason": "unauthorized"})
			return
		}
		user_data, auth_error := authservice.Check(SERVICE_NAME, value)
		if auth_error != nil {
			span.Finish()
			c.AbortWithStatusJSON(401, gin.H{"reason": "unauthorized"})
			return
		}
		var raw struct {
			Uid     string
			Email   string
			Role    string
			Created string
		}
		if json_err := json.Unmarshal([]byte(user_data), &raw); json_err != nil {
			span.Finish()
			panic(json_err.Error())
		}

		if raw.Email == "" {
			span.Finish()
			c.AbortWithStatusJSON(403, gin.H{"reason": "unauthorized"})
			return
		}

		post_err := postservice.GetAll(feed_range)

		if post_err != nil {
			span.Finish()
			c.AbortWithStatusJSON(502, gin.H{"reason": "post service error"})
			return
		}
		out, marshal_err := json.Marshal(postservice)
		if marshal_err != nil {
			span.Finish()
			panic(marshal_err.Error())
		} else {
			c.String(200, string(out))
			span.Finish()
		}

	})

	return router
}

func main() {

	authservice := services.NewUserAuthService()
	postservice := services.NewPostService()
	router := setupRouter(authservice, postservice)
	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}

}
