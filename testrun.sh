#!/bin/bash

#chmod 755 ./testrun.sh

function EPHEMERAL_PORT() {
	LOW_BOUND=49152
	RANGE=16384
	while true; do
		CANDIDATE=$(($LOW_BOUND + ($RANDOM % $RANGE)))
		(echo "" > /dev/tcp/127.0.0.1/${CANDIDATE}) > /dev/null 2>&1
		if [ $? -ne 0 ]; then
			echo $CANDIDATE
			break
		fi
	done
}

function CLEAN_AFTER_TEST() {
	unset ADDRESS
	unset RANDOM_PORT
	if [ -f "${TEMP_FILE}" ]; then
		echo "${TEMP_FILE} exists, deleting..."
		rm "${TEMP_FILE}"
	fi
	unset TEMP_FILE
	unset RESTORE
}

echo "go vet..."
go_vet_result=$( (go vet -vettool=$(which statictest-darwin-amd64) ./...) 2>&1)

if [ -z "$go_vet_result" ]; then
	echo "go vet passed"
else
	echo "go vet failed:"
	echo "$go_vet_result"
	exit
fi

echo "cleaning..."
rm ./ag ./ser
rm ./log*.txt

echo "building server"
go build -o ser ./cmd/server/*.go

echo "building agent"
go build -o ag ./cmd/agent/*.go

echo "Iter 1..."
metricstest-darwin-amd64 -test.v -test.run=^TestIteration1$ -agent-binary-path=./ag -binary-path=./ser > log1.txt
echo "Iter 1: $(tail -1 ./log1.txt)"

echo "Iter 2..."
metricstest-darwin-amd64 -test.v -test.run=^TestIteration2[AB]*$ -agent-binary-path=./ag -binary-path=./ser -source-path=. > log2.txt
echo "Iter 2: $(tail -1 ./log2.txt)"

echo "Iter 3..."
metricstest-darwin-amd64 -test.v -test.run=^TestIteration3[AB]*$ -agent-binary-path=./ag -binary-path=./ser -source-path=. > log3.txt
echo "Iter 3: $(tail -1 ./log3.txt)"

echo "Iter 4..."
RANDOM_PORT=$(EPHEMERAL_PORT)
echo RANDOM_PORT: "$RANDOM_PORT"
export ADDRESS="localhost:${RANDOM_PORT}"
metricstest-darwin-amd64 -test.v -test.run=^TestIteration4$ -server-port="$RANDOM_PORT" -agent-binary-path=./ag -binary-path=./ser -source-path=. > log4.txt
CLEAN_AFTER_TEST
echo "Iter 4: $(tail -1 ./log4.txt)"

echo "Iter 5..."
RANDOM_PORT=$(EPHEMERAL_PORT)
echo RANDOM_PORT: "$RANDOM_PORT"
export ADDRESS="localhost:${RANDOM_PORT}"
metricstest-darwin-amd64 -test.v -test.run=^TestIteration5$ -server-port="$RANDOM_PORT" -agent-binary-path=./ag -binary-path=./ser -source-path=. > log5.txt
CLEAN_AFTER_TEST
echo "Iter 5: $(tail -1 ./log5.txt)"

echo "Iter 6..."
RANDOM_PORT=$(EPHEMERAL_PORT)
echo RANDOM_PORT: "$RANDOM_PORT"
export ADDRESS="localhost:${RANDOM_PORT}"
metricstest-darwin-amd64 -test.v -test.run=^TestIteration6$ -server-port="$RANDOM_PORT" -agent-binary-path=./ag -binary-path=./ser -source-path=. > log6.txt
CLEAN_AFTER_TEST
echo "Iter 6: $(tail -1 ./log6.txt)"

echo "Iter 7..."
RANDOM_PORT=$(EPHEMERAL_PORT)
RANDOM_PORT="8080"
RESTORE=false
echo RANDOM_PORT: "$RANDOM_PORT"
export ADDRESS="localhost:${RANDOM_PORT}"
metricstest-darwin-amd64 -test.v -test.run=^TestIteration7$ -server-port="$RANDOM_PORT" -agent-binary-path=./ag -binary-path=./ser -source-path=. > log7.txt
CLEAN_AFTER_TEST
echo "Iter 7: $(tail -1 ./log7.txt)"

echo "Iter 8..."
rm /tmp/metrics-db.json
RANDOM_PORT=$(EPHEMERAL_PORT)
echo RANDOM_PORT: "$RANDOM_PORT"
export ADDRESS="localhost:${RANDOM_PORT}"
export TEMP_FILE="./tempfile${RANDOM_PORT}"
echo TEMP FILE: "$TEMP_FILE"
export RESTORE=true
metricstest-darwin-amd64 -test.v -test.run=^TestIteration8$ -server-port="$RANDOM_PORT" -agent-binary-path=./ag -binary-path=./ser -source-path=. > log8.txt
CLEAN_AFTER_TEST
echo "Iter 8: $(tail -1 ./log8.txt)"

echo "Iter 9..."
rm /tmp/metrics-db.json
RANDOM_PORT=$(EPHEMERAL_PORT)
echo RANDOM_PORT: "$RANDOM_PORT"
export ADDRESS="localhost:${RANDOM_PORT}"
export TEMP_FILE="./tempfile${RANDOM_PORT}"
echo TEMP FILE: "$TEMP_FILE"
export RESTORE=true
metricstest-darwin-amd64 -test.v -test.run=^TestIteration9$ -file-storage-path=$TEMP_FILE -server-port="$RANDOM_PORT" -agent-binary-path=./ag -binary-path=./ser -source-path=. > log9.txt
CLEAN_AFTER_TEST
echo "Iter 9: $(tail -1 ./log9.txt)"
