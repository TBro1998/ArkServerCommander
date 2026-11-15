import { useCallback, useSyncExternalStore } from 'react';

// 同一页面内的修改不会触发 storage 事件，用自定义事件通知其他订阅者
const EVENT = 'local-storage-change';

function subscribe(callback: () => void) {
  window.addEventListener('storage', callback);
  window.addEventListener(EVENT, callback);
  return () => {
    window.removeEventListener('storage', callback);
    window.removeEventListener(EVENT, callback);
  };
}

/**
 * 保存在 localStorage 里的选项（例如列表的显示方式），只接受 allowed 里的值。
 * 服务端渲染和首次水合时返回 fallback，避免水合不一致
 */
export function useLocalStorage<T extends string>(key: string, fallback: T, allowed: readonly T[]): [T, (value: T) => void] {
  const value = useSyncExternalStore(
    subscribe,
    () => {
      try {
        return window.localStorage.getItem(key);
      } catch {
        return null;
      }
    },
    () => null,
  );
  const set = useCallback(
    (next: T) => {
      try {
        window.localStorage.setItem(key, next);
      } catch {
        // 隐私模式等无法写入时忽略，保持原来的显示方式
      }
      window.dispatchEvent(new Event(EVENT));
    },
    [key],
  );
  return [allowed.includes(value as T) ? (value as T) : fallback, set];
}
