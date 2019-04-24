package bucket

import (
	"github.com/zzzhr1990/go-wcs-cloud-sdk/utility"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/core"
	log "github.com/sirupsen/logrus"
)

//Manager manager all
type Manager struct {
	auth        *utility.Auth
	config      *core.Config
	httpManager *utility.HTTPManager
}

func NewBucketManager(auth *utility.Auth,config  *core.Config) (bm *Manager) {
	if nil == auth {
		log.Errorln("Auth is nil!!!")
		return &Manager{auth, config, utility.NewHTTPManager()}
	}
	if nil == config {
		config = core.NewDefaultConfig()
	}
	return &Manager{auth, config, utility.NewHTTPManager()}
}