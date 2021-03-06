#! /bin/bash

rm -rf package-lock.json
rm -rf ./build/
rm -rf ./node_modules/

echo "start install node packages....."
# npm install  -g node-gyp 
# npm install  --unsafe-perm --allow-root --save-dev grunt
yarn global add node-gyp
yarn 

echo "start install pkg...."
#npm install -g pkg
yarn global add pkg
#cd ./node_mudules/scrypt

yarnBin=`yarn global bin`
export PATH=$PATH:${yarnBin}
echo $PATH

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

if [ -f "$buildDir/mse" ] ;then
    chmod +x "$buildDir/mse"
else
    echo "mse build error"
    exit
fi

echo "build successful, see build floder"


echo 'start building clmse...'
cd ../clmse

echo 'start install node packages.....'
yarn
yarn link

echo 'clmse build success'