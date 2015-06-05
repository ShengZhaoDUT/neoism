package neoism

import (
	"encoding/json"
	"strconv"
)

// curl -X POST http://localhost:7474/db/data/ext/ShortestDistance/node/23337/shortestDistance   -H "Content-Type: application/json"   -d '{"targets":"http://localhost:7474/db/data/node/23338","targets":"http://localhost:7474/db/data/node/23339", "depth":"2"}'
func (db *Database) NodeDistance(src int64, dst int64, relType []string, depth int) int {
	url := join(db.Url, "ext", "ShortestDistance", "node", strconv.FormatInt(src, 10), "shortestDistance")
	var result interface{}
	ne := NeoError{}
	type s struct {
		Target string   `json:"target"`
		Types  []string `json:"types"`
		Depth  int      `json:"depth"`
	}
	payload := s{
		Target: join(db.HrefNode, strconv.FormatInt(dst, 10)),
		Types:  relType,
		Depth:  depth,
	}
	_, err := db.Session.Post(url, payload, &result, &ne)
	if err != nil {
		return depth
	}

	return int(result.(float64))
}

func (db *Database) Profile(src int64, dst int64, relType []string, depth int) []byte {
	url := join(db.Url, "ext", "Profile", "node", strconv.FormatInt(src, 10), "profile")
	var result []interface{}
	ne := NeoError{}
	type s struct {
		Target string   `json:"target"`
		Types  []string `json:"types"`
		Depth  int      `json:"depth"`
	}
	payload := s{
		Target: join(db.HrefNode, strconv.FormatInt(dst, 10)),
		Types:  relType,
		Depth:  depth,
	}
	profile := make(map[string]string)
	_, err := db.Session.Post(url, payload, &result, &ne)
	if err != nil {
		return []byte{}
	}
	length := len(result)
	if length%2 == 0 {
		for i := 0; i < length; i += 2 {
			profile[result[i].(string)] = result[i+1].(string)
		}
	}
	b, errr := json.Marshal(profile)
	if errr != nil {
		return []byte{}
	}

	return b
}

func (db *Database) Profiles(src int64, dst []int64, relType []string, depth int) []byte {
	url := join(db.Url, "ext", "Profiles", "node", strconv.FormatInt(src, 10), "profiles")
	targets := make([]string, len(dst))
	for i, x := range dst {
		targets[i] = join(db.HrefNode, strconv.FormatInt(x, 10))
	}
	type s struct {
		Targets []string `json:"targets"`
		Types   []string `json:"types"`
		Depth   int      `json:"depth"`
	}

	payload := s{
		Targets: targets,
		Types:   relType,
		Depth:   depth,
	}
	return db.GetProfiles(url, payload)
}

func (db *Database) Props(src int64, dst []int64) []byte {
	url := join(db.Url, "ext", "Props", "node", strconv.FormatInt(src, 10), "props")
	targets := make([]string, len(dst))
	for i, x := range dst {
		targets[i] = join(db.HrefNode, strconv.FormatInt(x, 10))
	}
	type s struct {
		Targets []string `json:"targets"`
	}

	payload := s{
		Targets: targets,
	}
	return db.GetProfiles(url, payload)
}

func (db *Database) RelatedNode(src int64, relType string) []byte {
	url := join(db.Url, "ext", "RelatedNode", "node", strconv.FormatInt(src, 10), "all")
	type s struct {
		Type string `json:"type"`
	}

	payload := s{
		Type: relType,
	}
	return db.GetProfiles(url, payload)
}

func (db *Database) IncomingRelatedNode(src int64, relType string) []byte {
	url := join(db.Url, "ext", "RelatedNode", "node", strconv.FormatInt(src, 10), "incoming")
	type s struct {
		Type string `json:"type"`
	}

	payload := s{
		Type: relType,
	}
	return db.GetProfiles(url, payload)
}

func (db *Database) OutgoingRelatedNode(src int64, relType string) []byte {
	url := join(db.Url, "ext", "RelatedNode", "node", strconv.FormatInt(src, 10), "outgoing")
	type s struct {
		Type string `json:"type"`
	}

	payload := s{
		Type: relType,
	}
	return db.GetProfiles(url, payload)
}

func (db *Database) InCommonNode(src int64, dst int64, relType string) []byte {
	url := join(db.Url, "ext", "InCommonNode", "node", strconv.FormatInt(src, 10), "all")
	target := join(db.HrefNode, strconv.FormatInt(dst, 10))
	type s struct {
		Target string `json:"target"`
		Type   string `json:"type"`
	}

	payload := s{
		Target: target,
		Type:   relType,
	}
	return db.GetProfiles(url, payload)
}

// return value
//{
//"source": Marshal Object
//	"type1": Marshal Object
//	"type2": Marshal Object
//}
func (db *Database) FullProfile(src int64, relType []string) []byte {
	url := join(db.Url, "ext", "FullProfile", "node", strconv.FormatInt(src, 10), "all")
	type s struct {
		Types []string `json:"types"`
	}

	payload := s{
		Types: relType,
	}
	return db.GetProfiles(url, payload)
}

func (db *Database) GetProfiles(url string, payload interface{}) []byte {
	var result interface{}
	ne := NeoError{}
	_, err := db.Session.Post(url, payload, &result, &ne)
	if err != nil || result == nil {
		return []byte{}
	}

	return []byte(result.(string))
}

func (db *Database) MultiUniqRelate(src int64, dst []int64, relType string, p interface{}) {
	url := join(db.Url, "ext", "MultiUniqRelate", "node", strconv.FormatInt(src, 10), "relate")
	targets := make([]string, len(dst))
	for i, x := range dst {
		targets[i] = join(db.HrefNode, strconv.FormatInt(x, 10))
	}

	type s struct {
		Targets []string `json:"targets"`
		Type    string   `json:"type"`
		Props   string   `json:"props"`
	}

	payload := s{
		Targets: targets,
		Type:    relType,
	}

	if p != nil {
		prop, _ := json.Marshal(p)
		payload.Props = string(prop)
	}
	var result interface{}
	ne := NeoError{}
	_, err := db.Session.Post(url, payload, &result, &ne)
	if err != nil {
		panic(err)
	}

	return
}
