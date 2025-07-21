"use client";

import { useLocale } from 'next-intl';
import { setClientLocale, getClientLocale, type Locale } from '@/lib/locale';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Button } from '@/components/ui/button';
import { Languages } from 'lucide-react';
import { useState, useEffect } from 'react';

export function LanguageSwitcher() {
  const serverLocale = useLocale();
  const [currentLocale, setCurrentLocale] = useState<Locale>(serverLocale as Locale);

  useEffect(() => {
    // 同步客户端和服务端的语言设置
    const clientLocale = getClientLocale();
    if (clientLocale !== serverLocale) {
      setCurrentLocale(clientLocale);
    }
  }, [serverLocale]);

  const changeLocale = (nextLocale: Locale) => {
    // 设置cookie
    setClientLocale(nextLocale);
    setCurrentLocale(nextLocale);

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