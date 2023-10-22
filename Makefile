
testAll: clean vet server agent test1 test2 test3 test4 test5 test6 test7 test8

.PHONY: vet
vet:
	go vet -vettool=$(which statictest-darwin-amd64) ./...


server:
	@echo "building server"
	go build -o ./cmd/server/server ./cmd/server/*.go

agent:
	@echo "building agent"
	go build -o ./cmd/agent/agent ./cmd/agent/*.go

test1:
	../go/bin/metricstest-darwin-amd64 -test.v -test.run=^TestIteration1$$ -agent-binary-path=./cmd/agent/agent -binary-path=./cmd/server/server > log1.txt

test2:
	../go/bin/metricstest-darwin-amd64 -test.v -test.run=^TestIteration2[AB]*$$ -agent-binary-path=./cmd/agent/agent -binary-path=./cmd/server/server -source-path=. > log2.txt

test3:
	../go/bin/metricstest-darwin-amd64 -test.v -test.run=^TestIteration3[AB]*$$ -agent-binary-path=./cmd/agent/agent -binary-path=./cmd/server/server -source-path=. > log3.txt

test4:
	../go/bin/metricstest-darwin-amd64 -test.v -test.run=^TestIteration4$$ -server-port="8080" -agent-binary-path=./cmd/agent/agent -binary-path=./cmd/server/server -source-path=. > log4.txt

test5:
	../go/bin/metricstest-darwin-amd64 -test.v -test.run=^TestIteration5$$ -server-port="8080" -agent-binary-path=./cmd/agent/agent -binary-path=./cmd/server/server -source-path=. > log5.txt

test6:
	../go/bin/metricstest-darwin-amd64 -test.v -test.run=^TestIteration6$$ -server-port="8080" -agent-binary-path=./cmd/agent/agent -binary-path=./cmd/server/server -source-path=. > log6.txt

test7:
	../go/bin/metricstest-darwin-amd64 -test.v -test.run=^TestIteration7$$ -server-port="8080" -agent-binary-path=./cmd/agent/agent -binary-path=./cmd/server/server -source-path=. > log7.txt

test8:
	../go/bin/metricstest-darwin-amd64 -test.v -test.run=^TestIteration8$$ -server-port="8080" -agent-binary-path=./cmd/agent/agent -binary-path=./cmd/server/server -source-path=. > log8.txt


.PHONY: clean
clean:
	rm -f ./cmd/agent/agent ./cmd/server/server
	rm -f ./log*.txt
	rm -f ./tempfile*

.PHONY: lnt
lnt:
	# excluded 'paralleltest' by the reason - not now
	# excluded 'wsl' by the reason - 'wsl' and 'gofumpt' fights between each other
	# excluded 'gochecknoglobals' by reason - I need global variables sometimes
	# excluded 'exhaustivestruct' - deprecated
	# excluded 'depguard' - no need in it
	# golangci-lint run -v --enable-all --disable gochecknoglobals --disable paralleltest --disable exhaustivestruct --disable depguard --disable wsl
	golangci-lint run -v

.PHONY: fmt
fmt:
	# to install it:
	# go install mvdan.cc/gofumpt@latest
	gofumpt -l -w .

.PHONY: gci
gci:
	# to install it:
	# go install github.com/daixiang0/gci@latest
	gci write --skip-generated -s default .

.PHONY: gofmt
gofmt:
	gofmt -s -w .

.PHONY: fix
fix: gofmt gci fmt

.PHONY: cover
cover:
	rm -f ./cover.html cover.out coverage.txt
	go test -race -coverprofile cover.out  ./... ./internal/... -coverpkg=./...
	go tool cover -html=cover.out -o cover.html
	#https://blog.seriesci.com/how-to-measure-code-coverage-in-go/
	go tool cover -func cover.out | grep total:
