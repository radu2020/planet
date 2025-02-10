package storage

import (
	_ "github.com/lib/pq"
	"github.com/paulmach/orb/geojson"
)

type Storage interface {
	GetCollection() (*geojson.FeatureCollection, error)
	GetOrgIDs() ([]int, error)
}
