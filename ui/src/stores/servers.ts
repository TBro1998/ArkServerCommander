import { create } from 'zustand';
import axios from 'axios';

export type ServerStatus = 'running' | 'stopped' | 'starting' | 'stopping';
export type ServerPhase = 'copying' | 'downloading' | 'installing_mods' | 'launching' | 'ready' | '';

// 运行中服务器的实时信息，由后端每 10 秒采集一次
export interface RuntimeInfo {
    phase: ServerPhase;
    progress?: number;
    mod_current?: number;
    mod_total?: number;
    players: number;
    max_players?: number;
    server_name?: string;
    // copying 阶段（启动前复制服务端文件）的来源服务器名称
    copy_source?: string;
    cpu_percent: number;
    memory_bytes: number;
    memory_limit?: number;
    updated_at: string;
}

// 定义服务器对象类型
export interface Server {
    id: number;
    identifier?: string;
    cluster_id?: string;
    session_name: string;
    // 后端可能返回 restarting、unknown 等临时状态，只用于展示
    status: ServerStatus | string;
    last_error?: string;
    port: number;
    query_port: number;
    rcon_port: number;
    admin_password?: string;
    map: string;
    max_players: number;
    game_mod_ids?: string;
    auto_restart?: boolean;
    auto_update?: boolean;
    game_user_settings?: string;
    game_ini?: string;
    server_args?: {
        query_params: Record<string, string>;
        command_line_args: Record<string, string | number | boolean>;
        custom_args: string[];
    };
    image_outdated?: boolean;
    config_pending?: boolean;
    runtime?: RuntimeInfo;
    created_at: string;
    updated_at: string;
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

// 定义服务器状态类型
interface ServersState {
    servers: Server[];
    loaded: boolean;
    isLoading: boolean;
    isImageStatusLoading: boolean;
    imageStatus: ImageStatus | null;
    publicHost: string;
    actions: ServersActions;
}

// 定义服务器操作类型
interface ServersActions {
    fetchServers: () => Promise<void>;
    createServer: (serverData: Partial<Server>) => Promise<Server>;
    updateServer: (serverId: Server['id'], updateData: Partial<Server>) => Promise<Server>;
    deleteServer: (serverId: Server['id'], deleteData?: boolean) => Promise<void>;
    getServer: (serverId: Server['id']) => Promise<Server>;
    getServerPassword: (serverId: Server['id']) => Promise<string>;
    getImageStatus: () => Promise<void>;
    fetchPublicHost: () => Promise<void>;
    startServer: (serverId: Server['id']) => Promise<void>;
    stopServer: (serverId: Server['id']) => Promise<void>;
    restartServer: (serverId: Server['id']) => Promise<void>;
    updateGame: (serverId: Server['id'], game: boolean, mods: boolean) => Promise<void>;
    recreateServer: (serverId: Server['id']) => Promise<void>;
    updateServerStatus: (serverId: Server['id'], status: Server['status']) => void;
}

const useServersStore = create<ServersState>((set, get) => ({
    servers: [],
    loaded: false,
    isLoading: false,
    isImageStatusLoading: false,
    imageStatus: null,
    publicHost: '',
    actions: {
        fetchServers: async () => {
            if (get().isLoading) return;
            set({ isLoading: true });
            try {
                const response = await axios.get('/api/servers');
                set({ servers: response.data.data || [], loaded: true });
            } finally {
                set({ isLoading: false });
            }
        },
        createServer: async (serverData) => {
            const response = await axios.post('/api/servers', serverData);
            const newServer = response.data.data;
            set((state) => ({ servers: [...state.servers, newServer] }));
            return newServer;
        },
        updateServer: async (serverId, updateData) => {
            const response = await axios.put(`/api/servers/${serverId}`, updateData);
            const updatedServer = response.data.data;
            set((state) => ({
                servers: state.servers.map((s) => (s.id === serverId ? { ...s, ...updatedServer } : s)),
            }));
            return updatedServer;
        },
        deleteServer: async (serverId, deleteData = false) => {
            await axios.delete(`/api/servers/${serverId}`, { params: deleteData ? { delete_data: true } : undefined });
            set((state) => ({
                servers: state.servers.filter((s) => s.id !== serverId),
            }));
        },
        getServer: async (serverId) => {
            const response = await axios.get(`/api/servers/${serverId}`);
            return response.data.data;
        },
        getServerPassword: async (serverId) => {
            const response = await axios.get(`/api/servers/${serverId}/rcon`);
            return response.data.data.admin_password;
        },
        getImageStatus: async () => {
            // 镜像状态和服务器列表各自使用独立的加载标记，两者不会互相阻塞
            if (get().isImageStatusLoading) return;
            set({ isImageStatusLoading: true });
            try {
                const response = await axios.get('/api/images/status');
                set({ imageStatus: response.data.data });
            } finally {
                set({ isImageStatusLoading: false });
            }
        },
        fetchPublicHost: async () => {
            if (get().publicHost) return;
            try {
                const response = await axios.get('/api/system/info');
                set({ publicHost: response.data.data?.public_host || '' });
            } catch {
                // 拿不到公网地址时卡片只显示端口
            }
        },
        startServer: async (serverId) => {
            await axios.post(`/api/servers/${serverId}/start`, {});
            // 实际结果由轮询获取，这里只标记为启动中
            get().actions.updateServerStatus(serverId, 'starting');
        },
        stopServer: async (serverId) => {
            await axios.post(`/api/servers/${serverId}/stop`, {});
            get().actions.updateServerStatus(serverId, 'stopping');
        },
        restartServer: async (serverId) => {
            await axios.post(`/api/servers/${serverId}/restart`, {});
            const server = get().servers.find((s) => s.id === serverId);
            get().actions.updateServerStatus(serverId, server?.status === 'running' ? 'stopping' : 'starting');
        },
        updateGame: async (serverId, game, mods) => {
            await axios.post(`/api/servers/${serverId}/update-game`, { game, mods });
            const server = get().servers.find((s) => s.id === serverId);
            get().actions.updateServerStatus(serverId, server?.status === 'running' ? 'stopping' : 'starting');
        },
        recreateServer: async (serverId) => {
            await axios.post(`/api/servers/${serverId}/recreate`, {});
            // 重建在后台执行，期间状态为 stopping，结果由轮询获取
            get().actions.updateServerStatus(serverId, 'stopping');
        },
        updateServerStatus: (serverId, status) => {
            set((state) => ({
                servers: state.servers.map((s) => (s.id === serverId ? { ...s, status, last_error: undefined, runtime: undefined } : s)),
            }));
        },
    },
}));

export const useServers = () => useServersStore((state) => state.servers);
export const useServersLoaded = () => useServersStore((state) => state.loaded);
export const useServersIsLoading = () => useServersStore((state) => state.isLoading);
export const useImageStatus = () => useServersStore((state) => state.imageStatus);
export const usePublicHost = () => useServersStore((state) => state.publicHost);
export const serversActions = useServersStore.getState().actions;

export const isTransitioning = (status: Server['status']) => status === 'starting' || status === 'stopping';

export default useServersStore;
