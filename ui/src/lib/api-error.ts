import axios from 'axios';
import { useCallback } from 'react';
import { useLocale, useTranslations } from 'next-intl';

interface ErrorBody {
  error?: unknown;
  code?: unknown;
}

function errorBody(error: unknown): ErrorBody | null {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data;
    if (data && typeof data === 'object') return data as ErrorBody;
  }
  return null;
}

// 优先使用后端返回的错误信息，拿不到时回退到调用方提供的默认文案
export function getErrorMessage(error: unknown, fallback: string): string {
  const body = errorBody(error);
  if (body && typeof body.error === 'string') return body.error;
  return fallback;
}

// 按界面语言选择错误文案：后端的 error 字段是中文且带有具体细节（例如哪个端口冲突），
// 中文界面直接显示；其他语言按 code 显示翻译，没有对应翻译时再显示后端原文
export function useApiError() {
  const locale = useLocale();
  const t = useTranslations('apiErrors');
  return useCallback(
    (error: unknown, fallback: string): string => {
      const body = errorBody(error);
      if (!body) return fallback;
      const message = typeof body.error === 'string' ? body.error : '';
      const code = typeof body.code === 'string' ? body.code : '';
      if (locale !== 'zh' && code && t.has(code)) return t(code);
      return message || fallback;
    },
    [locale, t],
  );
}
