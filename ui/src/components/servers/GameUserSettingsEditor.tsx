"use client";

import { useTranslations } from 'next-intl';

export function GameUserSettingsEditor() {
  const t = useTranslations('servers.editor');
  return <div>{t('visualEditMode')} - GameUserSettings.ini</div>;
}