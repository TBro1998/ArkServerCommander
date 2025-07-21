"use client";

import { useTranslations } from 'next-intl';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { RefreshCw, Download, CheckCircle, AlertCircle, Eye, RotateCcw } from 'lucide-react';
import { ImageStatus as ImageStatusType } from '@/stores/servers';

interface ImageStatusProps {
  imageStatus: ImageStatusType;
  onRefresh: () => void;
  onManualDownload: () => void;
  onCheckUpdates: () => void;
  onDownloadImage: (imageName: string) => void;
  onUpdateImage: (imageName: string) => void;
}

export function ImageStatus({
  imageStatus,
  onRefresh,
  onManualDownload,
  onCheckUpdates,
  onDownloadImage,
  onUpdateImage
}: ImageStatusProps) {
  const t = useTranslations('servers.dockerImages');

  // 格式化字节大小
  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  // 获取镜像显示名称
  const getImageDisplayName = (imageName: string): string => {
    switch (imageName) {
      case 'tbro98/ase-server:latest': return t('arkServer');
      case 'alpine:latest': return t('alpineSystem');
      default: return imageName;
    }
  };

  // 获取状态文本
  const getStatusText = (status: { ready: boolean; pulling: boolean; has_update?: boolean }): string => {
    if (status.pulling) {
      return t('downloading');
    } else if (status.ready) {
      return t('ready');
    } else {
      return t('notReady');
    }
  };

  return (
    <div className="flex flex-col gap-4 text-sm w-full">
      {/* 总体状态 */}
      <div className="flex items-center gap-2">
        {imageStatus.any_pulling ? (
          <div className="flex items-center gap-2 text-yellow-600">
            <RefreshCw className="w-5 h-5 animate-spin" />
            <span className="font-medium">{imageStatus.overall_status}</span>
          </div>
        ) : !imageStatus.can_start_server ? (
          <div className="flex items-center gap-2 text-red-600">
            <AlertCircle className="w-5 h-5" />
            <span className="font-medium">{t('imageMissingManualDownload')}</span>
          </div>
        ) : (
          <div className="flex items-center gap-2 text-green-600">
            <CheckCircle className="w-5 h-5" />
            <span className="font-medium">{t('imageReady')}</span>
          </div>
        )}
        
        <Button
          variant="ghost"
          size="sm"
          onClick={onRefresh}
          className="ml-auto text-blue-600 hover:text-blue-800 hover:bg-blue-50"
          title={t('refreshStatus')}
        >
          <RefreshCw className="w-5 h-5" />
        </Button>
      </div>

      {/* 操作按钮区域 */}
      <div className="flex gap-2">
        {/* 手动下载按钮 */}
        {!imageStatus.can_start_server && !imageStatus.any_pulling && (
          <Button
            variant="default"
            size="sm"
            onClick={onManualDownload}
            className="bg-green-600 hover:bg-green-700"
          >
            <Download className="w-4 h-4 mr-2" />
            {t('manualDownload')}
          </Button>
        )}
        
        {/* 检查更新按钮 */}
        {imageStatus.can_start_server && (
          <Button
            variant="default"
            size="sm"
            onClick={onCheckUpdates}
            className="bg-purple-600 hover:bg-purple-700"
          >
            <RefreshCw className="w-4 h-4 mr-2" />
            {t('checkUpdates')}
          </Button>
        )}
      </div>

      {/* 镜像横向排列 */}
      {imageStatus.images && (
        <div className="flex flex-wrap gap-4">
          {Object.entries(imageStatus.images).map(([imageName, status]) => (
            <Card
              key={imageName}
              className={`p-4 flex-1 min-w-[300px] max-w-[400px] border-l-4 ${
                status.ready && !status.has_update
                  ? 'border-l-green-500'
                  : status.has_update || status.pulling
                  ? 'border-l-yellow-500'
                  : 'border-l-red-500'
              }`}
            >
              {/* 镜像标题和状态 */}
              <div className="flex justify-between items-center mb-3">
                <h3 className="text-sm font-semibold text-gray-900">
                  {getImageDisplayName(imageName)}
                </h3>
                <Badge
                  variant="secondary"
                  className={
                    status.ready && !status.has_update
                      ? 'text-green-600 bg-green-100'
                      : status.has_update || status.pulling
                      ? 'text-yellow-600 bg-yellow-100'
                      : 'text-red-600 bg-red-100'
                  }
                >
                  {getStatusText(status)}
                </Badge>
              </div>

              {/* 层级信息 - 使用 Popover 展示 */}
              {status.pulling && status.layers && (
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <span className="text-xs font-medium text-gray-700">
                      {t('layerProgress')}:
                    </span>
                    <Popover>
                      <PopoverTrigger asChild>
                        <Button
                          variant="ghost"
                          size="sm"
                          className="h-6 px-2 text-xs"
                        >
                          <Eye className="w-3 h-3 mr-1" />
                          {Object.keys(status.layers).length} {t('layers')}
                        </Button>
                      </PopoverTrigger>
                      <PopoverContent className="w-96 max-h-96 overflow-y-auto">
                        <div className="text-sm font-medium mb-3">{t('layerDetails')}</div>
                        <div className="space-y-3">
                          {Object.entries(status.layers).map(([layerId, layer]) => (
                            <div key={layerId} className="bg-white p-3 rounded border">
                              <div className="flex justify-between items-center mb-2">
                                <span className="text-gray-600 font-mono text-xs">
                                  {layer.id.substring(0, 12)}...
                                </span>
                                <Badge
                                  variant="secondary"
                                  className={`text-xs ${
                                    layer.status === 'downloading'
                                      ? 'text-blue-600'
                                      : layer.status === 'extracting'
                                      ? 'text-yellow-600'
                                      : layer.status === 'verifying'
                                      ? 'text-purple-600'
                                      : layer.status === 'complete'
                                      ? 'text-green-600'
                                      : 'text-gray-500'
                                  }`}
                                >
                                  {t(`layerStatus.${layer.status}`) || layer.status}
                                </Badge>
                              </div>
                              
                              <div className="text-xs text-gray-500 mb-2">
                                {layer.size > 0
                                  ? `${formatBytes(layer.progress)} / ${formatBytes(layer.size)}`
                                  : `${formatBytes(layer.progress)} / ${t('unknownSize')}`
                                }
                              </div>
                              
                              {/* 层级进度条 */}
                              <Progress
                                value={layer.size > 0 ? Math.min((layer.progress / layer.size) * 100, 100) : (layer.status === 'complete' ? 100 : 0)}
                                className="h-2"
                              />
                            </div>
                          ))}
                        </div>
                      </PopoverContent>
                    </Popover>
                  </div>
                </div>
              )}

              {/* 镜像就绪状态 */}
              {status.ready && !status.pulling && (
                <div className="text-center py-4">
                  {!status.has_update ? (
                    <CheckCircle className="w-8 h-8 text-green-500 mx-auto mb-2" />
                  ) : (
                    <RefreshCw className="w-8 h-8 text-yellow-500 mx-auto mb-2" />
                  )}
                  <p className={`text-sm font-medium ${!status.has_update ? 'text-green-600' : 'text-yellow-600'}`}>
                    {!status.has_update ? t('ready') : t('updateAvailable')}
                  </p>
                  {status.has_update && (
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => onUpdateImage(imageName)}
                      className="mt-2 text-purple-600 hover:text-purple-800 hover:bg-purple-50"
                      title={t('update')}
                    >
                      <RotateCcw className="w-4 h-4" />
                    </Button>
                  )}
                </div>
              )}

              {/* 镜像未就绪状态 */}
              {!status.ready && !status.pulling && (
                <div className="text-center py-4">
                  <AlertCircle className="w-8 h-8 text-gray-400 mx-auto mb-2" />
                  <p className="text-sm text-gray-600">{t('waitingDownload')}</p>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => onDownloadImage(imageName)}
                    className="mt-2 text-green-600 hover:text-green-800 hover:bg-green-50"
                    title={t('download')}
                  >
                    <Download className="w-4 h-4" />
                  </Button>
                </div>
              )}
            </Card>
          ))}
        </div>
      )}
      
      {/* 统计信息 */}
      {imageStatus.total_images && (
        <div className="mt-6 pt-4 border-t border-gray-200">
          <div className="flex justify-between items-center text-sm text-gray-600">
            <span>{t('totalImages')}: {imageStatus.total_images}</span>
            {imageStatus.pulling_count > 0 && (
              <span className="text-yellow-600 font-medium">
                {t('downloadingCount')}: {imageStatus.pulling_count}
              </span>
            )}
          </div>
        </div>
      )}
    </div>
  );
}