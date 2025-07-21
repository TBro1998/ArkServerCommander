import { NextResponse } from 'next/server';
import axios from 'axios';
import { headers } from 'next/headers';

const getApiBase = () => process.env.NEXT_PUBLIC_API_BASE;

async function proxyRequest(request: Request, method: 'GET' | 'POST') {
  const headersList = await headers();
  const authorization = headersList.get('authorization');

  if (!authorization) {
    return NextResponse.json({ error: '未授权' }, { status: 401 });
  }

  const config = {
    headers: { Authorization: authorization, 'Content-Type': 'application/json' },
  };

  try {
    let response;
    if (method === 'GET') {
      response = await axios.get(`${getApiBase()}/servers`, config);
    } else {
      const body = await request.json();
      response = await axios.post(`${getApiBase()}/servers`, body, config);
    }
    return NextResponse.json(response.data);
  } catch (error: unknown) {
    const axiosError = error as { response?: { data?: { error?: string }, status?: number } };
    return NextResponse.json({
      error: axiosError.response?.data?.error || '请求失败'
    }, { status: axiosError.response?.status || 500 });
  }
}

export async function GET(request: Request) {
  return proxyRequest(request, 'GET');
}

export async function POST(request: Request) {
  return proxyRequest(request, 'POST');
}