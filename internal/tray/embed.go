// oreon/defense · watchthelight <wtl>
// Icons by serval_serval_serval_serval

package assets

import _ "embed"

//go:embed icons/secure-16.png
var ProtectedIcon []byte

//go:embed icons/warning-16.png
var WarningIcon []byte

//go:embed icons/critical-16.png
var AlertIcon []byte

//go:embed icons/pending-16.png
var ScanningIcon []byte

//go:embed icons/paused-16.png
var PausedIcon []byte
