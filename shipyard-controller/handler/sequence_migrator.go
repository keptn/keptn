package handler

import "context"

type SequenceMigrator interface {
	Start(ctx context.Context)
}
