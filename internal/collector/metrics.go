package collector

import (
	"fmt"
	"strings"
)

// formatMetric combines a metric name with a tags mapping in a format consistent with that
// expected by Telegraf/InfluxDB.
func formatMetric(metric string, tags map[string]string) string {
	var tagComponents []string
	for key, value := range tags {
		tagComponents = append(tagComponents, fmt.Sprintf("%s=%s", key, value))
	}

	suffix := strings.Join(tagComponents, ",")

	if len(tagComponents) > 0 {
		return fmt.Sprintf("%s,%s", metric, suffix)
	}

	return metric
}
