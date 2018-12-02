package dns

import (
	"github.com/ivanthescientist/docker-dns-addon/internal/app/container"
	log "github.com/sirupsen/logrus"
	"sync"
)

// DomainRegistry is a thread-safe domain to address map, it also handles creation of domains through getDomain method
type DomainRegistry struct {
	domainIndex  map[string]container.Container
	indexMutex   *sync.RWMutex
	logger       *log.Logger
	domainSuffix string
}

// Create new registry using provided domainSuffix for domain creation
func NewDomainRegistry(logger *log.Logger, domainSuffix string) *DomainRegistry {
	return &DomainRegistry{
		domainIndex:  make(map[string]container.Container),
		indexMutex:   &sync.RWMutex{},
		logger:       logger,
		domainSuffix: domainSuffix,
	}
}

// Handle container event
func (r *DomainRegistry) HandleEvent(event container.Event) {
	switch event.Type {
	case container.EventContainerStarted:
		r.AddRecord(event.Container)
	case container.EventContainerStopped:
		r.RemoveRecord(event.Container)
	}
}

// Wait for all ongoing reads and then add domain record for a container
func (r *DomainRegistry) AddRecord(c container.Container) {
	r.indexMutex.Lock()
	defer r.indexMutex.Unlock()
	domain := r.getDomain(c)

	r.logger.Infof("Adding domain record: %s - %s", domain, c.String())
	r.domainIndex[domain] = c

}

// Wait for all ongoing reads and then remove domain record for a container
func (r *DomainRegistry) RemoveRecord(c container.Container) {
	r.indexMutex.Lock()
	defer r.indexMutex.Unlock()
	domain := r.getDomain(c)

	r.logger.Infof("Removing domain record: %s - %s", domain, c.String())
	delete(r.domainIndex, domain)
}

// Resolve domain address, if not present return an empty string
func (r *DomainRegistry) ResolveDomain(domain string) string {
	r.indexMutex.RLock()
	defer r.indexMutex.RUnlock()

	if c, ok := r.domainIndex[domain]; ok {
		return c.Addr
	}

	return ""
}

func (r *DomainRegistry) getDomain(c container.Container) string {
	return c.Name + r.domainSuffix
}
