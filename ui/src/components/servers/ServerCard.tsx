"use client";

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { Server } from '@/stores/servers';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import {
  Play,
  Square,
  Loader2,
  HelpCircle,
  Info,
  Edit,
  Trash2,
  Wifi,
  Lock,
  Eye,
  EyeOff,
  Copy,
} from 'lucide-react';

interface ServerCardProps {
  server: Server;
  canStartServer: boolean;
  onStart: (server: Server) => void;
  onStop: (server: Server) => void;
  onEdit: (server: Server) => void;
  onDelete: (server: Server) => void;
}

export function ServerCard({
  server,
  canStartServer,
  onStart,
  onStop,
  onEdit,
  onDelete,
}: ServerCardProps) {
  const t = useTranslations('servers');
  const [showPassword, setShowPassword] = useState(false);

  // 获取地图的显示名称（翻译名称或原始名称）
  const getMapDisplayName = (mapName: string) => {
    const mapKey = `edit.maps.${mapName}`;
    const translatedName = t(mapKey);
    // 如果翻译存在且不等于key本身，返回翻译；否则返回原始名称
    return translatedName !== mapKey ? translatedName : mapName;
  };

  const getStatusVariant = (status: Server['status']): 'default' | 'destructive' | 'secondary' | 'outline' => {
    switch (status) {
      case 'running':
        return 'default';
      case 'stopped':
        return 'destructive';
      case 'starting':
      case 'stopping':
        return 'secondary';
      default:
        return 'outline';
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    // You might want to add a toast notification here
  };

  const StartStopButton = () => {
    switch (server.status) {
      case 'running':
        return (
          <Button
            variant="ghost"
            size="sm"
            className="h-8 w-8 p-0 text-red-600 hover:text-red-700 hover:bg-red-50"
            onClick={() => onStop(server)}
          >
            <Square className="h-4 w-4" />
          </Button>
        );
      case 'stopped':
        return (
          <Button
            variant="ghost"
            size="sm"
            className="h-8 w-8 p-0 text-green-600 hover:text-green-700 hover:bg-green-50"
            onClick={() => onStart(server)}
            disabled={!canStartServer}
          >
            <Play className="h-4 w-4" />
          </Button>
        );
      case 'starting':
      case 'stopping':
        return (
          <Button
            variant="ghost"
            size="sm"
            className="h-8 w-8 p-0 text-blue-600"
            disabled
          >
            <Loader2 className="h-4 w-4 animate-spin" />
          </Button>
        );
      default:
        return (
          <Button
            variant="ghost"
            size="sm"
            className="h-8 w-8 p-0 text-gray-400"
            disabled
          >
            <HelpCircle className="h-4 w-4" />
          </Button>
        );
    }
  };

  return (
    <Card className="h-full hover:shadow-lg transition-all duration-200 border-0 shadow-sm">
      <CardHeader className="pb-3">
        <div className="flex justify-between items-start mb-2">
          <CardTitle className="text-base font-semibold truncate pr-2">{server.name}</CardTitle>
          <div className="flex items-center gap-1 flex-shrink-0">
            {/* 启动/停止按钮 */}
            {StartStopButton()}

            {/* 删除按钮 */}
            <Popover>
              <PopoverTrigger asChild>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 w-8 p-0 text-red-600 hover:text-red-700 hover:bg-red-50"
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </PopoverTrigger>
              <PopoverContent>
                <div className="space-y-2">
                  <p className="text-sm">{t('card.confirmDeleteMessage', { identifier: server.name })}</p>
                  <Button size="sm" variant="destructive" onClick={() => onDelete(server)}>{t('deleteServer')}</Button>
                </div>
              </PopoverContent>
            </Popover>

            {/* 编辑按钮 */}
            <Button
              variant="ghost"
              size="sm"
              className="h-8 w-8 p-0 text-blue-600 hover:text-blue-700 hover:bg-blue-50"
              onClick={() => onEdit(server)}
            >
              <Edit className="h-4 w-4" />
            </Button>
          </div>
        </div>
        <Badge
          variant={getStatusVariant(server.status)}
          className="w-fit text-xs px-2 py-1"
        >
          {server.status}
        </Badge>
      </CardHeader>
      <CardContent className="pt-0 space-y-3">
        {/* 端口配置 - 紧凑显示 */}
        <div className="bg-gray-50 rounded-lg p-3">
          <div className="flex items-center mb-2">
            <Wifi className="h-3.5 w-3.5 text-blue-600 mr-1.5" />
            <span className="text-xs font-medium text-gray-700">{t('card.portConfig')}</span>
          </div>
          <div className="grid grid-cols-3 gap-2 text-xs">
            <div className="text-center">
              <div className="text-gray-500">{t('card.gamePort')}</div>
              <div className="font-mono font-semibold">{server.port}</div>
            </div>
            <div className="text-center">
              <div className="text-gray-500">{t('card.queryPort')}</div>
              <div className="font-mono font-semibold">{server.query_port}</div>
            </div>
            <div className="text-center">
              <div className="text-gray-500">RCON</div>
              <div className="font-mono font-semibold">{server.rcon_port}</div>
            </div>
          </div>
        </div>

        {/* 服务器信息 */}
        <div className="space-y-2">
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-600 flex items-center">
              <Info className="h-3.5 w-3.5 mr-1.5 text-green-600" />
              {t('card.map')}
            </span>
            <span className="font-medium truncate ml-2">{getMapDisplayName(server.map)}</span>
          </div>
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-600">{t('card.maxPlayers')}</span>
            <span className="font-medium">{server.max_players}</span>
          </div>
        </div>

        {/* 管理员密码 */}
        <div className="bg-amber-50 rounded-lg p-2">
          <div className="flex items-center justify-between">
            <div className="flex items-center">
              <Lock className="h-3.5 w-3.5 text-amber-600 mr-1.5" />
              <span className="text-xs font-medium text-gray-700">{t('card.adminPassword')}</span>
            </div>
            <div className="flex items-center gap-1">
              <span className="font-mono text-xs">
                {showPassword ? server.admin_password : '••••••••'}
              </span>
              <Button
                variant="ghost"
                size="sm"
                className="h-6 w-6 p-0"
                onClick={() => setShowPassword(!showPassword)}
              >
                {showPassword ? <EyeOff className="h-3 w-3" /> : <Eye className="h-3 w-3" />}
              </Button>
              <Button
                variant="ghost"
                size="sm"
                className="h-6 w-6 p-0"
                onClick={() => copyToClipboard(server.admin_password)}
              >
                <Copy className="h-3 w-3" />
              </Button>
            </div>
          </div>
        </div>


      </CardContent>
    </Card>
  );
}