package db

import (
	"database/sql"

)

type DataStorage interface {
	GetFeedInfo() ([]*FeedInfo, error)
	PutFeedInfo(feedInfo *FeedInfo) error
}

func NewSQLDataStorage(sqlDB *sql.DB) SQLDataStorage {
	return SQLDataStorage{DB: sqlDB}
}

type SQLDataStorage struct {
	DB *sql.DB
}