import { getRequestConfig } from 'next-intl/server';
import { getLocale, locales } from './lib/locale-server';

// 导出locales以保持兼容性
export { locales };

export default getRequestConfig(async () => {
  // 从cookie获取语言设置
  const locale = await getLocale();

  return {
    locale,
    messages: (await import(`../messages/${locale}.ts`)).default
  };
});