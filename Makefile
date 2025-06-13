build:

test:
	go test -modfile test.mod -cover coverprofile coverage.out -json -race -memprofile=mem.out -cpuprofile=cpu.out -v ./...

bench: benchname:=.
bench:
	go test -modfile test.mod -bench="$(benchname)" -run="^$$" -benchmem -memprofile=mem.out -cpuprofile=cpu.out -v ./...

run:
