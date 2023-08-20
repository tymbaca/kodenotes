package spellcheck

import (
        "context"
        "encoding/json"
        "errors"
        "io"
        "net/http"
        "net/url"
        "time"

        "github.com/tymbaca/kodenotes/util"
)

const (
        yandexSpellerUrlEnvVar = "YANDEX_SPELLER_URL"
)

var (
        ErrYandexTooBigText = errors.New("text is too big: pass text with less than 10_000 characters")

        yandexSpellerUrl      = util.MustGetenv(yandexSpellerUrlEnvVar)
        yandexResponseTimeout = util.GetenvIntOrDefault("YANDEX_SPELLER_TIMEOUT", 10)
)

type ResponseWrapper struct {
        resp *http.Response
        err  error
}

type YandexSpeller struct{}

func NewYandexSpeller() *YandexSpeller {
        yandexSpeller := &YandexSpeller{}
        return yandexSpeller
}

// Check sends text to Yandex.Speller (https://yandex.ru/dev/speller) to check it for spelling mistakes.
// Returns:
//   - CheckResponse, nil - if spelling mistake found
//   - nil, error - if unexcpected error accured during the request. Error can be ErrTooBigText or another.
func (y *YandexSpeller) Check(text string) (CheckResponse, error) {

        if len(text) >= 10_000 {
                return CheckResponse{}, ErrYandexTooBigText
        }

        formData := url.Values{}
        formData.Add("text", text)

        resp, err := fetchCheckData(formData)
        if err != nil {
                return CheckResponse{}, err
        }

        data, err := io.ReadAll(resp.Body)
        if err != nil {
                return CheckResponse{}, err
        }

        var checkResponse CheckResponse

        err = json.Unmarshal(data, &checkResponse)
        if err != nil {
                return CheckResponse{}, err
        }

        return checkResponse, nil
}

func fetchCheckData(formData url.Values) (*http.Response, error) {
        ctx, cancel := context.WithTimeout(context.Background(), time.Duration(yandexResponseTimeout)*time.Second)
        defer cancel()
        respWrapCh := make(chan ResponseWrapper)

        go func() {
                resp, err := http.PostForm(yandexSpellerUrl, formData)
                respWrapCh <- ResponseWrapper{resp: resp, err: err}
        }()

        for {
                select {
                case <-ctx.Done():
                        return nil, ErrCheckTimeout
                case respWrap := <-respWrapCh:
                        return respWrap.resp, respWrap.err
                }

        }
}
