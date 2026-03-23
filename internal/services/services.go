package services

import (
	"fmt"
	repo "todo-golang/internal/repositories"
)

func GetDashboard(param string) string {
	metrics := repo.GetMetrics(param)
	topIPs := repo.GetTopIPs(param)

	fmt.Println(metrics, topIPs)
	return ""
}
