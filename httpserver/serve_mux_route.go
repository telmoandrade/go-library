package httpserver

import (
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type (
	serveMuxRoute struct {
		allowedMethods []string
		maxAge         int
		cors           *cors
	}
)

func (smr *serveMuxRoute) appendMethod(method string) {
	if !slices.Contains(smr.allowedMethods, method) {
		smr.allowedMethods = append(smr.allowedMethods, method)
	}
}

func (smr *serveMuxRoute) addMethod(method string) {
	smr.appendMethod("OPTIONS")

	if method == "GET" {
		smr.appendMethod("HEAD")
	}
	if method == "" {
		smr.appendMethod("CONNECT")
		smr.appendMethod("DELETE")
		smr.appendMethod("GET")
		smr.appendMethod("HEAD")
		smr.appendMethod("PATCH")
		smr.appendMethod("POST")
		smr.appendMethod("PUT")
		smr.appendMethod("TRACE")
	} else {
		smr.appendMethod(method)
	}

	sort.Strings(smr.allowedMethods)
}

func (smr *serveMuxRoute) headerAllow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", strings.Join(smr.allowedMethods, ", "))
	if smr.maxAge > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%v", smr.maxAge))
	}
}

func (smr *serveMuxRoute) handlerOptions(w http.ResponseWriter, r *http.Request) {
	if smr.cors == nil {
		smr.headerAllow(w, r)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (smr *serveMuxRoute) middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if smr.cors != nil {
			origin := r.Header.Get("Origin")

			if origin == "" || !smr.cors.isOriginAllowed(r, origin) {
				if r.Method == http.MethodOptions {
					smr.headerAllow(w, r)
				}
				next.ServeHTTP(w, r)
				return
			}

			if smr.cors.allowedOriginsAll && !smr.cors.allowCredentials {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			if smr.cors.allowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(smr.allowedMethods, ", "))

				if smr.cors.maxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", strconv.Itoa(smr.cors.maxAge))
				}
				if smr.cors.allowedHeadersAll && !smr.cors.allowCredentials {
					w.Header().Set("Access-Control-Allow-Headers", "*")
				} else {
					requestHeaders := strings.Split(r.Header.Get("Access-Control-Request-Headers"), ",")
					allowedHeaders := []string{}
					for _, v := range requestHeaders {
						header := http.CanonicalHeaderKey(strings.TrimSpace(v))
						if !smr.cors.allowedHeadersAll && !slices.Contains(smr.cors.allowedHeaders, header) {
							continue
						}

						allowedHeaders = append(allowedHeaders, header)
					}

					if len(allowedHeaders) > 0 {
						w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
					}
				}

				method := r.Header.Get("Access-Control-Request-Method")
				if method == "" {
					w.WriteHeader(http.StatusBadRequest)
				}
				if !slices.Contains(smr.allowedMethods, method) {
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
			} else {
				if len(smr.cors.exposedHeaders) > 0 {
					w.Header().Set("Access-Control-Expose-Headers", strings.Join(smr.cors.exposedHeaders, ", "))
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
