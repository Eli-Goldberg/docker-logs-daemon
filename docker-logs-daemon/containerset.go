package main

// ContainerSet preserves a state for a list of uniquely identified container IDs
// and has the ability to check for the existence of a specific container, list them, etc.
type ContainerSet struct {
	containerMap map[string]bool
}

// NewContainerSet Creates a new container
func NewContainerSet() *ContainerSet {
	containerMap := make(map[string]bool)

	return &ContainerSet{containerMap: containerMap}
}

// Exists checks if a specfic id exists
func (c *ContainerSet) Exists(id string) bool {
	_, ok := c.containerMap[id]
	return ok;
}

// Add adds a new id to the set
func (c *ContainerSet) Add(id string) {
	c.containerMap[id] = true;
}

// Remove removes a new id to the set
func (c *ContainerSet) Remove(id string) {
	delete(c.containerMap, id);
}

// GetAll returns a lisit of all keys 
func (c ContainerSet) GetAll() []string{
	keys := make([]string, 0, len(c.containerMap))
	for k := range c.containerMap {
		keys = append(keys, k)
	}
	return keys
}