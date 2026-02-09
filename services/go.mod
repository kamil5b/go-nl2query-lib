module github.com/kamil5b/go-nl2query-lib/services

go 1.25.6

replace github.com/kamil5b/go-nl2query-lib/domains => ../domains

replace github.com/kamil5b/go-nl2query-lib/ports => ../ports

replace github.com/kamil5b/go-nl2query-lib/testsuites => ../testsuites

require (
	github.com/kamil5b/go-nl2query-lib/domains v0.0.0-00010101000000-000000000000
	github.com/kamil5b/go-nl2query-lib/ports v0.0.0-00010101000000-000000000000
	github.com/kamil5b/go-nl2query-lib/testsuites v0.0.0-00010101000000-000000000000
	github.com/toon-format/toon-go v0.0.0-20251202084852-7ca0e27c4e8c
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
