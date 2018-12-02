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

// NewDomainRegistry creates new registry using provided domainSuffix for domain creation
func NewDomainRegistry(logger *log.Logger, domainSuffix string) *DomainRegistry {
	return &DomainRegistry{
		domainIndex:  make(map[string]container.Container),
		indexMutex:   &sync.RWMutex{},
		logger:       logger,
		domainSuffix: domainSuffix,
	}
}

// HandleEvent translates container events into appropriate actions e.g. adding or removing domain record for container
func (r *DomainRegistry) HandleEvent(event container.Event) {
	switch event.Type {
	case container.EventContainerStarted:
		r.AddRecord(event.Container)
	case container.EventContainerStopped:
		r.RemoveRecord(event.Container)
	}
}

// AddRecord constructs a domain and adds a domain record for it, blocks until all current reads are done
func (r *DomainRegistry) AddRecord(c container.Container) {
	r.indexMutex.Lock()
	defer r.indexMutex.Unlock()
	domain := r.getDomain(c)

	r.logger.Infof("Adding domain record: %s - %s", domain, c.String())
	r.domainIndex[domain] = c

}

// RemoveRecord constructs a domain and removes corresponding domain record, blocks until all current reads are done
func (r *DomainRegistry) RemoveRecord(c container.Container) {
	r.indexMutex.Lock()
	defer r.indexMutex.Unlock()
	domain := r.getDomain(c)

	r.logger.Infof("Removing domain record: %s - %s", domain, c.String())
	delete(r.domainIndex, domain)
}

// ResolveDomain resolves domain address, if not present returns an empty string
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
