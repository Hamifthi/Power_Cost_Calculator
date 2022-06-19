package entity

import "sync"

func CreateSyncPools(syncPoolSize int64) *SyncPools {
	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]byte, syncPoolSize)
		return lines
	}}
	stringsPool := sync.Pool{New: func() interface{} {
		strs := ""
		return strs
	}}
	return &SyncPools{LinesPool: &linesPool, StringsPool: &stringsPool}
}
