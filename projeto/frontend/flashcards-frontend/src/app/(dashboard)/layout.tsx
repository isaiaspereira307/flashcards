'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Header from '@/src/components/Header';
import { useAuth } from '@/src/hooks/useAuth';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { isAuthenticated, loadFromStorage } = useAuth();
  const router = useRouter();

  useEffect(() => {
    loadFromStorage();
  }, [loadFromStorage]);

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login');
    }
  }, [isAuthenticated, router]);

  if (!isAuthenticated) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <main className="max-w-7xl mx-auto px-4 py-8">{children}</main>
    </div>
  );
}
