package iface

type NodeAPI interface {
	HandlePut(requesterID int, key, value string) error
	HandleGet(key string) (string, error)
	HandleDelete(requesterID int, key string) error
}
