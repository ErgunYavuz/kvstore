package iface

type NodeAPI interface {
	HandlePut(key, value string) error
	HandleGet(key string) (string, error)
	Delete(key string) (bool, error)
}
