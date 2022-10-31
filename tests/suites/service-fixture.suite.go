package suites

import (
	"database/sql"
	"log"
	"path/filepath"
	"sort"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/ToggyO/dbshaker/pkg"
)

type ServiceFixtureSuite struct {
	suite.Suite

	MigrationRoot            string
	InitialMigrationVersions []int64

	Db   *sql.DB
	Conf DbConf
}

func (sf *ServiceFixtureSuite) Init(dialect, connectionString string) {
	db, err := dbshaker.OpenDbWithDriver(dialect, connectionString)
	if err != nil {
		panic(err)
	}

	sf.Db = db

	dir, err := filepath.Abs("./migrations")
	if err != nil {
		log.Fatalln(err)
	}

	sf.MigrationRoot = dir
	sf.InitialMigrationVersions = []int64{11092022001, 15102022005, 31102022003}

	sf.sortVersions(sf.InitialMigrationVersions)
	sf.setupInitialMigrations()
}

func (sf *ServiceFixtureSuite) TearDownSuite() {
	err := dbshaker.Down(sf.Db, sf.MigrationRoot)
	require.NoError(sf.Suite.T(), err)

	migrations, err := dbshaker.ListMigrations()
	require.NoError(sf.Suite.T(), err)
	require.Empty(sf.T(), migrations)
}

func (sf *ServiceFixtureSuite) setupInitialMigrations() {
	sf.testUpTo()
	sf.testUp()
}

func (sf *ServiceFixtureSuite) testUpTo() {
	target := sf.InitialMigrationVersions[1]
	err := dbshaker.UpTo(sf.Db, sf.MigrationRoot, target)
	require.NoError(sf.Suite.T(), err)

	migrations, err := dbshaker.ListMigrations()
	require.NoError(sf.Suite.T(), err)
	require.Len(sf.T(), migrations, 2)
}

func (sf *ServiceFixtureSuite) testUp() {
	err := dbshaker.Up(sf.Db, sf.MigrationRoot)
	require.NoError(sf.Suite.T(), err)

	migrations, err := dbshaker.ListMigrations()
	require.NoError(sf.Suite.T(), err)

	versions := make([]int64, 0, len(sf.InitialMigrationVersions))
	for _, m := range migrations {
		versions = append(versions, m.Version)
	}

	sf.sortVersions(versions)
	require.Equal(sf.Suite.T(), sf.InitialMigrationVersions, versions)
}

func (sf *ServiceFixtureSuite) sortVersions(versions []int64) {
	sort.Slice(sf.InitialMigrationVersions, func(i, j int) bool {
		return versions[i] < versions[j]
	})
}
