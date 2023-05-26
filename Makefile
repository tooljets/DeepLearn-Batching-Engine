
lint:
	revive -formatter friendly -config revive.toml *.go
	revive -formatter friendly -config revive.toml service/*.go

build: