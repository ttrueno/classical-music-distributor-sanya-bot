package models

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyID = errors.New("id cannot be empty string")
)

type Composer struct {
	ID          string
	FirstName   string
	LastName    string
	ImageLink   string
	Description string
	Version     string
}

func (c *Composer) Describe() string {
	return fmt.Sprintf(
		"%s %s\n",
		c.FirstName,
		c.LastName,
	)
}

func (c *Composer) DescribeVerbose() string {
	return fmt.Sprintf(
		"%s %s -\n%s",
		c.FirstName,
		c.LastName,
		c.Description,
	)
}
