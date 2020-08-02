#!/bin/bash
echo "Starting benchmark"
echo "Cleaning log files"
rm *.log
rm -rf logs

echo "Building go binaries"
go build array-access-100.go
go build array-access-1000.go
go build array-access-10000.go
go build array-while-create-100.go
go build array-while-create-1000.go
go build array-while-create-10000.go
go build class-instance-100.go
go build class-instance-100-micro.go
go build class-instance-1000.go
go build class-instance-1000-micro.go
go build class-instance-10000.go
go build class-instance-10000-micro.go
go build dynamic-instance-100.go
go build dynamic-instance-100-micro.go
go build dynamic-instance-1000.go
go build dynamic-instance-1000-micro.go
go build dynamic-instance-10000.go
go build dynamic-instance-10000-micro.go
go build dynamic-method-100.go
go build dynamic-method-100-micro.go
go build dynamic-method-1000.go
go build dynamic-method-1000-micro.go
go build dynamic-method-10000.go
go build dynamic-method-10000-micro.go

x=1
while [ $x -le 10 ]
do
    echo "Running tests $x times"
    # ./array-access-100
    # ./array-access-1000
    # ./array-access-10000
    ./array-while-create-100
    ./array-while-create-1000
    ./array-while-create-10000
    # ./class-instance-100
    # ./class-instance-100-micro
    # ./class-instance-1000
    # ./class-instance-1000-micro
    # ./class-instance-10000
    # ./class-instance-10000-micro
    # ./dynamic-instance-100
    # ./dynamic-instance-100-micro
    # ./dynamic-instance-1000
    # ./dynamic-instance-1000-micro
    # ./dynamic-instance-10000
    # ./dynamic-instance-10000-micro
    # ./dynamic-method-100
    # ./dynamic-method-100-micro
    # ./dynamic-method-1000
    # ./dynamic-method-1000-micro
    # ./dynamic-method-10000
    # ./dynamic-method-10000-micro
    x=$(( $x + 1 ))
done

mkdir logs
mv *.log logs
echo "Benchmark logs saved at ./log"