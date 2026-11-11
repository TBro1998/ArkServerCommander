"use client";

import { useTranslations } from 'next-intl';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { RefreshCw, Download, AlertCircle, Layers, ArrowUpCircle } from 'lucide-react';
import { ImageStatus as ImageStatusType } from '@/stores/servers';
import { formatBytes } from '@/lib/format';
import { cn } from '@/lib/utils';

interface ImageStatusProps {
  imageStatus: ImageStatusType;
  checkingUpdates: boolean;
  onRefresh: () => void;
  onManualDownload: () => void;
  onCheckUpdates: () => void;
  onDownloadImage: (imageName: string) => void;
  onUpdateImage: (imageName: string) => void;
}

type ImageEntry = ImageStatusType['images'][string];

// 用各层的已下载字节估算整体进度；层大小未知时按已完成层数估算
function pullProgress(status: ImageEntry): number {
  const layers = Object.values(status.layers ?? {});
  if (layers.length === 0) return 0;
  const total = layers.reduce((sum, l) => sum + (l.size || 0), 0);
  if (total > 0) {
    const done = layers.reduce((sum, l) => sum + (l.status === 'complete' ? l.size : Math.min(l.progress, l.size)), 0);
    return Math.min((done / total) * 100, 100);
  }
  return (layers.filter((l) => l.status === 'complete').length / layers.length) * 100;
}

// 镜像管理：每个镜像一行，状态 + 当前能做的操作放在同一行
export function ImageStatus({
  imageStatus,
  checkingUpdates,
  onRefresh,
  onManualDownload,
  onCheckUpdates,
  onDownloadImage,
  onUpdateImage,
}: ImageStatusProps) {
  const t = useTranslations('servers.dockerImages');

  const displayName = (imageName: string): string => {
    switch (imageName) {
      case 'tbro98/ase-server:latest':
        return t('arkServer');
      case 'alpine:latest':
        return t('alpineSystem');
      default:
        return imageName;
    }
  };

  const missing = !imageStatus.can_start_server && !imageStatus.any_pulling;

  return (
    <div className="space-y-4 text-sm">
      <div className="flex flex-wrap items-center gap-2">
        {missing ? (
          <Button size="sm" onClick={onManualDownload}>
            <Download className="h-4 w-4" />
            {t('manualDownload')}
          </Button>
        ) : (
          <Button variant="outline" size="sm" onClick={onCheckUpdates} disabled={checkingUpdates || imageStatus.any_pulling}>
            <ArrowUpCircle className={cn('h-4 w-4', checkingUpdates && 'animate-pulse')} />
            {checkingUpdates ? t('checkingUpdates') : t('checkUpdates')}
          </Button>
        )}
        <Button variant="ghost" size="icon" className="ml-auto h-8 w-8" onClick={onRefresh} title={t('refreshStatus')} aria-label={t('refreshStatus')}>
          <RefreshCw className="h-4 w-4" />
        </Button>
      </div>

      <ul className="divide-y rounded-md border">
        {Object.entries(imageStatus.images ?? {}).map(([imageName, status]) => {
          const tone = status.pulling || status.has_update ? 'busy' : status.ready ? 'ok' : 'off';
          const progress = status.pulling ? pullProgress(status) : 0;
          return (
            <li key={imageName} className="space-y-2 px-3 py-3">
              <div className="flex items-center gap-3">
                <span
                  aria-hidden
                  className={cn(
                    'h-2.5 w-2.5 flex-shrink-0 rounded-full',
                    tone === 'ok' && 'bg-green-500',
                    tone === 'busy' && 'bg-amber-500',
                    tone === 'off' && 'bg-red-500',
                  )}
                />
                <div className="min-w-0 flex-1">
                  <div className="font-medium">{displayName(imageName)}</div>
                  <div className="truncate font-mono text-xs text-muted-foreground">{imageName}</div>
                </div>

                {status.pulling ? (
                  <span className="num text-xs font-medium text-amber-800">
                    {t('downloading')} {progress.toFixed(0)}%
                  </span>
                ) : status.ready && status.has_update ? (
                  <Button size="sm" variant="outline" onClick={() => onUpdateImage(imageName)}>
                    <ArrowUpCircle className="h-4 w-4 text-amber-600" />
                    {t('update')}
                  </Button>
                ) : status.ready ? (
                  <span className="text-xs text-green-700">{t('ready')}</span>
                ) : (
                  <Button size="sm" variant="outline" onClick={() => onDownloadImage(imageName)}>
                    <Download className="h-4 w-4" />
                    {t('download')}
                  </Button>
                )}
              </div>

              {status.pulling && (
                <div className="flex items-center gap-2 pl-5">
                  <Progress value={progress} className="h-1.5 flex-1" />
                  {status.layers && (
                    <Popover>
                      <PopoverTrigger asChild>
                        <Button variant="ghost" size="sm" className="h-6 px-2 text-xs">
                          <Layers className="h-3 w-3" />
                          {Object.keys(status.layers).length} {t('layers')}
                        </Button>
                      </PopoverTrigger>
                      <PopoverContent className="max-h-96 w-96 overflow-y-auto">
                        <div className="mb-3 text-sm font-medium">{t('layerDetails')}</div>
                        <ul className="space-y-3">
                          {Object.entries(status.layers).map(([layerId, layer]) => (
                            <li key={layerId} className="space-y-1.5">
                              <div className="flex items-center justify-between text-xs">
                                <span className="font-mono text-muted-foreground">{layer.id.substring(0, 12)}</span>
                                <span className={layer.status === 'complete' ? 'text-green-700' : 'text-muted-foreground'}>
                                  {t(`layerStatus.${layer.status}`)}
                                </span>
                              </div>
                              <Progress
                                value={layer.size > 0 ? Math.min((layer.progress / layer.size) * 100, 100) : layer.status === 'complete' ? 100 : 0}
                                className="h-1.5"
                              />
                              <div className="num text-xs text-muted-foreground">
                                {formatBytes(layer.progress)} / {layer.size > 0 ? formatBytes(layer.size) : t('unknownSize')}
                              </div>
                            </li>
                          ))}
                        </ul>
                      </PopoverContent>
                    </Popover>
                  )}
                </div>
              )}

              {/* 最近一次拉取失败的原因，下次开始拉取时由后端清除 */}
              {status.error && !status.pulling && (
                <div className="ml-5 flex items-start gap-1.5 rounded-md bg-red-50 p-2 text-xs text-red-800">
                  <AlertCircle className="mt-0.5 h-3.5 w-3.5 flex-shrink-0" />
                  <span className="break-all">{status.error}</span>
                </div>
              )}
            </li>
          );
        })}
      </ul>
    </div>
  );
}
