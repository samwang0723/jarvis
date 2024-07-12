package remotetest

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func TestCreatePostgres(t *testing.T) {
	t.Parallel()

	container, err := CreatePostgres()
	assert.Nil(t, err)
	pool, err := container.CreateConnPool()
	assert.Nil(t, err)

	pool.Close()

	err = container.RunMigrations()
	assert.Nil(t, err)

	assert.Nil(t, container.Purge())
}

func TestSetupPostgresClient(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, SetupPostgresClient(t, true))
}
