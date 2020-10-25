package testutil

import (
	"os"
	"regexp"

	razorpay "github.com/jitendra-1217/razorpay-go"
)

// Sets api key, secret and host for tests, reading from env values.
func init() {
	razorpay.APIKey = os.Getenv("TEST_API_KEY")
	razorpay.APISecret = os.Getenv("TEST_API_SECRET")
	razorpay.APIHost = os.Getenv("TEST_API_HOST")
}

// idRegex matches any Razorpay format id.
var idRegex = regexp.MustCompile(`^[a-z]{0,10}_[\w]{14}$`)

// IsAnyID returns if value is any valid id.
func IsAnyID(id string) bool {
	return idRegex.MatchString(id)
}
