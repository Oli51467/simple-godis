package resp

// Connection represent a Connection from Client
type Connection interface {
	Write([]byte) error
	GetDBIndex() int
	SelectDB(int)
}
