cd $(dirname "$0")
DIR="$(pwd)"
TargetDir=$DIR/../confmanager/template
CONFIG_BRANCH="$1"
CONFIG_REPO_URL="git@codeup.aliyun.com:62b023a03e81781f3ad195c6/Server_V2_Config.git"

if [ -d "./Server_V2_Config" ];then
    rm -rf ./Server_V2_Config
fi

git clone \
  --branch "$CONFIG_BRANCH" \
  --depth 1  \
  "$CONFIG_REPO_URL"

cd Server_V2_Config || exit

RELEASE_VERSION=$(sed -n 1p RELEASE_VERSION)
echo "当前发行版本号：""$RELEASE_VERSION"

for f in ${TargetDir}/*
do
	fname=$(basename $f)
	cp $RELEASE_VERSION/$fname $TargetDir
done

echo  "配置文件夹"$TargetDir
