package spellcheck

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
        yandexSpellerUrlEnvVar = "YANDEX_SPELLER_URL"
)

var (
        yandexSpellerUrl = os.Getenv(yandexSpellerUrlEnvVar)
)

type YandexSpeller struct {}

func NewYandexSpeller() *YandexSpeller {
        yandexSpeller := &YandexSpeller{}
        return yandexSpeller
}

// Check sends text to Yandex.Speller (https://yandex.ru/dev/speller) to check it for spelling mistakes. 
// Returns: 
//      - true, nil, nil - if text is correct
//      - false, CheckResponse, nil - if spelling mistake found
//      - false, nil, error - if unexcpected error accured durring the request
//
// WARNING: Yandex.Speller is not working correctly via POST method. 
// Need to use GET method with 10KB URL limitaion (with text path parameter included)
//
func (y *YandexSpeller) Check(text string) (bool, CheckResponse, error) {
        urlWithText := yandexSpellerUrl + "?text=" + url.QueryEscape(text)

        if len(urlWithText) >= 10_000 {
                return false, nil, errors.New("final url with text is too big: pass smaller text")
        }

        resp, err := http.Get(urlWithText)
        if err != nil {
                return false, nil, err
        }
        
        data, err := io.ReadAll(resp.Body)
        if err != nil {
                return false, nil, err
        }

        var checkResponse CheckResponse

        err = json.Unmarshal(data, &checkResponse)
        if err != nil {
                return false, nil, err
        }


        if len(checkResponse) == 0 {
                return true, nil, nil
        } else {
                return false, checkResponse, nil
        }
}
