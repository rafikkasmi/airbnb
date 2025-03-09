package utils

import (
	"log"
	"math/rand"
	"net/url"
	"sync"
	"time"
)

// ProxyRotator manages a pool of proxies and provides methods to get proxies in a rotating fashion
type ProxyRotator struct {
	proxies     []*url.URL
	currentIdx  int
	mutex       sync.Mutex
	lastRotated time.Time
}

// NewProxyRotator creates a new ProxyRotator with the given proxy URLs
func NewProxyRotator(proxyURLs []string) (*ProxyRotator, error) {
	var proxies []*url.URL
	
	for _, proxyStr := range proxyURLs {
		if proxyStr == "" {
			continue
		}
		
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			return nil, err
		}
		proxies = append(proxies, proxyURL)
	}
	
	// Initialize with a random starting point
	startIdx := 0
	if len(proxies) > 0 {
		startIdx = rand.Intn(len(proxies))
	}
	
	return &ProxyRotator{
		proxies:     proxies,
		currentIdx:  startIdx,
		lastRotated: time.Now(),
	}, nil
}

// GetNextProxy returns the next proxy in the rotation
// If no proxies are available, returns nil
func (pr *ProxyRotator) GetNextProxy() *url.URL {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()
	
	if len(pr.proxies) == 0 {
		return nil
	}
	
	proxy := pr.proxies[pr.currentIdx]
	pr.currentIdx = (pr.currentIdx + 1) % len(pr.proxies)
	pr.lastRotated = time.Now()
	
	log.Printf("Rotating to proxy: %s", proxy.String())
	return proxy
}

// GetRandomProxy returns a random proxy from the pool
// If no proxies are available, returns nil
func (pr *ProxyRotator) GetRandomProxy() *url.URL {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()
	
	if len(pr.proxies) == 0 {
		return nil
	}
	
	idx := rand.Intn(len(pr.proxies))
	proxy := pr.proxies[idx]
	pr.lastRotated = time.Now()
	
	log.Printf("Using random proxy: %s", proxy.String())
	return proxy
}

// GetCurrentProxy returns the current proxy without rotating
// If no proxies are available, returns nil
func (pr *ProxyRotator) GetCurrentProxy() *url.URL {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()
	
	if len(pr.proxies) == 0 {
		return nil
	}
	
	return pr.proxies[pr.currentIdx]
}

// Count returns the number of proxies in the pool
func (pr *ProxyRotator) Count() int {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()
	
	return len(pr.proxies)
}
