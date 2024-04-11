package session

import (
	"github.com/mrzack99s/cocong/model/inmemory_model"
	"github.com/mrzack99s/cocong/network"
	"github.com/mrzack99s/cocong/vars"
)

func NewSession(session *inmemory_model.Session) (err error) {

	if !vars.SYS_DEBUG {
		err = network.AllowAccess(session)
		if err != nil {
			return
		}
	}

	err = vars.InMemoryDatabase.Create(session).Error
	return
}
