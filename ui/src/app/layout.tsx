import type { Metadata } from 'next'
import { NextIntlClientProvider } from 'next-intl';
import { getMessages } from 'next-intl/server';
import { getLocale } from '@/lib/locale-server';
import { Toaster } from '@/components/ui/toaster';
import '@fontsource-variable/saira/standard.css'
import './globals.css'

export const metadata: Metadata = {
  title: 'ARK Server Commander',
  description: 'Manage your ARK servers with ease',
}

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const locale = await getLocale();
  const messages = await getMessages();

  return (
    <html lang={locale}>
      <body>
        <NextIntlClientProvider locale={locale} messages={messages}>
          {children}
          <Toaster />
        </NextIntlClientProvider>
      </body>
    </html>
  )
}