
package route

import (
	"context"
	"io"
	"net"
	"net/http"
	urlpkg "net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/tunnel"

	"github.com/go-chi/render"
)

const (
	maxBodySize = 2 * 1024 // 2KB
)

type ContentCheckResult struct {
	Proxy   string `json:"proxy"`
	URL     string `json:"url"`
	Content string `json:"content,omitempty"`
	Error   string `json:"error,omitempty"`
}

func contentCheckRouter() http.Handler {
	return http.HandlerFunc(contentCheck)
}

func contentCheck(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, newError("url query parameter is required"))
		return
	}

	timeoutQuery := r.URL.Query().Get("timeout")
	timeout, err := time.ParseDuration(timeoutQuery)
	if err != nil || timeout == 0 {
		timeout = 5 * time.Second
	}

	proxies := tunnel.Proxies()
	results := make([]*ContentCheckResult, 0, len(proxies))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, proxy := range proxies {
		// Exclude proxy groups to avoid recursion and unnecessary checks
		proxyType := proxy.Type().String()
		if proxyType == "Selector" || proxyType == "URLTest" || proxyType == "Fallback" || proxyType == "LoadBalance" {
			continue
		}

		wg.Add(1)
		go func(p C.Proxy) {
			defer wg.Done()

			result := &ContentCheckResult{
				Proxy: p.Name(),
				URL:   url,
			}

			content, err := checkProxyContent(p, url, timeout)
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Content = strings.TrimSpace(content)
			}

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(proxy)
	}

	wg.Wait()

	render.JSON(w, r, render.M{"results": results})
}

func checkProxyContent(p C.Proxy, url string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a custom transport that uses the proxy's DialContext
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// The address from http.Request is used to create the metadata
			u, err := urlpkg.Parse("http://" + addr)
			if err != nil {
				return nil, err
			}
			port, err := strconv.ParseUint(u.Port(), 10, 16)
			if err != nil {
				port = 80
			}

			metadata := &C.Metadata{
				Host:    u.Hostname(),
				DstPort: uint16(port),
				NetWork: C.TCP,
			}
			return p.DialContext(ctx, metadata)
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))
	if err != nil {
		return "", err
	}

	return string(body), nil
}
