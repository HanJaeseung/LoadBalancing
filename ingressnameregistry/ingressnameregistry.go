package ingressnameregistry

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

var lock sync.RWMutex

// Common errors.
var (
	ErrServiceNotFound = errors.New("service name/version not found")
)

// Registry is an interface used to lookup the target host
// for a given service name / version pair.
type Registry interface {
	Add(ingressName, url string)
	Delete(url string)
	//Delete(host, path, endpoint string)             // Remove an endpoint to our registry
	Failure(host, path, endpoint string, err error) // Mark an endpoint as failed.
	Lookup(ingressName string) ([]string, error)
}


// DefaultRegistry is a basic registry using the following format:
// {
//   "IngressName": [
//       "keti.test.com/test",
//       "lb_test.com/service",
//     ],
// }

//type DefaultRegistry map[string]map[string]map[string]stringzmgma
type DefaultRegistry map[string][]string

// Lookup return the endpoint list for the given service name/version.

func (r DefaultRegistry) Add(ingressName, url string) {
	fmt.Println("*****Ingress Name Add*****")
	lock.Lock()
	defer lock.Unlock()

	service, ok := r[ingressName]
	if !ok {
		service = []string{}
		r[ingressName] = service
	}
	service = append(service, url)
}


func (r DefaultRegistry) Lookup(ingressName string) ([]string, error) {
	fmt.Println("----Lookup----")
	lock.RLock()
	targets, ok := r[ingressName]
	lock.RUnlock()
	if !ok {
		return nil, ErrServiceNotFound
	}
	return targets, nil
}

func (r DefaultRegistry) Failure(host, path, endpoint string, err error) {
	// Would be used to remove an endpoint from the rotation, log the failure, etc.
	//log.Printf("Error accessing %s/%s (%s): %s", path, endpoint, err)
	log.Printf("Error accessing %s %s (%s): %s", host, path, endpoint, err)
}

func (r DefaultRegistry) Delete(ingressName string) {
	fmt.Println("*****Delete*****")
	lock.Lock()
	defer lock.Unlock()

	_, ok := r[ingressName]
	if !ok {
		return
	}

	delete(r, ingressName)
}

//// Delete removes the given endpoit for the service name/version.
//func (r DefaultRegistry) Delete(host, path, endpoint string) {
//	fmt.Println("----Delete----")
//	lock.Lock()
//	defer lock.Unlock()
//
//	service, ok := r[host]
//	if !ok {
//		return
//	}
//
//begin:
//	for i, svc := range service[path] {
//		if svc == endpoint {
//			copy(service[path][i:], service[path][i+1:])
//			service[path] = service[path][:len(service[path])-1]
//			goto begin
//		}
//	}
//}
