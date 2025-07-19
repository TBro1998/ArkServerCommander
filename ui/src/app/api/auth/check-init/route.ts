import { NextResponse } from 'next/server';
import axios from 'axios';

export async function GET() {
  try {
    const response = await axios.get(`${process.env.NEXT_PUBLIC_API_BASE}/auth/check-init`);
    return NextResponse.json(response.data);
  } catch (error: any) {
    return NextResponse.json({ error: error.response?.data?.error || '检查初始化状态失败' }, { status: error.response?.status || 500 });
  }
}