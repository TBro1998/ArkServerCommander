"use client";

import { useLocale } from 'next-intl';
import { useRouter, usePathname } from '@/navigation';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Button } from '@/components/ui/button';
import { Languages } from 'lucide-react';

export function LanguageSwitcher() {
  const router = useRouter();
  const pathname = usePathname();
  const locale = useLocale();

  const changeLocale = (nextLocale: string) => {
    router.replace(pathname, { locale: nextLocale });
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon">
          <Languages className="h-5 w-5" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={() => changeLocale('en')} disabled={locale === 'en'}>
          English
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => changeLocale('zh')} disabled={locale === 'zh'}>
          简体中文
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}