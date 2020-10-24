package database

import (
	"fmt"

	// loading db driver
	_ "github.com/lib/pq"

	"covid-19/internal/config"

	"xorm.io/xorm"
)

// Orm for the database
type Orm struct {
	*xorm.Engine
}

// Session to separate db operations to session
type Session struct {
	*xorm.Session
}

// GetOrm returns an orm
func GetOrm(db config.DB) *Orm {
	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s/%s", db.User, db.Pass, db.Host, db.Name)
	if engine, err := xorm.NewEngine("postgres", dataSourceName); err != nil {
		panic(err)
	} else {
		return &Orm{engine}
	}
	return nil
}

// NewSession Create a new session
// TODO: Should start use session instead of engiine which uses the default session.
func (orm Orm) NewSession() *Session {
	return &Session{orm.Engine.NewSession()}
}
