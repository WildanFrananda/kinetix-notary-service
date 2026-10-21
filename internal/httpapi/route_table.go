package httpapi

import (
	"net/http"
	"slices"
	"strings"
)

type routeTable struct {
	mux        *http.ServeMux
	instrument Instrument
	methods    map[string][]string
}

func newRouteTable(instrument Instrument) *routeTable {
	return &routeTable{mux: http.NewServeMux(), instrument: instrument, methods: map[string][]string{}}
}

func (t *routeTable) handle(method, path string, handler http.Handler) {
	t.mux.Handle(method+" "+path, t.instrument(path, handler))
	t.methods[path] = append(t.methods[path], method)
}

func (t *routeTable) handler() http.Handler {
	for path, methods := range t.methods {
		t.mux.Handle(path, methodNotAllowed(allowed(methods)))
	}

	unmatched := t.instrument("unmatched", http.HandlerFunc(notFound))
	fallback := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			unmatched.ServeHTTP(w, r)

			return
		}

		notFound(w, r)
	}

	t.mux.HandleFunc("/", fallback)

	return t.mux
}

func allowed(methods []string) string {
	allow := slices.Clone(methods)
	if slices.Contains(allow, http.MethodGet) {
		allow = append(allow, http.MethodHead)
	}

	slices.Sort(allow)

	return strings.Join(allow, ", ")
}

func methodNotAllowed(allow string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Allow", allow)
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "this path answers "+allow)
	})
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	writeError(w, http.StatusNotFound, "not_found", "nothing is served at this path")
}
