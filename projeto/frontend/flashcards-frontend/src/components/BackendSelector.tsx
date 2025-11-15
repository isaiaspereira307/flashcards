'use client';

import React, { useState } from 'react';
import { useBackend } from '@/src/hooks/useBackend';
import { Backend } from '@/src/types';

export default function BackendSelector() {
  const { backend, setBackend } = useBackend();
  const [open, setOpen] = useState(false);

  const backends: Backend[] = ['golang', 'java', 'fastapi', 'django'];

  const getSwaggerUrl = (backend: Backend): string => {
    const urls: Record<Backend, string> = {
      golang: 'http://localhost:8001/swagger',
      java: 'http://localhost:8002/swagger-ui.html',
      fastapi: 'http://localhost:8003/docs',
      django: 'http://localhost:8004/api/docs',
    };
    return urls[backend];
  };

  return (
    <div className="relative">
      <button
        onClick={() => setOpen(!open)}
        className="px-3 py-1 border border-gray-300 rounded text-sm bg-white text-gray-800 font-semibold hover:bg-gray-50"
      >
        Backend: {backend?.toUpperCase() || 'FASTAPI'} ▼
      </button>

      {open && (
        <div className="absolute right-0 mt-2 w-48 bg-white border border-gray-200 rounded shadow-lg z-10">
          {backends.map((b) => (
            <div key={b} className="flex justify-between items-center p-2 hover:bg-blue-50">
              <button
                onClick={() => {
                  setBackend(b);
                  setOpen(false);
                }}
                className={`flex-1 text-left ${backend === b ? 'font-bold text-blue-600' : 'text-gray-700'}`}
              >
                {b.toUpperCase()}
              </button>
              <a
                href={getSwaggerUrl(b)}
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-500 hover:text-blue-700 text-sm ml-2"
                title="Swagger/Docs"
              >
                📖
              </a>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}