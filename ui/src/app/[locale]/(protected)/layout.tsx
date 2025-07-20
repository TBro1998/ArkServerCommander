"use client";

import { useAuthActions, useAuthUser } from "@/stores/auth";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Menu, Server, LogOut, Github, User as UserIcon, Sun, Moon } from "lucide-react";
import { useTranslations } from 'next-intl';
import { Link, usePathname } from '@/navigation';
import { LanguageSwitcher } from '@/components/LanguageSwitcher';

export default function ProtectedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const t = useTranslations('navigation');
  const { logout } = useAuthActions();
  const user = useAuthUser();
  const pathname = usePathname();

  const navLinks = [
    { href: '/servers', label: t('servers'), icon: Server },
  ];

  return (
    <div className="min-h-screen flex flex-col bg-gray-50">
      <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container flex h-16 items-center">
          <Link href="/" className="flex items-center mr-8">
            <Server className="h-6 w-6 text-blue-600 mr-2" />
            <span className="font-bold">ARK SM</span>
          </Link>
          
          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center space-x-6 text-sm font-medium">
            {navLinks.map(link => (
              <Link
                key={link.href}
                href={link.href}
                className={`transition-colors hover:text-foreground/80 ${pathname === link.href ? 'text-foreground' : 'text-foreground/60'}`}
              >
                {link.label}
              </Link>
            ))}
          </nav>

          <div className="flex flex-1 items-center justify-end space-x-4">
            {/* Mobile Navigation */}
            <div className="md:hidden">
              <Sheet>
                <SheetTrigger asChild>
                  <Button variant="ghost" size="icon">
                    <Menu className="h-5 w-5" />
                  </Button>
                </SheetTrigger>
                <SheetContent side="left">
                  <nav className="grid gap-6 text-lg font-medium mt-8">
                    <Link href="/" className="flex items-center gap-2 text-lg font-semibold mb-4">
                       <Server className="h-6 w-6 text-blue-600" />
                       <span className="font-bold">ARK SM</span>
                    </Link>
                    {navLinks.map(link => (
                      <Link
                        key={link.href}
                        href={link.href}
                        className="hover:text-foreground"
                      >
                        {link.label}
                      </Link>
                    ))}
                  </nav>
                </SheetContent>
              </Sheet>
            </div>

            <a href="https://github.com/luoxuhai/ARK-Server-Manager" target="_blank" rel="noopener noreferrer">
              <Button variant="ghost" size="icon">
                <Github className="h-5 w-5" />
              </Button>
            </a>

            <LanguageSwitcher />
            
            {/* User Menu */}
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="secondary" size="icon" className="rounded-full">
                  <UserIcon className="h-5 w-5" />
                  <span className="sr-only">Toggle user menu</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>{user?.username}</DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => logout()}>
                  <LogOut className="mr-2 h-4 w-4" />
                  <span>{t('logout')}</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </header>
      <main className="flex-1 container py-8">
        {children}
      </main>
    </div>
  );
}