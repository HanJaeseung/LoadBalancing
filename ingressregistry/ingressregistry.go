package ingressregistry

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

var lock sync.RWMutex

var (
	ErrServiceNotFound = errors.New("service name/version not found")
)


type Registry interface {
	Add(host, path, endpoint string)                // Add an endpoint to our registry
	Delete(host, path, endpoint string)             // Remove an endpoint to our registry
	Failure(host, path, endpoint string, err error) // Mark an endpoint as failed.
	Lookup(host, path string) ([]string, error)     // Return the endpoint list for the given service name/version
}

type DefaultRegistry map[string]map[string][]string


func (r DefaultRegistry) Lookup(host string, path string) ([]string, error) {
	fmt.Println("----Lookup----")
	fmt.Println(host)
	fmt.Println(path)
	lock.RLock()
	targets, ok := r[host][path]
	lock.RUnlock()
	if !ok {
		return nil, ErrServiceNotFound
	}
	return targets, nil
}


func (r DefaultRegistry) Failure(host, path, endpoint string, err error) {
	log.Printf("Error accessing %s %s (%s): %s", host, path, endpoint, err)
}

func (r DefaultRegistry) Add(host, path, endpoint string) {
	fmt.Println("----Add----")
	lock.Lock()
	defer lock.Unlock()

	service, ok := r[host]
	if !ok {
		service = map[string][]string{}
		r[host] = service
	}
	service[path] = append(service[path], endpoint)
}


func (r DefaultRegistry) Delete(host, path, endpoint string) {
	fmt.Println("----Delete----")
	lock.Lock()
	defer lock.Unlock()

	service, ok := r[host]
	if !ok {
		return
	}
begin:
	for i, svc := range service[path] {
		if svc == endpoint {
			copy(service[path][i:], service[path][i+1:])
			service[path][len(service)-1] = ""
			service[path] = service[path][:len(service)-1]
			goto begin
		}
	}
}
