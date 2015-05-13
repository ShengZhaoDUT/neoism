package neoism

import (
	"strconv"
)

func (db *Database) Recommendation(id int64, limit int) ([]Recommendation, error) {
	reco := []Recommendation{}

	uri := join(db.HrefReco, strconv.FormatInt(id, 10))
	ne := NeoError{}
	resp, err := db.Session.Get(uri, nil, &reco, &ne)
	if err != nil {
		return reco, err
	}
	switch resp.Status() {
	default:
		err = ne
	case 200:
		err = nil // Success!
	case 404:
		err = NotFound
	}
	return reco, err
}

type Recommendation struct {
	//	UUID  string      `json:"uuid"`
	ID    int64       `json:"id"`
	Score interface{} `json:"score"`
}
