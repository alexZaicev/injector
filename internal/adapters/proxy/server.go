package proxy

type Server interface {
	Serve() error
	Stop()
}
