package schema

type ConditionalStringSliceAppender struct {
	slice []string
}

func (c *ConditionalStringSliceAppender) AppendWithCondition(condition bool, s ...string) {
	if condition {
		c.slice = append(c.slice, s...)
	}
}

func (c *ConditionalStringSliceAppender) Append(s ...string) {
	c.slice = append(c.slice, s...)
}

func (c *ConditionalStringSliceAppender) Slice() []string {
	return c.slice
}
