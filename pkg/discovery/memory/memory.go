package memory

import (
	"context"
	"sync"
	"time"

	"movieexample.com/pkg/discovery"
)

type serviceName string
type instanceID string
type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

// Registry defines an in-memory service registry
type Registry struct {
	sync.RWMutex
	serviceAddrs map[serviceName]map[instanceID]*serviceInstance
}

// NewRegsitry creates a new in-memory registry
func NewRegistry() *Registry {
	return &Registry{serviceAddrs: map[serviceName]map[instanceID]*serviceInstance{}}
}

// Register registers a service with the in-memory registry
func (r *Registry) Register(ctx context.Context, instID string, srvName string, hostPort string) error {
	r.Lock()
	defer r.Unlock()
	sn := serviceName(srvName)
	iID := instanceID(instID)
	if _, ok := r.serviceAddrs[sn]; !ok {
		r.serviceAddrs[sn] = map[instanceID]*serviceInstance{}
	}
	r.serviceAddrs[sn][iID] = &serviceInstance{
		hostPort:   hostPort,
		lastActive: time.Now(),
	}
	return nil
}

// Deregister removes a service from the registry
func (r *Registry) Deregister(ctx context.Context, srvName string, instID string) error {
	r.Lock()
	defer r.Unlock()

	sn := serviceName(srvName)
	if _, ok := r.serviceAddrs[sn]; !ok {
		return nil
	}
	iID := instanceID(instID)
	delete(r.serviceAddrs[sn], iID)
	return nil
}

// ServiceAddresses returns the list of addresses of active instances of the given service
func (r *Registry) ServiceAddresses(ctx context.Context, srvName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	sn := serviceName(srvName)
	if _, ok := r.serviceAddrs[sn]; !ok {
		return nil, discovery.ErrNotFound
	}
	col := r.serviceAddrs[sn]
	if len(col) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, v := range col {
		if v.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		res = append(res, v.hostPort)
	}
	return res, nil
}

// ReportHealthStatus is a push mechanism for reporting healthy state to the registry
func (r *Registry) ReportHealthyStatus(instID string, srvName string) error {
	r.Lock()
	defer r.Unlock()
	sn := serviceName(srvName)
	srvcs, ok := r.serviceAddrs[sn]
	if !ok {
		return discovery.ErrSrvcNotRegistered
	}
	iID := instanceID(instID)
	inst, ok := srvcs[iID]
	if !ok {
		return discovery.ErrINstanceNotRegistered
	}
	inst.lastActive = time.Now()
	return nil
}
