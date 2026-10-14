package httpapi

var Routes = []string{"/health/live", "/health/ready", "/metrics", ProofsRoute, RegistrationsRoute}

var WatchRoutes = []string{"/health/live", "/health/ready", "/metrics"}
