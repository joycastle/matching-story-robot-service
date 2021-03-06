#!/bin/bash
cd $(dirname "$0")
DIR="$(pwd)"
TargetDir=$DIR/../confmanager/template
CONFIG_BRANCH="$1"
CONFIG_REPO_URL="git@codeup.aliyun.com:62b023a03e81781f3ad195c6/Server_V2_Config.git"

if [ -d "./Server_V2_Config" ];then
    rm -rf ./Server_V2_Config
fi

go env -w GO111MODULE=off
go build format.go
go env -w GO111MODULE=on

git clone \
  --branch "$CONFIG_BRANCH" \
  --depth 1  \
  "$CONFIG_REPO_URL"

cd Server_V2_Config || exit

RELEASE_VERSION=$(sed -n 1p RELEASE_VERSION)

if [ "$2" != "" ];then
	RELEASE_VERSION=$2
fi

echo "当前发行版本号：""$RELEASE_VERSION"


for f in ${TargetDir}/*
do
	fname=$RELEASE_VERSION/$(basename $f)
	../format -f $fname -ff $f
	echo "新文件:"$fname" 写入："$f
done

echo  "配置文件夹"$TargetDir



#bom
#Bytes	Encoding Form
#00 00 FE FF        UTF-32, big-endian
#FF FE 00 00        UTF-32, little-endian
#FE FF	                UTF-16, big-endian
#FF FE	                UTF-16, little-endian
#EF BB BF	        UTF-8
