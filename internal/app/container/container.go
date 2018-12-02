package container

// Container is generic container descriptor containing only business domain relevant fields
type Container struct {
	ID   string
	Name string
	Addr string
}

func (c Container) String() string {
	return "[" +
		c.ID + ", " +
		c.Name + ", " +
		c.Addr + ", " +
		"]"
}
