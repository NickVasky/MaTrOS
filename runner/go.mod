module github.com/NickVasky/MaTrOS/runner

go 1.24.4

require (
	github.com/NickVasky/MaTrOS/shared v0.1.0
	github.com/joho/godotenv v1.5.1
	github.com/segmentio/kafka-go v0.4.48
)

//replace github.com/NickVasky/MaTrOS/shared => ../shared

require (
	github.com/emersion/go-imap/v2 v2.0.0-beta.5 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
