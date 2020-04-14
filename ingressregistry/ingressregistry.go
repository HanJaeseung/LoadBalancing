// Package registry defines the Registry interface which can be used with goproxy.
package ingressregistry

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

//-----------------------------------------------------------------------------------------------------------------------------

//var lock sync.RWMutex
//
//// Common errors.
//var (
//	ErrServiceNotFound = errors.New("service name/version not found")
//)
//
//// Registry is an interface used to lookup the target host
//// for a given service name / version pair.
//type Registry interface {
//	Add(path, endpoint string)                // Add an endpoint to our registry
//	Delete(path, endpoint string)             // Remove an endpoint to our registry
//	Failure(path, endpoint string, err error) // Mark an endpoint as failed.
//	Lookup(path string) ([]string, error)     // Return the endpoint list for the given service name/version
//}
//
//// DefaultRegistry is a basic registry using the following format:
//// {
////   "serviceName": {
////     "serviceVersion": [
////       "endpoint1:port",
////       "endpoint2:port"
////     ],
////   },
//// }
//
//// DefaultRegistry is a basic registry using the following format:
//// {
////   "Host": {
////     "Path": [
////       "endpoint1:port",
////       "endpoint2:port"
////     ],
////   },
//// }
//
//
//type DefaultRegistry map[string][]string
//
//// Lookup return the endpoint list for the given service name/version.
//func (r DefaultRegistry) Lookup(path string) ([]string, error) {
//	fmt.Println("Lookup")
//	lock.RLock()
//	targets, ok := r[path]
//
//	fmt.Println(r[path])
//	lock.RUnlock()
//	fmt.Println("LookupCheack1")
//	fmt.Println(targets)
//	if !ok {
//		return nil, ErrServiceNotFound
//	}
//	return targets, nil
//}
//
//// Failure marks the given endpoint for service name/version as failed.
//func (r DefaultRegistry) Failure(path, endpoint string, err error) {
//	// Would be used to remove an endpoint from the rotation, log the failure, etc.
//	//log.Printf("Error accessing %s/%s (%s): %s", path, endpoint, err)
//	log.Printf("Error accessing %s (%s): %s", path, endpoint, err)
//}
//
//// Add adds the given endpoit for the service name/version.
//func (r DefaultRegistry) Add(path, endpoint string) {
//	fmt.Println("Add")
//	lock.Lock()
//	defer lock.Unlock()
//
//	service, ok := r[path]
//	if !ok {
//		service = []string{}
//		r[path] = service
//	}
//	service = append(service, endpoint)
//}
//
//// Delete removes the given endpoit for the service name/version.
//func (r DefaultRegistry) Delete(path, endpoint string) {
//	fmt.Println("Delete")
//	lock.Lock()
//	defer lock.Unlock()
//
//	service, ok := r[path]
//	if !ok {
//		return
//	}
//begin:
//	for i, svc := range service {
//		if svc == endpoint {
//			copy(service[i:], service[i+1:])
//			service[len(service)-1] = ""
//			service = service[:len(service)-1]
//			goto begin
//		}
//	}
//}

//---------------------------------------------------------------------------------------------------------
var lock sync.RWMutex

// Common errors.
var (
	ErrServiceNotFound = errors.New("service name/version not found")
)

// Registry is an interface used to lookup the target host
// for a given service name / version pair.
type Registry interface {
	Add(host, path, endpoint string)                // Add an endpoint to our registry
	Delete(host, path, endpoint string)             // Remove an endpoint to our registry
	Failure(host, path, endpoint string, err error) // Mark an endpoint as failed.
	Lookup(host, path string) ([]string, error)     // Return the endpoint list for the given service name/version
}

// DefaultRegistry is a basic registry using the following format:
// {
//   "Host": {
//     "Path": [
//       "cluster1",
//       "cluster2"
//     ],
//   },
// }

//type DefaultRegistry map[string]map[string]map[string]string
type DefaultRegistry map[string]map[string][]string


// Lookup return the endpoint list for the given service name/version.

func (r DefaultRegistry) Lookup(host string, path string) ([]string, error) {
	fmt.Println("Lookup")
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
	// Would be used to remove an endpoint from the rotation, log the failure, etc.
	//log.Printf("Error accessing %s/%s (%s): %s", path, endpoint, err)
	log.Printf("Error accessing %s %s (%s): %s", host, path, endpoint, err)
}

func (r DefaultRegistry) Add(host, path, endpoint string) {
	fmt.Println("Add")
	lock.Lock()
	defer lock.Unlock()

	service, ok := r[host]
	if !ok {
		service = map[string][]string{}
		r[host] = service
	}
	service[path] = append(service[path], endpoint)
}


// Delete removes the given endpoit for the service name/version.
func (r DefaultRegistry) Delete(host, path, endpoint string) {
	fmt.Println("Delete")
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