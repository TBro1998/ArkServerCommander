"use client";

import { useTranslations } from 'next-intl';
import { Button } from '@/components/ui/button';
import { Link } from '@/navigation';
import { useAuthUser } from '@/stores/auth';

export default function HomePage() {
    const t = useTranslations('home');
    const profile = useAuthUser();

    return (
        <div className="w-full max-w-none py-8">
            <div className="bg-white shadow-sm rounded-lg p-8">
                <h1 className="text-3xl font-bold text-gray-900 mb-2">{t('title')}</h1>
                <p className="text-gray-600 mb-6">{t('subtitle')}</p>

                <div className="mb-8 p-4 border rounded-lg">
                    <h2 className="text-xl font-semibold mb-2">{t('systemInfo')}</h2>
                    <p><strong>{t('username')}:</strong> {profile?.username}</p>
                    <p><strong>{t('userID')}:</strong> {profile?.id}</p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    <div className="border p-6 rounded-lg">
                        <h3 className="font-semibold text-lg mb-2">{t('serverManagement')}</h3>
                        <p className="text-sm text-gray-600 mb-4">{t('serverManagementDesc')}</p>
                        <Link href="/servers">
                            <Button>{t('startManage')}</Button>
                        </Link>
                    </div>
                    <div className="border p-6 rounded-lg bg-gray-50">
                        <h3 className="font-semibold text-lg mb-2">{t('playerManagement')}</h3>
                        <p className="text-sm text-gray-600 mb-4">{t('playerManagementDesc')}</p>
                        <Button disabled>{t('comingSoon')}</Button>
                    </div>
                    <div className="border p-6 rounded-lg bg-gray-50">
                        <h3 className="font-semibold text-lg mb-2">{t('logMonitoring')}</h3>
                        <p className="text-sm text-gray-600 mb-4">{t('logMonitoringDesc')}</p>
                        <Button disabled>{t('comingSoon')}</Button>
                    </div>
                </div>
                <p className="text-center text-sm text-gray-500 mt-8">{t('tip')}</p>
            </div>
        </div>
    );
}