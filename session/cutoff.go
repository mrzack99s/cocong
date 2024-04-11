package session

import (
	"github.com/mrzack99s/cocong/model/inmemory_model"
	"github.com/mrzack99s/cocong/network"
	"github.com/mrzack99s/cocong/vars"
)

func CutOffSession(ss inmemory_model.Session) (err error) {

	if !vars.SYS_DEBUG {
		err = network.DeleteAllowAccess(&ss)
		if err != nil {
			return
		}

	}

	vars.InMemoryDatabase.Unscoped().Where("id = ?", ss.ID).Delete(&inmemory_model.Session{})

	return
}
