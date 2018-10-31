#! /bin/bash -il

targetDir="bin"

pwd=`pwd`
targetFile=`basename $pwd`

rm -rf ./bin/*

echo "gopath $GOPATH"

if [ "$GOPATH" ] ;then
    echo "GOPATH found $GOPATH"
else
    echo "GOPATH not found"
    exit
fi

datxDir="${GOPATH}/src/datx/"
mkdir -p $datxDir
echo "datxDir $datxDir"

ln -s ${pwd} $datxDir

echo "start build $targetFile ..."

buildResult=`go build -o "${targetDir}/${targetFile}" "datx/${targetFile}" 2>&1`
echo "build command: go build -o ${targetDir}/${targetFile} datx/${targetFile}"
if [ -f "${targetDir}/${targetFile}" ] ;then
    chmod +x ${targetDir}/${targetFile}
else
    echo "build error $buildResult"
    exit
fi

lsdbin="${pwd}/${targetDir}/${targetFile}"
echo "lsd bin  $lsdbin"
if [ ! -f "${lsdbin}" ];then 
    printf "the %s is not exist.\n" ${lsdbin}
    exit
fi

echo "build success,you can start it."