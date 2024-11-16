package httpserver

import (
	"fmt"
	"maps"
	"strings"
)

type (
	patternRoute struct {
		host      string
		pattern   string
		endSlash  bool
		multiName string
		wildcard  map[string]bool
	}
)

func newPatternRoute(pattern string) *patternRoute {
	if pattern == "" {
		pattern = "/"
	}
	i := strings.IndexByte(pattern, '/')
	if i < 0 {
		pattern = pattern + "/"
		i = len(pattern) - 1
	}

	pr := &patternRoute{
		host:     pattern[:i],
		pattern:  pattern[i:],
		endSlash: strings.HasSuffix(pattern[i:], "/"),
		wildcard: map[string]bool{},
	}

	pr.extractWildcard()

	return pr
}

func (pr *patternRoute) extractWildcardSegment(seg string) {
	if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
		name := seg[1 : len(seg)-1]

		name, cut := strings.CutSuffix(name, "...")
		if cut {
			pr.multiName = name
		}

		if name != "" && name != "$" {
			if pr.wildcard[name] {
				panic(fmt.Errorf("httpserver: duplicate wildcard name %q", name))
			}
			pr.wildcard[name] = true
		}
	}
}

func (pr *patternRoute) extractWildcard() {
	rest := pr.pattern
	seg := ""
	for len(rest) > 0 {
		rest = rest[1:]

		i := strings.IndexByte(rest, '/')
		if i < 0 {
			i = len(rest)
		}
		seg, rest = rest[:i], rest[i:]

		pr.extractWildcardSegment(seg)
	}
}

func patternRouteJoinWildcard(prefixPr, suffixPr *patternRoute) map[string]bool {
	wildcard := map[string]bool{}

	maps.Copy(wildcard, prefixPr.wildcard)

	for k, v := range suffixPr.wildcard {
		if wildcard[k] {
			panic(fmt.Errorf("httpserver: duplicate wildcard name %q", k))
		}
		wildcard[k] = v
	}

	return wildcard
}

func patternRouteJoinHost(prefixHost, suffixHost string) string {
	if prefixHost != "" && suffixHost != "" && prefixHost != suffixHost {
		panic(fmt.Errorf("httpserver: host %s conflicts with host %s", prefixHost, suffixHost))
	}

	host := prefixHost
	if suffixHost != "" {
		host = suffixHost
	}

	return host
}

func (pr *patternRoute) removeEndSlash() string {
	pattern := pr.pattern
	if pr.endSlash {
		pattern = pattern[:len(pattern)-1]
	}
	return pattern
}

func (pr *patternRoute) join(pattern string) *patternRoute {
	suffixPr := newPatternRoute(pattern)

	host := patternRouteJoinHost(pr.host, suffixPr.host)
	wildcard := patternRouteJoinWildcard(pr, suffixPr)

	pattern = pr.removeEndSlash() + suffixPr.removeEndSlash()

	endSlash := suffixPr.endSlash
	if pr.pattern != "/" && suffixPr.pattern == "/" {
		endSlash = false
	}

	if endSlash {
		pattern = pattern + "/"
	}

	prResult := &patternRoute{
		host:      host,
		pattern:   pattern,
		endSlash:  endSlash,
		multiName: suffixPr.multiName,
		wildcard:  wildcard,
	}
	return prResult
}

func (pr *patternRoute) String() string {
	return fmt.Sprintf("%v%v", pr.host, pr.pattern)
}

func (pr *patternRoute) mountMethodNotAllowed() string {
	pattern := pr.pattern
	if pr.multiName != "" {
		pattern, _ = strings.CutSuffix(pattern, fmt.Sprintf("{%s...}", pr.multiName))
	}

	return fmt.Sprintf("%v%v", pr.host, pattern)
}
