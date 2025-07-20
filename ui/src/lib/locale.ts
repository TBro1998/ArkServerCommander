import Cookies from 'js-cookie';

export const locales = ['en', 'zh'] as const;
export type Locale = typeof locales[number];
export const defaultLocale: Locale = 'en';

const LOCALE_COOKIE_NAME = 'NEXT_LOCALE';

// 客户端获取语言
export function getClientLocale(): Locale {
  if (typeof window === 'undefined') {
    return defaultLocale;
  }
  
  const locale = Cookies.get(LOCALE_COOKIE_NAME);
  
  if (locale && locales.includes(locale as Locale)) {
    return locale as Locale;
  }
  
  return defaultLocale;
}

// 客户端设置语言
export function setClientLocale(locale: Locale) {
  if (typeof window === 'undefined') {
    return;
  }
  
  Cookies.set(LOCALE_COOKIE_NAME, locale, {
    expires: 365, // 1年
    path: '/',
    sameSite: 'lax'
  });
}