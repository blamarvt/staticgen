package vars

// Store holds a map of variables that can be used in page templates
type Store struct {
	variables map[string]string
}

// NewStore creates a new Store with an empty map of variables
func NewStore() *Store {
	return &Store{
		variables: make(map[string]string),
	}
}

// Set sets a variable value in the store
func (s *Store) Set(key, value string) {
	s.variables[key] = value
}

// Get retrieves a variable value from the store
func (s *Store) Get(key string) (string, bool) {
	value, ok := s.variables[key]
	return value, ok
}

// GetOrDefault retrieves a variable value from the store, or returns a default value if not found
func (s *Store) GetOrDefault(key, defaultValue string) string {
	if value, ok := s.variables[key]; ok {
		return value
	}
	return defaultValue
}
