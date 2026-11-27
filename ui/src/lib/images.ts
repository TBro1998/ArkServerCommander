// 本系统管理的镜像，名称与后端 utils/asm_info.go 一致
// ASE：原生 Linux 版（不用 ArkApi）和 Proton 版（启用 ArkApi）
export const ARK_LINUX_IMAGE = 'tbro98/ase-server:linux';
export const ARK_PROTON_IMAGE = 'tbro98/ase-server:proton';
// ASA 没有 Linux 版服务端，只有 Proton 版（是否加载 AsaApi 由容器的 ASAAPI 决定）
export const ASA_PROTON_IMAGE = 'tbro98/asa-server:proton';
// P2P 联机房的 EasyTier 节点
export const EASYTIER_IMAGE = 'easytier/easytier:v2.6.4';

export type ImageLabelKey = 'arkServerLinux' | 'arkServerProton' | 'asaServerProton' | 'easyTier';

// 镜像的显示名称对应的翻译键（servers.dockerImages 下），未知镜像返回 null
export function imageLabelKey(imageName: string): ImageLabelKey | null {
  switch (imageName) {
    case ARK_LINUX_IMAGE:
      return 'arkServerLinux';
    case ARK_PROTON_IMAGE:
      return 'arkServerProton';
    case ASA_PROTON_IMAGE:
      return 'asaServerProton';
    case EASYTIER_IMAGE:
      return 'easyTier';
    default:
      return null;
  }
}
