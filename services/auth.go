package services

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var SERVICE_URL = os.Getenv("AUTH_SERVICE_URL")

type AuthService interface {
	Login(service string, email string, password string) (string, error)
	Check(service string, token string) (string, error)
	Update() (bool, error)
	Create(service string, email string, password string) (bool, error)
	Delete(string) (bool, error)
}

type userAuthService struct {
	token string
}

func NewUserAuthService() AuthService {
	return &userAuthService{
		token: "",
	}
}

func (userAuth *userAuthService) Login(service string, email string, password string) (string, error) {
	tracer := opentracing.GlobalTracer()
	childSpan := tracer.StartSpan(
		"userauth login",
	)
	defer childSpan.Finish()
	endpoint := SERVICE_URL + "/login"
	data := url.Values{"service": {service}, "email": {email}, "password": {password}}
	req, _ := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	ext.SpanKindRPCClient.Set(childSpan)
	ext.HTTPUrl.Set(childSpan, endpoint)
	ext.HTTPMethod.Set(childSpan, "POST")

	tracer.Inject(childSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body), nil
	} else {
		return "", errors.New("unauthorized: login error")
	}

}

func (userAuth *userAuthService) Check(service string, token string) (string, error) {
	tracer := opentracing.GlobalTracer()
	childSpan := tracer.StartSpan(
		"userauth check",
	)
	defer childSpan.Finish()
	endpoint := SERVICE_URL + "/user?service=" + service + "&token=" + token
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	ext.SpanKindRPCClient.Set(childSpan)
	ext.HTTPUrl.Set(childSpan, endpoint)
	ext.HTTPMethod.Set(childSpan, "GET")

	tracer.Inject(childSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body), nil
	} else {
		return "", errors.New("unauthorized: check error")
	}
}

func (userAuth *userAuthService) Update() (bool, error) {
	return false, nil
}

func (userAuth *userAuthService) Create(service string, email string, password string) (bool, error) {
	tracer := opentracing.GlobalTracer()
	childSpan := tracer.StartSpan(
		"userauth create",
	)
	defer childSpan.Finish()
	resp, err := http.PostForm(SERVICE_URL+"/user",
		url.Values{"service": {service}, "email": {email}, "password": {password}})
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return true, nil
	} else {
		return false, errors.New("auth: create user error")
	}

}

func (userAuth *userAuthService) Delete(string) (bool, error) {
	return false, nil
}
