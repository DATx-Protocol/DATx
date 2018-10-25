#! /bin/bash

rm -rf ./build/*
rm -rf ./node_mudules/*

echo "start install node packages....."
npm install  -g node-gyp 
npm install  --unsafe-perm

echo "start install pkg...."
npm install -g pkg
#cd ./node_mudules/scrypt

echo "start compile ...."
pkg .
rm -f mse-win.exe

sysOS=`uname -s`
if [ $sysOS == "Darwin" ];then
    echo "Detected Mac OS"
	rm -f mse-linux
    mv mse-macos mse
elif [ $sysOS == "Linux" ];then
	echo "Detected Linux OS"
    rm -f mse-macos
    mv mse-linux mse
else
	echo "unsupported OS"
fi

pwd=`pwd`
buildDir="${pwd}/build"
mkdir -p $buildDir

mv -f  mse  $buildDir

echo "build successful, see build floder"