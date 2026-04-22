package testutil

import "flag"

var SkipDataApiIntegrationTest = flag.Bool("skip-dataapi", false, "skips data api tests")
var SkipExtractionTest = flag.Bool("skip-extract", false, "skips extraction tests")
