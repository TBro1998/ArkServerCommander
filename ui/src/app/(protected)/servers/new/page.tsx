"use client";

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { Server } from '@/stores/servers';
import { serversActions } from '@/stores/servers';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Loader2, ArrowLeft, Plus } from 'lucide-react';
import { GameUserSettingsEditor } from '@/components/servers/GameUserSettingsEditor';
import { GameIniEditor } from '@/components/servers/GameIniEditor';
import { ServerArgsEditor } from '@/components/servers/ServerArgsEditor';
import { Alert, AlertDescription } from '@/components/ui/alert';

export default function ServerNewPage() {
  const t = useTranslations('servers.edit');
  const tCommon = useTranslations('common');
  const router = useRouter();

  const { createServer } = serversActions;

  const [formData, setFormData] = useState<Partial<Server>>({
    name: 'ARK Server',
    port: 7777,
    query_port: 27015,
    rcon_port: 27020,
    map: 'TheIsland',
    max_players: 70,
  });
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    const isNumber = type === 'number';
    setFormData((prev) => ({ ...prev, [name]: isNumber ? Number(value) : value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setSaving(true);
      await createServer(formData);
      // 创建成功后跳转回服务器列表
      router.push('/servers');
    } catch (err) {
      setError(t('serverCreateError'));
    } finally {
      setSaving(false);
    }
  };

  const handleBack = () => {
    router.push('/servers');
  };

  return (
    <div className="w-full max-w-none py-8">
      <div className="mb-6">
        <div className="flex items-center gap-4 mb-4">
          <Button variant="outline" onClick={handleBack}>
            <ArrowLeft className="mr-2 h-4 w-4" />
            {tCommon('back')}
          </Button>
          <div>
            <h1 className="text-2xl lg:text-3xl font-bold text-gray-900">{t('createTitle')}</h1>
            <p className="text-gray-600">{t('createServerDesc')}</p>
          </div>
        </div>

        {error && (
          <Alert variant="destructive" className="mb-4">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t('createTitle')}</CardTitle>
        </CardHeader>
        <CardContent>
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
              <Plus className="mr-2 h-4 w-4" />
              {t('createServer')}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}