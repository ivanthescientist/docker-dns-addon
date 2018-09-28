package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/miekg/dns"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

type ContainerRecord struct {
	ID     string
	Name   string
	Domain string
	Addr   string // not net.Addr for readability purposes
}

func (record ContainerRecord) String() string {
	return fmt.Sprintf("[%s] %s %s %s", record.ID, record.Name, record.Domain, record.Addr)
}

// Domain to ContainerRecord
var containerMap = map[string]*ContainerRecord{}
var containerMapMutex = &sync.RWMutex{}

const DomainSuffix = ".docker."

func DockerHandler(w dns.ResponseWriter, r *dns.Msg) {
	defer w.Close()
	if len(r.Question) > 1 {
		w.Close()
		return
	}

	containerMapMutex.RLock()
	defer containerMapMutex.RUnlock()

	requestName := r.Question[0].Name
	log.Printf("Received resolution request for: %s", requestName)

	resp := new(dns.Msg)
	resp.SetReply(r)
	defer w.WriteMsg(resp)

	containerRecord, isPresent := containerMap[requestName]
	if !isPresent {
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
			A: net.ParseIP(containerRecord.Addr),
		},
	}
}

func DomainFromName(containerName string) string {
	return strings.Replace(containerName, "/", "", -1) + DomainSuffix
}

func main() {
	var bindAddr string
	var err error

	if os.Getenv("SUDO_USER") != "" {
		bindAddr = "127.0.0.1:53"
	} else {
		bindAddr = "127.0.0.1:5300"
	}

	log.Printf("DNS server bind address: %s", bindAddr)

	dockerClient, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("Failed to connect to docker: %s", err)
	}

	containers, err := dockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Fatalf("Failed to fetch list of containers: %s", err)
	}

	log.Printf("Populating initial contianer domain records:")
	log.Printf("[ID] Name Domain IPAddr")
	for _, container := range containers {
		record, err := GetContainerInfo(container.ID, dockerClient)
		if err != nil {
			log.Printf("Failed to fetch additional container info: %s", err)
		}
		containerMap[record.Domain] = record

		log.Print(record)
	}

	filter := filters.NewArgs()
	filters.NewArgs().Add("type", "container")

	eventCh, errCh := dockerClient.Events(context.Background(), types.EventsOptions{
		Filters: filter,
	})
	go DockerEventHandler(eventCh, errCh, dockerClient)

	dns.HandleFunc("docker.", DockerHandler)
	err = dns.ListenAndServe(bindAddr, "udp", nil)
	if err != nil {
		log.Println(err)
	}
}

func GetContainerInfo(id string, dockerClient client.APIClient) (*ContainerRecord, error) {
	containerInfo, err := dockerClient.ContainerInspect(context.Background(), id)
	if err != nil {
		return nil, err
	}
	domain := DomainFromName(containerInfo.Name)
	record := &ContainerRecord{
		Domain: domain,
		ID:     containerInfo.ID,
		Name:   containerInfo.Name,
		Addr:   containerInfo.NetworkSettings.IPAddress,
	}
	return record, nil
}

func DockerEventHandler(eventCh <-chan events.Message, errCh <-chan error, dockerClient client.APIClient) {
	for {
		select {
		case event := <-eventCh:
			switch event.Action {
			case "start":
				err := AddContainerRecord(event.ID, dockerClient)
				if err != nil {
					log.Printf("Failed to add container record: %s", err)
				}
				log.Println(containerMap)
			case "die", "stop", "kill":
				err := RemoveContainerRecord(event.ID)
				if err != nil {
					log.Printf("Failed to remove container record: %s", err)
				}
				log.Println(containerMap)
			}
		case err := <-errCh:
			log.Println(err)
			return
		}
	}
}

func FindContainerRecordByID(id string) *ContainerRecord {
	var result *ContainerRecord

	for _, record := range containerMap {
		if record.ID == id {
			result = record
		}
	}

	return result
}

func AddContainerRecord(id string, dockerClient client.APIClient) error {
	containerMapMutex.Lock()
	defer containerMapMutex.Unlock()

	containerRecord, err := GetContainerInfo(id, dockerClient)
	if err != nil {
		return err
	}

	containerMap[containerRecord.Domain] = containerRecord

	return nil
}

func RemoveContainerRecord(id string) error {
	containerMapMutex.Lock()
	defer containerMapMutex.Unlock()

	record := FindContainerRecordByID(id)
	if record != nil {
		delete(containerMap, record.Domain)
	}

	return nil
}
