package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Config structure for the configuration file
type Config struct {
	Domains       []DomainConfig `json:"domains"`
	CheckInterval int            `json:"check_interval_minutes"`
	LogFile       string         `json:"log_file"`
}

// DomainConfig structure for each domain
type DomainConfig struct {
	Domain   string `json:"domain"`
	Host     string `json:"host"` // e.g., "@" or "www"
	Password string `json:"password"`
}

// IPState to track current and previous IP
type IPState struct {
	CurrentIP  string
	PreviousIP string
}

const (
	namecheapURL = "https://dynamicdns.park-your-domain.com/update"
	ipCheckURL   = "https://api.ipify.org"
)

func main() {
	// Set up logging
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	logFile, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "DDNS: ", log.LstdFlags)

	fmt.Println("Starting Dynamic DNS updater...")
	logger.Println("Starting Dynamic DNS updater")

	// Initialize IP state
	ipState := &IPState{}

	interval := time.Duration(config.CheckInterval) * time.Minute

	// Main loop
	for {
		startTime := time.Now()
		
		fmt.Printf("\nChecking IP at %s\n", startTime.Format(time.RFC1123))
		
		currentIP, err := getPublicIP()
		if err != nil {
			msg := fmt.Sprintf("Error getting public IP: %v", err)
			fmt.Println(msg)
			logger.Printf(msg)
			countdown(interval, startTime)
			continue
		}

		ipState.CurrentIP = currentIP

		if ipState.CurrentIP != ipState.PreviousIP {
			fmt.Printf("IP changed from %s to %s\n", ipState.PreviousIP, ipState.CurrentIP)
			logger.Printf("IP changed from %s to %s", ipState.PreviousIP, ipState.CurrentIP)
			
			// Update all domains
			for i, domain := range config.Domains {
				fmt.Printf("Updating %d/%d: %s.%s...", 
					i+1, len(config.Domains), domain.Host, domain.Domain)
				err := updateDNS(domain, ipState.CurrentIP, logger)
				if err != nil {
					msg := fmt.Sprintf("FAILED: %v", err)
					fmt.Println(msg)
					logger.Printf("Error updating DNS for %s: %v", domain.Domain, err)
				} else {
					msg := "SUCCESS"
					fmt.Println(msg)
					logger.Printf("Successfully updated DNS for %s.%s to %s", 
						domain.Host, domain.Domain, ipState.CurrentIP)
				}
			}
			
			ipState.PreviousIP = ipState.CurrentIP
		} else {
			fmt.Println("No IP change detected")
		}

		countdown(interval, startTime)
	}
}

// countdown displays a timer until the next check
func countdown(interval time.Duration, startTime time.Time) {
	targetTime := startTime.Add(interval)
	
	for {
		remaining := time.Until(targetTime)
		if remaining <= 0 {
			break
		}
		
		fmt.Printf("\rNext check in: %s", remaining.Round(time.Second))
		time.Sleep(1 * time.Second)
	}
	fmt.Println() // New line after countdown
}

// loadConfig reads and parses the configuration file
func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	// Set default check interval if not specified
	if config.CheckInterval == 0 {
		config.CheckInterval = 5 // 5 minutes default
	}

	return &config, nil
}

// getPublicIP retrieves the current public IP address
func getPublicIP() (string, error) {
	resp, err := http.Get(ipCheckURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(ip)), nil
}

// updateDNS sends the update request to Namecheap
func updateDNS(domain DomainConfig, ip string, logger *log.Logger) error {
	url := fmt.Sprintf("%s?host=%s&domain=%s&password=%s&ip=%s",
		namecheapURL,
		domain.Host,
		domain.Domain,
		domain.Password,
		ip)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
