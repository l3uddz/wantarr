package web

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/l3uddz/wantarr/logger"
	"github.com/l3uddz/wantarr/utils/lists"

	"github.com/imroc/req"
	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
)

var (
	// Logging
	log = logger.GetLogger("web")
)

/* Structs */

// HTTPMethod - The HTTP request method to use
type HTTPMethod int
type Retry struct {
	backoff.Backoff
	MaxAttempts          float64
	RetryableStatusCodes []int
	ExpectedContentType  string
}

const (
	// GET - Use GET HTTP method
	GET HTTPMethod = iota + 1
	// POST - Use POST HTTP method
	POST
	// PUT - Use PUT HTTP method
	PUT
	// DELETE - Use DELETE HTTP method
	DELETE
)

/* Private */

func init() {
	req.SetJSONEscapeHTML(false)
}

/* Public */

func GetResponse(method HTTPMethod, requestUrl string, timeout int, v ...interface{}) (*req.Resp, error) {
	// prepare request
	inputs := make([]interface{}, 0)

	var retry Retry
	for _, vv := range v {
		switch vT := vv.(type) {
		case *Retry:
			retry = *vT
		case Retry:
			retry = vT
		default:
			inputs = append(inputs, vT)
		}
	}

	// add context if available
	if timeout > 0 {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		inputs = append(inputs, ctx)
	}

	// Response var
	var resp *req.Resp
	var err error

	// Exponential backoff
	for {
		switch method {
		case GET:
			resp, err = req.Get(requestUrl, inputs...)
		case POST:
			resp, err = req.Post(requestUrl, inputs...)
		default:
			log.Error("Request method has not been implemented")
			return nil, errors.New("request method has not been implemented")
		}

		// validate response
		if err != nil {
			log.WithError(err).Debugf("Failed requesting: %q", requestUrl)
			if os.IsTimeout(err) {
				if retry.MaxAttempts == 0 || retry.Attempt() >= retry.MaxAttempts {
					return nil, err
				}

				d := retry.Duration()
				log.Debugf("Retrying failed request in %s: %q", d, requestUrl)
				time.Sleep(d)
				continue
			}

			return nil, err
		}

		log.Tracef("Request URL: %s", resp.Request().URL)
		log.Tracef("Request Response: %s", resp.Response().Status)

		if retry.MaxAttempts == 0 || retry.Attempt() >= retry.MaxAttempts {
			break
		}

		// check status code vs retryable ones
		if lists.IntListContains(resp.Response().StatusCode, retry.RetryableStatusCodes) {
			// close response body
			_ = resp.Response().Body.Close()

			// retry
			d := retry.Duration()
			log.Debugf("Retrying failed request in %s: %d - %q", d, resp.Response().StatusCode, requestUrl)

			time.Sleep(d)
			continue
		}

		// check response content type vs expected one
		if retry.ExpectedContentType != "" {
			// check response content type
			contentType := resp.Response().Header.Get("Content-Type")
			if !strings.Contains(strings.ToLower(contentType), strings.ToLower(retry.ExpectedContentType)) &&
				!strings.EqualFold(contentType, retry.ExpectedContentType) {
				// close response body
				_ = resp.Response().Body.Close()

				// retry
				d := retry.Duration()
				log.Debugf("Retrying failed request in %s: %d %s - %q", d, resp.Response().StatusCode, contentType, requestUrl)

				time.Sleep(d)
				continue
			}
		}

		break
	}

	return resp, err
}

func GetBodyBytes(method HTTPMethod, requestUrl string, timeout int, v ...interface{}) ([]byte, error) {
	// send request
	resp, err := GetResponse(method, requestUrl, timeout, v...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Response().Body.Close(); err != nil {
			log.WithError(err).Errorf("Failed to close response body for url: %q", requestUrl)
		}
	}()

	// process response
	body, err := ioutil.ReadAll(resp.Response().Body)
	if err != nil {
		log.WithError(err).Errorf("Failed reading response body for url: %q", requestUrl)
		return nil, errors.Wrap(err, "failed reading url response body")
	}

	return body, nil
}

func GetBodyString(method HTTPMethod, requestUrl string, timeout int, v ...interface{}) (string, error) {
	bodyBytes, err := GetBodyBytes(method, requestUrl, timeout, v...)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
