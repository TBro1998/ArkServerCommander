import { NextRequest, NextResponse } from 'next/server';
import { getLocale, setServerLocale, defaultLocale } from './lib/locale-server';

export async function middleware(request: NextRequest) {
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

  return NextResponse.next();
}

export const config = {
  // 匹配所有路径，但排除静态文件和API路由
  matcher: [
    '/((?!api|_next/static|_next/image|favicon.ico).*)',
  ]
};