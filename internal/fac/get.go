package fac

import (
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func FacGet(url string) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		zap.L().Error("could not initialize new request", zap.Error(err))
	}

	req.Header = http.Header{
		"X-API-Key":      {viper.GetString("api.key")},
		"Accept-Profile": {"api_v1_1_0"},
	}

	try_count := 3
	for {
		if try_count > 0 {
			res, err := client.Do(req)
			if err != nil {
				try_count = try_count - 1
				zap.L().Error("error in client.Do", zap.Error(err), zap.Int("countdown", try_count))
				time.Sleep(1 * time.Second)
			}

			if res.Body != nil {
				defer res.Body.Close()
			}

			body, readErr := io.ReadAll(res.Body)
			if readErr != nil {
				log.Fatal(readErr)
			}

			return body, nil
		} else {
			// Only try a few times if we get errors.
			break
		}
	}

	return []byte{}, errors.New("could not fetch body")

}
