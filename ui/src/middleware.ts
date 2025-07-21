import { NextRequest, NextResponse } from 'next/server';
import { getLocale, setServerLocale, defaultLocale } from './lib/locale-server';

export async function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  
  // 检查是否有语言切换请求
  const searchParams = request.nextUrl.searchParams;
  const localeParam = searchParams.get('locale');

  if (localeParam && ['en', 'zh'].includes(localeParam)) {
    // 设置新的语言cookie
    const response = NextResponse.redirect(new URL(request.nextUrl.pathname, request.url));
    response.cookies.set('NEXT_LOCALE', localeParam, {
      maxAge: 365 * 24 * 60 * 60, // 1年
      path: '/',
      sameSite: 'lax'
    });
    return response;
  }

  // 确保有语言cookie
  const currentLocale = request.cookies.get('NEXT_LOCALE')?.value;
  if (!currentLocale || !['en', 'zh'].includes(currentLocale)) {
    const response = NextResponse.next();
    response.cookies.set('NEXT_LOCALE', defaultLocale, {
      maxAge: 365 * 24 * 60 * 60, // 1年
      path: '/',
      sameSite: 'lax'
    });
    return response;
  }

  // 认证检查
  const token = request.cookies.get('auth-token')?.value;
  const isAuthPage = pathname.startsWith('/login') || pathname.startsWith('/init');
  const isProtectedPage = pathname.startsWith('/home') || pathname.startsWith('/servers');
  const isRootPage = pathname === '/';

  // 如果访问受保护页面但没有token，重定向到登录页
  if (isProtectedPage && !token) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  // 如果访问根页面且有token，重定向到home页面
  if (isRootPage && token) {
    return NextResponse.redirect(new URL('/home', request.url));
  }

  // 如果访问根页面且没有token，重定向到登录页
  if (isRootPage && !token) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  // 如果已登录但访问认证页面，重定向到home页面
  if (isAuthPage && token) {
    return NextResponse.redirect(new URL('/home', request.url));
  }

  return NextResponse.next();
}

export const config = {
  // 匹配所有路径，但排除静态文件和API路由
  matcher: [
    '/((?!api|_next/static|_next/image|favicon.ico).*)',
  ]
};