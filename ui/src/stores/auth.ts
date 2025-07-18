import { create } from 'zustand';
import axios from 'axios';
import Cookies from 'js-cookie';

// 定义用户对象类型
interface User {
  id: number;
  username: string;
  is_admin: boolean;
  created_at: string;
}

// 定义认证状态类型
interface AuthState {
  token: string | null;
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  actions: AuthActions;
}

// 定义认证操作类型
interface AuthActions {
  checkInit: () => Promise<boolean>;
  init: (credentials: Credentials) => Promise<{ success: boolean; message: string }>;
  login: (credentials: Credentials) => Promise<{ success: boolean; message: string }>;
  getProfile: () => Promise<void>;
  logout: () => void;
  initFromStorage: () => void;
}

// 定义登录凭据类型
interface Credentials {
  username?: string;
  password?: string;
}

// 定义API响应类型
interface AuthResponse {
  token: string;
  user: User;
  message: string;
}

const useAuthStore = create<AuthState>((set, get) => ({
  token: null,
  user: null,
  isLoading: false,
  isAuthenticated: false,
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
    init: async (credentials) => {
      set({ isLoading: true });
      try {
        const response = await axios.post<AuthResponse>('/api/auth/init', credentials);
        const { token, user, message } = response.data;
        set({ token, user, isAuthenticated: true });
        Cookies.set('auth-token', token, { expires: 7 });
        return { success: true, message };
      } catch (error: any) {
        return { success: false, message: error.response?.data?.error || '初始化失败' };
      } finally {
        set({ isLoading: false });
      }
    },
    login: async (credentials) => {
      set({ isLoading: true });
      try {
        const response = await axios.post<AuthResponse>('/api/auth/login', credentials);
        const { token, user, message } = response.data;
        set({ token, user, isAuthenticated: true });
        Cookies.set('auth-token', token, { expires: 7 });
        return { success: true, message };
      } catch (error: any) {
        return { success: false, message: error.response?.data?.error || '登录失败' };
      } finally {
        set({ isLoading: false });
      }
    },
    getProfile: async () => {
      try {
        const token = get().token;
        if (!token) return;
        const response = await axios.get<{ user: User }>('/api/profile', {
          headers: { Authorization: `Bearer ${token}` },
        });
        set({ user: response.data.user });
      } catch (error) {
        console.error('获取用户信息失败:', error);
        get().actions.logout();
      }
    },
    logout: () => {
      set({ token: null, user: null, isAuthenticated: false });
      Cookies.remove('auth-token');
    },
    initFromStorage: () => {
      const token = Cookies.get('auth-token');
      if (token) {
        set({ token, isAuthenticated: true });
        get().actions.getProfile();
      }
    },
  },
}));

export const useAuthActions = () => useAuthStore((state) => state.actions);
export const useIsAuthenticated = () => useAuthStore((state) => state.isAuthenticated);
export const useAuthUser = () => useAuthStore((state) => state.user);
export const useAuthIsLoading = () => useAuthStore((state) => state.isLoading);

export default useAuthStore;