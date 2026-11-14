import { NextRequest, NextResponse } from 'next/server';
import { getApiBase } from '@/lib/backend-server';

// 每次请求都要转发，不能缓存
export const dynamic = 'force-dynamic';

// MCP 入口：Agent 访问 http(s)://<主机>/mcp，原样转发给后端的 /api/mcp。
// 认证用 Agent 自带的 Authorization: Bearer <API 令牌>，不读取登录 cookie。
// 和通配代理分开，是因为 MCP 需要透传 Authorization、Accept 和 Mcp-* 头，且令牌无效时不能清掉浏览器的登录 cookie
const PASS_REQUEST_HEADERS = [
  'authorization',
  'content-type',
  'content-length',
  'accept',
  'mcp-session-id',
  'mcp-protocol-version',
  'last-event-id',
  'x-forwarded-for',
];
const PASS_RESPONSE_HEADERS = [
  'content-type',
  'content-length',
  'cache-control',
  'mcp-session-id',
  'mcp-protocol-version',
  'www-authenticate',
  'allow',
];

async function forward(request: NextRequest) {
  const headers: Record<string, string> = {};
  for (const name of PASS_REQUEST_HEADERS) {
    const value = request.headers.get(name);
    if (value) headers[name] = value;
  }

  const hasBody = request.method !== 'GET' && request.method !== 'HEAD';
  try {
    const res = await fetch(`${getApiBase()}/mcp`, {
      method: request.method,
      headers,
      body: hasBody ? request.body : undefined,
      ...(hasBody ? { duplex: 'half' } : {}),
      cache: 'no-store',
      signal: request.signal,
    } as RequestInit);

    const responseHeaders = new Headers();
    for (const name of PASS_RESPONSE_HEADERS) {
      const value = res.headers.get(name);
      if (value) responseHeaders.set(name, value);
    }
    if (responseHeaders.get('content-type')?.startsWith('text/event-stream')) {
      responseHeaders.set('cache-control', 'no-cache, no-transform');
      responseHeaders.set('x-accel-buffering', 'no');
    }

    const body = [204, 205, 304].includes(res.status) ? null : res.body;
    return new NextResponse(body, { status: res.status, headers: responseHeaders });
  } catch {
    if (request.signal.aborted) return new NextResponse(null, { status: 499 });
    return NextResponse.json(
      { jsonrpc: '2.0', id: null, error: { code: -32603, message: 'backend unavailable' } },
      { status: 502 },
    );
  }
}

export const GET = forward;
export const POST = forward;
export const DELETE = forward;
