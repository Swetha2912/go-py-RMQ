package infrastructures

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/globalsign/mgo"
	"strconv"
	"sample_rmq/common/utilities"
)

type MongoDBCollection struct {
	// real mongodb connection object here
	Conn   *mgo.Session
	DBName string
	C      *mgo.Collection
	Query  *mgo.Query
}

// PaginateResult with docs structure
type PaginateResult struct {
	TotalCount  int                      `json:"total"`
	Limit       int                      `json:"limit"`
	CurrentPage int                      `json:"current_page"`
	Docs        []map[string]interface{} `json:"docs"`
}


func (collection MongoDBCollection) Insert(payload interface{}) error {
	return collection.C.Insert(payload)
}

func (collection MongoDBCollection) Find(payload interface{}) (interface{}, error) {

	hystrix.ConfigureCommand("collection_find", hystrix.CommandConfig{Timeout: 5000, RequestVolumeThreshold: 10})
	output := make(chan interface{}, 1)
	errors := hystrix.Go("collection_find", func() error {
		var row interface{}
		err := collection.C.Find(payload).One(&row)
		if err == nil {
			output <- row
		}
		return err
	}, nil)

	select {
	case response := <-output:
		return response, nil
	case err := <-errors:
		fmt.Println("hystrix error", err)
		return nil, err
	}

}

func (collection MongoDBCollection) FindAll(payload interface{}) (interface{}, error) {

	hystrix.ConfigureCommand("collection_find_all", hystrix.CommandConfig{
		Timeout:               1000,
		MaxConcurrentRequests: 50,
		//ErrorPercentThreshold: 25,
	})

	output := make(chan []interface{}, 1)
	errors := hystrix.Go("collection_find_all", func() error {
		var rows []interface{}
		err := collection.C.Find(payload).All(&rows)
		if err == nil {
			output <- rows
		}
		return err
	}, nil)

	select {
	case response := <-output:
		return response, nil
	case err := <-errors:
		hystrix.GetCircuit("collection_find")
		fmt.Println("hystrix error", err)
		return nil, err
	}

}

func (collection MongoDBCollection) Delete(payload interface{}) error {
	return collection.C.Remove(payload)
}

func (collection MongoDBCollection) Exists(payload interface{}) bool {
	count, err := collection.C.Find(payload).Count()

	if count > 0 {
		return true
	}
	if err != nil || count == 0 {
		return false
	}
	return false
}

func (collection MongoDBCollection) DeleteAll(payload interface{}) (interface{}, error) {
	data, err := collection.C.RemoveAll(payload)
	return data, err
}

func (collection MongoDBCollection) Update(filter interface{}, payload interface{}) error {
	return collection.C.Update(filter, payload)
}

func (collection MongoDBCollection) Aggregation(payload interface{}) (interface{}, error) {
	var row interface{}
	err := collection.C.Pipe(payload).One(&row)
	return row, err
}

func (collection MongoDBCollection) Paginate(query interface{}, selectFilter interface{}, _page interface{}, _limit interface{}, sort string, pRes *interface{}) error {
	var page int
	var limit int

	var results []map[string]interface{}

	if _page == nil {
		page = 1
	} else {
		switch _page.(type) {
		case int:
			page = _page.(int)
			break
		case float64:
			page = int(_page.(float64))
			break
		case string:
			page, _ = strconv.Atoi(_page.(string))
			break
		}
	}

	if _limit == nil {
		limit = 10
	} else {

		switch _limit.(type) {
		case int:
			limit = _limit.(int)
			break
		case float64:
			limit = int(_limit.(float64))
			break
		case string:
			limit, _ = strconv.Atoi(_limit.(string))
			break
		}
	}

	count, _ := collection.C.Find(query).Count()

	var pResult utilities.PaginateResult
	pResult.TotalCount = count
	pResult.Limit = limit
	pResult.CurrentPage = page

	mongoQuery := collection.C.Find(query).Select(selectFilter).Skip(limit * (page - 1)).Limit(limit)
	if sort != "" {
		mongoQuery = mongoQuery.Sort(sort)
	}

	err := mongoQuery.All(&results)
	pResult.Docs = results

	if results == nil {
		pResult.Docs = make([]map[string]interface{}, 0)
	}

	*pRes = pResult

	return err
}

// ProvideMongoDBCollection -
func ProvideMongoDBCollection(db MongoDB, dbName string, cName string) MongoDBCollection {
	coll := db.Conn.DB(dbName).C(cName)
	return MongoDBCollection{Conn: db.Conn, DBName: dbName, C: coll}
}
