import { NextRequest, NextResponse } from 'next/server';
import { clearTokenCookie, getApiBase, TOKEN_COOKIE } from '@/lib/backend-server';

// 每次请求都要转发，不能缓存
export const dynamic = 'force-dynamic';

// 原样转发给浏览器的响应头：下载备份需要 Content-Disposition，日志流需要关闭缓存
const PASS_RESPONSE_HEADERS = ['content-type', 'content-disposition', 'content-length', 'cache-control', 'retry-after'];

// 统一代理：从 HttpOnly cookie 取出令牌，转成 Bearer 头转发给后端。
// 请求体和响应体都以流的方式转发，支持上传备份、下载备份和 SSE 实时日志
async function proxy(request: NextRequest, { params }: { params: Promise<{ path: string[] }> }) {
  const { path } = await params;
  if (path.some((segment) => segment === '.' || segment === '..')) {
    return NextResponse.json({ error: '请求路径无效' }, { status: 400 });
  }
  // MCP 走单独的 /mcp 路由（用 API 令牌认证）；这里会带上登录 cookie，令牌校验失败时还会把 cookie 清掉
  if (path[0] === 'mcp') {
    return NextResponse.json({ error: 'MCP 地址是 /mcp', code: 'not_found' }, { status: 404 });
  }

  const token = request.cookies.get(TOKEN_COOKIE)?.value;
  const headers: Record<string, string> = {};
  const contentType = request.headers.get('content-type');
  if (contentType) headers['Content-Type'] = contentType;
  const contentLength = request.headers.get('content-length');
  if (contentLength) headers['Content-Length'] = contentLength;
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const forwardedFor = request.headers.get('x-forwarded-for');
  if (forwardedFor) headers['X-Forwarded-For'] = forwardedFor;

  const hasBody = request.method !== 'GET' && request.method !== 'HEAD';
  const url = `${getApiBase()}/${path.map(encodeURIComponent).join('/')}${request.nextUrl.search}`;

  try {
    const res = await fetch(url, {
      method: request.method,
      headers,
      body: hasBody ? request.body : undefined,
      // 以流的方式发送请求体时 Node 要求 duplex: 'half'
      ...(hasBody ? { duplex: 'half' } : {}),
      cache: 'no-store',
      // 浏览器断开（例如关闭日志窗口）时一并取消对后端的请求
      signal: request.signal,
    } as RequestInit);

    const responseHeaders = new Headers();
    for (const name of PASS_RESPONSE_HEADERS) {
      const value = res.headers.get(name);
      if (value) responseHeaders.set(name, value);
    }
    if (!responseHeaders.has('content-type')) responseHeaders.set('content-type', 'application/json');
    if (responseHeaders.get('content-type')?.startsWith('text/event-stream')) {
      responseHeaders.set('cache-control', 'no-cache, no-transform');
      responseHeaders.set('x-accel-buffering', 'no');
    }

    const body = [204, 205, 304].includes(res.status) ? null : res.body;
    const response = new NextResponse(body, { status: res.status, headers: responseHeaders });
    // 令牌过期或无效时清掉 cookie，proxy 会把用户带回登录页
    if (res.status === 401 && token) clearTokenCookie(response);
    return response;
  } catch {
    if (request.signal.aborted) return new NextResponse(null, { status: 499 });
    return NextResponse.json({ error: '后端服务不可用', code: 'backend_unavailable' }, { status: 502 });
  }
}

export const GET = proxy;
export const POST = proxy;
export const PUT = proxy;
export const DELETE = proxy;
