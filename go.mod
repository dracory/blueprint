module project

go 1.25.0

// replace github.com/dracory/base => ../base
// replace github.com/dracory/taskstore => ../../_modules_dracory/taskstore

// replace github.com/dracory/cachestore => ../../_modules_dracory/cachestore

// replace github.com/dracory/sessionstore => ../../_modules_dracory/sessionstore

require (
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2
	github.com/disintegration/imaging v1.6.2
	github.com/dracory/api v1.7.0
	github.com/dracory/auditstore v0.2.0
	github.com/dracory/auth v0.25.0
	github.com/dracory/base v0.26.0
	github.com/dracory/blindindexstore v1.6.0
	github.com/dracory/blockeditor v0.24.0
	github.com/dracory/blogstore v1.1.0
	github.com/dracory/bs v0.15.0
	github.com/dracory/cachestore v0.23.0
	github.com/dracory/cdn v1.8.0
	github.com/dracory/cmd v0.2.0
	github.com/dracory/cmsstore v0.34.0
	github.com/dracory/crud/v2 v2.0.0-20251030193142-403ea1e5e710
	github.com/dracory/csrf v0.2.0
	github.com/dracory/customstore v1.6.0
	github.com/dracory/dashboard v1.10.0
	github.com/dracory/database v0.5.0
	github.com/dracory/dataobject v1.6.0
	github.com/dracory/entitystore v1.1.0
	github.com/dracory/env v1.0.0
	github.com/dracory/envenc v1.1.0
	github.com/dracory/feedstore v0.5.0
	github.com/dracory/filesystem v1.0.0
	github.com/dracory/form v0.19.0
	github.com/dracory/geostore v0.15.0
	github.com/dracory/hb v1.88.0
	github.com/dracory/liveflux v0.11.0
	github.com/dracory/llm v0.9.0
	github.com/dracory/logstore v1.9.0
	github.com/dracory/metastore v1.3.0
	github.com/dracory/req v0.1.0
	github.com/dracory/rtr v1.1.0
	github.com/dracory/sb v0.12.0
	github.com/dracory/sessionstore v1.5.2
	github.com/dracory/settingstore v1.2.0
	github.com/dracory/shopstore v1.5.1
	github.com/dracory/statsstore v0.7.0
	github.com/dracory/str v0.17.0
	github.com/dracory/subscriptionstore v0.5.0
	github.com/dracory/taskstore v1.8.1
	github.com/dracory/test v0.4.0
	github.com/dracory/ui v0.17.0
	github.com/dracory/uid v1.8.0
	github.com/dracory/uncdn v0.9.0
	github.com/dracory/userstore v1.5.0
	github.com/dracory/vaultstore v0.26.0
	github.com/dracory/websrv v0.1.0
	github.com/dracory/wf v0.6.0
	github.com/dromara/carbon/v2 v2.6.14
	github.com/faabiosr/cachego v0.26.0
	github.com/flosch/pongo2/v6 v6.0.0
	github.com/glebarez/sqlite v1.11.0
	github.com/go-co-op/gocron v1.37.0
	github.com/go-sql-driver/mysql v1.9.3
	github.com/gouniverse/responses v0.6.0
	github.com/jellydator/ttlcache/v3 v3.4.0
	github.com/lmittmann/tint v1.1.2
	github.com/mileusna/useragent v1.3.5
	github.com/mingrammer/cfmt v1.1.0
	github.com/mocktools/go-smtp-mock v1.10.0
	github.com/robertkrimen/otto v0.5.1
	github.com/sfreiberg/simplessh v0.0.0-20220719182921-185eafd40485
	github.com/spf13/cast v1.10.0
	github.com/stretchr/testify v1.11.1
	github.com/tidwall/gjson v1.18.0
	github.com/yuin/goldmark v1.7.13
	gorm.io/gorm v1.31.0
)

require (
	cloud.google.com/go v0.123.0 // indirect
	cloud.google.com/go/aiplatform v1.108.0 // indirect
	cloud.google.com/go/auth v0.17.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	cloud.google.com/go/iam v1.5.3 // indirect
	cloud.google.com/go/longrunning v0.7.0 // indirect
	cloud.google.com/go/vertexai v0.15.0 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.39.5 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.2 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.18.20 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.89.1 // indirect
	github.com/aws/smithy-go v1.23.1 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/davidmz/go-pageant v1.0.2 // indirect
	github.com/doug-martin/goqu/v9 v9.19.0 // indirect
	github.com/dracory/arr v0.2.0 // indirect
	github.com/dracory/crypto v0.3.0 // indirect
	github.com/dracory/sqlfilestore v1.0.0 // indirect
	github.com/dracory/versionstore v0.5.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gabriel-vasile/mimetype v1.4.11 // indirect
	github.com/glebarez/go-sqlite v1.22.0 // indirect
	github.com/go-chi/chi/v5 v5.2.3 // indirect
	github.com/go-chi/cors v1.2.2 // indirect
	github.com/go-chi/httprate v0.15.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.6 // indirect
	github.com/googleapis/gax-go/v2 v2.15.0 // indirect
	github.com/goravel/framework v1.16.5 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/gouniverse/base v0.9.0 // indirect
	github.com/gouniverse/shortcode v0.4.0 // indirect
	github.com/gouniverse/validator v0.11.0 // indirect
	github.com/jedib0t/go-pretty/v6 v6.6.9 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/pkg/sftp v1.13.10 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/sashabaranov/go-openai v1.41.2 // indirect
	github.com/teambition/rrule-go v1.8.2 // indirect
	github.com/tidwall/match v1.2.0 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.63.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.63.0 // indirect
	go.opentelemetry.io/otel v1.38.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20251023183803-a4bb9ffd2546 // indirect
	golang.org/x/image v0.32.0 // indirect
	golang.org/x/oauth2 v0.32.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/term v0.36.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	google.golang.org/api v0.254.0 // indirect
	google.golang.org/genai v1.33.0 // indirect
	google.golang.org/genproto v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/grpc v1.76.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/libc v1.66.10 // indirect
	modernc.org/sqlite v1.40.0 // indirect
)

require (
	github.com/darkoatanasovski/htmltags v1.0.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/georgysavva/scany v1.2.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gouniverse/maputils v0.7.0
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible // indirect
	github.com/logrusorgru/aurora v2.0.3+incompatible // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/samber/lo v1.52.0
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
)
