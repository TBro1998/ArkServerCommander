import { NextResponse } from 'next/server';
import axios from 'axios';
import { headers } from 'next/headers';

const getApiBase = () => process.env.NEXT_PUBLIC_API_BASE;

async function proxyRequest(request: Request, method: 'GET' | 'PUT' | 'DELETE', id: string) {
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
    const url = `${getApiBase()}/servers/${id}`;

    if (method === 'GET') {
      response = await axios.get(url, config);
    } else if (method === 'PUT') {
      const body = await request.json();
      response = await axios.put(url, body, config);
    } else { // DELETE
      response = await axios.delete(url, config);
    }

    return NextResponse.json(response.data);
  } catch (error: any) {
    return NextResponse.json({ error: error.response?.data?.error || '请求失败' }, { status: error.response?.status || 500 });
  }
}

export async function GET(request: Request, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return proxyRequest(request, 'GET', id);
}

export async function PUT(request: Request, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return proxyRequest(request, 'PUT', id);
}

export async function DELETE(request: Request, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return proxyRequest(request, 'DELETE', id);
}