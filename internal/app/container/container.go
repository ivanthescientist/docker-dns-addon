package container

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
