package infrastructures

import (
	"github.com/globalsign/mgo"
	"sample_rmq/common/interfaces"
)

type MongoDB struct {
	Conn *mgo.Session
}

func (dbHandler MongoDB) GetTable(tableName string) interfaces.IDbTableHandler {
	dbName := "test"
	collection := dbHandler.Conn.DB(dbName).C(tableName)
	return MongoDBCollection{
		Conn:   dbHandler.Conn,
		DBName: dbName,
		C:      collection,
	}
}

// ProvideMongoDB returns a MongoDB
func ProvideMongoDB(connString string) interfaces.IDbHandler {
	conn, err := mgo.Dial(connString)
	if err != nil {
		panic("cannot connect to database")
	}
	return MongoDB{Conn: conn}
}
