package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	pool *ProxyPool
)

type ProxyItem struct {
	URL           *url.URL
	LastUsed      time.Time
	UseCount      int
	MaxUses       int
	ResetInterval time.Duration
	Disabled      bool
	DisableUntil  time.Time
}

type ProxyPool struct {
	proxies []*ProxyItem
	mu      sync.Mutex
}

func NewProxyPool(proxyURLs []string, maxUses int, resetInterval time.Duration) *ProxyPool {
	pool := &ProxyPool{}
	for _, urlString := range proxyURLs {
		parsedURL, err := url.Parse(urlString)
		if err != nil {
			log.Printf("Invalid proxy URL: %s", urlString)
			continue
		}
		pool.proxies = append(pool.proxies, &ProxyItem{
			URL:           parsedURL,
			MaxUses:       maxUses,
			ResetInterval: resetInterval,
		})
	}
	return pool
}

func (p *ProxyPool) ChooseProxy() *ProxyItem {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	var chosenProxy *ProxyItem
	for i, proxy := range p.proxies {
		if proxy.Disabled && now.Before(proxy.DisableUntil) {
			continue
		}
		if now.Sub(proxy.LastUsed) > proxy.ResetInterval {
			proxy.UseCount = 0
		}
		if proxy.UseCount < proxy.MaxUses || proxy.MaxUses == 0 {
			proxy.LastUsed = now
			proxy.UseCount++
			proxy.Disabled = false
			chosenProxy = proxy
			// 将选中的代理移动到队列末尾
			p.proxies = append(p.proxies[:i], p.proxies[i+1:]...)
			p.proxies = append(p.proxies, chosenProxy)
			break
		}
	}
	return chosenProxy
}

func (p *ProxyPool) DisableProxy(proxy *ProxyItem, duration time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	proxy.Disabled = true
	proxy.DisableUntil = time.Now().Add(duration)
}

func init() {
	// 初始化代理池，示例中的代理和次数限制是示意性的
	proxyURLs := []string{"43.153.90.69:443", "203.189.89.106:80"}
	pool = NewProxyPool(proxyURLs, 100, time.Second)
}

func crawler() {
	client := &http.Client{}

	for {
		proxy := pool.ChooseProxy()
		if proxy == nil {
			log.Println("No available proxies")
			break
		}

		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxy.URL),
		}

		response, err := client.Get(fmt.Sprintf("http://localhost:%d", port))
		if err != nil {
			log.Printf("Request failed: %v", err)
			continue
		}

		if response.StatusCode == http.StatusTooManyRequests {
			log.Println("Received 429 Too Many Requests")
			if response.Header.Get("Retry-After") == "" {
				log.Println("No Retry-After header, disabling proxy")
				pool.DisableProxy(proxy, 0)
			} else {
				retryAfter, _ := time.ParseDuration(response.Header.Get("Retry-After") + "s")
				pool.DisableProxy(proxy, retryAfter)
			}
		}

		body, _ := io.ReadAll(response.Body)
		response.Body.Close()
		log.Printf("Response body: %s", string(body))

		// 为了示例简单，我们只循环一次
		break
	}
}
