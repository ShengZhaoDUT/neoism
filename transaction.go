// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package neoism

import (
	"encoding/json"
	"errors"
)

// A Tx is an in-progress database transaction.
type Tx struct {
	db         *Database
	hrefCommit string
	Location   string
	Errors     []TxError
	Expires    string // Cannot unmarshall into time.Time :(
}

type txRequest struct {
	Statements []*CypherQuery `json:"statements"`
}

type txResponse struct {
	Commit  string
	Results []struct {
		Columns []string
		Data    []struct {
			Row []*json.RawMessage
		}
	}
	Transaction struct {
		Expires string
	}
	Errors []TxError
}

// unmarshal populates a slice of CypherQuery object with result data returned
// from the server.
func (tr *txResponse) unmarshal(qs []*CypherQuery) error {
	if len(tr.Results) != len(qs) {
		return errors.New("Result count does not match query count")
	}
	// NOTE: Beginning in 2.0.0-M05, the data format returned by transaction
	// endpoint diverged from the format returned by cypher batch.  At least
	// until final 2.0.0 release, we will work around this by munging the new
	// result format into the existing cypherResult struct.
	for i, res := range tr.Results {
		data := make([][]*json.RawMessage, len(res.Data))
		for n, d := range res.Data {
			data[n] = d.Row
		}
		q := qs[i]
		cr := cypherResult{
			Columns: res.Columns,
			Data:    data,
		}
		q.cr = cr
		if q.Result != nil {
			err := q.Unmarshal(q.Result)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (db *Database) CommitQueries(queries []string) error {
	qs := make([]*CypherQuery, 0, len(queries))
	for _, q := range queries {
		qs = append(qs, &CypherQuery{Statement: q})
	}
	return db.Commit(qs)
}

// Begin opens a new transaction, executing zero or more cypher queries
// inside the transaction.
func (db *Database) Commit(qs []*CypherQuery) error {
	payload := txRequest{Statements: qs}
	result := txResponse{}
	ne := NeoError{}
	resp, err := db.Session.Post(db.HrefTransaction+"/commit", payload, &result, &ne)
	if err != nil {
		return err
	}
	if resp.Status() != 201 {
		return ne
	}
	err = result.unmarshal(qs)
	if err != nil {
		return err
	}
	return nil
}
