package parser

import (
	"fmt"
	"sort"
	"strings"
)

func graphiteTags(path string) (string, map[string]string, error) {
	labels := make(map[string]string)

	if strings.IndexByte(path, ';') < 0 {
		return path, labels, nil
	}

	arr := strings.Split(path, ";")

	if len(arr[0]) == 0 {
		return "", labels, fmt.Errorf("cannot parse path %#v, no metric found", path)
	}

	metricPath, arr := arr[0], arr[1:]

	// check tags
	for _, label := range arr {
		if strings.Index(label, "=") < 1 {
			return "", labels, fmt.Errorf("cannot parse path %#v, invalid segment %#v", metricPath, label)
		}
	}

	sort.Stable(byKey(arr[1:]))

	// uniq
	toDel := 0
	prevKey := ""
	for i, label := range arr {
		key := label[:strings.Index(label, "=")]
		if key == prevKey {
			toDel++
		} else {
			prevKey = key
		}
		if toDel > 0 {
			arr[i-toDel] = label
		}
	}

	arr = arr[:len(arr)-toDel]

	for _, label := range arr {
		p := strings.Index(label, "=")
		labels[label[:p]] = label[p+1:]
	}

	return metricPath, labels, nil
}
