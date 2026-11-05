package jwks

import (
	"bytes"
	"context"
	"crypto"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

var (

	// ErrKIDNotFound indicates that the given key ID was not found in the JWKs.
	ErrKIDNotFound = errors.New("the given key ID was not found in the JWKs")

	// ErrMissingAssets indicates there are required assets missing to create a public key.
	ErrMissingAssets = errors.New("required assets are missing to create a public key")

	// ErrUnknownKeyType indicated that a key type is not implemented
	ErrUnknownKeyType = errors.New("the key has an unknown type")
)

// ErrorHandler is a function signature that consumes an error.
type ErrorHandler func(err error)

type JWK interface {
	ID() string
	Alg() string
	Type() string
	PublicKey() (crypto.PublicKey, error)
}

type JWKS interface {
	Key(kid string) (jwk JWK, err error)
	Cancel()
}

// JWKs represents a JSON Web Key Set.
type jwksImpl struct {
	keys                map[string]*jwkImpl
	cancel              context.CancelFunc
	client              *http.Client
	ctx                 context.Context
	jwksURL             string
	mux                 sync.RWMutex
	refreshErrorHandler ErrorHandler
	refreshInterval     time.Duration
	refreshRateLimit    time.Duration
	refreshRequests     chan context.CancelFunc
	refreshTimeout      time.Duration
	refreshUnknownKID   bool
}

// rawJWK represents a raw key inside a JWKs.
type jwkImpl struct {
	key         rawJWK
	precomputed interface{}
}

func (j *jwkImpl) ID() string {
	return j.key.Kid
}

func (j *jwkImpl) Alg() string {
	return j.key.Algorithm
}

func (j *jwkImpl) Type() string {
	return j.key.KeyType
}

func (j *jwkImpl) PublicKey() (crypto.PublicKey, error) {
	if j.key.KeyType == "RSA" {
		return j.rsa()
	} else if j.key.KeyType == "ecdsa" {
		return j.ecdsa()
	}

	return nil, ErrUnknownKeyType
}

type rawJWK struct {
	Algorithm string   `json:"alg"`
	KeyType   string   `json:"kty"`
	Use       string   `json:"use"`
	Curve     string   `json:"crv"`
	Exponent  string   `json:"e"`
	Kid       string   `json:"kid"`
	Modulus   string   `json:"n"`
	X         string   `json:"x"`
	Y         string   `json:"y"`
	X5C       []string `json:"x5c"`
}

// rawJWKs represents a JWKs in JSON format.
type rawJWKs struct {
	Keys []rawJWK `json:"keys"`
}

// New creates a new JWKs from a raw JSON message.
func NewFromJSON(jwksBytes []byte) (JWKS, error) {
	// Iterate through the keys in the raw JWKs. Add them to the JWKs.
	jwks := &jwksImpl{
		keys: map[string]*jwkImpl{},
	}

	if err := jwks.update(jwksBytes); err != nil {
		return nil, err
	}

	return jwks, nil
}

// NewFromURL loads the JWKs at the given URL.
func NewFromURL(jwksURL string, config Config) (JWKS, error) {
	// Apply the options to the JWKs.
	applyConfigDefaults(&config)

	// Create the JWKs.
	jwks := &jwksImpl{
		keys:                map[string]*jwkImpl{},
		jwksURL:             jwksURL,
		client:              config.Client,
		refreshTimeout:      config.RefreshTimeout,
		refreshErrorHandler: config.RefreshErrorHandler,
		refreshRateLimit:    config.RefreshRateLimit,
		refreshInterval:     config.RefreshInterval,
		refreshUnknownKID:   config.RefreshUnknownKID,
	}

	// Get the keys for the JWKs.
	if err := jwks.refresh(); err != nil {
		return nil, err
	}

	// Check to see if a background refresh of the JWKs should happen.
	if jwks.refreshInterval > 0 || jwks.refreshRateLimit > 0 || jwks.refreshUnknownKID {
		// Attach a context used to end the background goroutine.
		jwks.ctx, jwks.cancel = context.WithCancel(context.Background())

		// Create a channel that will accept requests to refresh the JWKs.
		jwks.refreshRequests = make(chan context.CancelFunc, 1)

		// Start the background goroutine for data refresh.
		go jwks.backgroundRefresh()
	}

	return jwks, nil
}

// Cancel ends the background goroutine to update the JWKs. It can only happen once and is only effective if the
// JWKs has a background goroutine refreshing the JWKs keys.
func (j *jwksImpl) Cancel() {
	if j.cancel != nil {
		j.cancel()
	}
}

// Key gets the JWK from the given KID from the JWKs. It may refresh the JWKs if configured to.
func (j *jwksImpl) Key(kid string) (JWK, error) {
	// Get the JSONKey from the JWKs.
	var key *jwkImpl
	var ok bool
	j.mux.RLock()
	if j.keys != nil {
		key, ok = j.keys[kid]
	}
	j.mux.RUnlock()

	// Check if the key was present.
	if !ok {
		// Check to see if configured to refresh on unknown kid.
		if j.refreshUnknownKID {
			// Create a context for refreshing the JWKs.
			ctx, cancel := context.WithCancel(j.ctx)

			// Refresh the JWKs.
			select {
			case <-j.ctx.Done():
				return key, nil
			case j.refreshRequests <- cancel:
			default:
				// If the j.refreshRequests channel is full, return the error early.
				return nil, ErrKIDNotFound
			}

			// Wait for the JWKs refresh to done.
			<-ctx.Done()

			// Lock the JWKs for async safe use.
			j.mux.RLock()
			defer j.mux.RUnlock()

			// Check if the JWKs refresh contained the requested key.
			if key, ok = j.keys[kid]; ok {
				return key, nil
			}
		}

		return nil, ErrKIDNotFound
	}

	return key, nil
}

// backgroundRefresh is meant to be a separate goroutine that will update the keys in a JWKs over a given interval of
// time.
func (j *jwksImpl) backgroundRefresh() {
	// Create some rate limiting assets.
	var lastRefresh time.Time
	var queueOnce sync.Once
	var refreshMux sync.Mutex

	if j.refreshRateLimit > 0 {
		lastRefresh = time.Now().Add(-j.refreshRateLimit)
	}

	// Create a channel that will never send anything unless there is a refresh interval.
	refreshInterval := make(<-chan time.Time)

	// Enter an infinite loop that ends when the background ends.
	for {
		// If there is a refresh interval, create the channel for it.
		if j.refreshInterval > 0 {
			refreshInterval = time.After(j.refreshInterval)
		}

		// Wait for a refresh to occur or the background to end.
		select {
		// Send a refresh request the JWKs after the given interval.
		case <-refreshInterval:
			select {
			case <-j.ctx.Done():
				return
			case j.refreshRequests <- func() {}:
			default: // If the j.refreshRequests channel is full, don't don't send another request.
			}

		// Accept refresh requests.
		case cancel := <-j.refreshRequests:
			// Rate limit, if needed.
			refreshMux.Lock()
			if j.refreshRateLimit > 0 && lastRefresh.Add(j.refreshRateLimit).After(time.Now()) {
				// Don't make the JWT parsing goroutine wait for the JWKs to refresh.
				cancel()

				// Only queue a refresh once.
				queueOnce.Do(func() {
					// Launch a goroutine that will get a reservation for a JWKs refresh or fail to and immediately return.
					go func() {
						// Wait for the next time to refresh.
						refreshMux.Lock()
						wait := time.Until(lastRefresh.Add(j.refreshRateLimit))
						refreshMux.Unlock()
						select {
						case <-j.ctx.Done():
							return
						case <-time.After(wait):
						}

						// Refresh the JWKs.
						refreshMux.Lock()
						defer refreshMux.Unlock()
						if err := j.refresh(); err != nil && j.refreshErrorHandler != nil {
							j.refreshErrorHandler(err)
						}

						// Reset the last time for the refresh to now.
						lastRefresh = time.Now()

						// Allow another queue.
						queueOnce = sync.Once{}
					}()
				})
			} else {
				// Refresh the JWKs.
				if err := j.refresh(); err != nil && j.refreshErrorHandler != nil {
					j.refreshErrorHandler(err)
				}

				// Reset the last time for the refresh to now.
				lastRefresh = time.Now()

				// Allow the JWT parsing goroutine to continue with the refreshed JWKs.
				cancel()
			}
			refreshMux.Unlock()

		// Clean up this goroutine when its context expires.
		case <-j.ctx.Done():
			return
		}
	}
}

// refresh does an HTTP GET on the JWKs URL to rebuild the JWKs.
func (j *jwksImpl) refresh() (err error) {
	// Create a context for the request.
	var ctx context.Context
	var cancel context.CancelFunc
	if j.ctx != nil {
		ctx, cancel = context.WithTimeout(j.ctx, j.refreshTimeout)
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), j.refreshTimeout)
	}
	defer cancel()

	// Create the HTTP request.
	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, j.jwksURL, bytes.NewReader(nil)); err != nil {
		return err
	}

	// Get the JWKs as JSON from the given URL.
	var resp *http.Response
	if resp, err = j.client.Do(req); err != nil {
		return err
	}
	defer resp.Body.Close() // Ignore any error.

	// Read the raw JWKs from the body of the response.
	var jwksBytes []byte
	if jwksBytes, err = io.ReadAll(resp.Body); err != nil {
		return err
	}

	if err = j.update(jwksBytes); err != nil {
		return err
	}

	return nil
}

func (j *jwksImpl) update(jwksBytes []byte) error {
	// Turn the raw JWKs into the correct Go type.
	var rawKS rawJWKs
	if err := json.Unmarshal(jwksBytes, &rawKS); err != nil {
		return err
	}

	keys := map[string]*jwkImpl{}

	for _, k := range rawKS.Keys {
		if k.Use != "sig" {
			continue
		}

		key := &jwkImpl{
			key: k,
		}

		keys[k.Kid] = key
	}

	// Lock the JWKs for async safe usage.
	j.mux.Lock()
	defer j.mux.Unlock()

	j.keys = keys

	return nil
}
