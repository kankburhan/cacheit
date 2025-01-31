package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/kankburhan/cacheit/pkg/cache"
)

var app_version = "dev"

func main() {
	cm := NewCacheManager()
	defer cm.Close()
	handleFlags(cm)
}

type CacheManager struct {
	cache *cache.Manager
}

func NewCacheManager() *CacheManager {
	cm, err := cache.NewManager()
	if err != nil {
		log.Fatal("Failed to initialize cache manager:", err)
	}

	return &CacheManager{
		cache: cm,
	}
}

func (c *CacheManager) Close() {
	// Placeholder for any future cleanup
}

func handleFlags(cm *CacheManager) {
	id := flag.String("id", "", "Retrieve cached data by ID")
	label := flag.String("l", "", "Label for the cached data (required when piping)")
	show := flag.Bool("show", false, "Show all cached entries")
	clearAll := flag.Bool("clear-all", false, "Clear all cached entries")
	clearOne := flag.String("clear-one", "", "Clear specific cache entry")
	output := flag.String("o", "", "Output to file")
	help := flag.Bool("help", false, "Show help")
	version := flag.Bool("version", false, "Show version")
	flag.Parse()

	if v := os.Getenv("TAKEIT_VERSION"); v != "" {
		app_version = v
	}

	switch {
	case *help:
		printHelp()
	case *version:
		fmt.Printf("cacheit version %s\n", app_version)
	case isPipedInput():
		handlePipeInput(cm, *label)
	case *show:
		handleShow(cm)
	case *clearAll:
		handleClearAll(cm)
	case *clearOne != "":
		handleClearOne(cm, *clearOne)
	case *id != "":
		handleRetrieve(cm, *id, *output)
	default:
		printHelp()
		os.Exit(0)
	}
}

func isPipedInput() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func handlePipeInput(cm *CacheManager, label string) {
	if label == "" {
		log.Fatal("Error: -l label is required when piping data")
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Error reading stdin:", err)
	}

	id, err := cm.cache.Save(label, data)
	if err != nil {
		log.Fatal("Error saving cache:", err)
	}

	fmt.Println(id)
}

func handleShow(cm *CacheManager) {
	entries, err := cm.cache.LoadMetadata()
	if err != nil {
		log.Fatal("Error loading metadata:", err)
	}

	fmt.Println("ID\tCommand\tLast Used")
	for _, entry := range entries {
		fmt.Printf("%s\t%s\t%s\n",
			entry.ID,
			truncateString(entry.Label, 40),
			entry.LastUsed.Format(time.RFC3339),
		)
	}
}

func handleClearAll(cm *CacheManager) {
	if err := cm.cache.ClearAll(); err != nil {
		log.Fatal("Error clearing cache:", err)
	}
	fmt.Println("All cache entries cleared")
}

func handleClearOne(cm *CacheManager, id string) {
	if err := cm.cache.ClearOne(id); err != nil {
		log.Fatal("Error clearing entry:", err)
	}
	fmt.Printf("Cache entry %s cleared\n", id)
}

func handleRetrieve(cm *CacheManager, id, output string) {
	data, err := cm.cache.Retrieve(id)
	if err != nil {
		log.Fatal("Error retrieving cache:", err)
	}

	if output != "" {
		if err := os.WriteFile(output, data, 0644); err != nil {
			log.Fatal("Error writing output:", err)
		}
		fmt.Printf("Data written to %s\n", output)
	} else {
		os.Stdout.Write(data)
	}
}

func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func printHelp() {
	fmt.Println(`
    _____            _          _____ _   
   / ____|          | |        |_   _| |  
   | |     __ _  ___| |__   ___  | | | |_ 
   | |    / _| |/ __| '_ \ / _ \ | | | __|
   | |___| (_| | (__| | | |  __/_| |_| |_ 
    \_____\__,_|\___|_| |_|\___|_____|\__|
                           by kankburhan

cacheit - CLI Data Caching Tool
Usage:
  command | cacheit -l "label"  Save piped data with label
  cacheit [options]

Options:
  -l 		      	 Label for piped data (required when piping)
  -id       		 Retrieve cached data by ID
  -show            	 Show all cached entries
  -clear-all       	 Clear all cached entries
  -clear-one   		 Clear specific cache entry
  -o         		 Output to file
  -version         	 Show version
  -help            	 Show this help

Examples:
  subfinder -d example.com | cacheit -l "subfinder scan"
  cacheit -id abc123 -o results.txt
  cacheit -show
  cacheit -clear-one abc123`)
}
