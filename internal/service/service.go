package service

import (
	"github.com/paulmach/orb/geojson"
	"github.com/radu2020/planet/internal/storage"
)

type DataService struct {
	storage storage.Storage
}

func NewDataService(storage storage.Storage) *DataService {
	return &DataService{storage: storage}
}

// Get all features from the database and return geojson FeatureCollection
func (s DataService) GetCollection() (*geojson.FeatureCollection, error) {
	fc, err := s.storage.GetCollection()
	if err != nil {
		return nil, err
	}
	return fc, nil
}

type OrgIDList struct {
	OrgIDs []int `json:"org_ids"`
}

func (s DataService) GetOrgIDs() (*OrgIDList, error) {
	orgIDs, err := s.storage.GetOrgIDs()
	if err != nil {
		return nil, err
	}
	return &OrgIDList{OrgIDs: orgIDs}, nil
}
