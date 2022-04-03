#!/bin/bash

CURR_PATH=$(pwd)

cd $CURR_PATH/bilibili-API-collect/grpc_api/bilibili/app/archive/v1
protoc --proto_path=. \
    --go_out=$CURR_PATH/proto/archive \
    --go_opt=Marchive.proto="./;archive" \
    archive.proto

cd $CURR_PATH/bilibili-API-collect/grpc_api/bilibili/community/service/dm/v1
protoc --proto_path=. \
    --go_out=$CURR_PATH/proto/dm \
    --go_opt=Mdm.proto="./;dm" \
    dm.proto

cd $CURR_PATH
go mod tidy