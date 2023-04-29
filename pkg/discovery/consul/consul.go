package consul

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	consul "github.com/hashicorp/consul/api"
	"movieexample.com/pkg/discovery"
)

// Registry defines a COnsul-base service regstry
type Registry struct {
	client *consul.Client
}

// NewRegistry creates a new Consul-based service registry instance
func NewRegistry(addr string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{client: client}, nil
}

// Register creates aservice instance record in the regitry
func (r *Registry) Register(ctx context.Context, instanceID string, serviceName string, hostPost string) error {
	parts := strings.Split(hostPost, ":")
	if len(parts) != 2 {
		return errors.New("hostPort must be in the form <host>:<port>, example localhost:8081")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Address: parts[0],
		ID:      instanceID,
		Name:    serviceName,
		Port:    port,
		Check: &consul.AgentServiceCheck{
			CheckID: instanceID,
			TTL:     "5s",
		},
	})
}

// Deregister removes aservice instance record form the registry
func (r *Registry) Deregister(ctx context.Context, instanceID string, _ string) error {
	return r.client.Agent().ServiceDeregister(instanceID)
}

// ServiceAddresses returns the list of addresses of active instances of the given service
func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	} else if len(entries) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, v := range entries {
		res = append(res, fmt.Sprintf("%s:%d", v.Service.Address, v.Service.Port))
	}
	return res, nil
}

// ReportHealthStatus is a push mechanism for reporting healthy state to the registry
func (r *Registry) ReportHealthyStatus(instanceID string, serviceName string) error {
	return r.client.Agent().PassTTL(instanceID, "")
}
