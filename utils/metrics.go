package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetMetricName(spider int, section string) string {
	if spider <= 0 {
		return fmt.Sprintf("spider_%s", section)
	}
	return fmt.Sprintf("spider_%s_%d", section, spider)
}

func GetMetricFileName(spider int, section string) string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("spider-%s.%d", section, spider))
}
