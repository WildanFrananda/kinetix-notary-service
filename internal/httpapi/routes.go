package httpapi

var Routes = []string{
	"/health/live",
	"/health/ready",
	"/metrics",
	ProofsRoute,
	RegistrationsRoute,
	OrderEvidenceRoute,
	ReviewEvidenceRoute,
}

var WatchRoutes = []string{"/health/live", "/health/ready", "/metrics"}
