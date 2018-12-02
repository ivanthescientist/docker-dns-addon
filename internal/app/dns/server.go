package dns

import (
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
)

type Server struct {
	dnsServer *dns.Server
	registry  *DomainRegistry
	logger    *log.Logger
}

func NewServer(logger *log.Logger, host string, port int, protocol string, registry *DomainRegistry) *Server {
	server := &Server{}

	dnsServer := &dns.Server{
		Addr:    host + ":" + strconv.Itoa(port),
		Net:     protocol,
		Handler: server,
	}

	server.dnsServer = dnsServer
	server.registry = registry
	server.logger = logger

	return server
}

func (s *Server) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	var err error
	resp := new(dns.Msg)
	resp.SetReply(r)
	// Return an empty response in worst case scenario
	defer func() {
		err = w.WriteMsg(resp)
		if err != nil {
			s.logger.Errorf("Failed to write dns response: %s", err)
		}
	}()

	if len(r.Question) == 0 {
		return
	}

	requestName := r.Question[0].Name
	s.logger.Infof("Received resolution request for: %s", requestName)

	addr := s.registry.ResolveDomain(requestName)
	if addr == "" {
		return
	}

	resp.Authoritative = true
	resp.Answer = []dns.RR{
		&dns.A{
			Hdr: dns.RR_Header{
				Name:   r.Question[0].Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    0,
			},
			A: net.ParseIP(addr),
		},
	}
}

func (s *Server) ListenAndServe() error {
	s.logger.Info("Starting DNS server")
	return s.dnsServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	s.logger.Info("Shutting down DNS server")
	return s.dnsServer.Shutdown()
}
