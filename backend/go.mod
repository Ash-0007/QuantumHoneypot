module github.com/pqcd/backend

go 1.22.0

toolchain go1.24.3

require (
	github.com/cloudflare/circl v1.3.3
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/rs/cors v1.9.0
)

require (
	github.com/pqcd/backend/crypto v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/sys v0.10.0 // indirect
)

replace github.com/pqcd/backend/crypto => ./crypto
