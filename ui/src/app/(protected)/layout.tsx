"use client";

import { useAuthActions, useAuthUser } from "@/stores/auth";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { Menu, Server, LogOut, Github, User as UserIcon } from "lucide-react";
import Link from "next/link";
import { usePathname } from 'next/navigation';

export default function ProtectedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { logout } = useAuthActions();
  const user = useAuthUser();
  const pathname = usePathname();

  const navLinks = [
    { href: '/servers', label: '服务器管理', icon: Server },
  ];

  return (
    <div className="flex h-screen bg-gray-100">
      {/* 侧边栏 */}
      <aside className="hidden md:flex w-64 flex-col bg-white border-r">
          <div className="h-16 flex items-center px-6 border-b">
              <Server className="h-6 w-6 text-blue-600 mr-2" />
              <h1 className="text-lg font-bold">ARK SM</h1>
          </div>
          <nav className="flex-1 py-6 space-y-2">
            {navLinks.map(({ href, label, icon: Icon }) => (
                <Link key={href} href={href} className={`flex items-center px-6 py-3 text-gray-700 hover:bg-gray-100 ${pathname === href ? 'bg-gray-100 font-semibold' : ''}`}>
                    <Icon className="h-5 w-5 mr-3" />
                    {label}
                </Link>
            ))}
          </nav>
          <div className="p-6 border-t">
              <div className="flex items-center mb-4">
                  <UserIcon className="w-8 h-8 rounded-full bg-gray-200 p-1 mr-3" />
                  <div>
                      <p className="font-semibold">{user?.username}</p>
                      <p className="text-xs text-gray-500">{user?.is_admin ? '管理员' : '用户'}</p>
                  </div>
              </div>
              <Button onClick={() => logout()} variant="ghost" className="w-full justify-start">
                <LogOut className="mr-2 h-4 w-4" />
                退出登录
              </Button>
          </div>
      </aside>

      {/* 主内容区 */}
      <div className="flex-1 flex flex-col overflow-hidden">
        <header className="h-16 flex items-center justify-between md:justify-end px-6 bg-white border-b">
           {/* 移动端菜单 */}
          <Sheet>
            <SheetTrigger asChild>
              <Button variant="outline" size="icon" className="md:hidden">
                <Menu className="h-6 w-6" />
              </Button>
            </SheetTrigger>
            <SheetContent side="left" className="w-64 p-0">
               <div className="h-16 flex items-center px-6 border-b">
                  <Server className="h-6 w-6 text-blue-600 mr-2" />
                  <h1 className="text-lg font-bold">ARK SM</h1>
              </div>
              <nav className="py-6 space-y-2">
                {navLinks.map(({ href, label, icon: Icon }) => (
                    <Link key={href} href={href} className={`flex items-center px-6 py-3 text-gray-700 hover:bg-gray-100 ${pathname === href ? 'bg-gray-100 font-semibold' : ''}`}>
                        <Icon className="h-5 w-5 mr-3" />
                        {label}
                    </Link>
                ))}
              </nav>
            </SheetContent>
          </Sheet>

          <div className="flex items-center gap-4">
             <Link href="https://github.com/luoxuhai/ARK-Server-Manager" target="_blank">
                <Button variant="ghost" size="icon">
                    <Github className="h-5 w-5" />
                </Button>
             </Link>
             <Button onClick={() => logout()} variant="ghost">
                <LogOut className="mr-2 h-4 w-4" />
                退出登录
              </Button>
          </div>
        </header>
        <main className="flex-1 overflow-y-auto p-6">
            {children}
        </main>
      </div>
    </div>
  );
}