package datastore

// Datastore is responsible for in-memory data
// All data is retained until the service is running. On shutdown, the data will be lost.
// The data that load balancer stores does not need persistence.

// IsHealthyRegistry stores live health data of available services. The key is the unique ID of each host generated at load balancer startup.
type IsHealthyRegistry map[string]bool