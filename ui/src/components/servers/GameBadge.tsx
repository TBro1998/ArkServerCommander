import { GAME_NAMES, type Game } from '@/lib/games';
import { cn } from '@/lib/utils';

// 游戏标记：ASE / ASA，悬停显示游戏全名
export function GameBadge({ game, className }: { game: Game; className?: string }) {
  return (
    <span
      title={GAME_NAMES[game]}
      className={cn(
        'num rounded-full border px-2 py-0.5 text-xs font-medium',
        game === 'asa' ? 'border-orange-200 bg-orange-50 text-orange-800' : 'border-gray-200 bg-gray-50 text-gray-700',
        className,
      )}
    >
      {game.toUpperCase()}
    </span>
  );
}
