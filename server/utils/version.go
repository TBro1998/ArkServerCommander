package utils

// Version 管理器的版本号，构建时通过 -ldflags "-X ark-server-commander/utils.Version=x.y.z" 注入；
// 本地直接运行或构建时没有传入时为 dev
var Version = "dev"

// ReleaseRepo 发布版本所在的 GitHub 仓库，用于检查新版本
const ReleaseRepo = "TBro1998/ArkServerCommander"

// CommanderImageRepo 管理端发布的镜像（不含镜像仓库地址和标签）。只有用这个镜像的 latest 标签运行时才能在界面上一键更新
const CommanderImageRepo = "tbro98/arkservercommander"
