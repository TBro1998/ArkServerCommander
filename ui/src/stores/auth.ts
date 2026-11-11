import { create } from 'zustand';
import axios from 'axios';

// 令牌由 Next 服务端写入 HttpOnly cookie，前端脚本读不到，这里只保存用户信息和登录状态

// 定义用户对象类型
interface User {
  id: number;
  username: string;
  is_admin: boolean;
  created_at: string;
}

// 定义认证状态类型
interface AuthState {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  isInitialized: boolean; // 是否已向服务端确认过登录状态
  actions: AuthActions;
}

// 定义认证操作类型
interface AuthActions {
  checkInit: () => Promise<boolean>;
  init: (credentials: Credentials) => Promise<AuthResult>;
  login: (credentials: Credentials) => Promise<AuthResult>;
  logout: () => Promise<void>;
  restoreSession: () => Promise<void>;
}

// 登录 / 初始化结果；code 是后端的错误码（如 too_many_attempts），界面据此显示对应语言的提示
export interface AuthResult {
  success: boolean;
  message: string;
  code?: string;
}

// 定义登录凭据类型
interface Credentials {
  username?: string;
  password?: string;
}

// 定义API响应类型
interface AuthResponse {
  user: User;
  message: string;
}

const useAuthStore = create<AuthState>((set) => {
  // 登录和初始化共用：成功后服务端已写入 cookie，这里只更新状态
  const submit = async (url: string, credentials: Credentials, fallback: string): Promise<AuthResult> => {
    set({ isLoading: true });
    try {
      const response = await axios.post<AuthResponse>(url, credentials);
      const { user, message } = response.data;
      set({ user, isAuthenticated: true, isInitialized: true });
      return { success: true, message };
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        return { success: false, message: error.response.data?.error || fallback, code: error.response.data?.code };
      }
      return { success: false, message: fallback };
    } finally {
      set({ isLoading: false });
    }
  };

  return {
    user: null,
    isLoading: false,
    isAuthenticated: false,
    isInitialized: false,
    actions: {
      checkInit: async () => {
        try {
          const response = await axios.get('/api/auth/check-init');
          return response.data.initialized;
        } catch (error) {
          console.error('检查初始化状态失败:', error);
          return false;
        }
      },
      init: (credentials) => submit('/api/auth/init', credentials, '初始化失败'),
      login: (credentials) => submit('/api/auth/login', credentials, '登录失败'),
      logout: async () => {
        try {
          await axios.post('/api/auth/logout');
        } catch (error) {
          console.error('退出登录请求失败:', error);
        }
        set({ user: null, isAuthenticated: false });
      },
      // 页面加载时向服务端确认 cookie 是否仍然有效
      restoreSession: async () => {
        try {
          const response = await axios.get<{ user: User }>('/api/profile');
          set({ user: response.data.user, isAuthenticated: true, isInitialized: true });
        } catch {
          set({ user: null, isAuthenticated: false, isInitialized: true });
        }
      },
    },
  };
});

export const useAuthActions = () => useAuthStore((state) => state.actions);
export const useIsAuthenticated = () => useAuthStore((state) => state.isAuthenticated);
export const useAuthUser = () => useAuthStore((state) => state.user);
export const useAuthIsLoading = () => useAuthStore((state) => state.isLoading);
export const useAuthIsInitialized = () => useAuthStore((state) => state.isInitialized);

export default useAuthStore;
