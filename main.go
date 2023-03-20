package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/barrett370/cloudflare_dynamicdns/cron"
)

type IPUpdateService struct {
	target string
	client *http.Client
}

const (
	envNextDNSTarget          = "NEXTDNS_TARGET"
	envNextDNSProfile         = "NEXTDNS_PROFILE"
	envNextDNSIntervalSeconds = "NEXTDNS_INTERVAL_SECONDS"
)

func NewIPUpdateService(target string) *IPUpdateService {
	return &IPUpdateService{
		target: target,
		client: http.DefaultClient,
	}
}

func (i IPUpdateService) Run(logger *log.Logger) error {
	resp, err := i.client.Get(i.target)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logger.Print(string(respBytes))
	return nil
}

func main() {
	nextDNSTarget := os.Getenv(envNextDNSTarget)
	if nextDNSTarget == "" {
		log.Fatalf("must set nextdns taget with env variable %s", envNextDNSTarget)
	}
	nextDNSProfile := os.Getenv(envNextDNSProfile)
	if nextDNSProfile == "" {
		nextDNSProfile = "unknown"
	}
	nextDNSIntervalSecondsStr := os.Getenv(envNextDNSIntervalSeconds)
	var (
		nextDNSIntervalSeconds = 60
		err                    error
	)
	if nextDNSIntervalSecondsStr != "" {
		nextDNSIntervalSeconds, err = strconv.Atoi(nextDNSIntervalSecondsStr)
		if err != nil {
			log.Fatalf("error parsing interval, err: %v", err)
		}
	}
	log.Printf("Starting update service: profile: %s, target: %s, interval seconds: %d\n", nextDNSProfile, nextDNSTarget, nextDNSIntervalSeconds)
	svc := NewIPUpdateService(nextDNSTarget)
	cronSvc := cron.New(fmt.Sprintf("DynamicNextDNS [%s]", nextDNSProfile), svc, time.Duration(nextDNSIntervalSeconds)*time.Second)
	cronSvc.Start()
	interruptC := make(chan os.Signal, 1)
	signal.Notify(interruptC, os.Interrupt, syscall.SIGTERM)
	<-interruptC
	log.Println("received os interrupt or kill, stopping update processes..")
	cronSvc.Stop()
}
