'use client';

import { useEffect, useState } from 'react';
import { Flashcard, Collection } from '@/src/types';
import { flashcardsService, collectionsService } from '@/src/services/collections.service';
import Link from 'next/link';
import { useParams } from 'next/navigation';

export default function CollectionDetailPage() {
  const params = useParams();
  const id = Array.isArray(params.id) ? params.id[0] : params.id;
  const [collectionId, setCollectionId] = useState<string>('');
  const [collection, setCollection] = useState<Collection | null>(null);
  const [flashcards, setFlashcards] = useState<Flashcard[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (id) {
      setCollectionId(String(id));
    }
  }, [id]);

  useEffect(() => {
    if (!collectionId) return;

    const fetchData = async () => {
      try {
        const [collData, cardsData] = await Promise.all([
          collectionsService.getById(collectionId),
          flashcardsService.listByCollection(collectionId),
        ]);
        setCollection(collData);
        setFlashcards(cardsData);
      } catch (err: any) {
        setError(err.response?.data?.message || 'Failed to load collection');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [collectionId]);

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-96">
        <p className="text-lg text-gray-600">Loading collection...</p>
      </div>
    );
  }

  if (!collection) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-600 text-lg">Collection not found</p>
        <Link href="/collections" className="text-blue-600 hover:text-blue-500">
          Back to collections
        </Link>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-start">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">{collection.name}</h1>
          {collection.description && (
            <p className="mt-2 text-gray-600">{collection.description}</p>
          )}
          <div className="mt-4 flex gap-4 text-sm text-gray-500">
            <span>{flashcards.length} cards</span>
            <span>{collection.is_public ? '🌍 Public' : '🔒 Private'}</span>
          </div>
        </div>
        <div className="flex gap-2">
          <Link
            href={`/collections/${collection.id}/generate`}
            className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
          >
            🤖 Generate Cards
          </Link>
          <Link
            href={`/collections/${collection.id}/study`}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            📖 Study
          </Link>
        </div>
      </div>

      {error && (
        <div className="rounded-md bg-red-50 p-4">
          <p className="text-sm font-medium text-red-800">{error}</p>
        </div>
      )}

      {flashcards.length === 0 ? (
        <div className="text-center py-12 bg-white rounded-lg">
          <p className="text-gray-600 text-lg">No flashcards yet</p>
          <div className="mt-4 flex gap-2 justify-center">
            <Link
              href={`/collections/${collection.id}/generate`}
              className="text-blue-600 hover:text-blue-500 font-medium"
            >
              Generate cards with AI
            </Link>
            <span className="text-gray-400">or</span>
            <Link
              href={`/collections/${collection.id}/new-card`}
              className="text-blue-600 hover:text-blue-500 font-medium"
            >
              Create manually
            </Link>
          </div>
        </div>
      ) : (
        <div className="space-y-2">
          {flashcards.map((card) => (
            <div key={card.id} className="p-4 bg-white rounded-lg shadow">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-xs text-gray-500 uppercase font-semibold">Front</p>
                  <p className="text-gray-900">{card.front}</p>
                </div>
                <div>
                  <p className="text-xs text-gray-500 uppercase font-semibold">Back</p>
                  <p className="text-gray-900">{card.back}</p>
                </div>
              </div>
              <div className="mt-3 flex justify-end gap-2 text-xs text-gray-500">
                {card.created_by_ia && <span>🤖 AI Generated</span>}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
