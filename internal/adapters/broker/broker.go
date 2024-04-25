package broker

const (
	MetadataKeyError = "error"
)

type Broker interface {
	Start() error
	Stop()
}
