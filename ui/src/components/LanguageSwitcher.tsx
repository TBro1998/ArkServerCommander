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
import { Check, Languages } from 'lucide-react';
import { cn } from '@/lib/utils';

const LANGUAGES: { value: Locale; label: string; short: string }[] = [
  { value: 'zh', label: '简体中文', short: '中文' },
  { value: 'en', label: 'English', short: 'EN' },
];

// 按钮上直接显示当前语言，不用点开才知道
export function LanguageSwitcher({ className }: { className?: string }) {
  // 服务端渲染使用的语言来自 cookie，切换语言会整页刷新，所以直接以它作为当前语言
  const currentLocale = useLocale() as Locale;
  const current = LANGUAGES.find((l) => l.value === currentLocale) ?? LANGUAGES[1];

  const changeLocale = (nextLocale: Locale) => {
    if (nextLocale === currentLocale) return;
    setClientLocale(nextLocale);
    window.location.reload();
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="sm" className={cn('gap-1.5', className)} aria-label={`Language: ${current.label}`}>
          <Languages className="h-4 w-4" />
          <span className="text-xs font-medium">{current.short}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {LANGUAGES.map((l) => (
          <DropdownMenuItem key={l.value} onClick={() => changeLocale(l.value)}>
            <Check className={cn('mr-2 h-4 w-4', l.value === currentLocale ? 'opacity-100' : 'opacity-0')} />
            {l.label}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
