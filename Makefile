build:

test:
	go test -modfile test.mod -cover -coverprofile coverage.out -race -memprofile=mem.out -cpuprofile=cpu.out -v ./...

bench: benchname:=.
bench:
	go test -modfile test.mod -bench="$(benchname)" -count 10 -ldflags '-s -w' -run="^$$" -benchmem -memprofile=mem.out -cpuprofile=cpu.out -v ./...

run:
