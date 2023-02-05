
cd ./usecase1
go run ./alice/cmd/worker/main.go & go run ./bob/cmd/worker/main.go  & go run ./cmd/worker/main.go & wait


