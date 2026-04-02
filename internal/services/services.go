package services

import (
	"sort"
	"strings"
	"time"
	repo "todo-golang/internal/repositories"
	"todo-golang/internal/utils"
)

type DashboardSummary struct {
	TotalRequests       int `json:"totalRequests"`
	UniqueIPs           int `json:"uniqueIPs"`
	BlockedRequests     int `json:"blockedRequests"`
	RateLimitedRequests int `json:"rateLimitedRequests"`
	AllowedRequests     int `json:"allowedRequests"`
}

type IPMetric struct {
	IP       string `json:"ip"`
	Requests int    `json:"requests"`
	Country  string `json:"country"`
	Status   string `json:"status"`
}

type EndpointMetric struct {
	Endpoint string `json:"endpoint"`
	Count    int    `json:"count"`
	Status   string `json:"status"`
}

type DashboardEvent struct {
	Timestamp string `json:"timestamp"`
	IP        string `json:"ip"`
	Endpoint  string `json:"endpoint"`
	Method    string `json:"method"`
	Action    string `json:"action"`
	Country   string `json:"country"`
	Score     int    `json:"score"`
}

type DashboardResponse struct {
	Range        string           `json:"range"`
	Summary      DashboardSummary `json:"summary"`
	TopIPs       []IPMetric       `json:"topIPs"`
	TopEndpoints []EndpointMetric `json:"topEndpoints"`
	RecentEvents []DashboardEvent `json:"recentEvents"`
}

type ipAccumulator struct {
	count    int
	country  string
	severity int
	status   string
}

type endpointAccumulator struct {
	count    int
	severity int
	status   string
}

func GetDashboard(param string) DashboardResponse {
	selectedRange := normalizeRange(param)
	since := time.Now().Add(-durationForRange(selectedRange))
	events := repo.GetEventsSince(since)

	summary := DashboardSummary{}
	uniqueIPs := make(map[string]struct{})
	ipStats := make(map[string]*ipAccumulator)
	endpointStats := make(map[string]*endpointAccumulator)

	for _, event := range events {
		summary.TotalRequests++
		uniqueIPs[event.IP] = struct{}{}

		switch event.Action {
		case utils.ActionBlocked:
			summary.BlockedRequests++
		case utils.ActionRateLimited:
			summary.RateLimitedRequests++
		default:
			summary.AllowedRequests++
		}

		if ipStats[event.IP] == nil {
			ipStats[event.IP] = &ipAccumulator{}
		}
		ipStats[event.IP].count++
		if event.Country != "" {
			ipStats[event.IP].country = event.Country
		}
		updateStatus(&ipStats[event.IP].severity, &ipStats[event.IP].status, event.Action)

		endpoint := event.Endpoint
		if endpoint == "" {
			endpoint = "/"
		}
		if endpointStats[endpoint] == nil {
			endpointStats[endpoint] = &endpointAccumulator{}
		}
		endpointStats[endpoint].count++
		updateStatus(&endpointStats[endpoint].severity, &endpointStats[endpoint].status, event.Action)
	}

	summary.UniqueIPs = len(uniqueIPs)

	topIPs := make([]IPMetric, 0, len(ipStats))
	for ip, stat := range ipStats {
		country := stat.country
		if country == "" {
			country = "--"
		}
		topIPs = append(topIPs, IPMetric{
			IP:       ip,
			Requests: stat.count,
			Country:  country,
			Status:   stat.status,
		})
	}
	sort.Slice(topIPs, func(i, j int) bool {
		if topIPs[i].Requests == topIPs[j].Requests {
			return topIPs[i].IP < topIPs[j].IP
		}
		return topIPs[i].Requests > topIPs[j].Requests
	})
	if len(topIPs) > 5 {
		topIPs = topIPs[:5]
	}

	topEndpoints := make([]EndpointMetric, 0, len(endpointStats))
	for endpoint, stat := range endpointStats {
		topEndpoints = append(topEndpoints, EndpointMetric{
			Endpoint: endpoint,
			Count:    stat.count,
			Status:   stat.status,
		})
	}
	sort.Slice(topEndpoints, func(i, j int) bool {
		if topEndpoints[i].Count == topEndpoints[j].Count {
			return topEndpoints[i].Endpoint < topEndpoints[j].Endpoint
		}
		return topEndpoints[i].Count > topEndpoints[j].Count
	})
	if len(topEndpoints) > 5 {
		topEndpoints = topEndpoints[:5]
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.After(events[j].Timestamp)
	})

	recentEvents := make([]DashboardEvent, 0, min(len(events), 10))
	for _, event := range events {
		recentEvents = append(recentEvents, DashboardEvent{
			Timestamp: event.Timestamp.Format(time.RFC3339),
			IP:        event.IP,
			Endpoint:  event.Endpoint,
			Method:    event.Method,
			Action:    event.Action,
			Country:   event.Country,
			Score:     event.Score,
		})
		if len(recentEvents) == 10 {
			break
		}
	}

	return DashboardResponse{
		Range:        selectedRange,
		Summary:      summary,
		TopIPs:       topIPs,
		TopEndpoints: topEndpoints,
		RecentEvents: recentEvents,
	}
}

func normalizeRange(param string) string {
	normalized := strings.ToUpper(strings.TrimSpace(param))
	switch normalized {
	case "1H", "24H", "7D", "30D":
		return normalized
	default:
		return utils.DefaultDashboardRange
	}
}

func durationForRange(selectedRange string) time.Duration {
	switch selectedRange {
	case "1H":
		return time.Hour
	case "7D":
		return 7 * 24 * time.Hour
	case "30D":
		return 30 * 24 * time.Hour
	default:
		return 24 * time.Hour
	}
}

func updateStatus(currentSeverity *int, currentStatus *string, action string) {
	severity := actionSeverity(action)
	if severity > *currentSeverity {
		*currentSeverity = severity
		*currentStatus = action
	}
}

func actionSeverity(action string) int {
	switch action {
	case utils.ActionBlocked:
		return 3
	case utils.ActionRateLimited:
		return 2
	default:
		return 1
	}
}
