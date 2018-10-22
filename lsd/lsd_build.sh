#! /bin/bash

targetDir="bin"

pwd=`pwd`
targetFile=`basename $pwd`

rm -rf ./build/*
rm -rf ./bin/*

datxDir="${pwd}/build/src/datx/${targetFile}"
mkdir -p $datxDir

cp -rf ./* $datxDir

export GOPATH=${pwd}/build

echo "gopath $GOPATH"

echo "start build $targetFile ..."

buildResult=`go build -o "${targetDir}/${targetFile}" "${datxDir}" 2>&1`
echo "build command: ${targetDir}/${targetFile} ${datxDir}"
if [ -z "$buildResult" ] ;then
    echo "build success $buildResult"
    chmod +x ${targetDir}/${targetFile}
else
    echo "build error $buildResult"
    exit
fi

echo "build success,you can start it."