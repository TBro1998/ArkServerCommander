import { NextResponse } from 'next/server';
import axios from 'axios';
import { headers } from 'next/headers';

const getApiBase = () => process.env.NEXT_PUBLIC_API_BASE;

export async function POST(request: Request, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const headersList = await headers();
  const authorization = headersList.get('authorization');

  if (!authorization) {
    return NextResponse.json({ error: '未授权' }, { status: 401 });
  }

  const config = {
    headers: { Authorization: authorization, 'Content-Type': 'application/json' },
  };

  try {
    const url = `${getApiBase()}/servers/${id}/start`;
    const response = await axios.post(url, {}, config);
    return NextResponse.json(response.data);
  } catch (error: unknown) {
    const axiosError = error as { response?: { data?: { error?: string }, status?: number } };
    return NextResponse.json({
      error: axiosError.response?.data?.error || '请求失败'
    }, { status: axiosError.response?.status || 500 });
  }
}