package service

import (
	"errors"
	"github.com/paulmach/orb/geojson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// Mock Storage
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetCollection() (*geojson.FeatureCollection, error) {
	args := m.Called()
	return args.Get(0).(*geojson.FeatureCollection), args.Error(1)
}

func (m *MockStorage) GetOrgIDs() ([]int, error) {
	args := m.Called()
	return args.Get(0).([]int), args.Error(1)
}

func TestGetCollection_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	expectedFC := geojson.NewFeatureCollection()
	mockStorage.On("GetCollection").Return(expectedFC, nil)

	service := NewDataService(mockStorage)
	fc, err := service.GetCollection()

	assert.NoError(t, err)
	assert.Equal(t, expectedFC, fc)
	mockStorage.AssertExpectations(t)
}

func TestGetCollection_Error(t *testing.T) {
	mockStorage := new(MockStorage)
	mockStorage.On("GetCollection").Return((*geojson.FeatureCollection)(nil), errors.New("database error"))

	service := NewDataService(mockStorage)
	fc, err := service.GetCollection()

	assert.Error(t, err)
	assert.Nil(t, fc)
	mockStorage.AssertExpectations(t)
}

func TestGetOrgIDs_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	expectedIDs := []int{1, 2, 3}
	mockStorage.On("GetOrgIDs").Return(expectedIDs, nil)

	service := NewDataService(mockStorage)
	result, err := service.GetOrgIDs()

	assert.NoError(t, err)
	assert.Equal(t, &OrgIDList{OrgIDs: expectedIDs}, result)
	mockStorage.AssertExpectations(t)
}

func TestGetOrgIDs_Error(t *testing.T) {
	mockStorage := new(MockStorage)
	mockStorage.On("GetOrgIDs").Return([]int(nil), errors.New("database error"))

	service := NewDataService(mockStorage)
	result, err := service.GetOrgIDs()

	assert.Error(t, err)
	assert.Nil(t, result)
	mockStorage.AssertExpectations(t)
}
