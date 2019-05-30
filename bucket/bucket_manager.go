package bucket

import (
	log "github.com/sirupsen/logrus"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/core"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/utility"
)

//Manager manager all
type Manager struct {
	auth        *utility.Auth
	config      *core.Config
	httpManager *utility.HTTPManager
}

//NewBucketManager bb
func NewBucketManager(auth *utility.Auth, config *core.Config) *Manager {
	if nil == auth {
		log.Errorln("Auth is nil!!!")
		return &Manager{auth, config, utility.NewHTTPManager()}
	}
	if nil == config {
		config = core.NewDefaultConfig()
	}
	return &Manager{auth, config, utility.NewHTTPManager()}
}

// TestManagerHost need test
func (manager *Manager) TestManagerHost(bucket string) error {
	key := "newSys/do-not-delete.jpg"
	stat, err := manager.Stat(bucket, key)
	if err != nil {
		log.Errorf("stat %v:%v err, %v", bucket, key, err)
		return err
	}
	log.Infof("stat %v:%v success, hash: %v", bucket, key, stat.Hash)
	return nil
}
