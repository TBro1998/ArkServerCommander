import { NextRequest, NextResponse } from 'next/server';
import { defaultLocale } from './lib/locale-server';

const LOCALES = ['en', 'zh'];
const LOCALE_COOKIE = {
  maxAge: 365 * 24 * 60 * 60, // 1年
  path: '/',
  sameSite: 'lax' as const,
};

// 需要登录的页面前缀
const PROTECTED_PREFIXES = ['/home', '/servers', '/settings'];

export async function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // 语言切换请求：写入 cookie 后去掉参数重新访问
  const localeParam = request.nextUrl.searchParams.get('locale');
  if (localeParam && LOCALES.includes(localeParam)) {
    const response = NextResponse.redirect(new URL(pathname, request.url));
    response.cookies.set('NEXT_LOCALE', localeParam, LOCALE_COOKIE);
    return response;
  }

  // 缺少语言 cookie 时补上，但仍然继续做登录检查
  const currentLocale = request.cookies.get('NEXT_LOCALE')?.value;
  const needLocale = !currentLocale || !LOCALES.includes(currentLocale);
  const withLocale = (response: NextResponse) => {
    if (needLocale) response.cookies.set('NEXT_LOCALE', defaultLocale, LOCALE_COOKIE);
    return response;
  };

  // 认证检查：这里只看 cookie 是否存在，令牌真伪由后端校验
  const token = request.cookies.get('auth-token')?.value;
  const isAuthPage = pathname.startsWith('/login') || pathname.startsWith('/init');
  const isProtectedPage = PROTECTED_PREFIXES.some((p) => pathname === p || pathname.startsWith(`${p}/`));

  if (pathname === '/') {
    return withLocale(NextResponse.redirect(new URL(token ? '/home' : '/login', request.url)));
  }
  if (isProtectedPage && !token) {
    return withLocale(NextResponse.redirect(new URL('/login', request.url)));
  }
  if (isAuthPage && token) {
    return withLocale(NextResponse.redirect(new URL('/home', request.url)));
  }

  return withLocale(NextResponse.next());
}

export const config = {
  // 匹配所有路径，但排除静态文件、API 路由和 MCP 入口
  matcher: [
    '/((?!api|mcp|_next/static|_next/image|favicon.ico).*)',
  ]
};
