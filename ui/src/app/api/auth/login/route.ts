import { NextResponse } from 'next/server';
import axios from 'axios';

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const response = await axios.post(`${process.env.NEXT_PUBLIC_API_BASE}/auth/login`, body);
    return NextResponse.json(response.data);
  } catch (error: unknown) {
    const axiosError = error as { response?: { data?: { error?: string }, status?: number } };
    return NextResponse.json({
      error: axiosError.response?.data?.error || '登录失败'
    }, { status: axiosError.response?.status || 500 });
  }
}