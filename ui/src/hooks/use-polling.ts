import { useEffect, useRef } from 'react';

// 按间隔执行 callback，页面不可见（切到其他标签页）时暂停，重新可见时立即执行一次。
// interval 为 0 或 enabled 为 false 时不轮询
export function usePolling(callback: () => unknown, interval: number, enabled = true) {
  const saved = useRef(callback);
  useEffect(() => {
    saved.current = callback;
  }, [callback]);

  useEffect(() => {
    if (!enabled || interval <= 0) return;
    let timer: ReturnType<typeof setInterval> | null = null;

    const run = () => {
      try {
        const result = saved.current();
        if (result instanceof Promise) result.catch(() => {});
      } catch {
        // 轮询失败不打断后续轮询
      }
    };
    const start = () => {
      if (timer === null) timer = setInterval(run, interval);
    };
    const stop = () => {
      if (timer !== null) {
        clearInterval(timer);
        timer = null;
      }
    };
    const onVisibility = () => {
      if (document.hidden) {
        stop();
      } else {
        run();
        start();
      }
    };

    if (!document.hidden) start();
    document.addEventListener('visibilitychange', onVisibility);
    return () => {
      stop();
      document.removeEventListener('visibilitychange', onVisibility);
    };
  }, [interval, enabled]);
}
