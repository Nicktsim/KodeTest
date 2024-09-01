package response

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"strings"

	"github.com/pkg/errors"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}


type YandexSpellerResponse []struct {
        Code int      `json:"code"`
        Pos  int      `json:"pos"`
        Row  int      `json:"row"`
        Col  int      `json:"col"`
        Len  int      `json:"len"`
        Word string   `json:"word"`
        S    []string `json:"s"`
    }

func ValidateNote(title,description string) error {
    u, err := url.Parse("https://speller.yandex.net/services/spellservice.json/checkText")
    if err != nil {
        return errors.Wrap(err, "failed to parse URL")
    }

    q := u.Query()
    q.Set("text", title+" "+description) 
    u.RawQuery = q.Encode()

    resp, err := http.Get(u.String())
    if err != nil {
        return errors.Wrap(err, "failed to make GET request to Yandex Speller")
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return errors.Wrap(err, "failed to read response body from Yandex Speller")
    }

    var spellerResponse YandexSpellerResponse
    err = json.Unmarshal([]byte(body), &spellerResponse)
    if err != nil {
        return errors.Wrap(err, "failed unmarshall response from Yandex Speller")
    }

    if len(spellerResponse) > 0 {
        var validationErrors []string
        for _, word := range spellerResponse {
            validationErrors = append(validationErrors, fmt.Sprintf("Mistake in word %s, may be %s?", word.Word, word.S[0]))
        }
        return fmt.Errorf("found some mistakes: %s", strings.Join(validationErrors, ","))

    }
    return nil
}