import { create } from 'zustand';

export type ToastVariant = 'success' | 'error' | 'info';

export interface ToastItem {
  id: number;
  variant: ToastVariant;
  title?: string;
  message: string;
}

interface ToastState {
  toasts: ToastItem[];
  push: (toast: Omit<ToastItem, 'id'>, duration?: number) => void;
  dismiss: (id: number) => void;
}

let nextId = 1;

export const useToastStore = create<ToastState>((set, get) => ({
  toasts: [],
  push: (toast, duration) => {
    const id = nextId++;
    // 同时最多显示 5 条，旧的先移除
    set((state) => ({ toasts: [...state.toasts.slice(-4), { ...toast, id }] }));
    const ms = duration ?? (toast.variant === 'error' ? 8000 : 4000);
    if (ms > 0) setTimeout(() => get().dismiss(id), ms);
  },
  dismiss: (id) => set((state) => ({ toasts: state.toasts.filter((t) => t.id !== id) })),
}));

// 在任何地方（包括事件处理器和 store）都可以调用
export const toast = {
  success: (message: string, title?: string) => useToastStore.getState().push({ variant: 'success', message, title }),
  error: (message: string, title?: string) => useToastStore.getState().push({ variant: 'error', message, title }),
  info: (message: string, title?: string) => useToastStore.getState().push({ variant: 'info', message, title }),
};
