"use client";

import { useTranslations } from 'next-intl';

export function GameIniEditor() {
  const t = useTranslations('servers.editor');
  return <div>{t('visualEditMode')} - Game.ini</div>;
}