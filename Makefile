
testAll: clean vet server agent test1 test2 test3 test4 test5 test6 test7 test8

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


clean:
	rm -f ./cmd/agent/agent ./cmd/server/server
	rm -f ./log*.txt
	rm -f ./tempfile*

lnt:
	golangci-lint run --enable-all --disable gochecknoglobals --disable paralleltest --disable exhaustivestruct --disable depguard --disable ifshort


fmt:
	# to install it:
	# go install mvdan.cc/gofumpt@latest
	gofumpt -l -w .