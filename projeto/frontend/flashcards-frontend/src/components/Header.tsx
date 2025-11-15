'use client';

import Link from 'next/link';
import React from 'react';
import { useAuth } from '@/src/hooks/useAuth';
import { useRouter } from 'next/navigation';
import BackendSelector from './BackendSelector';

export default function Header() {
  const { user, logout, isAuthenticated } = useAuth();
  const router = useRouter();

  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  return (
    <header className="bg-white shadow">
      <nav className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
        <Link href="/" className="text-2xl font-bold text-blue-600">
          📚 Flashcards
        </Link>

        <div className="flex gap-4 items-center">
          <BackendSelector />

          {isAuthenticated && user ? (
            <div className="flex gap-4 items-center">
              <span className="text-sm text-gray-800">{user.email}</span>
              <span className="px-3 py-1 bg-blue-100 text-blue-800 rounded text-sm">
                {user.plan}
              </span>
              <button
                onClick={handleLogout}
                className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600"
              >
                Logout
              </button>
            </div>
          ) : (
            <div className="flex gap-4">
              <Link href="/login" className="px-4 py-2 border border-gray-300 rounded text-gray-800 font-medium hover:bg-gray-50">
                Login
              </Link>
              <Link href="/register" className="px-4 py-2 bg-blue-600 text-white rounded font-medium hover:bg-blue-700">
                Register
              </Link>
            </div>
          )}
        </div>
      </nav>
    </header>
  );
}
