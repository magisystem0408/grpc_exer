#!/bin/bash

python -m grpc_tools.protoc -I protos --python_out=client --grpc_python_out=client protos/mamushi.proto