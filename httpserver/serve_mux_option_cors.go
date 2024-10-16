package httpserver

import (
	"net/http"
	"slices"
	"strings"
)

type (
	wildcard struct {
		prefix string
		suffix string
		len    int
	}

	cors struct {
		allowedOrigins         []string
		allowedWildcardOrigins []wildcard
		allowOriginFunc        func(r *http.Request, origin string) bool
		allowedHeaders         []string
		exposedHeaders         []string
		maxAge                 int
		allowedOriginsAll      bool
		allowedHeadersAll      bool
		allowCredentials       bool
	}

	// OptionCors is used to apply configurations to a cors when creating it with [WithCors].
	OptionCors func(*cors)
)

// WithCors is an [OptionServeMux] that defines information that will be used in the handlers for CORS processing.
// A variadic set of [OptionCors] used to configure CORS processing.
//
// Default:
//   - The default maximum age is 86400 seconds.
//
// Default Behavior:
//   - Wildcard handling (e.g. http://*.bar.com).
//   - If the not defined [WithAllowedOrigins] and [WithAllowOriginFunc], any origin is allowed.
//   - If the not defined [WithAllowedHeaders], any header is allowed.
//   - If the defined [WithAllowOriginFunc], will have the responsibility to allow origin, ignoring the [WithAllowedOrigins].
func WithCors(opts ...OptionCors) OptionServeMux {
	return func(mux *ServeMux) {
		c := &cors{
			maxAge:                 86400,
			allowedOrigins:         []string{},
			allowedWildcardOrigins: []wildcard{},
			allowedHeaders:         []string{},
			exposedHeaders:         []string{},
		}

		for _, opt := range opts {
			opt(c)
		}

		if len(c.allowedOrigins) == 0 && c.allowOriginFunc == nil {
			c.allowedOriginsAll = true
		} else if c.allowOriginFunc != nil {
			c.allowedOriginsAll = false
			c.allowedOrigins = []string{}
		} else {
			allowedOrigins := []string{}

			for _, origin := range c.allowedOrigins {
				if origin == "*" {
					c.allowedOriginsAll = true
					c.allowedOrigins = []string{}
					break
				} else if i := strings.IndexByte(origin, '*'); i >= 0 {
					w := wildcard{
						prefix: origin[0:i],
						suffix: origin[i+1:],
						len:    len(origin) - 1,
					}
					c.allowedWildcardOrigins = append(c.allowedWildcardOrigins, w)
				} else {
					allowedOrigins = append(allowedOrigins, origin)
				}
			}
			c.allowedOrigins = allowedOrigins
		}

		if len(c.allowedHeaders) == 0 {
			c.allowedHeadersAll = true
		} else {
			for _, h := range c.allowedHeaders {
				if h == "*" {
					c.allowedHeadersAll = true
					c.allowedHeaders = []string{}
					break
				}
			}
		}

		mux.options.cors = c
	}
}

// WithAllowedOrigins is an [OptionCors] that defines the list of allowed origins.
func WithAllowedOrigins(allowedOrigins ...string) OptionCors {
	return func(c *cors) {
		for _, origin := range allowedOrigins {
			origin = strings.ToLower(origin)
			if !slices.Contains(c.allowedOrigins, origin) {
				c.allowedOrigins = append(c.allowedOrigins, origin)
			}
		}
	}
}

// WithAllowOriginFunc is an [OptionCors] that defines a function to allow origin.
func WithAllowOriginFunc(fn func(r *http.Request, origin string) bool) OptionCors {
	return func(c *cors) {
		if fn != nil {
			c.allowOriginFunc = fn
		}
	}
}

// WithAllowedHeaders is an [OptionCors] that defines the list of allowed headers.
func WithAllowedHeaders(allowedHeaders ...string) OptionCors {
	return func(c *cors) {
		for _, header := range allowedHeaders {
			header = http.CanonicalHeaderKey(header)
			if !slices.Contains(c.allowedHeaders, header) {
				c.allowedHeaders = append(c.allowedHeaders, header)
			}
		}
	}
}

// WithExposedHeaders is an [OptionCors] that defines the list of exposed headers.
func WithExposedHeaders(exposedHeaders ...string) OptionCors {
	return func(c *cors) {
		for _, header := range exposedHeaders {
			header = http.CanonicalHeaderKey(header)
			if !slices.Contains(c.exposedHeaders, header) {
				c.exposedHeaders = append(c.exposedHeaders, header)
			}
		}
	}
}

// WithAllowCredentials is an [OptionCors] that allows credentials.
func WithAllowCredentials(allowCredentials bool) OptionCors {
	return func(c *cors) {
		c.allowCredentials = allowCredentials
	}
}

// WithCorsMaxAge is an [OptionCors] that defines a maximum age in seconds.
func WithCorsMaxAge(seconds int) OptionCors {
	return func(c *cors) {
		c.maxAge = seconds
	}
}

func (c *cors) isOriginAllowed(r *http.Request, origin string) bool {
	if c.allowOriginFunc != nil {
		return c.allowOriginFunc(r, origin)
	}
	if c.allowedOriginsAll {
		return true
	}
	origin = strings.ToLower(origin)
	for _, o := range c.allowedOrigins {
		if o == origin {
			return true
		}
	}
	for _, w := range c.allowedWildcardOrigins {
		if w.match(origin) {
			return true
		}
	}
	return false
}

func (w wildcard) match(s string) bool {
	return len(s) > w.len && strings.HasPrefix(s, w.prefix) && strings.HasSuffix(s, w.suffix)
}
