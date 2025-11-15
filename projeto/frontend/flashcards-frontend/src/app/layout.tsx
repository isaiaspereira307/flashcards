import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'Flashcards - Learn with AI',
  description: 'Generate and study flashcards with AI',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="bg-gray-50">
        {children}
      </body>
    </html>
  );
}