package auth

import (
	"fmt"
	"net/http"
	"net/url"
	//"io/ioutil"
	"log"
	"net/http/cookiejar"
	
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/gin-gonic/gin"
)

type AuthAdapter struct {
	authURL string `env:"AUTH_URL" envDefault:"http://localhost:3500/auth/api/v1"`
	client *http.Client
}

func New() *AuthAdapter {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}
	return &AuthAdapter{
		client: &http.Client{
			Jar: jar,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func (a *AuthAdapter) Verify(ctx *gin.Context) error {
	accessToken, err := ctx.Request.Cookie("access_token")
	if err != nil {
		return err
	}

	refreshToken, err := ctx.Request.Cookie("refresh_token")
	if err != nil {
		return err
	}

	urlObj, _ := url.Parse("http://localhost:8000/")
	a.client.Jar.SetCookies(urlObj, []*http.Cookie{accessToken, refreshToken})

	authR, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8000/auth/api/v1/verify", http.NoBody)
	if err != nil {
		return err
	}

	authR.Close = true

	resp, err := a.client.Do(authR)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("unexpected error")
}