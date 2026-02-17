run:
	go run ./cmd/analyst/main.go

wire:
	cd pkg/wire && go run -mod=mod github.com/google/wire/cmd/wire .
