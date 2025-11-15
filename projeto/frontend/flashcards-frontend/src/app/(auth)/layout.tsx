import Header from '@/src/components/Header';

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      {children}
    </div>
  );
}
