package tools

import "time"

type CustomDate struct {
	time.Time
	Initial string
}

func (c *CustomDate) UnmarshalJSON(b []byte) error {
	t, err := time.Parse(`"02.01.2006"`, string(b))
	if err != nil {
		return err
	}

	c.Time = t
	c.Initial = string(b)

	return nil
}
