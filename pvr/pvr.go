package pvr

import (
	"fmt"
	"github.com/l3uddz/wantarr/config"
	"strings"
)

type Interface interface {
	GetWanted() error
	GetQueueSize() (int, error)
}

func Get(pvrName string, pvrType string, pvrConfig *config.Pvr) (Interface, error) {

	switch strings.ToLower(pvrType) {
	case "sonarr":
		return NewSonarr(pvrName, pvrConfig), nil
	default:
		break
	}

	return nil, fmt.Errorf("unsupported pvr type provided: %q", pvrType)
}
