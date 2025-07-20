"use client";

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { Server } from '@/stores/servers';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
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
  MoreVertical,
  Wifi,
  Lock,
  Clock,
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
          <Button variant="ghost" size="icon" onClick={() => onStop(server)}>
            <Square className="h-4 w-4" />
          </Button>
        );
      case 'stopped':
        return (
          <Button variant="ghost" size="icon" onClick={() => onStart(server)} disabled={!canStartServer}>
            <Play className="h-4 w-4" />
          </Button>
        );
      case 'starting':
      case 'stopping':
        return (
          <Button variant="ghost" size="icon" disabled>
            <Loader2 className="h-4 w-4 animate-spin" />
          </Button>
        );
      default:
        return (
          <Button variant="ghost" size="icon" disabled>
            <HelpCircle className="h-4 w-4" />
          </Button>
        );
    }
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex justify-between items-start">
            <CardTitle className="truncate">{server.name}</CardTitle>
            <div className="flex items-center">
                {StartStopButton()}
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon">
                        <MoreVertical className="h-4 w-4" />
                    </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                    <DropdownMenuItem onClick={() => onEdit(server)}>
                        <Edit className="mr-2 h-4 w-4" />
                        {t('editServer')}
                    </DropdownMenuItem>
                    <Popover>
                        <PopoverTrigger asChild>
                        <DropdownMenuItem onSelect={(e) => e.preventDefault()}>
                            <Trash2 className="mr-2 h-4 w-4" />
                            {t('deleteServer')}
                        </DropdownMenuItem>
                        </PopoverTrigger>
                        <PopoverContent>
                        <div className="space-y-2">
                            <p className="text-sm">{t('card.confirmDeleteMessage', { identifier: server.name })}</p>
                            <Button size="sm" variant="destructive" onClick={() => onDelete(server)}>{t('deleteServer')}</Button>
                        </div>
                        </PopoverContent>
                    </Popover>
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
        </div>
        <Badge variant={getStatusVariant(server.status)}>{server.status}</Badge>
      </CardHeader>
      <CardContent className="space-y-4">
        <div>
            <h4 className="flex items-center text-sm font-semibold mb-2"><Wifi className="mr-2 h-4 w-4" /> {t('card.portConfig')}</h4>
            <div className="text-sm space-y-1">
                <div className="flex justify-between"><span>{t('card.gamePort')}</span> <span className="font-mono">{server.port}</span></div>
                <div className="flex justify-between"><span>{t('card.queryPort')}</span> <span className="font-mono">{server.query_port}</span></div>
                <div className="flex justify-between"><span>{t('card.rconPortLabel')}</span> <span className="font-mono">{server.rcon_port}</span></div>
            </div>
        </div>
        <div>
            <h4 className="flex items-center text-sm font-semibold mb-2"><Lock className="mr-2 h-4 w-4" /> {t('card.authInfo')}</h4>
            <div className="text-sm flex justify-between items-center">
                <span>{t('card.adminPassword')}</span>
                <div className="flex items-center gap-1">
                    <span className="font-mono">{showPassword ? server.admin_password : '••••••••'}</span>
                    <Button variant="ghost" size="icon" onClick={() => setShowPassword(!showPassword)}>
                        {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                    </Button>
                    <Button variant="ghost" size="icon" onClick={() => copyToClipboard(server.admin_password)}>
                        <Copy className="h-4 w-4" />
                    </Button>
                </div>
            </div>
        </div>
        <div>
            <h4 className="flex items-center text-sm font-semibold mb-2"><Info className="mr-2 h-4 w-4" /> {t('card.serverName')}</h4>
             <div className="text-sm space-y-1">
                <div className="flex justify-between"><span>{t('card.map')}</span> <span>{server.map}</span></div>
                <div className="flex justify-between"><span>{t('card.maxPlayers')}</span> <span>{server.max_players}</span></div>
             </div>
        </div>
        <div>
            <h4 className="flex items-center text-sm font-semibold mb-2"><Clock className="mr-2 h-4 w-4" /> {t('card.timeInfo')}</h4>
            <div className="text-xs text-gray-500 space-y-1">
                <div className="flex justify-between"><span>{t('card.createdAt')}</span> <span>{new Date(server.created_at).toLocaleString()}</span></div>
                <div className="flex justify-between"><span>{t('card.updatedAt')}</span> <span>{new Date(server.updated_at).toLocaleString()}</span></div>
            </div>
        </div>
      </CardContent>
    </Card>
  );
}