module main

go 1.23.8

require ecc v0.0.0

require transaction v0.0.0

require (
	github.com/tsuna/endian v0.0.0-20250821203744-206f48965e13 // indirect
	golang.org/x/crypto v0.41.0 // indirect
)

replace ecc => ./ecc

replace transaction => ./transaction
