CONFIG_BRANCH="$1"
#CONFIG_REPO_URL="git@codeup.teambition.com:joycastle/mm/Server_V2_Config.git"
CONFIG_REPO_URL="git@codeup.aliyun.com:62b023a03e81781f3ad195c6/Server_V2_Config.git"

https://codeup.teambition.com/joycastle/mm/Server_V2_Config/blob/dev/README.md
DIR="$(pwd)"

echo "当前地址""$DIR"

echo "从""$CONFIG_REPO_URL""的""$CONFIG_BRANCH""分支拉取文件"

git clone \
  --branch "$CONFIG_BRANCH" \
  --depth 1  \
  "$CONFIG_REPO_URL"

cd Server_V2_Config || exit

RELEASE_VERSION=$(sed -n 1p RELEASE_VERSION)

echo "========================================================================================================="
echo "当前发行版本号：""$RELEASE_VERSION"
echo "满级玩法配置如下："

tail -n +3 "$RELEASE_VERSION"/maxleveltournament.csv

echo "========================================================================================================="
echo "公会活动配置如下："

tail -n 2 "$RELEASE_VERSION"/teamchesttime.csv

echo "========================================================================================================="

cd "$RELEASE_VERSION" && cp ./* "$DIR"/confmanager/template
