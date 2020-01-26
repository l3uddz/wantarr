package pvr

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/l3uddz/wantarr/config"
	"github.com/l3uddz/wantarr/logger"
	"github.com/l3uddz/wantarr/utils/web"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
)

/* Structs */

type Sonarr struct {
	cfg        *config.Pvr
	log        *logrus.Entry
	apiUrl     string
	reqHeaders req.Header
}

type Queue struct {
	Size int `json:"totalRecords"`
}

/* Initializer */

func NewSonarr(name string, c *config.Pvr) *Sonarr {
	// set api url
	apiUrl := ""
	if strings.Contains(c.URL, "/api") {
		apiUrl = c.URL
	} else {
		apiUrl = web.JoinURL(c.URL, "/api/v3")
	}

	// set headers
	reqHeaders := req.Header{
		"X-Api-Key": c.ApiKey,
	}

	return &Sonarr{
		cfg:        c,
		log:        logger.GetLogger(name),
		apiUrl:     apiUrl,
		reqHeaders: reqHeaders,
	}
}

/* Interface Implements */

func (p *Sonarr) GetWanted() error {
	return nil
}

func (p *Sonarr) GetQueueSize() (int, error) {
	// send request
	resp, err := web.GetResponse(web.GET, web.JoinURL(p.apiUrl, "/queue"), 15, p.reqHeaders)
	if err != nil {
		return 0, errors.WithMessage(err, "failed retrieving queue api response from sonarr")
	}
	defer resp.Response().Body.Close()

	// validate response
	if resp.Response().StatusCode != 200 {
		return 0, fmt.Errorf("failed retrieving valid queue api response from sonarr: %s", resp.Response().Status)
	}

	// decode response
	var q Queue
	if err := resp.ToJSON(&q); err != nil {
		return 0, errors.WithMessage(err, "failed decoding queue api response from sonarr")
	}

	return q.Size, nil
}
