package repository

import (
	"be-dashboard-nba/internal/db"
)

type Query struct {
	db db.Query
}

func NewQuery(db db.Query) *Query {
	return &Query{
		db: db,
	}
}
