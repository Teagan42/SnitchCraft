package heuristics

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type BadActorsCheck struct {
	BadIPs []string // This should be populated with known bad IPs
}

func (h BadActorsCheck) Name() string {
	return "known_bad_actor"
}

func (h BadActorsCheck) Check(r *http.Request) (string, bool) {
	if len(h.BadIPs) == 0 {
		h.BadIPs = GetBadIPs(badIPsURL)
	}
	ip := r.RemoteAddr
	for _, badIP := range h.BadIPs {
		if strings.HasPrefix(ip, badIP) {
			return "Request from known bad IP", true
		}
	}
	return "", false
}

var badIPsURL = "https://raw.githubusercontent.com/ramit-mitra/blocklist-ipsets/main/rottenIPs.json"

func GetBadIPs(url string) []string {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("[heuristics] failed to fetch bad IPs:", err)
		return []string{}
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("[heuristics] failed to close response body:", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[heuristics] failed to read response:", err)
		return []string{}
	}

	var result []string
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("[heuristics] failed to parse JSON:", err)
		return []string{}
	}
	return result
}

func init() {
	fmt.Println("[heuristics] registering BadActorsCheck heuristic...")
	RegisterHeuristic(BadActorsCheck{})
}
