package testutils

import (
	"database/sql"
	"testing"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/postgresqldb"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Создает DBCluster с mock master и репликами для тестов
func NewMockDBCluster(t *testing.T, numReplicas int) (*postgresqldb.DBCluster, sqlmock.Sqlmock, []sqlmock.Sqlmock) {
	// Создаем mock для master
	masterDB, masterMock, err := sqlmock.New()
	assert.NoError(t, err)

	// Создаем список для mock'ов реплик
	var replicaMocks []sqlmock.Sqlmock
	var replicaDBs []*sql.DB

	// Если указано количество реплик, создаем моки для реплик
	for i := 0; i < numReplicas; i++ {
		replicaDB, replicaMock, err := sqlmock.New()
		assert.NoError(t, err)
		replicaDBs = append(replicaDBs, replicaDB)
		replicaMocks = append(replicaMocks, replicaMock)
	}

	// Возвращаем DBCluster с master и репликами
	dbCluster := postgresqldb.NewDBCluster(masterDB, replicaDBs)

	return dbCluster, masterMock, replicaMocks
}
