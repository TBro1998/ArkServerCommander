// Package assets 随管理端一起发布的文件
package assets

import _ "embed"

// BuiltinArkApiVersion 内置的 ArkApi 版本。GitHub（Michidu/ARK-Server-API）上的发布只到 3.54，
// 3.56 来自 gameservershub，原先打包在服务器镜像里，现在随管理端发布，离线也能安装
const BuiltinArkApiVersion = "3.56"

// BuiltinArkApiZip 内置 ArkApi 的压缩包，根目录下是 version.dll、msdia140.dll、version.pdb、config.json
//
//go:embed arkapi/ArkApi_3.56.zip
var BuiltinArkApiZip []byte
