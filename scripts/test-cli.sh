#!/bin/bash

# =================================================================
#
# Work of the U.S. Department of Defense, Defense Digital Service.
# Released as open source under the MIT License.  See LICENSE file.
#
# =================================================================

set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export testdata_local="${DIR}/../testdata"

export temp="${DIR}/../temp"


# examples/texas_hurricanes.csv is copied from https://en.wikipedia.org/wiki/List_of_United_States_hurricanes#Texas
testTexasHurricanes() {
  local expected='{"August":23,"July":9,"June":11,"October":5,"September":17}'
  local output=$("${DIR}/../bin/timebucket" -v '{{.ClosestApproach}} {{.Year}}' --layouts 'January 2 2006' -k 'January' -o 'json' 'examples/texas_hurricanes.csv'  2>&1)
  assertEquals "unexpected output" "${expected}" "${output}"
}

oneTimeSetUp() {
  echo "Using temporary directory at ${SHUNIT_TMPDIR}"
  echo "Reading testdata from ${testdata_local}"
}

oneTimeTearDown() {
  echo "Tearing Down"
}

# Load shUnit2.
. "${DIR}/shunit2"
