import { authProxy } from '@/lib/backend-server';

export async function POST(request: Request) {
  return authProxy(request, 'login', '登录失败');
}
