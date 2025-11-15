'use client';

import Link from 'next/link';
import { useAuth } from '@/src/hooks/useAuth';
import Header from '@/src/components/Header';

export default function Home() {
  const { isAuthenticated } = useAuth();

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <Header />
      <main className="max-w-7xl mx-auto px-4 py-16 sm:py-24">
        <div className="text-center space-y-8">
          <h1 className="text-5xl sm:text-6xl font-bold text-gray-900">
            📚 Learn Better with Flashcards
          </h1>
          <p className="text-xl text-gray-600 max-w-2xl mx-auto">
            Create, organize, and study flashcards powered by AI. Master any subject with
            intelligent spaced repetition.
          </p>

          <div className="flex gap-4 justify-center pt-4">
            {isAuthenticated ? (
              <>
                <Link
                  href="/collections"
                  className="px-8 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium"
                >
                  Go to Collections
                </Link>
              </>
            ) : (
              <>
                <Link
                  href="/register"
                  className="px-8 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium"
                >
                  Get Started
                </Link>
                <Link
                  href="/login"
                  className="px-8 py-3 border border-gray-300 text-gray-900 rounded-lg hover:bg-gray-50 font-medium"
                >
                  Login
                </Link>
              </>
            )}
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mt-16">
            <div className="bg-white p-6 rounded-lg shadow">
              <div className="text-4xl mb-4">🤖</div>
              <h2 className="text-xl font-semibold text-gray-900">AI-Powered</h2>
              <p className="text-gray-600 mt-2">Generate flashcards automatically using AI</p>
            </div>
            <div className="bg-white p-6 rounded-lg shadow">
              <div className="text-4xl mb-4">📊</div>
              <h2 className="text-xl font-semibold text-gray-900">Track Progress</h2>
              <p className="text-gray-600 mt-2">Monitor your learning with detailed statistics</p>
            </div>
            <div className="bg-white p-6 rounded-lg shadow">
              <div className="text-4xl mb-4">🌍</div>
              <h2 className="text-xl font-semibold text-gray-900">Share & Collaborate</h2>
              <p className="text-gray-600 mt-2">Share collections with friends and study together</p>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
