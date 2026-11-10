"use client";

import { useTranslations } from 'next-intl';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ClosableAlert } from '@/components/ui/closable-alert';
import { Link } from '@/navigation';
import { useAuthUser } from '@/stores/auth';
import useServersStore, { useImageStatus, serversActions, Server as ServerInfo } from '@/stores/servers';
import { ImageStatus } from '@/components/docker/ImageStatus';
import { ImageUpdateConfirmModal } from '@/components/docker/ImageUpdateConfirmModal';
import { useEffect, useRef, useCallback, useState } from 'react';
import { Server } from 'lucide-react';

interface Notice {
    type: 'success' | 'error';
    text: string;
}

interface UpdateTarget {
    imageName: string;
    servers: ServerInfo[];
}

export default function HomePage() {
    const t = useTranslations('home');
    const tImages = useTranslations('servers.dockerImages');
    const tCommon = useTranslations('common');
    const profile = useAuthUser();
    const imageStatus = useImageStatus();
    const { getImageStatus } = serversActions;
    const pollingIntervalRef = useRef<NodeJS.Timeout | null>(null);
    const [notice, setNotice] = useState<Notice | null>(null);
    const [checkingUpdates, setCheckingUpdates] = useState(false);
    const [updateTarget, setUpdateTarget] = useState<UpdateTarget | null>(null);
    const [submittingUpdate, setSubmittingUpdate] = useState(false);

    // 获取镜像状态
    const refreshImageStatus = useCallback(async () => {
        try {
            await getImageStatus();
        } catch (error) {
            console.error('获取镜像状态失败:', error);
        }
    }, [getImageStatus]);

    // 获取镜像显示名称
    const getImageLabel = (imageName: string): string => {
        switch (imageName) {
            case 'tbro98/ase-server:latest': return tImages('arkServer');
            case 'alpine:latest': return tImages('alpineSystem');
            default: return imageName;
        }
    };

    // 从响应中取后端返回的错误信息
    const readError = async (response: Response, fallback: string): Promise<string> => {
        try {
            const body = await response.json();
            return body.error || fallback;
        } catch {
            return fallback;
        }
    };

    // 处理手动下载
    const handleManualDownload = async () => {
        if (!imageStatus?.images) return;

        // 找到第一个未就绪的镜像
        const notReadyImage = Object.entries(imageStatus.images).find(([, status]) => !status.ready);

        if (!notReadyImage) {
            console.log('所有镜像都已就绪');
            return;
        }

        const imageName = notReadyImage[0];
        await handleDownloadImage(imageName);
    };

    // 处理单个镜像下载
    const handleDownloadImage = async (imageName: string) => {
        try {
            const response = await fetch('/api/images/pull', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    image_name: imageName
                })
            });

            if (!response.ok) {
                throw new Error('下载镜像失败');
            }

            console.log(`镜像 ${imageName} 开始下载`);

            // 开始轮询状态
            startPolling();
        } catch (error) {
            console.error('下载镜像失败:', error);
        }
    };

    // 点击更新：先取受影响的服务器，再弹出确认框
    const handleUpdateImage = async (imageName: string) => {
        try {
            const response = await fetch(`/api/images/affected?image_name=${encodeURIComponent(imageName)}`);
            if (!response.ok) {
                setNotice({ type: 'error', text: await readError(response, tImages('loadAffectedFailed')) });
                return;
            }
            const body = await response.json();
            setUpdateTarget({ imageName, servers: body.data?.servers || [] });
        } catch {
            setNotice({ type: 'error', text: tImages('loadAffectedFailed') });
        }
    };

    // 确认更新
    const handleConfirmUpdate = async (restartRunning: boolean) => {
        if (!updateTarget) return;
        setSubmittingUpdate(true);
        try {
            const response = await fetch('/api/images/update', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    image_name: updateTarget.imageName,
                    restart_running: restartRunning,
                })
            });

            if (!response.ok) {
                setNotice({ type: 'error', text: await readError(response, tImages('updateFailed')) });
                return;
            }

            setNotice({ type: 'success', text: tImages('updateStarted') });
            setUpdateTarget(null);

            // 开始轮询状态
            startPolling();
        } catch {
            setNotice({ type: 'error', text: tImages('updateFailed') });
        } finally {
            setSubmittingUpdate(false);
        }
    };

    // 处理检查更新：后端会记录检查结果，刷新镜像状态后即可显示"有新版本"
    const handleCheckUpdates = async () => {
        setCheckingUpdates(true);
        setNotice(null);
        try {
            const response = await fetch('/api/images/check-updates');

            if (!response.ok) {
                setNotice({ type: 'error', text: await readError(response, tImages('checkUpdatesFailed')) });
                return;
            }

            const body = await response.json();
            const results: Record<string, boolean> = body.data || {};
            const anyUpdate = Object.values(results).some(Boolean);
            setNotice({ type: 'success', text: anyUpdate ? tImages('updatesFound') : tImages('upToDate') });

            await refreshImageStatus();
        } catch {
            setNotice({ type: 'error', text: tImages('checkUpdatesFailed') });
        } finally {
            setCheckingUpdates(false);
        }
    };

    // 开始轮询状态
    const startPolling = () => {
        if (pollingIntervalRef.current) {
            clearInterval(pollingIntervalRef.current);
        }

        pollingIntervalRef.current = setInterval(async () => {
            await refreshImageStatus();

            // 从 store 读取最新状态，避免闭包里拿到旧的 imageStatus；没有镜像在拉取中时停止轮询
            if (!useServersStore.getState().imageStatus?.any_pulling) {
                stopPolling();
            }
        }, 2000);
    };

    // 停止轮询
    const stopPolling = () => {
        if (pollingIntervalRef.current) {
            clearInterval(pollingIntervalRef.current);
            pollingIntervalRef.current = null;
        }
    };

    useEffect(() => {
        refreshImageStatus();

        return () => {
            stopPolling();
        };
    }, [refreshImageStatus]);

    return (
        <div className="w-full max-w-none py-8 space-y-8">
            {/* 欢迎标题区域 */}
            <Card>
                <CardContent className="p-8">
                    <div className="text-center space-y-6">
                        <div className="mx-auto w-16 h-16 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-full flex items-center justify-center">
                            <Server className="w-8 h-8 text-white" />
                        </div>
                        <h1 className="text-4xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">
                            {t('title')}
                        </h1>
                        <p className="text-xl text-gray-600 max-w-3xl mx-auto leading-relaxed">
                            {t('subtitle')}
                        </p>
                    </div>

                    <div className="mt-8 p-4 border rounded-lg bg-gray-50">
                        <h2 className="text-xl font-semibold mb-2">{t('systemInfo')}</h2>
                        <p><strong>{t('username')}:</strong> {profile?.username}</p>
                        <p><strong>{t('userID')}:</strong> {profile?.id}</p>
                    </div>
                </CardContent>
            </Card>

            {/* 镜像管理区域 */}
            <Card>
                <CardHeader>
                    <CardTitle className="text-2xl font-semibold text-gray-900">
                        {t('imageManagement')}
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    {notice && (
                        <ClosableAlert
                            key={notice.text}
                            variant={notice.type === 'error' ? 'destructive' : 'default'}
                            className="mb-4"
                            title={notice.type === 'error' ? tCommon('error') : tCommon('success')}
                            onClose={() => setNotice(null)}
                        >
                            {notice.text}
                        </ClosableAlert>
                    )}
                    {imageStatus ? (
                        <ImageStatus
                            imageStatus={imageStatus}
                            checkingUpdates={checkingUpdates}
                            onRefresh={refreshImageStatus}
                            onManualDownload={handleManualDownload}
                            onCheckUpdates={handleCheckUpdates}
                            onDownloadImage={handleDownloadImage}
                            onUpdateImage={handleUpdateImage}
                        />
                    ) : (
                        <div className="text-center py-8">
                            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
                            <p className="mt-2 text-gray-600">{tImages('loadingStatus')}</p>
                        </div>
                    )}
                    {updateTarget && (
                        <ImageUpdateConfirmModal
                            key={updateTarget.imageName}
                            open
                            imageName={updateTarget.imageName}
                            imageLabel={getImageLabel(updateTarget.imageName)}
                            affectedServers={updateTarget.servers}
                            submitting={submittingUpdate}
                            onCancel={() => setUpdateTarget(null)}
                            onConfirm={handleConfirmUpdate}
                        />
                    )}
                </CardContent>
            </Card>

            {/* 功能卡片网格 */}
            <Card>
                <CardHeader>
                    <CardTitle className="text-2xl font-semibold text-gray-900">
                        {t('features')}
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        <Link href="/servers" className="block">
                            <Card className="h-full hover:shadow-xl transition-all duration-300 cursor-pointer group border-2 border-transparent hover:border-blue-200">
                                <CardContent className="text-center space-y-6 p-6">
                                    <div className="mx-auto w-12 h-12 bg-blue-100 group-hover:bg-blue-200 rounded-full flex items-center justify-center transition-colors">
                                        <Server className="w-6 h-6 text-blue-600" />
                                    </div>
                                    <div>
                                        <h3 className="text-xl font-semibold text-gray-900 mb-2">
                                            {t('serverManagement')}
                                        </h3>
                                        <p className="text-gray-600 leading-relaxed">
                                            {t('serverManagementDesc')}
                                        </p>
                                    </div>
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        className="group-hover:bg-blue-50 text-blue-600"
                                    >
                                        {t('startManage')}
                                        <span className="ml-1">→</span>
                                    </Button>
                                </CardContent>
                            </Card>
                        </Link>

                        <Card className="h-full opacity-60 border-2 border-gray-100">
                            <CardContent className="text-center space-y-6 p-6">
                                <div className="mx-auto w-12 h-12 bg-gray-100 rounded-full flex items-center justify-center">
                                    <span className="text-2xl">👥</span>
                                </div>
                                <div>
                                    <h3 className="text-xl font-semibold text-gray-500 mb-2">
                                        {t('playerManagement')}
                                    </h3>
                                    <p className="text-gray-400 leading-relaxed">
                                        {t('playerManagementDesc')}
                                    </p>
                                </div>
                                <Button variant="secondary" size="sm" disabled>
                                    {t('comingSoon')}
                                </Button>
                            </CardContent>
                        </Card>

                        <Card className="h-full opacity-60 border-2 border-gray-100">
                            <CardContent className="text-center space-y-6 p-6">
                                <div className="mx-auto w-12 h-12 bg-gray-100 rounded-full flex items-center justify-center">
                                    <span className="text-2xl">📊</span>
                                </div>
                                <div>
                                    <h3 className="text-xl font-semibold text-gray-500 mb-2">
                                        {t('logMonitoring')}
                                    </h3>
                                    <p className="text-gray-400 leading-relaxed">
                                        {t('logMonitoringDesc')}
                                    </p>
                                </div>
                                <Button variant="secondary" size="sm" disabled>
                                    {t('comingSoon')}
                                </Button>
                            </CardContent>
                        </Card>
                    </div>
                </CardContent>
            </Card>

            <p className="text-center text-sm text-gray-500">{t('tip')}</p>
        </div>
    );
}