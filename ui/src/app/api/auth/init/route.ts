import { authProxy } from '@/lib/backend-server';

export async function POST(request: Request) {
  return authProxy(request, 'init', '初始化失败');
}
