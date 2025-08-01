"use client";

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter, useParams } from 'next/navigation';
import { Server } from '@/stores/servers';
import { serversActions } from '@/stores/servers';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Loader2, ArrowLeft, Save } from 'lucide-react';
import { GameUserSettingsEditor } from '@/components/servers/GameUserSettingsEditor';
import { GameIniEditor } from '@/components/servers/GameIniEditor';
import { ServerArgsEditor } from '@/components/servers/ServerArgsEditor';
import { MapSelector } from '@/components/servers/MapSelector';
import { Alert, AlertDescription } from '@/components/ui/alert';

export default function ServerEditPage() {
  const tServers = useTranslations('servers');
  const tServersEdit = useTranslations('servers.edit');
  const tCommon = useTranslations('common');
  const router = useRouter();
  const params = useParams();
  const serverId = params.id as string;

  const { getServer, updateServer } = serversActions;

  const [server, setServer] = useState<Server | null>(null);
  const [formData, setFormData] = useState<Partial<Server>>({});
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  useEffect(() => {
    const loadServer = async () => {
      try {
        setLoading(true);
        const serverData = await getServer(serverId);
        setServer(serverData);
        setFormData(serverData);
      } catch (err) {
        setError(tServersEdit('loadServerInfoFailed'));
      } finally {
        setLoading(false);
      }
    };

    if (serverId) {
      loadServer();
    }
  }, [serverId, getServer, tServersEdit]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    const isNumber = type === 'number';
    setFormData((prev) => ({ ...prev, [name]: isNumber ? Number(value) : value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setSaving(true);
      await updateServer(serverId, formData);
      setSuccess(tServers('serverUpdateSuccess'));
      // 可选：保存成功后跳转回服务器列表
      // router.push('/servers');
    } catch (err) {
      setError(tServers('serverUpdateError'));
    } finally {
      setSaving(false);
    }
  };

  const handleBack = () => {
    router.push('/servers');
  };

  if (loading) {
    return (
      <div className="w-full max-w-none py-8">
        <div className="flex justify-center items-center h-64">
          <Loader2 className="h-8 w-8 animate-spin" />
          <span className="ml-2">{tCommon('loading')}</span>
        </div>
      </div>
    );
  }

  if (error && !server) {
    return (
      <div className="w-full max-w-none py-8">
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
        <Button onClick={handleBack} className="mt-4">
          <ArrowLeft className="mr-2 h-4 w-4" />
          {tCommon('back')}
        </Button>
      </div>
    );
  }

  return (
    <div className="w-full max-w-none py-8">
      <div className="mb-6">
        <div className="flex items-center gap-4 mb-4">
          <Button variant="outline" onClick={handleBack}>
            <ArrowLeft className="mr-2 h-4 w-4" />
            {tCommon('back')}
          </Button>
          <div>
            <h1 className="text-2xl lg:text-3xl font-bold text-gray-900">{tServersEdit('editTitle')}</h1>
            <p className="text-gray-600">{server?.session_name}</p>
          </div>
        </div>

        {error && (
          <Alert variant="destructive" className="mb-4">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {success && (
          <Alert className="mb-4">
            <AlertDescription>{success}</AlertDescription>
          </Alert>
        )}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{tServersEdit('editTitle')}</CardTitle>
        </CardHeader>
        <CardContent>
          <Tabs defaultValue="basic">
            <TabsList>
              <TabsTrigger value="basic">{tServersEdit('basicParams')}</TabsTrigger>
              <TabsTrigger value="game_user_settings">{tServersEdit('gameUserSettings')}</TabsTrigger>
              <TabsTrigger value="game_ini">{tServersEdit('gameIni')}</TabsTrigger>
              <TabsTrigger value="server_args">{tServersEdit('serverArgs')}</TabsTrigger>
            </TabsList>

            <TabsContent value="basic">
              <form onSubmit={handleSubmit} className="space-y-4 py-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label htmlFor="session_name">{tServersEdit('serverName')}</Label>
                    <Input id="session_name" name="session_name" value={formData.session_name || ''} onChange={handleChange} />
                  </div>
                  <div>
                    <MapSelector
                      value={formData.map || ''}
                      onChange={(value) => setFormData(prev => ({ ...prev, map: value }))}
                      label={tServersEdit('map')}
                    />
                  </div>
                </div>
                <div className="grid grid-cols-3 gap-4">
                  <div>
                    <Label htmlFor="port">{tServersEdit('gamePort')}</Label>
                    <Input id="port" name="port" type="number" value={formData.port || ''} onChange={handleChange} />
                  </div>
                  <div>
                    <Label htmlFor="query_port">{tServersEdit('queryPort')}</Label>
                    <Input id="query_port" name="query_port" type="number" value={formData.query_port || ''} onChange={handleChange} />
                  </div>
                  <div>
                    <Label htmlFor="rcon_port">{tServersEdit('rconPort')}</Label>
                    <Input id="rcon_port" name="rcon_port" type="number" value={formData.rcon_port || ''} onChange={handleChange} />
                  </div>
                </div>
                <div>
                  <Label htmlFor="admin_password">{tServersEdit('adminPassword')}</Label>
                  <Input id="admin_password" name="admin_password" type="password" value={formData.admin_password || ''} onChange={handleChange} />
                </div>
                <div>
                  <Label htmlFor="max_players">{tServersEdit('maxPlayers')}</Label>
                  <Input id="max_players" name="max_players" type="number" value={formData.max_players || ''} onChange={handleChange} />
                </div>
              </form>
            </TabsContent>

            <TabsContent value="game_user_settings">
              <GameUserSettingsEditor
                value={formData.game_user_settings}
                onChange={(v) => setFormData(p => ({ ...p, game_user_settings: v }))}
              />
            </TabsContent>

            <TabsContent value="game_ini">
              <GameIniEditor />
            </TabsContent>

            <TabsContent value="server_args">
              {/* @ts-expect-error: Prop 'value' is not available on type 'IntrinsicAttributes' */}
              <ServerArgsEditor value={formData.server_args} onChange={(v) => setFormData(p => ({ ...p, server_args: v }))} />
            </TabsContent>
          </Tabs>

          <div className="flex justify-end gap-2 mt-6 pt-6 border-t">
            <Button variant="outline" onClick={handleBack}>
              {tCommon('cancel')}
            </Button>
            <Button onClick={handleSubmit} disabled={saving}>
              {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              <Save className="mr-2 h-4 w-4" />
              {tServersEdit('saveChanges')}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}