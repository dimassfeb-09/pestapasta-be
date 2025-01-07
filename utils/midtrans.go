package utils

import "fmt"

var url = GetENV().MidtransUrl

var (
	EndpointMidtransCharge = fmt.Sprintf("%s/charge", url)
)
