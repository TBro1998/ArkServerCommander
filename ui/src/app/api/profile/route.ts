import { NextResponse } from 'next/server';
import axios from 'axios';
import { headers } from 'next/headers';

export async function GET() {
  try {
    const headersList = await headers();
    const authorization = headersList.get('authorization');

    if (!authorization) {
      return NextResponse.json({ error: '未授权' }, { status: 401 });
    }

    const response = await axios.get(`${process.env.NEXT_PUBLIC_API_BASE}/profile`, {
      headers: {
        Authorization: authorization,
      },
    });
    return NextResponse.json(response.data);
  } catch (error: unknown) {
    const axiosError = error as { response?: { data?: { error?: string }, status?: number } };
    return NextResponse.json({
      error: axiosError.response?.data?.error || '获取用户信息失败'
    }, { status: axiosError.response?.status || 500 });
  }
}