import { getRequestConfig } from 'next-intl/server';
import { notFound } from 'next/navigation';

export const locales = ['en', 'zh'];

export default getRequestConfig(async ({ locale }) => {
  if (!locales.includes(locale)) {
    notFound();
  }
  return {
    messages: (await import(`../messages/${locale}.ts`)).default
  }
});