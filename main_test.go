package main

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "github.com/vinhut/feed-service/mocks"

	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestPing(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock_auth := mocks.NewMockAuthService(ctrl)
	mock_post := mocks.NewMockPostService(ctrl)

	router := setupRouter(mock_auth, mock_post)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestGetFeed(t *testing.T) {

	now := time.Now()
	token := "852a37a34b727c0e0b331806-7af4bdfdcc60990d427f383efecc8529289d040dd67e0753b9e2ee5a1e938402186f28324df23f6faa4e2bbf43f584ae228c55b00143866215d6e92805d470a1cc2a096dcca4d43527598122313be412e17fbefdcdab2fae02e06a405791d936862d4fba688b3c7fd784d4"
	user_data := "{\"uid\": \"1\", \"email\": \"test@email.com\", \"role\": \"standard\", \"created\": \"" + now.Format("2006-01-02T15:04:05") + "\"}"

	os.Setenv("KEY", "12345678901234567890123456789012")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock_auth := mocks.NewMockAuthService(ctrl)
	mock_post := mocks.NewMockPostService(ctrl)

	mock_auth.EXPECT().Check(gomock.Any(), gomock.Any()).Return(user_data, nil)
	mock_post.EXPECT().GetAll(gomock.Any()).Return(nil)

	router := setupRouter(mock_auth, mock_post)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", SERVICE_NAME+"/feed", nil)
	req.Header.Set("Cookie", "token="+token+";")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	fmt.Println(w.Body.String())
}

func TestEmailEmpty(t *testing.T) {

	now := time.Now()
	token := "852a37a34b727c0e0b331806-7af4bdfdcc60990d427f383efecc8529289d040dd67e0753b9e2ee5a1e938402186f28324df23f6faa4e2bbf43f584ae228c55b00143866215d6e92805d470a1cc2a096dcca4d43527598122313be412e17fbefdcdab2fae02e06a405791d936862d4fba688b3c7fd784d4"
	user_data := "{\"uid\": \"1\", \"email\": \"\", \"role\": \"standard\", \"created\": \"" + now.Format("2006-01-02T15:04:05") + "\"}"

	os.Setenv("KEY", "12345678901234567890123456789012")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock_auth := mocks.NewMockAuthService(ctrl)
	mock_post := mocks.NewMockPostService(ctrl)

	mock_auth.EXPECT().Check(gomock.Any(), gomock.Any()).Return(user_data, nil)
	mock_post.EXPECT().GetAll(gomock.Any()).Return(nil)

	router := setupRouter(mock_auth, mock_post)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", SERVICE_NAME+"/feed", nil)
	req.Header.Set("Cookie", "token="+token+";")
	router.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
}
