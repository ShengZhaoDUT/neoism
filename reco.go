package neoism

import (
	"github.com/jmcvetta/napping"
	"strconv"
)

func (db *Database) Recommendation(id int64, limit int) ([]Recommendation, error) {
	reco := []Recommendation{}

	uri := join(db.HrefReco, strconv.FormatInt(id, 10))
	ne := NeoError{}
	resp, err := db.Session.Get(uri, &napping.Params{"limit": strconv.Itoa(limit)}, &reco, &ne)
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

func (db *Database) ContactRecommendation(id int64, limit int) ([]ConcactRecommendation, error) {
	reco := []ConcactRecommendation{}

	uri := join(db.HrefCReco, strconv.FormatInt(id, 10))
	ne := NeoError{}
	resp, err := db.Session.Get(uri, &napping.Params{"limit": strconv.Itoa(limit)}, &reco, &ne)
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

func (db *Database) GroupRecommendation(src int64, limit int) ([]int64, error) {
	type s struct {
		ID int64 `json:"id"`
	}
	result := []s{}
	cq := CypherQuery{
		Statement: `
start src = node({src})
match (n:Group)
where not src-[:Join]-n return id(n) as id limit 10`,
		Parameters: Props{"src": src},
		Result:     &result,
	}
	db.Cypher(&cq)
	IDList := make([]int64, 0)
	if result == nil {
		return IDList, nil
	}
	for _, x := range result {
		IDList = append(IDList, x.ID)
	}
	return IDList, nil
}

type Recommendation struct {
	//	UUID  string      `json:"uuid"`
	ID              int64       `json:"id"`
	Score           interface{} `json:"score"`
	inContact       bool        `json:"inContact"`
	friendsInCommon []int64     `json:"friendsInCommon"`
}

type ConcactRecommendation struct {
	//	UUID  string      `json:"uuid"`
	ID              int64       `json:"id"`
	Phone           string      `json:"phone"`
	Score           interface{} `json:"score"`
	friendsInCommon []int64     `json:"friendsInCommon"`
}
