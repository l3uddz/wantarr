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
	"time"
)

/* Structs */

type RadarrV3 struct {
	cfg        *config.Pvr
	log        *logrus.Entry
	apiUrl     string
	reqHeaders req.Header
	version    int
}

type RadarrV3Movie struct {
	Id          int
	AirDateUtc  time.Time `json:"inCinemas"`
	Status      string
	HasFile     bool
	IsAvailable bool
	Monitored   bool
}

type RadarrV3SystemStatus struct {
	Version string
}

type RadarrV3CommandStatus struct {
	Name    string
	Message string
	Started time.Time
	Ended   time.Time
	Status  string
}

type RadarrV3CommandResponse struct {
	Id int
}

type RadarrV3MovieSearch struct {
	Name   string `json:"name"`
	Movies []int  `json:"movieIds"`
}

/* Initializer */

func NewRadarrV3(name string, c *config.Pvr) *RadarrV3 {
	// set api url
	apiUrl := ""
	if strings.Contains(c.URL, "/api") {
		apiUrl = c.URL
	} else {
		apiUrl = web.JoinURL(c.URL, "/api")
	}

	// set headers
	reqHeaders := req.Header{
		"X-Api-Key": c.ApiKey,
	}

	return &RadarrV3{
		cfg:        c,
		log:        logger.GetLogger(name),
		apiUrl:     apiUrl,
		reqHeaders: reqHeaders,
	}
}

/* Private */

func (p *RadarrV3) getSystemStatus() (*RadarrV3SystemStatus, error) {
	// send request
	resp, err := web.GetResponse(web.GET, web.JoinURL(p.apiUrl, "/system/status"), 15, p.reqHeaders,
		&pvrDefaultRetry)
	if err != nil {
		return nil, errors.New("failed retrieving system status api response from radarr")
	}
	defer resp.Response().Body.Close()

	// validate response
	if resp.Response().StatusCode != 200 {
		return nil, fmt.Errorf("failed retrieving valid system status api response from radarr: %s",
			resp.Response().Status)
	}

	// decode response
	var s RadarrV3SystemStatus
	if err := resp.ToJSON(&s); err != nil {
		return nil, errors.WithMessage(err, "failed decoding system status api response from radarr")
	}

	return &s, nil
}

func (p *RadarrV3) getCommandStatus(id int) (*RadarrV3CommandStatus, error) {
	// send request
	resp, err := web.GetResponse(web.GET, web.JoinURL(p.apiUrl, fmt.Sprintf("/command/%d", id)), 15,
		p.reqHeaders, &pvrDefaultRetry)
	if err != nil {
		return nil, errors.New("failed retrieving command status api response from radarr")
	}
	defer resp.Response().Body.Close()

	// validate response
	if resp.Response().StatusCode != 200 {
		return nil, fmt.Errorf("failed retrieving valid command status api response from radarr: %s",
			resp.Response().Status)
	}

	// decode response
	var s RadarrV3CommandStatus
	if err := resp.ToJSON(&s); err != nil {
		return nil, errors.WithMessage(err, "failed decoding command status api response from radarr")
	}

	return &s, nil
}

/* Interface Implements */

func (p *RadarrV3) Init() error {
	// retrieve system status
	status, err := p.getSystemStatus()
	if err != nil {
		return errors.Wrap(err, "failed initializing radarr pvr")
	}

	// determine version
	switch status.Version[0:1] {
	case "3":
		p.version = 1
	default:
		return fmt.Errorf("unsupported version of radarr pvr: %s", status.Version)
	}
	return nil
}

func (p *RadarrV3) GetQueueSize() (int, error) {
	// send request
	resp, err := web.GetResponse(web.GET, web.JoinURL(p.apiUrl, "/queue"), 15, p.reqHeaders,
		&pvrDefaultRetry)
	if err != nil {
		return 0, errors.WithMessage(err, "failed retrieving queue api response from radarr")
	}
	defer resp.Response().Body.Close()

	// validate response
	if resp.Response().StatusCode != 200 {
		return 0, fmt.Errorf("failed retrieving valid queue api response from radarr: %s",
			resp.Response().Status)
	}

	// decode response
	var q []interface{}
	if err := resp.ToJSON(&q); err != nil {
		return 0, errors.WithMessage(err, "failed decoding queue api response from radarr")
	}

	queueSize := len(q)
	p.log.WithField("queue_size", queueSize).Debug("Queue retrieved")
	return queueSize, nil
}

func (p *RadarrV3) GetWantedMissing() ([]MediaItem, error) {
	// logic vars
	totalRecords := 0
	var wantedMissing []MediaItem

	// retrieve all page results
	p.log.Info("Retrieving wanted missing media...")

	// send request
	resp, err := web.GetResponse(web.GET, web.JoinURL(p.apiUrl, "/movie"), 15,
		p.reqHeaders, &pvrDefaultRetry)
	if err != nil {
		return nil, errors.WithMessage(err, "failed retrieving wanted missing api response from radarr")
	}
	defer resp.Response().Body.Close()

	// validate response
	if resp.Response().StatusCode != 200 {
		return nil, fmt.Errorf("failed retrieving valid wanted missing api response from radarr: %s",
			resp.Response().Status)
	}

	// decode response
	var m []RadarrV3Movie
	if err := resp.ToJSON(&m); err != nil {
		return nil, errors.WithMessage(err, "failed decoding wanted missing api response from radarr")
	}

	// process response
	for _, movie := range m {
		// is movie monitored?
		if !movie.Monitored {
			continue
		}
		// is movie available?
		if !movie.IsAvailable {
			continue
		}
		// is the status released?
		if movie.Status != "released" {
			continue
		}
		// do we have the file?
		if movie.HasFile {
			continue
		}

		// store this movie
		airDate := movie.AirDateUtc
		wantedMissing = append(wantedMissing, MediaItem{
			ItemId:     movie.Id,
			AirDateUtc: airDate,
			LastSearch: time.Time{},
		})
	}
	totalRecords = len(wantedMissing)

	p.log.WithField("media_items", totalRecords).Info("Finished")
	return wantedMissing, nil
}

func (p *RadarrV3) GetWantedCutoff() ([]MediaItem, error) {
	return nil, errors.New("currently not implemented")
}

func (p *RadarrV3) SearchMediaItems(mediaItemIds []int) (bool, error) {
	// set request data
	payload := RadarrMovieSearch{
		Name:   "moviesSearch",
		Movies: mediaItemIds,
	}

	// send request
	resp, err := web.GetResponse(web.POST, web.JoinURL(p.apiUrl, "/command"), 15, p.reqHeaders,
		&pvrDefaultRetry, req.BodyJSON(&payload))
	if err != nil {
		return false, errors.WithMessage(err, "failed retrieving command api response from radarr")
	}
	defer resp.Response().Body.Close()

	// validate response
	if resp.Response().StatusCode != 201 {
		return false, fmt.Errorf("failed retrieving valid command api response from radarr: %s",
			resp.Response().Status)
	}

	// decode response
	var q RadarrCommandResponse
	if err := resp.ToJSON(&q); err != nil {
		return false, errors.WithMessage(err, "failed decoding command api response from radarr")
	}

	// monitor search status
	p.log.WithField("command_id", q.Id).Debug("Monitoring search status")

	for {
		// retrieve command status
		searchStatus, err := p.getCommandStatus(q.Id)
		if err != nil {
			return false, errors.Wrapf(err, "failed retrieving command status from radarr for: %d", q.Id)
		}

		p.log.WithFields(logrus.Fields{
			"command_id": q.Id,
			"status":     searchStatus.Status,
		}).Debug("Status retrieved")

		// is status complete?
		if searchStatus.Status == "completed" {
			break
		}

		time.Sleep(10 * time.Second)
	}

	return true, nil
}
