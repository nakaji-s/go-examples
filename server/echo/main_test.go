package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/go-resty/resty"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	go func() {
		main()
	}()

	os.Exit(m.Run())
}

func TestHello(t *testing.T) {
	resp, err := resty.R().Get("http://localhost:8081/hello")
	assert.Nil(t, err)
	assert.Regexp(t, "hello[0-9]*", string(resp.Body()))
}

func TestHelloId(t *testing.T) {
	e := httpexpect.New(t, "http://localhost:8081")

	// is it working?
	e.GET("/hello/" + "me").
		Expect().
		Status(http.StatusOK).Body().Equal("hello me")
}

func TestEcho(t *testing.T) {
	resp, err := resty.R().
		SetBody("hogehoge").
		Post("http://localhost:8081/echo")
	assert.Nil(t, err)
	assert.Equal(t, "hogehoge", string(resp.Body()))
}

func TestEcho2(t *testing.T) {
	resp, err := resty.R().
		SetBody("hogehoge").
		Post("http://localhost:8081/blog/echo")
	assert.Nil(t, err)
	assert.Equal(t, "hogehoge"+"dummy", string(resp.Body()))
}

func TestMe(t *testing.T) {
	resp, err := resty.R().
		Get("http://localhost:8081/blog/me")
	assert.Nil(t, err)
	assert.Equal(t, "me", string(resp.Body()))
}

func TestIn(t *testing.T) {
	resp, err := resty.R().
		Get("http://localhost:8081/blog/in")
	assert.Nil(t, err)
	assert.Equal(t, "me", string(resp.Body()))
}
