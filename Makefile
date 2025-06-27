build:

test:
	go test -modfile devdeps.mod -cover -coverprofile coverage.out -race -memprofile=mem.out -cpuprofile=cpu.out -v ./...

bench: benchname:=.
bench:
	go test -modfile devdeps.mod -bench="$(benchname)" -count 5 -run="^$$" -benchmem -memprofile=mem.out -cpuprofile=cpu.out -v ./...

run:
