import { NextResponse } from 'next/server';
import axios from 'axios';
import { headers } from 'next/headers';

const getApiBase = () => process.env.NEXT_PUBLIC_API_BASE;

export async function POST(request: Request) {
  const headersList = await headers();
  const authorization = headersList.get('authorization');

  if (!authorization) {
    return NextResponse.json({ error: '未授权' }, { status: 401 });
  }

  try {
    const body = await request.json();
    const { image_name } = body;

    if (!image_name) {
      return NextResponse.json({ error: '镜像名称不能为空' }, { status: 400 });
    }

    const config = {
      headers: { Authorization: authorization, 'Content-Type': 'application/json' },
    };

    const url = `${getApiBase()}/servers/images/pull`;
    const response = await axios.post(url, { image_name }, config);
    return NextResponse.json(response.data);
  } catch (error: any) {
    return NextResponse.json({ error: error.response?.data?.error || '拉取镜像失败' }, { status: error.response?.status || 500 });
  }
}