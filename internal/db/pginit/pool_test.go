package pginit

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog"
	"go.uber.org/goleak"
)

var testHost, testPort string //nolint:gochecknoglobals // postgres dockertest info

func TestMain(m *testing.M) {
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	err = CheckAndDownloadImg("public.ecr.aws/docker/library/postgres", "14.8")
	if err != nil {
		log.Fatalf("Could not run CheckAndDownloadImg: %v", err)
	}

	resource, err := dockerPool.Run("public.ecr.aws/docker/library/postgres", "14.8", []string{
		"POSTGRES_PASSWORD=postgres",
		"POSTGRES_USER=postgres",
		"POSTGRES_DB=datawarehouse",
		"listen_addresses = '*'",
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	databaseURL := fmt.Sprintf(
		"postgres://postgres:%s@%s/datawarehouse?sslmode=disable",
		"postgres",
		getHostPort(resource, "5432/tcp"),
	)

	if err := dockerPool.Retry(func() error {
		ctx := context.Background()
		db, err := pgx.Connect(ctx, databaseURL)
		if err != nil {
			return fmt.Errorf("pgx connect: %w", err)
		}
		if err := db.Ping(ctx); err != nil {
			return fmt.Errorf("ping: %w", err)
		}
		db.Close(ctx)

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
				code = 1

				log.Fatalf("goleak: Errors on successful test run: %v\n", err)
			}
		}
	}

	// You can't defer this because os.Exit doesn't care for defer
	if err := dockerPool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

// CheckAndDownloadImg is a helper function to check and download image before run the concurrent test
// It prevent the api ratelimit due to concurrent calling pull image api in different test cases. This
// function should be called in TestMain function if you want to start docker instances in different
// testcases concurrently.
func CheckAndDownloadImg(imageName, imageTag string) error {
	retry := 3
	sleepTime := 100 * time.Millisecond

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

func getHostPort(resource *dockertest.Resource, id string) string {
	hostAndPort := resource.GetHostPort(id)
	hp := strings.Split(hostAndPort, ":")
	testHost = hp[0]
	testPort = hp[1]

	return testHost + ":" + testPort
}

func TestConnPool(t *testing.T) {
	t.Parallel()

	type args struct {
		Config Config
	}

	tests := []struct {
		name       string
		args       args
		wantConfig Config
		wantErr    bool
	}{
		{
			name: "expecting no error with default connection setting",
			args: args{
				Config{
					Host:         testHost,
					Port:         testPort,
					User:         "postgres",
					Password:     "postgres",
					Database:     "datawarehouse",
					MaxConns:     2,
					MaxIdleConns: 2,
				},
			},
			wantConfig: Config{
				MaxConns:     2,
				MaxIdleConns: 2,
				MaxLifeTime:  5 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "expecting no error with custom connection setting",
			args: args{
				Config{
					Host:         testHost,
					Port:         testPort,
					User:         "postgres",
					Password:     "postgres",
					Database:     "datawarehouse",
					MaxConns:     3,
					MaxIdleConns: 3,
					MaxLifeTime:  10 * time.Minute,
				},
			},
			wantConfig: Config{
				MaxConns:     3,
				MaxIdleConns: 3,
				MaxLifeTime:  10 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "expecting error with wrong user setting",
			args: args{
				Config{
					Host:     testHost,
					Port:     testPort,
					User:     "wrong",
					Password: "postgres",
					Database: "datawarehouse",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			logger := zerolog.New(io.Discard)

			db, err := NewTestConnPool(
				context.Background(),
				t,
				&tt.args.Config,
				[]Option{WithLogger(&logger, "")},
			)
			if tt.wantErr && err == nil {
				t.Errorf("expects err but nil returned")
			}

			defer func() {
				if db != nil {
					db.Close()
				}
			}()

			if err != nil {
				if !tt.wantErr {
					t.Errorf("expect no err but err returned: %s", err)
				}

				return
			}

			if err := db.Ping(ctx); err != nil {
				t.Errorf("cannot ping db: %s", err)
			}

			if db.Config().MaxConns != tt.wantConfig.MaxConns {
				t.Errorf("expected (%v) but got (%v)", tt.wantConfig.MaxConns, db.Config().MaxConns)
			}

			if db.Config().MaxConnLifetime != tt.wantConfig.MaxLifeTime {
				t.Errorf(
					"expected (%v) but got (%v)",
					tt.wantConfig.MaxLifeTime,
					db.Config().MaxConnLifetime,
				)
			}

			if db.Config().MinConns != tt.wantConfig.MaxIdleConns {
				t.Errorf(
					"expected (%v) but got (%v)",
					tt.wantConfig.MaxIdleConns,
					db.Config().MinConns,
				)
			}
		})
	}
}

func TestConnPoolWithLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		lvl       zerolog.Level
		wantedLvl tracelog.LogLevel
	}{
		{
			name:      "level debug",
			lvl:       zerolog.DebugLevel,
			wantedLvl: tracelog.LogLevelDebug,
		},
		{
			name:      "level info",
			lvl:       zerolog.InfoLevel,
			wantedLvl: tracelog.LogLevelInfo,
		},
		{
			name:      "level warn",
			lvl:       zerolog.WarnLevel,
			wantedLvl: tracelog.LogLevelWarn,
		},
		{
			name:      "level error",
			lvl:       zerolog.ErrorLevel,
			wantedLvl: tracelog.LogLevelError,
		},
		{
			name:      "level none",
			lvl:       zerolog.NoLevel,
			wantedLvl: tracelog.LogLevelNone,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			type contextKey string

			requestIDKey := contextKey("x-request-id")

			// set the request id
			reqID := uuid.Must(uuid.NewV7()).String()
			ctx := context.WithValue(context.Background(), requestIDKey, reqID)

			out := &SafeBuffer{
				b: bytes.NewBuffer([]byte{}),
			}
			logger := zerolog.New(out)

			db, err := NewTestConnPool(ctx, t, &Config{
				Host:     testHost,
				Port:     testPort,
				User:     "postgres",
				Password: "postgres",
				Database: "datawarehouse",
				MaxConns: 1,
			}, []Option{
				WithLogLevel(tt.lvl),
				WithLogger(&logger, requestIDKey),
				WithUUIDType(),
			})
			if err != nil {
				t.Error("expected no error")
			}

			defer func() {
				if db != nil {
					db.Close()
				}
			}()

			if err := db.Ping(ctx); err != nil {
				t.Error("expected no error")
			}

			if db.Config().ConnConfig.Tracer == nil {
				t.Error("expected logger not nil")
			}

			if _, err := db.Exec(ctx, "SELECT * FROM ERROR"); err == nil {
				t.Error("expected return error")
			}

			actualLog := out.String()
			if tt.name != "level none" && !strings.Contains(actualLog, reqID) {
				t.Errorf("expected log contains request id")
			}
		})
	}
}

func TestConnPool_mapCustomDataTypes(t *testing.T) {
	t.Parallel()

	type contextKey string

	requestIDKey := contextKey("x-request-id")

	logger := zerolog.New(io.Discard)

	tests := []struct {
		name             string
		opts             []Option
		expectErrDecimal bool
		expectErrUUID    bool
	}{
		{
			name: "decimal + uuid",
			opts: []Option{
				WithLogLevel(zerolog.DebugLevel),
				WithLogger(&logger, requestIDKey),
				WithUUIDType(),
			},
			expectErrDecimal: false,
			expectErrUUID:    false,
		},
		{
			name: "uuid + decimal",
			opts: []Option{
				WithLogLevel(zerolog.DebugLevel),
				WithLogger(&logger, requestIDKey),
				WithUUIDType(),
			},
			expectErrDecimal: false,
			expectErrUUID:    false,
		},
		{
			name: "decimal",
			opts: []Option{
				WithLogLevel(zerolog.DebugLevel),
				WithLogger(&logger, requestIDKey),
			},
			expectErrDecimal: false,
			expectErrUUID:    true,
		},
		{
			name: "uuid",
			opts: []Option{
				WithLogLevel(zerolog.DebugLevel),
				WithLogger(&logger, requestIDKey),
				WithUUIDType(),
			},
			expectErrDecimal: true,
			expectErrUUID:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			db, err := NewTestConnPool(context.Background(), t, &Config{
				Host:     testHost,
				Port:     testPort,
				User:     "postgres",
				Password: "postgres",
				Database: "datawarehouse",
				MaxConns: 1,
			}, tt.opts)
			if err != nil {
				t.Error("expected no error")
			}

			defer func() {
				if db != nil {
					db.Close()
				}
			}()

			err = db.Ping(ctx)
			if err != nil {
				t.Error("expected no error")
			}

			u := &uuid.UUID{}
			err = db.QueryRow(context.Background(), "select 'b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5'").
				Scan(u)

			if err != nil && !tt.expectErrUUID {
				t.Errorf("expected no err: %s", err)
			}
		})
	}
}

func TestConnPool_mapCustomTypes_CRUD(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
	}{
		{
			name: "CRUD operation with custom type uuid and decimal",
		},
	}

	logger := zerolog.New(io.Discard)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pool, err := NewTestConnPool(context.Background(), t, &Config{
				Host:     testHost,
				Port:     testPort,
				User:     "postgres",
				Password: "postgres",
				Database: "datawarehouse",
				MaxConns: 1,
			}, []Option{
				WithLogger(&logger, "cryptoRequestID"),
				WithUUIDType(),
			})
			if err != nil {
				t.Error("expected no error")
			}

			defer func() {
				if pool != nil {
					pool.Close()
				}
			}()

			conn, err := pool.Acquire(ctx)
			if err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			defer conn.Release()

			tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
			if err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			defer tx.Rollback(ctx)

			_, err = tx.Exec(
				ctx,
				"CREATE TABLE IF NOT EXISTS uuid_decimal(uuid uuid, PRIMARY KEY (uuid))",
			)
			if err != nil {
				t.Errorf("expected no error but got: %v", err)
			}

			// create
			id, _ := uuid.FromString("b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5")

			row := tx.QueryRow(
				ctx,
				"INSERT INTO uuid_decimal(uuid) VALUES($1) RETURNING uuid",
				id,
			)

			r := struct {
				uuid uuid.UUID
			}{}

			if err = row.Scan(&r.uuid); err != nil { //nolint:govet // inline err is within scope
				t.Errorf("expected no error but got: %v, (%+v)", err, row)
			}

			if r.uuid.String() != id.String() {
				t.Error("inserted data doesn't match with input")
			}

			// read
			rows, err := tx.Query(ctx, "SELECT * FROM uuid_decimal")
			if err != nil {
				t.Errorf("expected no error but got: %v", err)
			}

			defer rows.Close()

			var results []struct {
				uuid uuid.UUID
			}

			for rows.Next() {
				r = struct { //nolint:govet // r is within loop scope
					uuid uuid.UUID
				}{}

				if err := rows.Scan(&r.uuid); err != nil {
					t.Errorf("expected no error but got: %v", err)
				}

				if r.uuid.String() != id.String() {
					t.Error("inserted data doesn't match with input")
				}

				results = append(results, r)
			}

			if len(results) != 1 {
				t.Errorf("expected 1 result but got: %v", len(results))
			}

			// delete
			row = tx.QueryRow(
				ctx,
				"DELETE FROM uuid_decimal WHERE uuid = $1 RETURNING uuid",
				id,
			)

			if err := row.Scan(&r.uuid); err != nil {
				t.Errorf("expected no error but got: %v, (%+v)", err, row)
			}

			if r.uuid.String() != id.String() {
				t.Error("inserted data doesn't match with input")
			}

			row = tx.QueryRow(ctx, "SELECT * FROM uuid_decimal WHERE uuid = $1", id)
			if err := row.Scan(&r.uuid); err != nil && !errors.Is(err, pgx.ErrNoRows) {
				t.Errorf("expected no error but got: %v, (%+v)", err, row)
			}
		})
	}
}

type lightState string

const (
	lightStateOn  lightState = "on"
	lightStateOff lightState = "off"
)

func (e *lightState) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = lightState(s)
	case string:
		*e = lightState(s)
	default:
		return fmt.Errorf("unsupported scan type for lightState: %T", src)
	}

	return nil
}

func TestConnPool_WithCustomTypes(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	logger := zerolog.New(io.Discard)

	// Prepare the db
	pool, err := NewTestConnPool(context.Background(), t, &Config{
		Host:     testHost,
		Port:     testPort,
		User:     "postgres",
		Password: "postgres",
		Database: "datawarehouse",
		MaxConns: 1,
	}, []Option{
		WithLogger(&logger, "cryptoRequestID"),
	})
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}

	t.Cleanup(func() {
		conn.Release()
	})

	_, err = conn.Exec(ctx, `
	BEGIN;
	CREATE TYPE light_state AS ENUM ('on', 'off');
	CREATE TABLE IF NOT EXISTS light_events(id serial, state light_state, PRIMARY KEY (id));
	INSERT INTO light_events(state) VALUES ('on'), ('on'), ('off'), ('off');
	COMMIT;
	`)
	if err != nil {
		t.Fatalf("expected no error but got: %v", err)
	}

	type record struct {
		State lightState
		ID    int
	}

	scanRows := func(t *testing.T, rows pgx.Rows) []record {
		t.Helper()

		results := []record{}

		for rows.Next() {
			r := record{}

			if scanErr := rows.Scan(&r.ID, &r.State); scanErr != nil {
				t.Errorf("expected no error but got: %v", scanErr)
			}

			results = append(results, r)
		}

		return results
	}

	tests := []struct {
		name           string
		expectedResult []record
		typeNames      []string
		wantErr        bool
	}{
		{
			name:           "without custom types options",
			wantErr:        true,
			expectedResult: []record{},
			typeNames:      []string{},
		},
		{
			typeNames: []string{"light_state", "_light_state"},
			wantErr:   false,
			expectedResult: []record{
				{ID: 1, State: lightStateOn},
				{ID: 2, State: lightStateOn},
				{ID: 3, State: lightStateOff},
				{ID: 4, State: lightStateOff},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pool, err := NewTestConnPool(context.Background(), t, &Config{
				Host:     testHost,
				Port:     testPort,
				User:     "postgres",
				Password: "postgres",
				Database: "datawarehouse",
				MaxConns: 3,
			}, []Option{WithCustomTypes(tt.typeNames)})
			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}

			t.Cleanup(func() { pool.Close() })

			conn, err := pool.Acquire(ctx)
			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}

			t.Cleanup(func() { conn.Release() })

			rows, err := conn.Query(
				ctx,
				"SELECT * FROM light_events WHERE state = ANY($1::light_state[]) order by id",
				[]lightState{lightStateOn, lightStateOff},
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr: %v, got: %v", tt.wantErr, err)
			}

			defer rows.Close()

			records := scanRows(t, rows)
			if !cmp.Equal(records, tt.expectedResult) {
				t.Errorf("Query diff = %v", cmp.Diff(records, tt.expectedResult))
			}
		})
	}
}

func NewTestConnPool(
	ctx context.Context,
	t *testing.T,
	config *Config,
	opts []Option,
) (*pgxpool.Pool, error) {
	t.Helper()

	pgi, err := New(config, opts...)
	if err != nil {
		t.Fatalf("error pginit: %s", err)
	}

	pool, err := pgi.ConnPool(ctx)
	if err != nil {
		return pool, err
	}

	return pool, nil
}

type SafeBuffer struct {
	b *bytes.Buffer
	m sync.Mutex
}

func (b *SafeBuffer) Read(p []byte) (int, error) {
	b.m.Lock()
	defer b.m.Unlock()

	return b.b.Read(p)
}

func (b *SafeBuffer) Write(p []byte) (int, error) {
	b.m.Lock()
	defer b.m.Unlock()

	return b.b.Write(p)
}

func (b *SafeBuffer) String() string {
	b.m.Lock()
	defer b.m.Unlock()

	return b.b.String()
}

func BenchmarkConnPool(b *testing.B) {
	logger := zerolog.New(io.Discard)

	for i := 0; i <= b.N; i++ {
		ctx := context.Background()

		b.StartTimer()

		pgi, _ := New(
			&Config{
				Port:     testPort,
				User:     "postgres",
				Password: "postgres",
				Database: "datawarehouse",
				MaxConns: 1,
			},
			WithLogger(&logger, "cryptoRequestID"),
			WithUUIDType(),
		)

		pool, err := pgi.ConnPool(ctx)
		if err != nil {
			b.Errorf("expected no error but got: %v", err)
		}

		b.StopTimer()

		pool.Close()
	}
}
