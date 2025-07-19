import { NextResponse } from 'next/server';
import axios from 'axios';

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const response = await axios.post(`${process.env.NEXT_PUBLIC_API_BASE}/auth/init`, body);
    return NextResponse.json(response.data);
  } catch (error: any) {
    return NextResponse.json({ error: error.response?.data?.error || '初始化失败' }, { status: error.response?.status || 500 });
  }
}