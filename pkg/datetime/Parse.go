// =================================================================
//
// Work of the U.S. Department of Defense, Defense Digital Service.
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package datetime

import (
	"fmt"
	"strings"
	"time"
)

// Parse parses a string using multiple time layouts.  The first layout that is successfully parsed is used.  If no layout is successfully parsed, then returns an error.
func Parse(value string, layouts []string) (time.Time, error) {
	for _, layout := range layouts {
		t, err := time.Parse(layout, value)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("could not parse time.Time from string %q using layouts %q", value, strings.Join(layouts, ","))
}
