package store

type WithdrawnUpdater = func(userId string, currency string, addValue float64)
