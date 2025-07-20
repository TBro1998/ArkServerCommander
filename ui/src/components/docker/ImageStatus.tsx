"use client";

import { useTranslations } from 'next-intl';

export function ImageStatus() {
  const t = useTranslations('servers.dockerImages');
  return <div>{t('title')}</div>;
}