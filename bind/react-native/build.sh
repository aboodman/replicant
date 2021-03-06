#!/bin/sh
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
ORIG=`pwd`
cd $DIR
echo "Building React Native bindings..."
set -x
rm -rf build
mkdir build
mkdir build/replicant-react-native
ls | grep -v build | xargs -I{} cp -R {} build/replicant-react-native/
rm -rf build/replicant-react-native/ios/Frameworks
rm build/replicant-react-native/android/repm.aar
mkdir build/replicant-react-native/ios/Frameworks
cp -R $DIR/../../repm/build/Repm.framework $DIR/build/replicant-react-native/ios/Frameworks/
cp $DIR/../../repm/build/repm.aar $DIR/build/replicant-react-native/android/
cd build
tar -czvf replicant-react-native.tar.gz replicant-react-native
cd $ORIG

