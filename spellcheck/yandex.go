package spellcheck

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tymbaca/kodenotes/util"
)

const (
        yandexSpellerUrlEnvVar = "YANDEX_SPELLER_URL"
)

var (
        yandexSpellerUrl = os.Getenv(yandexSpellerUrlEnvVar)
        yandexResponseTimeout = util.GetenvOrDefault("YANDEX_SPELLER_TIMEOUT", "10")
)

type ResponseWrapper struct {
        resp *http.Response
        err   error
}

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
func (y *YandexSpeller) Check(text string) (CheckResponse, error) {

        if len(text) >= 10_000 {
                return false, nil, errors.New("final url with text is too big: pass smaller text")
        }
        var formData url.Values
        formData.Add("text", text)


        
        data, err := io.ReadAll(resp.Body)
        if err != nil {
                return nil, err
        }

        var checkResponse CheckResponse

        err = json.Unmarshal(data, &checkResponse)
        if err != nil {
                return nil, err
        }

        return 
}

func checkData(formData url.Values) (*http.Response, error) {
        ctx, cancel := context.WithTimeout(context.Background(), yandexResponseTimeout * time.Second)
        respWrapCh      := make(chan ResponseWrapper)

        go func() {
                resp, err := http.PostForm(yandexSpellerUrl, formData)
                respWrapCh <- ResponseWrapper{resp: resp, err: err}
        }()
        
        for {
                select {
                case <-ctx.Done():
                        return nil, http.ErrHandlerTimeout
                case respWrap := <-respWrapCh:
                        return respWrap.resp, respWrap.err
                }

        }
}

