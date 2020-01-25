package pvr

import (
	"github.com/l3uddz/wantarr/config"
	"github.com/l3uddz/wantarr/logger"
	"github.com/sirupsen/logrus"
)

type Sonarr struct {
	cfg *config.Pvr
	log *logrus.Entry
}

func NewSonarr(name string, c *config.Pvr) *Sonarr {
	return &Sonarr{
		cfg: c,
		log: logger.GetLogger(name),
	}
}

/* Interface Implements */

func (p *Sonarr) GetWanted() error {
	p.log.Info("Getting wanted!")
	return nil
}
