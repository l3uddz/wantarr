package pvr

import (
	"fmt"
	"github.com/jpillora/backoff"
	"github.com/l3uddz/wantarr/config"
	"github.com/l3uddz/wantarr/utils/web"
	"strings"
	"time"
)

var (
	pvrDefaultPageSize      = 1000
	pvrDefaultRetry         = web.Retry{
		MaxAttempts:          5,
		RetryableStatusCodes: []int{},
		Backoff: backoff.Backoff{
			Jitter: true,
			Min:    1 * time.Second,
			Max:    5 * time.Second,
		},
	}
)

type MediaItem struct {
	ItemId     int
	AirDateUtc time.Time
	LastSearch time.Time
}

type Interface interface {
	Init() error
	GetQueueSize() (int, error)
	GetWantedMissing() ([]MediaItem, error)
	GetWantedCutoff() ([]MediaItem, error)
	SearchMediaItems([]int) (bool, error)
}

/* Public */

func Get(pvrName string, pvrType string, pvrConfig *config.Pvr) (Interface, error) {

	switch strings.ToLower(pvrType) {
	case "sonarr_v3":
		return NewSonarrV3(pvrName, pvrConfig), nil
	case "radarr_v2":
		return NewRadarrV2(pvrName, pvrConfig), nil
	case "radarr_v3":
		return NewRadarrV3(pvrName, pvrConfig), nil
	default:
		break
	}

	return nil, fmt.Errorf("unsupported pvr type provided: %q", pvrType)
}
