package suite

import (
	"database/sql"
	"sync"
)

func NewDBs() DBStore {
	return &DBs{
		dbs: make(map[string]*sql.DB),
		m:   sync.RWMutex{},
	}
}

type DBs struct {
	dbs map[string]*sql.DB
	m   sync.RWMutex
}

func (d *DBs) StoreDB(name string, db *sql.DB) {
	d.m.Lock()
	d.dbs[name] = db
	d.m.Unlock()
}
func (d *DBs) RetrieveDB(name string) *sql.DB {
	var db *sql.DB
	d.m.RLock()
	db = d.dbs[name]
	d.m.RUnlock()
	return db
}
