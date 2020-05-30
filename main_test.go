package main

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "github.com/vinhut/feed-service/mocks"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock_auth := mocks.NewMockAuthService(ctrl)
	mock_post := mocks.NewMockPostService(ctrl)

	router := setupRouter(mock_auth, mock_post)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", SERVICE_NAME+"/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}
