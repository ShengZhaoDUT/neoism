package neoism

import (
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

func (db *Database) Profile(src int64, dst int64, relType []string, depth int) map[string]string {
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
		return profile
	}
	length := len(result)
	if length%2 == 0 {
		for i := 0; i < length; i += 2 {
			profile[result[i].(string)] = result[i+1].(string)
		}
	}
	return profile
}

func (db *Database) Profiles(src int64, dst []int64, relType []string, depth int) []map[string]string {
	url := join(db.Url, "ext", "Profiles", "node", strconv.FormatInt(src, 10), "profiles")
	var result []interface{}
	ne := NeoError{}
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
	profiles := make([]map[string]string, 0)
	_, err := db.Session.Post(url, payload, &result, &ne)
	if err != nil {
		return profiles
	}
	length := len(result)
	profile := make(map[string]string)
	if length%2 == 0 {
		for i := 0; i < length; i += 2 {
			key := result[i].(string)
			profile[key] = result[i+1].(string)
			if key == "id" {
				profiles = append(profiles, profile)
				profile = make(map[string]string)
			}
		}
	}
	return profiles
}
