package utils

// Version 管理器的版本号，构建时通过 -ldflags "-X ark-server-commander/utils.Version=x.y.z" 注入；
// 本地直接运行或构建时没有传入时为 dev
var Version = "dev"

// ReleaseRepo 发布版本所在的 GitHub 仓库，用于检查新版本
const ReleaseRepo = "TBro1998/ArkServerCommander"
