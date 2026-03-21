module project

go 1.26.1

replace github.com/dracory/base => ../../_modules_dracory/base

// replace github.com/dracory/test => ../../_modules_dracory/test

replace github.com/dracory/rtr => ../../_modules_dracory/rtr

require (
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2
	github.com/aws/smithy-go v1.24.2
	github.com/disintegration/imaging v1.6.2
	github.com/dracory/api v1.7.0
	github.com/dracory/auditstore v0.3.0
	github.com/dracory/auth v0.29.0
	github.com/dracory/base v0.31.0
	github.com/dracory/blindindexstore v1.8.0
	github.com/dracory/blockeditor v0.24.0
	github.com/dracory/blogstore v1.6.0
	github.com/dracory/bs v0.16.0
	github.com/dracory/cachestore v0.23.0
	github.com/dracory/cdn v1.10.0
	github.com/dracory/chatstore v0.8.0
	github.com/dracory/cmd v0.2.0
	github.com/dracory/cmsstore v1.13.0
	github.com/dracory/crud/v2 v2.0.0-20260318082437-20bdfa2008c2
	github.com/dracory/csrf v0.2.0
	github.com/dracory/customstore v1.7.0
	github.com/dracory/dashboard v1.11.0
	github.com/dracory/database v0.6.0
	github.com/dracory/dataobject v1.6.0
	github.com/dracory/entitystore v1.3.0
	github.com/dracory/env v1.0.0
	github.com/dracory/envenc v1.2.0
	github.com/dracory/feedstore v0.7.0
	github.com/dracory/filesystem v1.1.0
	github.com/dracory/form v0.21.0
	github.com/dracory/geostore v1.0.0
	github.com/dracory/hb v1.88.0
	github.com/dracory/liveflux v0.25.0
	github.com/dracory/llm v1.3.0
	github.com/dracory/logstore v1.13.0
	github.com/dracory/metastore v1.4.0
	github.com/dracory/req v0.1.0
	github.com/dracory/rtr v1.5.0
	github.com/dracory/sb v0.20.0
	github.com/dracory/sessionstore v1.6.0
	github.com/dracory/settingstore v1.5.0
	github.com/dracory/shopstore v1.7.0
	github.com/dracory/statsstore v0.11.0
	github.com/dracory/str v0.17.0
	github.com/dracory/subscriptionstore v0.7.0
	github.com/dracory/taskstore v1.19.0
	github.com/dracory/test v0.9.0
	github.com/dracory/ui v0.17.0
	github.com/dracory/uid v1.9.0
	github.com/dracory/uncdn v0.9.0
	github.com/dracory/userstore v1.6.0
	github.com/dracory/vaultstore v0.34.0
	github.com/dracory/versionstore v0.6.0
	github.com/dracory/websrv v0.1.0
	github.com/dracory/wf v0.6.0
	github.com/dromara/carbon/v2 v2.6.16
	github.com/faabiosr/cachego v0.26.0
	github.com/flosch/pongo2/v6 v6.0.0
	github.com/go-co-op/gocron v1.37.0
	github.com/go-sql-driver/mysql v1.9.3
	github.com/jellydator/ttlcache/v3 v3.4.0
	github.com/lmittmann/tint v1.1.3
	github.com/mileusna/useragent v1.3.5
	github.com/mingrammer/cfmt v1.1.0
	github.com/robertkrimen/otto v0.5.1
	github.com/sfreiberg/simplessh v0.0.0-20220719182921-185eafd40485
	github.com/spf13/cast v1.10.0
	github.com/stretchr/testify v1.11.1
	github.com/tidwall/gjson v1.18.0
	github.com/yuin/goldmark v1.7.17
	modernc.org/sqlite v1.47.0
)

require (
	atomicgo.dev/cursor v0.2.0 // indirect
	atomicgo.dev/keyboard v0.2.9 // indirect
	atomicgo.dev/schedule v0.1.0 // indirect
	cloud.google.com/go v0.123.0 // indirect
	cloud.google.com/go/aiplatform v1.120.0 // indirect
	cloud.google.com/go/auth v0.18.2 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	cloud.google.com/go/iam v1.5.3 // indirect
	cloud.google.com/go/longrunning v0.8.0 // indirect
	cloud.google.com/go/vertexai v0.17.0 // indirect
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.41.4 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.7 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.20 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.20 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.21 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.20 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.20 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.97.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/containerd/console v1.0.5 // indirect
	github.com/dave/dst v0.27.3 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/davidmz/go-pageant v1.0.2 // indirect
	github.com/doug-martin/goqu/v9 v9.19.0 // indirect
	github.com/dracory/arr v0.2.0 // indirect
	github.com/dracory/crypto v0.3.0 // indirect
	github.com/dracory/shortcode v0.5.0 // indirect
	github.com/dracory/sqlfilestore v1.2.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gabriel-vasile/mimetype v1.4.13 // indirect
	github.com/go-chi/chi/v5 v5.2.5 // indirect
	github.com/go-chi/cors v1.2.2 // indirect
	github.com/go-chi/httprate v0.15.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.14 // indirect
	github.com/googleapis/gax-go/v2 v2.19.0 // indirect
	github.com/gookit/color v1.6.0 // indirect
	github.com/goravel/framework v1.17.2 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.8.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jedib0t/go-pretty/v6 v6.7.8 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/lithammer/fuzzysearch v1.1.8 // indirect
	github.com/mattn/go-runewidth v0.0.21 // indirect
	github.com/mattn/go-sqlite3 v2.0.1+incompatible // indirect
	github.com/mocktools/go-smtp-mock v1.10.0 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/pkg/sftp v1.13.10 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/pterm/pterm v0.12.83 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/sashabaranov/go-openai v1.41.2 // indirect
	github.com/teambition/rrule-go v1.8.2 // indirect
	github.com/tidwall/match v1.2.0 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	github.com/zeebo/xxh3 v1.1.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.67.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.67.0 // indirect
	go.opentelemetry.io/otel v1.42.0 // indirect
	go.opentelemetry.io/otel/metric v1.42.0 // indirect
	go.opentelemetry.io/otel/trace v1.42.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20260312153236-7ab1446f8b90 // indirect
	golang.org/x/image v0.37.0 // indirect
	golang.org/x/oauth2 v0.36.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/term v0.41.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	google.golang.org/api v0.272.0 // indirect
	google.golang.org/genai v1.51.0 // indirect
	google.golang.org/genproto v0.0.0-20260319201613-d00831a3d3e7 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260319201613-d00831a3d3e7 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260319201613-d00831a3d3e7 // indirect
	google.golang.org/grpc v1.79.3 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/mysql v1.6.0 // indirect
	gorm.io/driver/postgres v1.6.0 // indirect
	gorm.io/driver/sqlite v1.6.0 // indirect
	gorm.io/gorm v1.31.1 // indirect
	modernc.org/libc v1.70.0 // indirect
)

require (
	github.com/darkoatanasovski/htmltags v1.0.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/georgysavva/scany v1.2.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible // indirect
	github.com/logrusorgru/aurora v2.0.3+incompatible // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/samber/lo v1.53.0
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
)
