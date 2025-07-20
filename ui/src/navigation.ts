// 使用标准的Next.js导航，不再依赖next-intl的路径导航
export { redirect } from 'next/navigation';
export { usePathname, useRouter } from 'next/navigation';

// 重新导出Next.js的Link组件
import NextLink from 'next/link';
export const Link = NextLink;

export const locales = ['en', 'zh'] as const;