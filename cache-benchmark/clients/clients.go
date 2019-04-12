package clients

// Cmd is the operation command
type Cmd struct {
	Name  string
	Key   string
	Value string
	Error error
}

// Client is the cache benchmark client
type Client interface {
	Run(*Cmd)
	PipelineRun([]*Cmd)
}

func New(typ, server string) Client {
	switch typ {
	case "redis":
		panic("not implemented")
		//return newRedisClient(server)
	case "http":
		return newHTTPClient(server)
	case "tcp":
		return newTCPClient(server)
	default:
		panic("unkown client type: " + typ)
	}
}
