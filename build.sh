#!/bin/bash
Version="1.0.0"
ToolName="ProtoGenExt"
 
export GOARCH=amd64
export GOPROXY=http://goproxy.cn

BuildBinary()
{
  set -e
  TargetDir=bin/"${1}"
  mkdir -p "${TargetDir}"
  export GOOS=${1}

  go build -o "${TargetDir}"
  tar zcvf bin/${ToolName}-${Version}-"${1}"-x86_64.tar.gz "${TargetDir}"
}

if [[ ${1} == "" ]]; then
  BuildBinary windows
  BuildBinary linux
  BuildBinary darwin
else
  BuildBinary "${1}"
fi