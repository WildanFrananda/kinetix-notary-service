package httpapi

import "net/http"

var Routes = []Route{
	{Method: http.MethodGet, Path: "/health/live"},
	{Method: http.MethodGet, Path: "/health/ready"},
	{Method: http.MethodGet, Path: "/metrics"},
	{Method: http.MethodGet, Path: ProofsRoute},
	{Method: http.MethodGet, Path: RegistrationsRoute},
	{Method: http.MethodPost, Path: RegistrationsRoute},
	{Method: http.MethodGet, Path: OrderEvidenceRoute},
	{Method: http.MethodGet, Path: ReviewEvidenceRoute},
}

var WatchRoutes = []Route{
	{Method: http.MethodGet, Path: "/health/live"},
	{Method: http.MethodGet, Path: "/health/ready"},
	{Method: http.MethodGet, Path: "/metrics"},
}
