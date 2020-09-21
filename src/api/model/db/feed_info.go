package db

import (
	"database/sql"

	uuid "github.com/satori/go.uuid"
)

type FeedInfo struct {
	ID			uuid.UUID	`json:"id"`
	Title		string		`json:"title"`
	URL 		string		`json:"rss_url"`
	XPath		string		`json:"xpath"`
}

func (data SQLDataStorage) GetFeedInfo() ([]*FeedInfo, error) {
	query := `select * from FeedInfoTable`
	rows, err := data.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedInfo []*FeedInfo
	for rows.Next() {
		f, err := feedInfoFromRows(rows)
		if err != nil {
			return nil, err
		}
		feedInfo = append(feedInfo, f)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return feedInfo, nil
}

func (data SQLDataStorage) PutFeedInfo(feedInfo *FeedInfo) error {
	query := `
	INSERT into feedinfotable (title, rss_url, xpath)
	VALUES ($1, $2, $3)
	RETURNING id
	`
	arguments := []interface{}{
		feedInfo.Title,
		feedInfo.URL,
		feedInfo.XPath,
	}
	err := data.DB.QueryRow(
		query, arguments...,
	).Scan(&feedInfo.ID)
	if err != nil {
		return err
	}
	return nil
} 

func feedInfoFromRows(rows *sql.Rows) (*FeedInfo, error) {
	f := &FeedInfo{}
	err := rows.Scan(
		&f.ID,
		&f.Title,
		&f.URL,
		&f.XPath,
	)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (data SQLDataStorage) DeleteFeedInfo(id uuid.UUID) error {
	query := `
	DELETE FROM feedinfotable
	WHERE id = $1
	`
	_, err := data.DB.Exec(query, sql.Named("id", id))
	return err
}