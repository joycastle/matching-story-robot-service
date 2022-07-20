cd $(dirname "$0")
CONFIG_BRANCH="$1"
CONFIG_REPO_URL="git@codeup.aliyun.com:62b023a03e81781f3ad195c6/Server_V2_Config.git"
DIR="$(pwd)"
echo "当前地址""$DIR"
echo "从""$CONFIG_REPO_URL""的""$CONFIG_BRANCH""分支拉取文件"

git clone \
  --branch "$CONFIG_BRANCH" \
  --depth 1  \
  "$CONFIG_REPO_URL"

cd Server_V2_Config || exit

RELEASE_VERSION=$(sed -n 1p RELEASE_VERSION)
echo "当前发行版本号：""$RELEASE_VERSION"

cd "$RELEASE_VERSION" && cp ./* "$DIR"/../confmanager/template

echo "$DIR"/../confmanager/template
