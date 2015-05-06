package neoism

import (
	"errors"
	"fmt"
	"strconv"
)

// Basic operation names
var (
	BatchGet          = "get"
	BatchCreate       = "create"
	BatchDelete       = "delete"
	BatchUpdate       = "update"
	BatchCreateUnique = "createUnique"
)

// Batcher is the interface for structs for making them compatible with Batch.
type Batcher interface {
	getBatchQuery(operation string, db *Database) map[string]interface{}
}

// Batch Base struct to support request
type Batch struct {
	DB    *Database
	Stack []*BatchRequest
}

// BatchRequest All batch request structs will be encapslated in this struct
type BatchRequest struct {
	Operation string
	Data      Batcher
}

// BatchResponse All returning results from Neo4j will be in BatchResponse format
type BatchResponse struct {
	ID       int         `json:"id"`
	Location string      `json:"location"`
	Body     interface{} `json:"body"`
	From     string      `json:"from"`
}

// GetLastIndex Returns last index of current stack
// This method can be used to obtain the latest index for creating
// manual batch requests or injecting the order number of pre-added request(s) id
func (batch *Batch) GetLastIndex() string {

	return strconv.Itoa(len(batch.Stack) - 1)
}

// NewBatch Creates New Batch request handler
func (db *Database) NewBatch() *Batch {
	return &Batch{
		DB:    db,
		Stack: make([]*BatchRequest, 0),
	}
}

// Get request to Neo4j as batch
func (batch *Batch) Get(obj Batcher) *Batch {
	batch.addToStack(BatchGet, obj)

	return batch
}

// Create request to Neo4j as batch
func (batch *Batch) Create(obj Batcher) *Batch {
	batch.addToStack(BatchCreate, obj)

	return batch
}

// Delete request to Neo4j as batch
func (batch *Batch) Delete(obj Batcher) *Batch {
	batch.addToStack(BatchDelete, obj)

	return batch
}

// Update request to Neo4j as batch
func (batch *Batch) Update(obj Batcher) *Batch {
	batch.addToStack(BatchUpdate, obj)

	return batch
}

// Adds requests to stack
// Used internally to pile up the batch request
func (batch *Batch) addToStack(operation string, obj Batcher) {
	batchRequest := &BatchRequest{
		Operation: operation,
		Data:      obj,
	}

	batch.Stack = append(batch.Stack, batchRequest)
}

// Execute Prepares and sends the request to Neo4j
// If the request is successful then parses the response
func (batch *Batch) Execute() error {

	// if Neo4j instance is not created return an error
	if batch.DB == nil {
		return errors.New("Batch request is not created by NewBatch method!")
	}
	// cache batch stack length
	stackLength := len(batch.Stack)

	//create result array
	//	response := make([]*BatchResponse, stackLength)

	if stackLength == 0 {
		panic("Empty batch")
	}

	// prepare request
	ne := NeoError{}
	var result interface{}

	resp, err := batch.DB.Session.Post(batch.DB.HrefBatch, prepareRequest(batch.Stack, batch.DB), result, ne)
	//spew.Fdump(os.Stdout, resp)

	if err != nil {
		return err
	}

	if resp.Status() != 200 {
		return ne
	}

	return nil
}

// prepares batch request as slice of map
func prepareRequest(stack []*BatchRequest, db *Database) []map[string]interface{} {
	request := make([]map[string]interface{}, len(stack))
	for i, value := range stack {
		// interface has this method getBatchQuery()
		query := value.Data.getBatchQuery(value.Operation, db)
		query["id"] = i
		request[i] = query
	}

	return request
}

func (n *Node) getBatchQuery(operation string, db *Database) map[string]interface{} {

	query := make(map[string]interface{})
	switch operation {
	case BatchGet:
		query["method"] = "GET"
		query["to"] = db.getRelativePath(n.HrefSelf)

	case BatchUpdate:

		query["method"] = "PUT"
		query["to"] = db.getRelativePath(n.HrefProperties)
		query["body"] = n.Data

	case BatchCreate:
		query = map[string]interface{}{
			"method": "POST",
			"to":     "/node",
			"body":   n.Data,
		}

	case BatchDelete:
		query["method"] = "DELETE"
		query["to"] = db.getRelativePath(n.HrefSelf)

	}
	return query

}

func (r *Relationship) getBatchQuery(operation string, db *Database) map[string]interface{} {

	if r.Type == "" {
		panic("No Type on relationship")
	}

	switch operation {
	case BatchGet:

		return map[string]interface{}{
			"method": "GET",
			"to":     db.getRelativePath(r.HrefSelf),
		}

	case BatchUpdate:
		return map[string]interface{}{
			"method": "PUT",
			"to":     db.getRelativePath(r.HrefSelf),
			"body":   r.Data,
		}
	case BatchCreate:

		return map[string]interface{}{
			"method": "POST",
			"to":     fmt.Sprintf("/node/%d/relationships", r.StartID()),
			"body": map[string]interface{}{
				"to":   fmt.Sprintf("/node/%d", r.EndID()),
				"type": r.Type,
				"data": r.Data,
			},
		}

	case BatchDelete:
		return map[string]interface{}{
			"method": "DELETE",
			"to":     db.getRelativePath(r.HrefSelf),
		}

	}
	return nil
}
