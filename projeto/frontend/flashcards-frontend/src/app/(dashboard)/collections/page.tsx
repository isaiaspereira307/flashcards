'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Collection } from '@/src/types';
import { collectionsService } from '@/src/services/collections.service';

export default function CollectionsPage() {
  const [collections, setCollections] = useState<Collection[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchCollections = async () => {
      try {
        const data = await collectionsService.list();
        setCollections(data);
      } catch (err: any) {
        setError(err.response?.data?.message || 'Failed to load collections');
      } finally {
        setLoading(false);
      }
    };

    fetchCollections();
  }, []);

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-96">
        <p className="text-lg text-gray-600">Loading collections...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold text-gray-900">My Collections</h1>
        <Link
          href="/collections/new"
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          ➕ New Collection
        </Link>
      </div>

      {error && (
        <div className="rounded-md bg-red-50 p-4">
          <p className="text-sm font-medium text-red-800">{error}</p>
        </div>
      )}

      {collections.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-600 text-lg">No collections yet.</p>
          <Link
            href="/collections/new"
            className="text-blue-600 hover:text-blue-500 font-medium"
          >
            Create your first collection
          </Link>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {collections.map((collection) => (
            <Link
              key={collection.id}
              href={`/collections/${collection.id}`}
              className="block p-6 bg-white rounded-lg shadow hover:shadow-lg transition"
            >
              <h2 className="text-xl font-semibold text-gray-900">
                {collection.name}
              </h2>
              {collection.description && (
                <p className="mt-2 text-gray-600 text-sm">
                  {collection.description}
                </p>
              )}
              <div className="mt-4 flex justify-between text-sm text-gray-500">
                <span>{collection.card_count || 0} cards</span>
                <span>
                  {collection.is_public ? '🌍 Public' : '🔒 Private'}
                </span>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
