"use client";

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { Server } from '@/stores/servers';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogDescription,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Loader2 } from 'lucide-react';
import { GameUserSettingsEditor } from './GameUserSettingsEditor';
import { GameIniEditor } from './GameIniEditor';
import { ServerArgsEditor } from './ServerArgsEditor';

interface ServerEditModalProps {
  show: boolean;
  mode: 'create' | 'edit';
  server: Partial<Server> | null;
  loading: boolean;
  saving: boolean;
  onClose: () => void;
  onSave: (data: Partial<Server>) => void;
}

export function ServerEditModal({
  show,
  mode,
  server,
  loading,
  saving,
  onClose,
  onSave,
}: ServerEditModalProps) {
  const t = useTranslations('servers.edit');
  const [formData, setFormData] = useState<Partial<Server>>({});

  useEffect(() => {
    if (show) {
      if (mode === 'create') {
        setFormData({
            name: 'ARK Server',
            port: 7777,
            query_port: 27015,
            rcon_port: 27020,
            map: 'TheIsland',
            max_players: 70,
        });
      } else {
        setFormData(server || {});
      }
    }
  }, [show, mode, server]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    const isNumber = type === 'number';
    setFormData((prev) => ({ ...prev, [name]: isNumber ? Number(value) : value }));
  };
  
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSave(formData);
  }

  return (
    <Dialog open={show} onOpenChange={onClose}>
      <DialogContent className="max-w-4xl">
        <DialogHeader>
          <DialogTitle>{mode === 'create' ? t('createTitle') : t('editTitle')}</DialogTitle>
          <DialogDescription>{server?.name}</DialogDescription>
        </DialogHeader>
        
        {loading ? (
            <div className="flex justify-center items-center h-64">
                <Loader2 className="h-8 w-8 animate-spin" />
            </div>
        ) : (
            <Tabs defaultValue="basic">
                <TabsList>
                    <TabsTrigger value="basic">{t('basicParams')}</TabsTrigger>
                    <TabsTrigger value="game_user_settings">{t('gameUserSettings')}</TabsTrigger>
                    <TabsTrigger value="game_ini">{t('gameIni')}</TabsTrigger>
                    <TabsTrigger value="server_args">{t('serverArgs')}</TabsTrigger>
                </TabsList>
                <TabsContent value="basic">
                    <form onSubmit={handleSubmit} className="space-y-4 py-4">
                        <div className="grid grid-cols-2 gap-4">
                            <div>
                                <Label htmlFor="name">{t('serverName')}</Label>
                                <Input id="name" name="name" value={formData.name || ''} onChange={handleChange} />
                            </div>
                             <div>
                                <Label htmlFor="map">{t('map')}</Label>
                                <Input id="map" name="map" value={formData.map || ''} onChange={handleChange} />
                            </div>
                        </div>
                        <div className="grid grid-cols-3 gap-4">
                            <div>
                                <Label htmlFor="port">{t('gamePort')}</Label>
                                <Input id="port" name="port" type="number" value={formData.port || ''} onChange={handleChange} />
                            </div>
                            <div>
                                <Label htmlFor="query_port">{t('queryPort')}</Label>
                                <Input id="query_port" name="query_port" type="number" value={formData.query_port || ''} onChange={handleChange} />
                            </div>
                             <div>
                                <Label htmlFor="rcon_port">{t('rconPort')}</Label>
                                <Input id="rcon_port" name="rcon_port" type="number" value={formData.rcon_port || ''} onChange={handleChange} />
                            </div>
                        </div>
                        <div>
                            <Label htmlFor="admin_password">{t('adminPassword')}</Label>
                            <Input id="admin_password" name="admin_password" type="password" value={formData.admin_password || ''} onChange={handleChange} />
                        </div>
                         <div>
                            <Label htmlFor="max_players">{t('maxPlayers')}</Label>
                            <Input id="max_players" name="max_players" type="number" value={formData.max_players || ''} onChange={handleChange} />
                        </div>
                    </form>
                </TabsContent>
                <TabsContent value="game_user_settings">
                    <GameUserSettingsEditor />
                </TabsContent>
                <TabsContent value="game_ini">
                    <GameIniEditor />
                </TabsContent>
                <TabsContent value="server_args">
                   {/* @ts-ignore */}
                    <ServerArgsEditor value={formData.server_args} onChange={(v) => setFormData(p => ({...p, server_args: v}))} />
                </TabsContent>
            </Tabs>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={onClose}>{t('common.cancel')}</Button>
          <Button onClick={handleSubmit} disabled={saving}>
            {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {t('saveChanges')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}