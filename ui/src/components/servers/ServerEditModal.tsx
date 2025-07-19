import { Server } from '@/stores/servers';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useEffect, useState } from 'react';

interface ServerEditModalProps {
  show: boolean;
  mode: 'create' | 'edit';
  server: Partial<Server> | null;
  loading: boolean;
  saving: boolean;
  onClose: () => void;
  onSave: (data: Partial<Server>) => void;
}

export function ServerEditModal({ show, mode, server, loading, saving, onClose, onSave }: ServerEditModalProps) {
  const [formData, setFormData] = useState<Partial<Server>>({});

  useEffect(() => {
    if (server) {
      setFormData(server);
    } else {
      setFormData({ name: '' });
    }
  }, [server]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSave(formData);
  };

  return (
    <Dialog open={show} onOpenChange={(open) => !open && onClose()}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{mode === 'create' ? '创建新服务器' : '编辑服务器'}</DialogTitle>
          <DialogDescription>
            {mode === 'create' ? '填写以下信息以创建一个新的 ARK 服务器。' : `正在编辑服务器 ${server?.name}。`}
          </DialogDescription>
        </DialogHeader>
        {loading ? (
          <div>加载中...</div>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">服务器名称</Label>
              <Input
                id="name"
                name="name"
                value={formData.name || ''}
                onChange={handleChange}
                placeholder="我的 ARK 服务器"
                required
              />
            </div>
            {/* 在未来可以添加更多服务器配置项 */}
            <DialogFooter>
              <Button type="button" variant="ghost" onClick={onClose}>取消</Button>
              <Button type="submit" disabled={saving}>
                {saving ? '保存中...' : '保存'}
              </Button>
            </DialogFooter>
          </form>
        )}
      </DialogContent>
    </Dialog>
  );
}