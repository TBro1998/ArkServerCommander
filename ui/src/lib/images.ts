// 本系统管理的镜像，名称与后端 utils/asm_info.go 一致
export const ARK_LINUX_IMAGE = 'tbro98/ase-server:linux';
export const ARK_PROTON_IMAGE = 'tbro98/ase-server:proton';
export const ALPINE_IMAGE = 'alpine:latest';

// 镜像的显示名称对应的翻译键（servers.dockerImages 下），未知镜像返回 null
export function imageLabelKey(imageName: string): 'arkServerLinux' | 'arkServerProton' | 'alpineSystem' | null {
  switch (imageName) {
    case ARK_LINUX_IMAGE:
      return 'arkServerLinux';
    case ARK_PROTON_IMAGE:
      return 'arkServerProton';
    case ALPINE_IMAGE:
      return 'alpineSystem';
    default:
      return null;
  }
}
