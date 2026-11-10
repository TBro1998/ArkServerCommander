import { NextResponse } from 'next/server';
import { clearTokenCookie } from '@/lib/backend-server';

export async function POST() {
  const response = NextResponse.json({ message: '已退出登录' });
  clearTokenCookie(response);
  return response;
}
