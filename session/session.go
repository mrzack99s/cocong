package session

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mrzack99s/cocong/network"
	"github.com/mrzack99s/cocong/types"
	"github.com/mrzack99s/cocong/vars"
)

func (i *sessionType) New() {
	i.ipmap = make(map[string]*types.SessionInfo)
}

func (i *sessionType) GetByIP(ipaddr string) (*types.SessionInfo, error) {
	keys := i.searchKey(ipaddr)
	if len(keys) > 0 {
		if s, ok := i.ipmap[keys[0]]; ok {
			return s, nil
		}
	}

	return nil, fmt.Errorf("not found ip address %s in session", ipaddr)
}

func (i *sessionType) GetByUsername(username string) ([]types.SessionInfo, error) {
	keys := i.searchKey(username)

	if len(keys) == 0 {
		return nil, fmt.Errorf("not found username %s in session", username)
	}

	var sessions []types.SessionInfo

	for _, k := range keys {
		if s, ok := i.ipmap[k]; ok {
			sessions = append(sessions, *s)
		}
	}

	return sessions, nil
}

func (i *sessionType) GetByID(id string) (*types.SessionInfo, error) {
	if s, ok := i.ipmap[id]; ok {
		return s, nil
	}

	return nil, errors.New("not found")
}

func (i *sessionType) Create(info types.SessionInfo) error {

	if !vars.SYS_DEBUG {
		err := network.AllowAccess(&info)
		if err != nil {
			return err
		}
	}

	key := fmt.Sprintf(key_pattern, info.User, info.IPAddress)
	info.ID = key
	info.LastSeen = time.Now().In(vars.TZ)
	i.ipmap[key] = &info

	return nil
}

func (i *sessionType) Delete(ipaddr string) error {
	session, err := i.GetByIP(ipaddr)
	if err != nil {
		return err
	}

	if !vars.SYS_DEBUG {
		network.DeleteAllowAccess(session)

	}

	delete(i.ipmap, session.ID)

	return nil
}

func (i *sessionType) DeleteByID(id string) error {
	session, ok := i.ipmap[id]
	if !ok {
		return errors.New("not found key")
	}

	if !vars.SYS_DEBUG {
		network.DeleteAllowAccess(session)

	}

	delete(i.ipmap, id)

	return nil
}

func (i *sessionType) UpdateLastSeen(ipaddr string) error {
	session, err := i.GetByIP(ipaddr)
	if err != nil {
		return err
	}

	session.LastSeen = time.Now().In(vars.TZ)
	i.ipmap[session.ID] = session

	return nil
}

func (i *sessionType) IsExpired(id string) (bool, error) {
	session, ok := i.ipmap[id]
	if !ok {
		return false, errors.New("not found key")
	}

	now := time.Now().In(vars.TZ)
	diff := now.Sub(session.LastSeen)
	if diff.Seconds() > float64(vars.Config.SessionIdle*60) {
		return true, nil
	}

	return false, nil
}

func (i *sessionType) searchKey(searchKey string) []string {
	var results []string

	for key := range i.ipmap {
		if strings.Contains(key, searchKey) {
			results = append(results, key)
		}
	}

	return results
}

func (i *sessionType) Search(searchKey string, offset, limit int) (*SearchResult, error) {
	var results []types.SessionInfo

	for key, value := range i.ipmap {
		if strings.Contains(key, searchKey) {
			results = append(results, *value)
		}
	}

	start := offset
	end := offset + limit

	resp := SearchResult{}

	if len(results) == 0 {
		return &resp, nil
	}

	if start > len(results) {
		return nil, errors.New("offset is out of legth")
	}

	if end > len(results) {
		end = len(results)
	}

	if len(results) > 0 {
		resp.Data = results[start:end]
		resp.Count = len(results)
	}

	return &resp, nil

}

func (i *sessionType) GetAllSession() []types.SessionInfo {
	var results []types.SessionInfo

	for _, value := range i.ipmap {
		results = append(results, *value)
	}

	return results

}
