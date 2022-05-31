package util

import (
	"fmt"
	"sync"

	"github.com/pborman/uuid"
)

var uuidLock sync.Mutex

func NewUUID() string {
	uuidLock.Lock()
	defer uuidLock.Unlock()
	return uuid.NewUUID().String()
}

func GetCRIEndpoint(runtime string) string {
	return fmt.Sprintf("unix:///var/run/%s/%s.sock", runtime, runtime)
}
