// =================================================================
//
// Work of the U.S. Department of Defense, Defense Digital Service.
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package datetime

import (
	"time"
)

var DefaultLayouts = []string{
	// timestamps
	"1/2/06 15:04:05 PM (MST)",
	"1/2/06 15:04:05 PM",
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05",
	"02 Jan 2006T15:04:05.999999999",
	// date only
	"2006",
	"2006-01-02",
	// time only
	"15:04:05.999999999",
	"15:04:05",
}
