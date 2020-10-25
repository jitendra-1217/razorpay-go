package razorpay

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sort"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	// clientVersion is sent to requests as User-Agent.
	clientVersion = "1.0.0-rc"

	// defaultBackendHost will be used as backend's host by default.
	defaultBackendHost = "https://api.razorpay.com"
)

var (
	// APIVersion is prefixed to all request paths.
	APIVersion string = "v1"

	// APIKey is the public part of Razorpay credential.
	APIKey string

	// APISecret is the secret part of Razorpay credential.
	APISecret string

	// APIHost will be used as api's host. It will default to `defaultBackendHost`.
	APIHost string

	// HTTPClient is pre configured http client. This can be set to any value
	// satisfying Doer interface.
	HTTPClient = &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: 1 * time.Second,
			DialContext:           (&net.Dialer{KeepAlive: 60 * time.Second, Timeout: 1 * time.Second}).DialContext,
			MaxIdleConns:          10,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   1 * time.Second,
			MaxIdleConnsPerHost:   10,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	// HTTPClientPrometheusCollector is prometheus collector for HTTPClient.
	HTTPClientPrometheusCollector = NewPrometheusCollector(HTTPClient, "default")

	// DefaultAPIBackend is if set will be used. It is helpful for unit
	// testability of integration code.
	DefaultAPIBackend Backend
)

// Doer interface has the method required to use a type as custom http client.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client is a configured backend to access apis.
type Client struct {
	apiVersion string
	apiKey     string
	apiSecret  string
	apiBackend Backend
}

// Call sets context and invokes' backend's call.
func (c *Client) Call(ctx context.Context, method string, path string, params RequestParams, v ResponseHolder) error {
	if params == nil {
		params = &Params{}
	}

	// Prefixes path with api version.
	path = c.apiVersion + path

	// Sets 'Authorization' header in `params`.
	authorizationValue := "Basic " + base64.StdEncoding.EncodeToString([]byte(c.apiKey+":"+c.apiSecret))
	params.SetHeader("Authorization", authorizationValue)

	return c.apiBackend.Call(ctx, method, path, params, v)
}

// IsValidPaymentSignature returns if payment signature is valid.
// Ref: https://razorpay.com/docs/payment-gateway/quick-integration/#step-4-verify-the-signature.
func (c *Client) IsValidPaymentSignature(_ context.Context, params map[string]string) (bool, error) {
	// Sample value of params:
	// {
	//   "razorpay_signature": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
	//   "razorpay_payment_id": "pay_00000000000001",
	//   "razorpay_order_id": "order_00000000000001"
	// }

	signature, ok := params["razorpay_signature"]
	if !ok {
		return false, fmt.Errorf("razorpay_signature is missing in params")
	}

	payload := ""
	keys := []string{}
	for k := range params {
		if k != "razorpay_signature" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		if payload != "" {
			payload += "|"
		}
		payload += params[k]
	}

	return isPayloadSignatureValid([]byte(payload), signature, c.apiSecret), nil
}

// GetDefaultClient returns client configured with defaults.
func GetDefaultClient() *Client {
	return NewClient(APIKey, APISecret, DefaultAPIBackend)
}

// NewClient returns new client.
func NewClient(apiKey string, apiSecret string, apiBackend Backend) *Client {
	if apiBackend == nil {
		apiBackend = &APIBackend{APIHost, HTTPClient}
	}

	return &Client{APIVersion, apiKey, apiSecret, apiBackend}
}

// Backend provides Call function to make request to remote host.
// It helps mocking for unit tests.
type Backend interface {
	Call(ctx context.Context, method string, path string, params RequestParams, v ResponseHolder) error
}

// APIBackend implements Backend.
type APIBackend struct {
	Host       string
	HTTPClient Doer
}

// Call builds, make requests, and unmarshals resp body into holder.
func (b *APIBackend) Call(_ context.Context, method string, path string, params RequestParams, v ResponseHolder) error {
	host := defaultBackendHost
	if b.Host != "" {
		host = b.Host
	}

	// Builds URL.
	// Appends `params` as URL query params for GET requests.
	url := host + "/" + path
	if isMethodGet(method) {
		queryParams, err := query.Values(params)
		if err != nil {
			return err
		}
		url = url + "?" + queryParams.Encode()
	}

	// Builds json body for non-GET requests.
	var body io.Reader
	if !isMethodGet(method) {
		jsonBody, err := json.Marshal(params)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	// Appends header params to request, if any set from client layer.
	for k, v := range params.Headers() {
		req.Header.Set(k, v)
	}
	// Appends rest of headers...
	if !isMethodGet(method) {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", "jitendra-1217/razorpay-go/"+clientVersion)

	// Makes request and reads resp body.
	resp, err := b.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// If resp is not success then unmarshals body into new error type and
	// returns. This way it is uniform and forces to handle error responses.
	if !isStatusCodeSuccess(resp.StatusCode) {
		v := &struct{ Error *Error }{}
		err := json.Unmarshal(respBody, v)
		if err != nil {
			return err
		}
		return v.Error
	}

	if v == nil {
		v = &Response{}
	}
	v.SetBody(respBody)
	// Unmarshals resp body into holder.
	err = json.Unmarshal(respBody, v)
	if err != nil {
		return err
	}

	return nil
}

func isMethodGet(method string) bool {
	return method == http.MethodGet
}

func isStatusCodeSuccess(statusCode int) bool {
	return statusCode == http.StatusOK
}

// RequestParams is request context i.e. query, body, and headers.
type RequestParams interface {
	SetHeader(key string, value string)
	Headers() map[string]string
}

// ResponseHolder holds response.
type ResponseHolder interface {
	SetBody([]byte)
}

// Entity is common part of entity representation.
type Entity struct {
	ID        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
}

// EntityList is collection of entities.
type EntityList struct {
	Count int64 `json:"count"`
}

// Notes represent https://razorpay.com/docs/api/#notes.
type Notes map[string]string

// UnmarshalJSON unmarshals raw notes into Notes type.
func (n *Notes) UnmarshalJSON(data []byte) error {
	// In json response, notes will appear slice i.e. `[]` when it is empty,
	// and so only attempts unmarshal when it appears object.
	if len(data) > 0 && data[0] == '{' {
		// Also it is possible that notes is not truly map[string]string, for
		// those case it will simply return error, for now.
		alias := map[string]string{}
		err := json.Unmarshal(data, &alias)
		if err != nil {
			return err
		}
		*n = alias
	}
	return nil
}

// ListParams is common list params that can be used when listing entities.
type ListParams struct {
	Params
	From   *int64   `url:"from,omitempty"`
	To     *int64   `url:"to,omitempty"`
	Count  *int64   `url:"count,omitempty"`
	Skip   *int64   `url:"skip,omitempty"`
	Expand []string `url:"expand[],omitempty"`
}

// GetParams is common params that can be used when getting entities.
type GetParams struct {
	Params
	Expand []string `url:"expand[],omitempty"`
}

// Params is common params.
type Params struct {
	headers map[string]string
}

// Headers returns set headers.
func (p *Params) Headers() map[string]string {
	return p.headers
}

// SetHeader sets a new header to be put onto request later.
func (p *Params) SetHeader(key string, value string) {
	if p.headers == nil {
		p.headers = map[string]string{}
	}
	p.headers[key] = value
}

// Response is common part of response.
type Response struct {
	Body []byte
}

// SetBody sets raw response body.
func (r *Response) SetBody(body []byte) {
	r.Body = body
}

// Error represents an error response.
type Error struct {
	Response
	Code        string            `json:"code"`
	Description string            `json:"description"`
	Field       string            `json:"field"`
	Source      string            `json:"source"`
	Step        string            `json:"step"`
	Reason      string            `json:"reason"`
	Metadata    map[string]string `json:"metadata"`
}

// Error returns one-liner error string.
func (e *Error) Error() string {
	return fmt.Sprintf("code: %s, description: %s", e.Code, e.Description)
}

// IsValidPaymentSignature returns if payment signature is valid.
// Ref: https://razorpay.com/docs/payment-gateway/quick-integration/#step-4-verify-the-signature.
func IsValidPaymentSignature(ctx context.Context, params map[string]string) (bool, error) {
	return GetDefaultClient().IsValidPaymentSignature(ctx, params)
}

// IsValidWebhookRequest returns if webhook request is valid.
// Ref: https://razorpay.com/docs/webhooks/#validation
func IsValidWebhookRequest(ctx context.Context, r *http.Request, secret string) (bool, error) {
	// Reads request body and hence closing it, but resets to not have side effect.
	// Ref: https://stackoverflow.com/a/23077519/1325949
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false, err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	signature := r.Header.Get("X-Razorpay-Signature")

	return isPayloadSignatureValid(reqBody, signature, secret), nil
}

func isPayloadSignatureValid(payload []byte, signature string, secret string) bool {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload) //nolint
	expectedSignature := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// Int64 returns pointer to an int64 value.
func Int64(v int64) *int64 {
	return &v
}

// String returns pointer to a string value.
func String(v string) *string {
	return &v
}

// Bool returns pointer to a bool value.
func Bool(v bool) *bool {
	return &v
}
