RUNENV=dev
if [[ "$1" == "pre" ]];then
	RUNENV=pre
elif  [[ "$1" == "prod" ]];then
	RUNENV=prod
else
	RUNENV=$1
fi
basepath=$(cd `dirname $0`; pwd)
DockerfilePath=$basepath/Dockerfile
CodePath=`dirname $basepath`
CodeConfPath=$CodePath/conf/$RUNENV
DockerImageName=matching-story-robot-service:latest
ContainerName=$RUNENV"-matching-story-robot-service"
ConfmanagerPath=$basepath/Server_V2_Config
LogPath=$CodePath/var/logs
ConfManagerBranchName=dev

if [[ "$RUNENV" == "pre" || "$RUNENV" == "prod" || "$RUNENV" == "dev" || "$RUNENV" == "dev1" ]];then
	LogPath="/home/ec2-user/run/matching-story-robot-service/logs"
fi

if [[ "$RUNENV" == "prod" ]];then
	ConfManagerBranchName=master
fi

if [ ! -d "$CodeConfPath" ]; then
    echo "not found conf path $CodeConfPath"
	exit
fi

echo "DockerfilePath: "$DockerfilePath
echo "CodePath: "$CodePath
echo "CodeConfPath: "$CodeConfPath
echo "DockerImageName: "$DockerImageName
echo "ContainerName: "$ContainerName
echo "ConfmanagerPath: "$ConfmanagerPath
echo "RunEnv: "$RUNENV
echo "LogPath: "$LogPath
echo "ConfManagerBranchName: "$ConfManagerBranchName

if [ -d "$ConfmanagerPath" ]; then
	rm -rf $ConfmanagerPath
fi

echo "更新confmanager配置"
git clone --branch $ConfmanagerPath --depth 1 git@codeup.aliyun.com:62b023a03e81781f3ad195c6/Server_V2_Config.git $ConfmanagerPath
cd $ConfmanagerPath || exit
RELEASE_VERSION=$(sed -n 1p $ConfmanagerPath/RELEASE_VERSION)
echo "当前发行版本号：""$RELEASE_VERSION"
cp $ConfmanagerPath/$RELEASE_VERSION/* $CodePath/confmanager/template

echo "docker build ....."
docker build -f $DockerfilePath -t $DockerImageName $CodePath 
echo "docker stop $ContainerName"
docker stop $ContainerName
echo "docker rm $ContainerName"
docker rm $ContainerName

echo "docker run ......"
#docker run 
docker run -d --restart=always \
--name $ContainerName \
--entrypoint /app/matching-story-robot-service -env=$RUNENV \
-p 18088:8088 \
-v $LogPath:/app/var/logs \
$DockerImageName \

echo "运行日志路径:$LogPath"
