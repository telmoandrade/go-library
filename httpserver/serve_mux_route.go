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
		allowedMethods       []string
		handlerOptionsMaxAge int
		cors                 *cors
	}
)

func (smr *serveMuxRoute) appendMethod(method string) {
	if !slices.Contains(smr.allowedMethods, method) {
		smr.allowedMethods = append(smr.allowedMethods, method)
	}
}

func (smr *serveMuxRoute) addMethod(method string) {
	if method == http.MethodGet {
		smr.appendMethod(http.MethodHead)
	}
	if method != "" {
		smr.appendMethod(method)
	}

	sort.Strings(smr.allowedMethods)
}

func (smr *serveMuxRoute) headerOptions(w http.ResponseWriter) {
	if len(smr.allowedMethods) > 0 {
		w.Header().Set("Allow", strings.Join(smr.allowedMethods, ", "))
		if smr.handlerOptionsMaxAge > 0 {
			w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%v", smr.handlerOptionsMaxAge))
		}
	}
}

func (smr *serveMuxRoute) middlewareMethodNotAllowed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		smr.headerOptions(w)
		next.ServeHTTP(w, r)
	})
}

func (smr *serveMuxRoute) headerCorsAll(w http.ResponseWriter, origin string) {
	if smr.cors.allowedOriginsAll && !smr.cors.allowCredentials {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	if smr.cors.allowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

func (smr *serveMuxRoute) headerCorsPreflightAllowHeaders(w http.ResponseWriter, r *http.Request) {
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

		w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
	}
}

func (smr *serveMuxRoute) headerCorsPreflight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(smr.allowedMethods, ", "))

	if smr.cors.maxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(smr.cors.maxAge))
	}
	smr.headerCorsPreflightAllowHeaders(w, r)

	method := r.Header.Get("Access-Control-Request-Method")
	if method == "" {
		w.Header().Set("Error", "Access-Control-Request-Method not found")
		w.WriteHeader(http.StatusBadRequest)
	} else if !slices.Contains(smr.allowedMethods, method) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (smr *serveMuxRoute) headerCorsActualRequest(w http.ResponseWriter) {
	if len(smr.cors.exposedHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(smr.cors.exposedHeaders, ", "))
	}
}

func (smr *serveMuxRoute) middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer next.ServeHTTP(w, r)

		origin := r.Header.Get("Origin")

		w.Header().Add("Vary", "Origin")
		if r.Method == http.MethodOptions {
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Add("Vary", "Access-Control-Request-Headers")
		}

		if origin == "" || !smr.cors.isOriginAllowed(r, origin) {
			if r.Method == http.MethodOptions || !slices.Contains(smr.allowedMethods, r.Method) {
				smr.headerOptions(w)
			}
			return
		}

		smr.headerCorsAll(w, origin)

		if r.Method == http.MethodOptions {
			smr.headerCorsPreflight(w, r)
		} else {
			smr.headerCorsActualRequest(w)
		}
	})
}
