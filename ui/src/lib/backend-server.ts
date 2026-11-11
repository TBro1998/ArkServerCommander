import { NextResponse } from 'next/server';

export const TOKEN_COOKIE = 'auth-token';
const TOKEN_MAX_AGE = 60 * 60 * 24; // 与后端 JWT 的 24 小时有效期一致

// 后端地址：优先使用运行时的 API_BASE，其次 NEXT_PUBLIC_API_BASE（开发环境），最后回退到本机回环地址
export const getApiBase = () =>
  process.env.API_BASE || process.env.NEXT_PUBLIC_API_BASE || 'http://127.0.0.1:8080/api';

// HttpOnly：浏览器脚本读不到令牌；仅在 HTTPS 部署时设置 COOKIE_SECURE=true
export const cookieOptions = () => ({
  httpOnly: true,
  sameSite: 'lax' as const,
  secure: process.env.COOKIE_SECURE === 'true',
  path: '/',
});

export function setTokenCookie(response: NextResponse, token: string) {
  response.cookies.set(TOKEN_COOKIE, token, { ...cookieOptions(), maxAge: TOKEN_MAX_AGE });
}

export function clearTokenCookie(response: NextResponse) {
  response.cookies.set(TOKEN_COOKIE, '', { ...cookieOptions(), maxAge: 0 });
}

// 登录/初始化：转发到后端，把返回的令牌写入 HttpOnly cookie，响应体里不再包含令牌
export async function authProxy(request: Request, endpoint: 'login' | 'init', fallbackError: string) {
  let body: unknown;
  try {
    body = await request.json();
  } catch {
    return NextResponse.json({ error: '请求参数错误' }, { status: 400 });
  }

  try {
    // 转发客户端地址，后端按来源 IP 做登录失败限流
    const headers: Record<string, string> = { 'Content-Type': 'application/json' };
    const forwardedFor = request.headers.get('x-forwarded-for');
    if (forwardedFor) headers['X-Forwarded-For'] = forwardedFor;
    const res = await fetch(`${getApiBase()}/auth/${endpoint}`, {
      method: 'POST',
      headers,
      body: JSON.stringify(body),
      cache: 'no-store',
    });
    const data = await res.json().catch(() => ({}));
    if (!res.ok || typeof data.token !== 'string') {
      return NextResponse.json({ error: data?.error || fallbackError, code: data?.code }, { status: res.ok ? 502 : res.status });
    }
    const { token, ...rest } = data;
    const response = NextResponse.json(rest);
    setTokenCookie(response, token);
    return response;
  } catch {
    return NextResponse.json({ error: '后端服务不可用' }, { status: 502 });
  }
}
