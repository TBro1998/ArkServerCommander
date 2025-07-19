import { Server } from '@/stores/servers';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';
import { MoreVertical, Play, StopCircle, Trash, Edit } from 'lucide-react';

interface ServerCardProps {
  server: Server;
  canStartServer: boolean;
  onStart: (server: Server) => void;
  onStop: (server: Server) => void;
  onEdit: (server: Server) => void;
  onDelete: (server: Server) => void;
}

export function ServerCard({ server, canStartServer, onStart, onStop, onEdit, onDelete }: ServerCardProps) {
  const statusClasses: { [key: string]: string } = {
    running: 'text-green-500',
    stopped: 'text-red-500',
    starting: 'text-yellow-500 animate-pulse',
    stopping: 'text-yellow-500 animate-pulse',
  };

  const statusTexts: { [key: string]: string } = {
    running: '运行中',
    stopped: '已停止',
    starting: '启动中',
    stopping: '停止中',
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex justify-between items-center">
          {server.name}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon">
                <MoreVertical className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem onClick={() => onEdit(server)}>
                <Edit className="mr-2 h-4 w-4" />
                编辑
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => onDelete(server)} disabled={server.status === 'running'}>
                <Trash className="mr-2 h-4 w-4" />
                删除
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </CardTitle>
        <CardDescription className={statusClasses[server.status]}>
          {statusTexts[server.status]}
        </CardDescription>
      </CardHeader>
      <CardFooter className="flex justify-end gap-2">
        {server.status === 'stopped' && (
          <Button onClick={() => onStart(server)} disabled={!canStartServer || server.status !== 'stopped'}>
            <Play className="mr-2 h-4 w-4" />
            启动
          </Button>
        )}
        {server.status === 'running' && (
          <Button onClick={() => onStop(server)} variant="destructive" disabled={server.status !== 'running'}>
            <StopCircle className="mr-2 h-4 w-4" />
            停止
          </Button>
        )}
      </CardFooter>
    </Card>
  );
}