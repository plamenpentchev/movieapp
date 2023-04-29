package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Registry defines a service registry
type Registry interface {

	//Register creates aservice instance record in the regitry
	Register(ctx context.Context, instanceID string, serviceName string, hostPost string) error
	//Deregister removes aservice instance record form the registry
	Deregister(ctx context.Context, instanceID string, serviceName string) error
	//ServiceAddresses returns the list of addresses of active instances of the given service
	ServiceAddresses(ctx context.Context, serviceID string) ([]string, error)
	//ReportHealthStatus is a push mechanism for reporting healthy state to the registry
	ReportHealthyStatus(instanceID string, serviceName string) error
}

// ErrNotFound is returnes when no service addresses are found
var ErrNotFound = errors.New("no service address found")
var ErrSrvcNotRegistered = errors.New("service is not registered")
var ErrINstanceNotRegistered = errors.New("service instance is not registered")

// GenerateInstanceID generates apseudo-random service instance identifier, using a service name
// suffixed by a random generated number
func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
