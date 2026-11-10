"use client";

import { useLocale } from 'next-intl';
import { setClientLocale, type Locale } from '@/lib/locale';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Button } from '@/components/ui/button';
import { Languages } from 'lucide-react';

export function LanguageSwitcher() {
  // 服务端渲染使用的语言来自 cookie，切换语言会整页刷新，所以直接以它作为当前语言
  const currentLocale = useLocale() as Locale;

  const changeLocale = (nextLocale: Locale) => {
    // 设置cookie
    setClientLocale(nextLocale);

    // 刷新页面以应用新语言
    window.location.reload();
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon">
          <Languages className="h-5 w-5" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={() => changeLocale('en')} disabled={currentLocale === 'en'}>
          English
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => changeLocale('zh')} disabled={currentLocale === 'zh'}>
          简体中文
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
