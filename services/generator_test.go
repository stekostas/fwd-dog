package services

import (
	"github.com/stekostas/fwd-dog/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
	"time"
)

type CacheAdapterMock struct {
	mock.Mock
}

func (m *CacheAdapterMock) Set(key string, value interface{}, ttl time.Duration) (string, error) {
	args := m.Called(key, value, ttl)

	return args.String(0), args.Error(1)
}

func (m *CacheAdapterMock) SetOrFail(key string, value interface{}, ttl time.Duration) (bool, error) {
	args := m.Called(key, value, ttl)

	return args.Bool(0), args.Error(1)
}

func (m *CacheAdapterMock) Get(key string) (string, error) {
	args := m.Called(key)

	return args.String(0), args.Error(1)
}

func (m *CacheAdapterMock) Delete(key string) error {
	args := m.Called(key)

	return args.Error(0)
}

type GeneratorTestSuite struct {
	suite.Suite
	generator    *LinkGenerator
	cacheAdapter *CacheAdapterMock
}

func TestLinkGenerator(t *testing.T) {
	suite.Run(t, new(GeneratorTestSuite))
}

func (s *GeneratorTestSuite) SetupTest() {
	s.cacheAdapter = new(CacheAdapterMock)
	s.generator = NewLinkGenerator(s.cacheAdapter)
}

func (s *GeneratorTestSuite) TestGenerateLinkNoOptions() {
	createLink := &models.CreateLink{
		TargetUrl: "https://fwd.dog",
		Ttl:       time.Second,
	}

	s.cacheAdapter.On("SetOrFail", mock.Anything, mock.Anything, time.Second).Return(true, nil)

	key, err := s.generator.Generate(createLink)

	s.cacheAdapter.AssertExpectations(s.T())
	s.NotEmpty(key)
	s.Nil(err)
}

func (s *GeneratorTestSuite) TestGenerateLinkSingleUse() {
	createLink := &models.CreateLink{
		TargetUrl: "https://fwd.dog",
		Ttl:       time.Second,
		SingleUse: true,
	}

	s.cacheAdapter.On("SetOrFail", mock.Anything, mock.Anything, time.Second).Return(true, nil)

	key, err := s.generator.Generate(createLink)

	s.cacheAdapter.AssertExpectations(s.T())
	s.NotEmpty(key)
	s.Nil(err)
	s.True(strings.HasPrefix(key, "."))
}

func (s *GeneratorTestSuite) TestGenerateLinkWithPassword() {
	createLink := &models.CreateLink{
		TargetUrl:         "https://fwd.dog",
		Ttl:               time.Second,
		PasswordProtected: true,
		Password:          "test",
	}

	s.cacheAdapter.On("SetOrFail", mock.Anything, mock.Anything, time.Second).Return(true, nil)

	key, err := s.generator.Generate(createLink)

	s.cacheAdapter.AssertExpectations(s.T())
	s.NotEmpty(key)
	s.Nil(err)
}

func (s *GeneratorTestSuite) TestGenerateLinkAllOptions() {
	createLink := &models.CreateLink{
		TargetUrl:         "https://fwd.dog",
		Ttl:               time.Second,
		SingleUse:         true,
		PasswordProtected: true,
		Password:          "test",
	}

	s.cacheAdapter.On("SetOrFail", mock.Anything, mock.Anything, time.Second).Return(true, nil)

	key, err := s.generator.Generate(createLink)

	s.cacheAdapter.AssertExpectations(s.T())
	s.NotEmpty(key)
	s.Nil(err)
	s.True(strings.HasPrefix(key, "."))
}

func (s *GeneratorTestSuite) TestGeneratePasswordProtectedLinkWithNoPassword() {
	createLink := &models.CreateLink{
		TargetUrl:         "https://fwd.dog",
		Ttl:               time.Second,
		PasswordProtected: true,
	}

	key, err := s.generator.Generate(createLink)

	s.Empty(key)
	s.NotNil(err)
}

func (s *GeneratorTestSuite) TestGenerateLinkUnsuccessfulTooManyTries() {
	createLink := &models.CreateLink{
		TargetUrl: "https://fwd.dog",
		Ttl:       time.Second,
	}

	s.cacheAdapter.On("SetOrFail", mock.Anything, mock.Anything, time.Second).Return(false, nil)

	key, err := s.generator.Generate(createLink)

	s.cacheAdapter.AssertExpectations(s.T())
	s.Empty(key)
	s.NotNil(err)
}

func (s *GeneratorTestSuite) TestGenerateLinkUnsuccessfulAfterFiveTries() {
	createLink := &models.CreateLink{
		TargetUrl: "https://fwd.dog",
		Ttl:       time.Second,
	}

	s.cacheAdapter.On("SetOrFail", mock.Anything, mock.Anything, time.Second).Times(5).Return(false, nil)
	s.cacheAdapter.On("SetOrFail", mock.Anything, mock.Anything, time.Second).Once().Return(true, nil)

	key, err := s.generator.Generate(createLink)

	s.cacheAdapter.AssertExpectations(s.T())
	s.NotEmpty(key)
	s.Nil(err)
	s.Len(key, 6)
}

func (s *GeneratorTestSuite) TestGenerateLinkEmptyData() {
	createLink := &models.CreateLink{}

	key, err := s.generator.Generate(createLink)

	s.Empty(key)
	s.NotNil(err)
}
