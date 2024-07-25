package ericlagergren_test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxtest"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	pgxdecimal "github.com/samwang0723/jarvis/internal/db/pginit/ext/decimal/ericlagergren"
)

var (
	databaseURL       string                                     //nolint:gochecknoglobals // test conn
	postgresImageName = "public.ecr.aws/docker/library/postgres" //nolint:gochecknoglobals // postgres dockertest info
	postgresImageTag  = "14.8"                                   //nolint:gochecknoglobals // postgres dockertest info
)

func TestMain(m *testing.M) {
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	err = CheckAndDownloadImg(postgresImageName, postgresImageTag)
	if err != nil {
		log.Fatalf("Could not run CheckAndDownloadImg: %v", err)
	}

	resource, err := dockerPool.Run(postgresImageName, postgresImageTag, []string{
		"POSTGRES_PASSWORD=postgres",
		"POSTGRES_USER=postgres",
		"POSTGRES_DB=datawarehouse",
		"listen_addresses = '*'",
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	databaseURL = fmt.Sprintf(
		"postgres://postgres:%s@%s/datawarehouse?sslmode=disable",
		"postgres",
		getHostPort(resource, "5432/tcp"),
	)

	if err := dockerPool.Retry(func() error {
		db, err := pgx.Connect(context.Background(), databaseURL)
		if err != nil {
			return fmt.Errorf("pgx connect: %w", err)
		}
		if err := db.Ping(context.Background()); err != nil {
			return fmt.Errorf("ping: %w", err)
		}
		db.Close(context.Background())

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker(%s): %s", databaseURL, err)
	}

	leak := flag.Bool("leak", false, "use leak detector")
	flag.Parse()

	code := m.Run()

	if *leak {
		if code == 0 {
			if err := goleak.Find(); err != nil {
				log.Fatalf("goleak: Errors on successful test run: %v\n", err)

				code = 1
			}
		}
	}

	// You can't defer this because os.Exit doesn't care for defer
	if err := dockerPool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func getHostPort(resource *dockertest.Resource, id string) string {
	hostAndPort := resource.GetHostPort(id)
	hp := strings.Split(hostAndPort, ":")
	testHost := hp[0]
	testPort := hp[1]

	return testHost + ":" + testPort
}

// CheckAndDownloadImg is a helper function to check and download image before run the concurrent test
// It prevent the api ratelimit due to concurrent calling pull image api in different test cases. This
// function should be called in TestMain function if you want to start docker instances in different
// testcases concurrently.
func CheckAndDownloadImg(imageName, imageTag string) error {
	retry := 3
	sleepTime := 100 * time.Millisecond //nolint:gomnd // default sleep time

	pool, err := dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("dockertest.NewPool: %w", err)
	}

	_, err = pool.Client.InspectImage(fmt.Sprintf("%s:%s", imageName, imageTag))
	if err == nil {
		return nil
	}

	for retry > 0 {
		err = pool.Client.PullImage(docker.PullImageOptions{
			Repository: imageName,
			Tag:        imageTag,
		}, docker.AuthConfiguration{})
		if err == nil {
			return nil
		}

		retry--

		time.Sleep(sleepTime)
	}

	if err != nil {
		return fmt.Errorf("pool.Client.PullImage: %w", err)
	}

	return nil
}

func TestCodecDecodeValue(t *testing.T) {
	t.Parallel()

	connTestRunner := pgxtest.ConnTestRunner{
		CreateConfig: func(ctx context.Context, tb testing.TB) *pgx.ConnConfig { //nolint:thelper // ref shopspring
			config, err := pgx.ParseConfig(databaseURL)
			if err != nil {
				tb.Fatalf("ParseConfig failed: %v", err)
			}

			return config
		},
		AfterConnect: func(ctx context.Context, tb testing.TB, conn *pgx.Conn) { //nolint:thelper // ref shopspring
			pgxdecimal.Register(conn.TypeMap())
		},
		AfterTest: func(ctx context.Context, tb testing.TB, conn *pgx.Conn) {}, //nolint:thelper // ref shopspring
		CloseConn: func(ctx context.Context, tb testing.TB, conn *pgx.Conn) { //nolint:thelper // ref shopspring
			err := conn.Close(ctx)
			if err != nil {
				tb.Errorf("Close failed: %v", err)
			}
		},
	}

	connTestRunner.RunTest(
		context.Background(),
		t,
		func(ctx context.Context, tb testing.TB, conn *pgx.Conn) { //nolint:thelper // ref shopspring
			original, _ := new(decimal.Big).SetString("1.234")

			rows, err := conn.Query(context.Background(), `select $1::numeric`, original)
			require.NoError(t, err)

			for rows.Next() {
				values, errV := rows.Values()
				require.NoError(t, errV)

				require.Len(t, values, 1)
				v0, ok := values[0].(decimal.Big)
				require.True(t, ok)
				require.Equal(t, 0, original.Cmp(&v0))
			}

			require.NoError(t, rows.Err())

			rows, err = conn.Query(context.Background(), `select $1::numeric`, nil)
			require.NoError(t, err)

			for rows.Next() {
				values, err := rows.Values()
				require.NoError(t, err)

				require.Len(t, values, 1)
				require.Equal(t, nil, values[0])
			}

			require.NoError(t, rows.Err())
		},
	)
}

func TestNaN(t *testing.T) {
	t.Parallel()

	connTestRunner := pgxtest.ConnTestRunner{
		CreateConfig: func(ctx context.Context, tb testing.TB) *pgx.ConnConfig { //nolint:thelper // ref shopspring
			config, err := pgx.ParseConfig(databaseURL)
			if err != nil {
				tb.Fatalf("ParseConfig failed: %v", err)
			}

			return config
		},
		AfterConnect: func(ctx context.Context, tb testing.TB, conn *pgx.Conn) { //nolint:thelper // ref shopspring
			pgxdecimal.Register(conn.TypeMap())
		},
		AfterTest: func(ctx context.Context, tb testing.TB, conn *pgx.Conn) {}, //nolint:thelper // ref shopspring
		CloseConn: func(ctx context.Context, tb testing.TB, conn *pgx.Conn) { //nolint:thelper // ref shopspring
			err := conn.Close(ctx)
			if err != nil {
				tb.Errorf("Close failed: %v", err)
			}
		},
	}

	connTestRunner.RunTest(
		context.Background(),
		t,
		func(ctx context.Context, tb testing.TB, conn *pgx.Conn) { //nolint:thelper // ref shopspring
			var d decimal.Big
			err := conn.QueryRow(context.Background(), `select 'NaN'::numeric`).Scan(&d)
			require.EqualError(t, err, `can't scan into dest[0]: cannot scan NaN into *decimal.Big`)
		},
	)
}

func TestArray(t *testing.T) {
	t.Parallel()

	connTestRunner := pgxtest.ConnTestRunner{
		CreateConfig: func(ctx context.Context, tb testing.TB) *pgx.ConnConfig { //nolint:thelper // ref shopspring
			config, err := pgx.ParseConfig(databaseURL)
			if err != nil {
				tb.Fatalf("ParseConfig failed: %v", err)
			}

			return config
		},
		AfterConnect: func(ctx context.Context, tb testing.TB, conn *pgx.Conn) { //nolint:thelper // ref shopspring
			pgxdecimal.Register(conn.TypeMap())
		},
		AfterTest: func(ctx context.Context, tb testing.TB, conn *pgx.Conn) {}, //nolint:thelper // ref shopspring
		CloseConn: func(ctx context.Context, tb testing.TB, conn *pgx.Conn) { //nolint:thelper // ref shopspring
			err := conn.Close(ctx)
			if err != nil {
				tb.Errorf("Close failed: %v", err)
			}
		},
	}

	connTestRunner.RunTest(
		context.Background(),
		t,
		func(ctx context.Context, tb testing.TB, conn *pgx.Conn) { //nolint:thelper // ref shopspring
			inputSlice := []*decimal.Big{}

			for i := 0; i < 10; i++ {
				d := decimal.New(int64(i), 0)
				inputSlice = append(inputSlice, d)
			}

			var outputSlice []decimal.Big
			err := conn.QueryRow(context.Background(), `select $1::numeric[]`, inputSlice).
				Scan(&outputSlice)
			require.NoError(t, err)

			require.Equal(t, len(inputSlice), len(outputSlice))
			for i := 0; i < len(inputSlice); i++ {
				require.Equal(t, 0, outputSlice[i].Cmp(inputSlice[i]))
			}
		},
	)
}

func isExpectedEqDecimal(a decimal.Big) func(interface{}) bool {
	return func(v interface{}) bool {
		decBig := v.(decimal.Big) //nolint:forcetypeassert // skip check

		return a.Cmp(&decBig) == 0
	}
}

func newFromString(val string) decimal.Big {
	decBig, _ := new(decimal.Big).SetString(val)

	return *decBig
}

func TestValueRoundTrip(t *testing.T) {
	t.Parallel()

	connTestRunner := pgxtest.ConnTestRunner{
		CreateConfig: func(ctx context.Context, tb testing.TB) *pgx.ConnConfig { //nolint:thelper // ref shopspring
			config, err := pgx.ParseConfig(databaseURL)
			if err != nil {
				tb.Fatalf("ParseConfig failed: %v", err)
			}

			return config
		},
		AfterConnect: func(ctx context.Context, tb testing.TB, conn *pgx.Conn) { //nolint:thelper // ref shopspring
			pgxdecimal.Register(conn.TypeMap())
		},
		AfterTest: func(ctx context.Context, tb testing.TB, conn *pgx.Conn) {}, //nolint:thelper // ref shopspring
		CloseConn: func(ctx context.Context, tb testing.TB, conn *pgx.Conn) { //nolint:thelper // ref shopspring
			err := conn.Close(ctx)
			if err != nil {
				tb.Errorf("Close failed: %v", err)
			}
		},
	}

	pgxtest.RunValueRoundTripTests(context.Background(),
		t, connTestRunner, nil, "numeric", []pgxtest.ValueRoundTripTest{
			{
				Param:  newFromString("1"),
				Result: new(decimal.Big),
				Test:   isExpectedEqDecimal(newFromString("1")),
			},
			{
				Param:  newFromString("0.000012345"),
				Result: new(decimal.Big),
				Test:   isExpectedEqDecimal(newFromString("0.000012345")),
			},
			{
				Param:  newFromString("123456.123456"),
				Result: new(decimal.Big),
				Test:   isExpectedEqDecimal(newFromString("123456.123456")),
			},
			{
				Param:  newFromString("-1"),
				Result: new(decimal.Big),
				Test:   isExpectedEqDecimal(newFromString("-1")),
			},
			{
				Param:  newFromString("-0.000012345"),
				Result: new(decimal.Big),
				Test:   isExpectedEqDecimal(newFromString("-0.000012345")),
			},
			{
				Param:  newFromString("-123456.123456"),
				Result: new(decimal.Big),
				Test:   isExpectedEqDecimal(newFromString("-123456.123456")),
			},
		})
}
