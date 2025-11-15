'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { CreateCollectionRequest } from '@/src/types';
import { collectionsService } from '@/src/services/collections.service';

export default function NewCollectionPage() {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [isPublic, setIsPublic] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const data: CreateCollectionRequest = {
        name,
        description: description || undefined,
        is_public: isPublic,
      };
      const collection = await collectionsService.create(data);
      router.push(`/collections/${collection.id}`);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create collection');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto">
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Create Collection</h1>

      <form onSubmit={handleSubmit} className="space-y-6 bg-white p-6 rounded-lg shadow">
        {error && (
          <div className="rounded-md bg-red-50 p-4">
            <p className="text-sm font-medium text-red-800">{error}</p>
          </div>
        )}

        <div>
          <label htmlFor="name" className="block text-sm font-medium text-gray-700">
            Collection Name *
          </label>
          <input
            id="name"
            type="text"
            required
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            placeholder="e.g., Spanish Vocabulary"
          />
        </div>

        <div>
          <label htmlFor="description" className="block text-sm font-medium text-gray-700">
            Description
          </label>
          <textarea
            id="description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={4}
            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            placeholder="Add a description..."
          />
        </div>

        <div className="flex items-center">
          <input
            id="public"
            type="checkbox"
            checked={isPublic}
            onChange={(e) => setIsPublic(e.target.checked)}
            className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
          />
          <label htmlFor="public" className="ml-2 block text-sm text-gray-700">
            Make this collection public
          </label>
        </div>

        <div className="flex gap-4">
          <button
            type="submit"
            disabled={loading}
            className="flex-1 py-2 px-4 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
          >
            {loading ? 'Creating...' : 'Create Collection'}
          </button>
          <button
            type="button"
            onClick={() => router.back()}
            className="flex-1 py-2 px-4 border border-gray-300 rounded hover:bg-gray-50"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
}
