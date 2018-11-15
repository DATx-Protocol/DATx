#! /bin/bash

rm -rf package-lock.json
rm -rf ./node_modules/

echo "start install node packages....."
npm install  --unsafe-perm

echo "start link"
npm link