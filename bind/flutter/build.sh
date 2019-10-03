ORIG=`pwd`
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

echo "Building Flutter SDK..."

cd $DIR
set -x
rm -rf build
mkdir build
cd build
mkdir replicant-flutter-sdk
ls ../ | grep -v build | grep -v sample | xargs -I{} cp -R ../{} replicant-flutter-sdk/{}
rm -rf replicant-flutter-sdk/ios/Repm.framework
rm replicant-flutter-sdk/android/repm.aar
cp -R ../../../repm/build/Repm.framework replicant-flutter-sdk/ios/
cp ../../../repm/build/repm.aar replicant-flutter-sdk/android/
tar -czvf replicant-flutter-sdk.tar.gz replicant-flutter-sdk
cd $ORIG

