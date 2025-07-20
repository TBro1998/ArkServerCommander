"use client";

import { useTranslations } from 'next-intl';

export function ImageUpdateConfirmModal() {
  const t = useTranslations('servers.dockerImages');
  return <div>{t('updateConfirm')}</div>;
}