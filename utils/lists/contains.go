package lists

import "strings"

/* Public */

func StringListContains(list []string, key string, caseSensitive bool) bool {
	for _, listKey := range list {
		switch caseSensitive {
		case false:
			if strings.EqualFold(listKey, key) {
				return true
			}
		default:
			if listKey == key {
				return true
			}
		}
	}

	return false
}

func IntListContains(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}
