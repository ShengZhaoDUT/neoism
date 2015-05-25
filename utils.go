package neoism

// this file works a helper class for other files

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// do a simple json.Marshall
// returns string
func jsonEncode(value interface{}) (string, error) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(jsonValue), nil
}

// decoding json here
func jsonDecode(data string, result *interface{}) error {
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return err
	}

	return nil
}

func getID(url string) int64 {
	parts := strings.Split(url, "/")
	s := parts[len(parts)-1]
	id, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		// Are both r.Info and r.Node valid?
		panic(err)
	}
	return id

}

// Obtain id from incoming URL
func getIDFromURL(base, url string) (string, error) {
	//add slash to end of base url,
	//because before id there is a slash
	target := base + "/"

	result := strings.SplitAfter(url, target)

	if len(result) > 1 {
		return result[1], nil
	}

	return "", errors.New("URL not valid")
}

//
func relationshipURLByNode(nodeID int64) string {
	return fmt.Sprintf("/node/%d/relationships", nodeID)
}

func nodeURLByNode(nodeID int64) string {
	return fmt.Sprintf("/node/%d", nodeID)
}

// to-do combine this method with doRequest function
func (db *Database) doBatchRequest(requestType, url, data string) (string, error) {

	//convert string into bytestream
	dataByte := strings.NewReader(data)
	req, err := http.NewRequest(requestType, url, dataByte)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := db.Session.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf(res.Status)
	}

	// read response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//A successful call returns err == nil
		return "", err
	}

	return string(body), nil

}
