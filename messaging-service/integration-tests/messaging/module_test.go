package integrationtests

import (
	"testing"

	redisClient "messaging-service/src/redis"

	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	redis *redisClient.RedisClient
}

func (s *IntegrationTestSuite) SetupSuite() {

	rClient := redisClient.New()
	s.redis = rClient
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
