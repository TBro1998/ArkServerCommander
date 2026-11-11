import { useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { prettyMapName } from '@/lib/maps';

// 地图显示名：有翻译用翻译，自定义地图做简单的格式化
export function useMapName() {
  const tMaps = useTranslations('servers.edit.maps');
  return useCallback((map: string) => (tMaps.has(map) ? tMaps(map) : prettyMapName(map)), [tMaps]);
}
