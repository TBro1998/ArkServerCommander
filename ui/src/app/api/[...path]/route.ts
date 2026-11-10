import { NextRequest, NextResponse } from 'next/server';
import { clearTokenCookie, getApiBase, TOKEN_COOKIE } from '@/lib/backend-server';

// 统一代理：从 HttpOnly cookie 取出令牌，转成 Bearer 头转发给后端
async function proxy(request: NextRequest, { params }: { params: Promise<{ path: string[] }> }) {
  const { path } = await params;
  if (path.some((segment) => segment === '.' || segment === '..')) {
    return NextResponse.json({ error: '请求路径无效' }, { status: 400 });
  }

  const token = request.cookies.get(TOKEN_COOKIE)?.value;
  const headers: Record<string, string> = {};
  const contentType = request.headers.get('content-type');
  if (contentType) headers['Content-Type'] = contentType;
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const hasBody = request.method !== 'GET' && request.method !== 'HEAD';
  const url = `${getApiBase()}/${path.map(encodeURIComponent).join('/')}${request.nextUrl.search}`;

  try {
    const res = await fetch(url, {
      method: request.method,
      headers,
      body: hasBody ? await request.text() : undefined,
      cache: 'no-store',
    });
    const text = await res.text();
    const body = [204, 205, 304].includes(res.status) ? null : text;
    const response = new NextResponse(body, {
      status: res.status,
      headers: { 'Content-Type': res.headers.get('content-type') ?? 'application/json' },
    });
    // 令牌过期或无效时清掉 cookie，proxy 会把用户带回登录页
    if (res.status === 401 && token) clearTokenCookie(response);
    return response;
  } catch {
    return NextResponse.json({ error: '后端服务不可用' }, { status: 502 });
  }
}

export const GET = proxy;
export const POST = proxy;
export const PUT = proxy;
export const DELETE = proxy;
