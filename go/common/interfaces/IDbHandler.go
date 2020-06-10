package interfaces

type IDbHandler interface {
	GetTable(string) IDbTableHandler
}