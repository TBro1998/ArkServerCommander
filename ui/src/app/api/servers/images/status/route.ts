import { NextResponse } from 'next/server';
import axios from 'axios';
import { headers } from 'next/headers';

const getApiBase = () => process.env.NEXT_PUBLIC_API_BASE;

export async function GET() {
  const headersList = await headers();
  const authorization = headersList.get('authorization');

  if (!authorization) {
    return NextResponse.json({ error: '未授权' }, { status: 401 });
  }

  const config = {
    headers: { Authorization: authorization, 'Content-Type': 'application/json' },
  };

  try {
    const url = `${getApiBase()}/servers/images/status`;
    const response = await axios.get(url, config);
    return NextResponse.json(response.data);
  } catch (error: any) {
    return NextResponse.json({ error: error.response?.data?.error || '请求失败' }, { status: error.response?.status || 500 });
  }
}