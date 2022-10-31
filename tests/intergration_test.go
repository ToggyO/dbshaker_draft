//go:build integration
// +build integration

package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"

	_ "github.com/ToggyO/dbshaker/tests/migrations"
	"github.com/ToggyO/dbshaker/tests/postgres"
)

func TestIntegration(t *testing.T) {
	suite.Run(t, new(postgres.PgTestSuite))
}
