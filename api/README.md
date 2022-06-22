# proto

## Usage

你可以使用`git submodule`在主项目下，依赖此repo。

Tip: 在Sourcetree的`子模块`下可以对该repo进行管理。

```bash
cd app
# 主项目
git remote add origin git@gitlab.com:joycastle/app.git

git submodule add https://gitlab.com/joycastle/backend/proto dest_folder
cd dest_folder
git remote get-url origin

```

## NOTICE

由于多proto处同一package下，需确保message name不冲突。