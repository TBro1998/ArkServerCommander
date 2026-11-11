"use client";

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import axios from 'axios';
import { Loader2, KeyRound } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useAuthUser } from '@/stores/auth';
import { PageHeader } from '@/components/PageHeader';
import { FormError } from '@/components/AuthShell';
import { useApiError } from '@/lib/api-error';
import { toast } from '@/lib/toast';

export default function SettingsPage() {
  const t = useTranslations('settings');
  const user = useAuthUser();
  const apiError = useApiError();
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirm, setConfirm] = useState('');
  const [error, setError] = useState('');
  const [saving, setSaving] = useState(false);

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    const bytes = new TextEncoder().encode(newPassword).length;
    if (bytes < 8) return setError(t('tooShort'));
    if (bytes > 72) return setError(t('tooLong'));
    if (newPassword !== confirm) return setError(t('mismatch'));
    setSaving(true);
    try {
      await axios.put('/api/auth/password', { old_password: oldPassword, new_password: newPassword });
      toast.success(t('changed'));
      setOldPassword('');
      setNewPassword('');
      setConfirm('');
    } catch (err) {
      setError(apiError(err, t('failed')));
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="w-full max-w-2xl space-y-6">
      <PageHeader title={t('title')} />

      <Card>
        <CardHeader>
          <CardTitle className="text-lg">{t('account')}</CardTitle>
        </CardHeader>
        <CardContent className="text-sm">
          <div className="flex justify-between py-1">
            <span className="text-muted-foreground">{t('username')}</span>
            <span className="font-medium">{user?.username}</span>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-lg">
            <KeyRound className="h-5 w-5" />
            {t('changePassword')}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={submit} className="space-y-4">
            <div className="space-y-1">
              <Label htmlFor="old_password">{t('oldPassword')}</Label>
              <Input
                id="old_password"
                type="password"
                autoComplete="current-password"
                value={oldPassword}
                onChange={(e) => setOldPassword(e.target.value)}
                required
              />
            </div>
            <div className="space-y-1">
              <Label htmlFor="new_password">{t('newPassword')}</Label>
              <Input
                id="new_password"
                type="password"
                autoComplete="new-password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                required
              />
              <p className="text-xs text-gray-500">{t('passwordHint')}</p>
            </div>
            <div className="space-y-1">
              <Label htmlFor="confirm_password">{t('confirmPassword')}</Label>
              <Input
                id="confirm_password"
                type="password"
                autoComplete="new-password"
                value={confirm}
                onChange={(e) => setConfirm(e.target.value)}
                required
              />
            </div>
            <FormError message={error} />
            <Button type="submit" disabled={saving}>
              {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {t('submit')}
            </Button>
            <p className="text-xs text-gray-500">{t('sessionNote')}</p>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
