package server

import (
	"backend/api/inet"
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
)

var Server *server

type server struct {
	PgPool *pgxpool.Pool
	ChPool driver.Conn
	Config *ServerConfig
}

type PostgresConfig struct {
	/* connection string of the postgres instance */
	DSN string
}

type ClickhouseConfig struct {
	/* connection string of the clickhouse instance */
	DSN string
}

type ServerConfig struct {
	PG                         PostgresConfig
	CH                         ClickhouseConfig
	MappingFileMaxSize         uint64
	SymbolsBucket              string
	SymbolsBucketRegion        string
	SymbolsAccessKey           string
	SymbolsSecretAccessKey     string
	AttachmentsBucket          string
	AttachmentsBucketRegion    string
	AttachmentsAccessKey       string
	AttachmentsSecretAccessKey string
	AWSEndpoint                string
	AttachmentOrigin           string
	SiteOrigin                 string
	APIOrigin                  string
	OAuthGitHubKey             string
	OAuthGitHubSecret          string
	OAuthGoogleKey             string
	AccessTokenSecret          []byte
	RefreshTokenSecret         []byte
	OtelServiceName            string
}

func NewConfig() *ServerConfig {
	mappingFileMaxSize, err := strconv.ParseUint(os.Getenv("MAPPING_FILE_MAX_SIZE"), 10, 64)
	if err != nil {
		log.Println("using default value of MAPPING_FILE_MAX_SIZE")
		mappingFileMaxSize = 524_288_000
	}

	symbolsBucket := os.Getenv("SYMBOLS_S3_BUCKET")
	if symbolsBucket == "" {
		log.Println("SYMBOLS_S3_BUCKET env var not set, mapping file uploads won't work")
	}

	symbolsBucketRegion := os.Getenv("SYMBOLS_S3_BUCKET_REGION")
	if symbolsBucketRegion == "" {
		log.Println("SYMBOLS_S3_BUCKET_REGION env var not set, mapping file uploads won't work")
	}

	symbolsAccessKey := os.Getenv("SYMBOLS_ACCESS_KEY")
	if symbolsAccessKey == "" {
		log.Println("SYMBOLS_ACCESS_KEY env var not set, mapping file uploads won't work")
	}

	symbolsSecretAccessKey := os.Getenv("SYMBOLS_SECRET_ACCESS_KEY")
	if symbolsSecretAccessKey == "" {
		log.Println("SYMBOLS_SECRET_ACCESS_KEY env var not set, mapping file uploads won't work")
	}

	attachmentsBucket := os.Getenv("ATTACHMENTS_S3_BUCKET")
	if attachmentsBucket == "" {
		log.Println("ATTACHMENTS_S3_BUCKET env var not set, event attachment uploads won't work")
	}

	attachmentsBucketRegion := os.Getenv("ATTACHMENTS_S3_BUCKET_REGION")
	if attachmentsBucketRegion == "" {
		log.Println("ATTACHMENTS_S3_BUCKET_REGION env var not set, event attachment uploads won't work")
	}

	attachmentsAccessKey := os.Getenv("ATTACHMENTS_ACCESS_KEY")
	if attachmentsAccessKey == "" {
		log.Println("ATTACHMENTS_ACCESS_KEY env var not set, event attachment uploads won't work")
	}

	attachmentsSecretAccessKey := os.Getenv("ATTACHMENTS_SECRET_ACCESS_KEY")
	if attachmentsSecretAccessKey == "" {
		log.Println("ATTACHMENTS_SECRET_ACCESS_KEY env var not set, event attachment uploads won't work")
	}

	attachmentOrigin := os.Getenv("ATTACHMENTS_S3_ORIGIN")
	if attachmentOrigin == "" {
		log.Println("ATTACHMENTS_S3_ORIGIN env var not set, event attachment downloads will be proxied")
	}

	siteOrigin := os.Getenv("SITE_ORIGIN")
	if siteOrigin == "" {
		log.Fatal("SITE_ORIGIN env var not set. Need for Cross Origin Resource Sharing (CORS) to work.")
	}

	apiOrigin := os.Getenv("API_ORIGIN")
	if apiOrigin == "" {
		log.Fatal("API_ORIGIN env var not set. Need for proxying session attachments.")
	}

	oauthGitHubKey := os.Getenv("OAUTH_GITHUB_KEY")
	if oauthGitHubKey == "" {
		log.Println("OAUTH_GITHUB_KEY env var is not set, dashboard authn won't work")
	}

	oauthGitHubSecret := os.Getenv("OAUTH_GITHUB_SECRET")
	if oauthGitHubSecret == "" {
		log.Println("OAUTH_GITHUB_SECRET env var is not set, dashboard authn won't work")
	}

	oauthGoogleKey := os.Getenv("OAUTH_GOOGLE_KEY")
	if oauthGoogleKey == "" {
		log.Println("OAUTH_GOOGLE_KEY env var is not set, dashboard authn won't work")
	}

	atSecret := os.Getenv("SESSION_ACCESS_SECRET")
	if atSecret == "" {
		log.Println("SESSION_ACCESS_SECRET env var is not set, dashboard authn won't work")
	}

	rtSecret := os.Getenv("SESSION_REFRESH_SECRET")
	if rtSecret == "" {
		log.Println("SESSION_REFRESH_SECRET env var is not set, dashboard authn won't work")
	}

	postgresDSN := os.Getenv("POSTGRES_DSN")
	if postgresDSN == "" {
		log.Fatal("POSTGRES_DSN env var is not set, cannot start server")
	}

	clickhouseDSN := os.Getenv("CLICKHOUSE_DSN")
	if clickhouseDSN == "" {
		log.Fatal("CLICKHOUSE_DSN env var is not set, cannot start server")
	}

	otelServiceName := os.Getenv("OTEL_SERVICE_NAME")
	if otelServiceName == "" {
		log.Println("OTEL_SERVICE_NAME env var is not set, o11y will not work")
	}

	endpoint := os.Getenv("AWS_ENDPOINT_URL")

	return &ServerConfig{
		PG: PostgresConfig{
			DSN: postgresDSN,
		},
		CH: ClickhouseConfig{
			DSN: clickhouseDSN,
		},
		MappingFileMaxSize:         mappingFileMaxSize,
		SymbolsBucket:              symbolsBucket,
		SymbolsBucketRegion:        symbolsBucketRegion,
		SymbolsAccessKey:           symbolsAccessKey,
		SymbolsSecretAccessKey:     symbolsSecretAccessKey,
		AttachmentsBucket:          attachmentsBucket,
		AttachmentsBucketRegion:    attachmentsBucketRegion,
		AttachmentsAccessKey:       attachmentsAccessKey,
		AttachmentsSecretAccessKey: attachmentsSecretAccessKey,
		AWSEndpoint:                endpoint,
		AttachmentOrigin:           attachmentOrigin,
		SiteOrigin:                 siteOrigin,
		APIOrigin:                  apiOrigin,
		OAuthGitHubKey:             oauthGitHubKey,
		OAuthGitHubSecret:          oauthGitHubSecret,
		OAuthGoogleKey:             oauthGoogleKey,
		AccessTokenSecret:          []byte(atSecret),
		RefreshTokenSecret:         []byte(rtSecret),
		OtelServiceName:            otelServiceName,
	}
}

func Init(config *ServerConfig) {
	pgPool, err := pgxpool.New(context.Background(), config.PG.DSN)
	if err != nil {
		log.Fatalf("Unable to create PG connection pool: %v\n", err)
	}

	chOpts, err := clickhouse.ParseDSN(config.CH.DSN)
	if err != nil {
		log.Fatalf("Unable to parse CH connection string: %v\n", err)
	}

	if gin.Mode() == gin.ReleaseMode {
		chOpts.Settings = clickhouse.Settings{
			"use_query_cache":                1,
			"query_cache_ttl":                600,
			"query_cache_min_query_duration": 1000,
		}
	}

	chPool, err := clickhouse.Open(chOpts)
	if err != nil {
		log.Fatalf("Unable to create CH connection pool: %v", err)
	}

	if err := inet.Init(); err != nil {
		log.Fatalf("Unable to initialize geo ip lookup system: %v", err)
	}

	Server = &server{
		PgPool: pgPool,
		ChPool: chPool,
		Config: config,
	}
}

func (sc ServerConfig) InitTracer() func(context.Context) error {
	otelCollectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	otelInsecureMode := os.Getenv("OTEL_INSECURE_MODE")
	otelServiceName := sc.OtelServiceName

	var secureOption otlptracegrpc.Option

	if strings.ToLower(otelInsecureMode) == "false" || otelInsecureMode == "0" || strings.ToLower(otelInsecureMode) == "f" {
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(otelCollectorURL),
		),
	)

	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", otelServiceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
	}

	otel.SetTracerProvider(
		trace.NewTracerProvider(
			trace.WithSampler(trace.AlwaysSample()),
			trace.WithBatcher(exporter),
			trace.WithResource(resources),
		),
	)

	return exporter.Shutdown
}
