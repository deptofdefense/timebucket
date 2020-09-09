// =================================================================
//
// Work of the U.S. Department of Defense, Defense Digital Service.
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	expected := time.Date(2020, time.Month(1), 01, 0, 0, 0, 1000000, time.UTC)
	//
	datetime, err := Parse("01 Jan 2020T00:00:00.001", DefaultLayouts)
	assert.NoError(t, err)
	assert.Equal(t, expected, datetime)
	//
}
