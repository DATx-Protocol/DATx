#!/bin/bash

if [ $# -ne 1 ]; then
    echo Usage: $0 my_new_accessory
    echo ... where my_new_accessory is the name of the accessory you want to create
    exit 1
fi

accessoryName=$1

echo Copying template...
cp -r template_accessory $accessoryName

echo Renaming files/directories...
mv $accessoryName/include/datxos/template_accessory $accessoryName/include/datxos/$accessoryName
for file in `find $accessoryName -type f -name '*template_accessory*'`; do mv $file `sed s/template_accessory/$accessoryName/g <<< $file`; done;

echo Renaming in files...
find $accessoryName -type f -exec sed -i "s/template_accessory/$accessoryName/g" {} \;

echo "Done! $accessoryName is ready. Don't forget to add it to CMakeLists.txt!"
