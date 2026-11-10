import { create } from 'zustand';
import axios from 'axios';

// 定义服务器对象类型
export interface Server {
    id: number;
    identifier?: string;
    cluster_id?: string;
    session_name: string;
    status: 'running' | 'stopped' | 'starting' | 'stopping';
    last_error?: string;
    port: number;
    query_port: number;
    rcon_port: number;
    admin_password?: string;
    map: string;
    max_players: number;
    game_mod_ids?: string;
    auto_restart?: boolean;
    game_user_settings?: string;
    game_ini?: string;
    server_args?: {
        query_params: Record<string, string>;
        command_line_args: Record<string, string | number | boolean>;
        custom_args: string[];
    };
    image_outdated?: boolean;
    config_pending?: boolean;
    created_at: string;
    updated_at: string;
}

// 定义服务器状态类型
interface ServersState {
    servers: Server[];
    isLoading: boolean;
    isImageStatusLoading: boolean;
    error: string | null;
    imageStatus: ImageStatus | null;
    actions: ServersActions;
}

export interface ImageStatus {
    can_create_server: boolean;
    can_start_server: boolean;
    any_pulling: boolean;
    any_not_ready: boolean;
    overall_status: string;
    pulling_count: number;
    total_images: number;
    images: {
        [imageName: string]: {
            ready: boolean;
            pulling: boolean;
            has_update?: boolean;
            error?: string;
            layers?: {
                [layerId: string]: {
                    id: string;
                    status: 'pending' | 'downloading' | 'extracting' | 'verifying' | 'complete';
                    progress: number;
                    size: number;
                };
            };
        };
    };
}

// 定义服务器操作类型
interface ServersActions {
    fetchServers: () => Promise<void>;
    createServer: (serverData: Partial<Server>) => Promise<Server>;
    updateServer: (serverId: Server['id'], updateData: Partial<Server>) => Promise<Server>;
    deleteServer: (serverId: Server['id']) => Promise<void>;
    getServer: (serverId: Server['id']) => Promise<Server>;
    getServerPassword: (serverId: Server['id']) => Promise<string>;
    getImageStatus: () => Promise<void>;
    startServer: (serverId: Server['id']) => Promise<void>;
    stopServer: (serverId: Server['id']) => Promise<void>;
    updateServerStatus: (serverId: Server['id'], status: Server['status']) => void;
}

const useServersStore = create<ServersState>((set, get) => ({
    servers: [],
    isLoading: false,
    isImageStatusLoading: false,
    error: null,
    imageStatus: null,
    actions: {
        fetchServers: async () => {
            if (get().isLoading) return;
            set({ isLoading: true, error: null });
            try {
                const response = await axios.get('/api/servers');
                set({ servers: response.data.data || [] });
            } catch (error) {
                set({ error: '获取服务器列表失败' });
                throw error;
            } finally {
                set({ isLoading: false });
            }
        },
        createServer: async (serverData) => {
            try {
                const response = await axios.post('/api/servers', serverData);
                const newServer = response.data.data;
                set((state) => ({ servers: [...state.servers, newServer] }));
                return newServer;
            } catch (error) {
                set({ error: '创建服务器失败' });
                throw error;
            }
        },
        updateServer: async (serverId, updateData) => {
            try {
                const response = await axios.put(`/api/servers/${serverId}`, updateData);
                const updatedServer = response.data.data;
                set((state) => ({
                    servers: state.servers.map((s) => (s.id === serverId ? updatedServer : s)),
                }));
                return updatedServer;
            } catch (error) {
                set({ error: '更新服务器失败' });
                throw error;
            }
        },
        deleteServer: async (serverId) => {
            try {
                await axios.delete(`/api/servers/${serverId}`);
                set((state) => ({
                    servers: state.servers.filter((s) => s.id !== serverId),
                }));
            } catch (error) {
                set({ error: '删除服务器失败' });
                throw error;
            }
        },
        getServer: async (serverId) => {
            try {
                const response = await axios.get(`/api/servers/${serverId}`);
                return response.data.data;
            } catch (error) {
                set({ error: '获取服务器信息失败' });
                throw error;
            }
        },
        getServerPassword: async (serverId) => {
            try {
                const response = await axios.get(`/api/servers/${serverId}/rcon`);
                return response.data.data.admin_password;
            } catch (error) {
                set({ error: '获取管理员密码失败' });
                throw error;
            }
        },
        getImageStatus: async () => {
            // 镜像状态和服务器列表各自使用独立的加载标记，两者不会互相阻塞
            if (get().isImageStatusLoading) return;

            set({ isImageStatusLoading: true, error: null });
            try {
                const response = await axios.get('/api/images/status');
                set({ imageStatus: response.data.data });
            } catch (error) {
                set({ error: '获取镜像状态失败' });
                throw error;
            } finally {
                set({ isImageStatusLoading: false });
            }
        },
        startServer: async (serverId) => {
            try {
                await axios.post(`/api/servers/${serverId}/start`, {});
                // 实际结果由服务器列表页轮询获取，这里只标记为启动中
                get().actions.updateServerStatus(serverId, 'starting');
            } catch (error) {
                set({ error: '启动服务器失败' });
                throw error;
            }
        },
        stopServer: async (serverId) => {
            try {
                await axios.post(`/api/servers/${serverId}/stop`, {});
                get().actions.updateServerStatus(serverId, 'stopping');
            } catch (error) {
                set({ error: '停止服务器失败' });
                throw error;
            }
        },
        updateServerStatus: (serverId, status) => {
            set((state) => ({
                servers: state.servers.map((s) => (s.id === serverId ? { ...s, status, last_error: undefined } : s)),
            }));
        },
    },
}));

export const useServers = () => useServersStore((state) => state.servers);
export const useServersIsLoading = () => useServersStore((state) => state.isLoading);
export const useServersError = () => useServersStore((state) => state.error);
export const useImageStatus = () => useServersStore((state) => state.imageStatus);
export const serversActions = useServersStore.getState().actions;

export default useServersStore;