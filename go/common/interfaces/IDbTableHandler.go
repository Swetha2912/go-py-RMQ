package interfaces

type IDbTableHandler interface {
	Find(interface{}) (interface{}, error)
	FindAll(interface{}) (interface{}, error)
	Insert(interface{}) error
	Delete(interface{}) error
	Exists(interface{}) bool
	DeleteAll(interface{}) (interface{}, error)
	Update(interface{}, interface{}) error
	Aggregation(interface{}) (interface{},error)
	Paginate(interface{}, interface{}, interface{}, interface{}, string, *interface{}) error
	
}
