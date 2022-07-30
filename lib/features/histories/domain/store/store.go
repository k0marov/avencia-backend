package store

import "github.com/k0marov/avencia-backend/lib/core/fs_facade"

type HistoryGetter = func(userId string) (fs_facade.Documents, error)
