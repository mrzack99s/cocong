package session

import "github.com/mrzack99s/cocong/types"

type sessionType struct {
	ipmap map[string]*types.SessionInfo
}

type SearchResult struct {
	Data  []types.SessionInfo
	Count int
}
