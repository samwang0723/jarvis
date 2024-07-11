package pginit

import (
	"context"
	"fmt"
	"net"
	"time"

	gofrs "github.com/jackc/pgx-gofrs-uuid"
	zerologadapter "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

const (
	defaultMaxConns     = 25
	defaultMaxIdleConns = 25
	defaultMaxLifeTime  = 5 * time.Minute
)

// Option configures PGInit behavior.
type Option func(*PGInit)

// PGInit provides capabilities for connect to postgres with pgx.pool.
type PGInit struct {
	tracer                pgx.QueryTracer
	pgxConf               *pgxpool.Config
	registerDataTypeFuncs []func(typeMap *pgtype.Map)
	customDataTypeNames   []string
	logLvl                tracelog.LogLevel
}

type Config struct {
	User         string
	Password     string
	Host         string
	Port         string
	Database     string
	MaxConns     int32
	MaxIdleConns int32
	MaxLifeTime  time.Duration
}

// New initializes a PGInit using the provided Config and options. If
// opts is not provided it will initializes PGInit with default configuration.
func New(conf *Config, opts ...Option) (*PGInit, error) {
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		conf.User, conf.Password, net.JoinHostPort(conf.Host, conf.Port), conf.Database,
	)

	pgxConf, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	pgxConf.MaxConns = defaultMaxConns
	if conf.MaxConns != 0 {
		pgxConf.MaxConns = conf.MaxConns
	}

	pgxConf.MinConns = defaultMaxIdleConns
	if conf.MaxIdleConns != 0 && conf.MaxConns >= conf.MaxIdleConns {
		pgxConf.MinConns = conf.MaxIdleConns
	} else {
		pgxConf.MinConns = pgxConf.MaxConns
	}

	pgxConf.MaxConnLifetime = defaultMaxLifeTime
	if conf.MaxLifeTime != 0 {
		pgxConf.MaxConnLifetime = conf.MaxLifeTime
	}

	pgi := &PGInit{
		pgxConf:             pgxConf,
		logLvl:              tracelog.LogLevelWarn,
		customDataTypeNames: []string{},
	}

	for _, opt := range opts {
		opt(pgi)
	}

	if pgi.tracer != nil {
		pgi.pgxConf.ConnConfig.Tracer = pgi.tracer
	}

	pgi.pgxConf.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		for _, fn := range pgi.registerDataTypeFuncs {
			fn(conn.TypeMap())
		}

		return registerTypes(ctx, pgi.customDataTypeNames, conn)
	}

	return pgi, nil
}

func registerTypes(ctx context.Context, types []string, conn *pgx.Conn) error {
	// https://pkg.go.dev/github.com/jackc/pgx/v5/pgtype#hdr-New_PostgreSQL_Type_Support
	for _, typeName := range types {
		dataType, err := conn.LoadType(ctx, typeName)
		if err != nil {
			return fmt.Errorf("load type: %w", err)
		}

		conn.TypeMap().RegisterType(dataType)
	}

	return nil
}

// ConnPool initiates connection to database and return a pgxpool.Pool.
func (pgi *PGInit) ConnPool(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.NewWithConfig(ctx, pgi.pgxConf)
	if err != nil {
		return nil, fmt.Errorf("connect config: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

// WithLogger Add logger to pgx. if the request context contains request id,
// can pass in the request id context key to reqIDKeyFromCtx and logger will
// log with the request id. Only will log if the log level is equal and above pgx.LogLevelWarn.
func WithLogger(logger *zerolog.Logger, reqIDKeyFromCtx interface{}) Option {
	return func(pgi *PGInit) {
		pgi.pgxConf.ConnConfig.Tracer = &tracelog.TraceLog{
			LogLevel: pgi.logLvl,
			Logger: zerologadapter.NewLogger(*logger, zerologadapter.WithContextFunc(
				func(ctx context.Context, logWith zerolog.Context) zerolog.Context {
					if ctxValue, ok := ctx.Value(reqIDKeyFromCtx).(string); ok {
						logWith = logWith.Str("x-request-id", ctxValue)
					}

					return logWith
				},
			)),
		}
	}
}

// WithTracer set pgx tracer.
func WithTracer(tracer pgx.QueryTracer) Option {
	return func(pgi *PGInit) {
		pgi.tracer = tracer
	}
}

// WithLogLevel set pgx log level.
func WithLogLevel(zLvl zerolog.Level) Option {
	return func(pgi *PGInit) {
		switch {
		case zLvl == zerolog.DebugLevel:
			pgi.logLvl = tracelog.LogLevelDebug
		case zLvl == zerolog.InfoLevel:
			pgi.logLvl = tracelog.LogLevelInfo
		case zLvl == zerolog.WarnLevel:
			pgi.logLvl = tracelog.LogLevelWarn
		case zLvl == zerolog.ErrorLevel:
			pgi.logLvl = tracelog.LogLevelError
		case zLvl == zerolog.NoLevel:
			pgi.logLvl = tracelog.LogLevelNone
		}
	}
}

// WithUUIDType set pgx uuid type to gofrs/uuid.
func WithUUIDType() Option {
	return func(pgi *PGInit) {
		pgi.registerDataTypeFuncs = append(pgi.registerDataTypeFuncs, gofrs.Register)
	}
}

func WithCustomTypes(types []string) Option {
	return func(pgi *PGInit) {
		pgi.customDataTypeNames = append(pgi.customDataTypeNames, types...)
	}
}
