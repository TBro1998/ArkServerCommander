import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
  const token = request.cookies.get('auth-token');
  const { pathname } = request.nextUrl;

  // 定义公共路由和受保护路由
  const publicPaths = ['/login', '/init'];
  const isPublicPath = publicPaths.some(path => pathname.startsWith(path));

  if (!token && !isPublicPath) {
    // 如果没有token并且访问的不是公共页面，重定向到登录页
    return NextResponse.redirect(new URL('/login', request.url));
  }

  if (token && isPublicPath) {
    // 如果有token并且访问的是公共页面，重定向到服务器列表页
    return NextResponse.redirect(new URL('/servers', request.url));
  }

  return NextResponse.next();
}

export const config = {
  // 匹配器，用于指定中间件应该在哪些路径上运行
  // 排除 _next, /api, /static, /favicon.ico 等
  matcher: [
    '/((?!api|_next/static|_next/image|favicon.ico).*)',
  ],
};