package janitor

import (
	"context"
)

type Janitor interface {
	Run(context.Context) error
}
