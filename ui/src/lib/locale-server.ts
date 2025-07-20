import { cookies } from 'next/headers';

export const locales = ['en', 'zh'] as const;
export type Locale = typeof locales[number];
export const defaultLocale: Locale = 'en';

const LOCALE_COOKIE_NAME = 'NEXT_LOCALE';

// 服务端获取语言
export async function getLocale(): Promise<Locale> {
  const cookieStore = await cookies();
  const locale = cookieStore.get(LOCALE_COOKIE_NAME)?.value;
  
  if (locale && locales.includes(locale as Locale)) {
    return locale as Locale;
  }
  
  return defaultLocale;
}

// 服务端设置语言
export async function setServerLocale(locale: Locale) {
  const cookieStore = await cookies();
  cookieStore.set(LOCALE_COOKIE_NAME, locale, {
    maxAge: 365 * 24 * 60 * 60, // 1年
    path: '/',
    sameSite: 'lax'
  });
}