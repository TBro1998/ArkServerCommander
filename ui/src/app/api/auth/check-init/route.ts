import { NextResponse } from 'next/server';
import axios from 'axios';

export async function GET() {
  try {
    const response = await axios.get(`${process.env.NEXT_PUBLIC_API_BASE}/auth/check-init`);
    return NextResponse.json(response.data);
  } catch (error: unknown) {
    const axiosError = error as { response?: { data?: { error?: string }, status?: number } };
    return NextResponse.json({
      error: axiosError.response?.data?.error || '检查初始化状态失败'
    }, { status: axiosError.response?.status || 500 });
  }
}