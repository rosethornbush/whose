package output

import "time"

func formatDate(value string) string {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return value
	}

	return t.UTC().Format("2006-01-02 15:04:05 UTC")
}
