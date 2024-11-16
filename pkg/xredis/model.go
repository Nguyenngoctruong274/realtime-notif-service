package xredis

type StringSliceCmd struct {
	err error

	val []string
}

func (c *StringSliceCmd) Result() ([]string, error) {
	return c.val, c.err
}
