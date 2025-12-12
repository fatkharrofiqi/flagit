import { useStore } from '@/store/useStore';
import { Sidebar } from './sidebar';
import { Header } from './header';
import { useEffect } from 'react';
import type { ReactNode } from 'react';

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const { sidebarOpen, setSidebarOpen } = useStore();

  useEffect(() => {
    // Connect to SSE for real-time updates
    // sseClient.connect();
    
    return () => {
      // sseClient.disconnect();
    };
  }, []);

  return (
    <div className="flex h-screen bg-gray-100">
      <Sidebar />
      
      <div className="flex flex-col flex-1 overflow-hidden">
        <Header 
          onMenuToggle={() => setSidebarOpen(!sidebarOpen)} 
        />
        
        <main className="flex-1 overflow-y-auto p-4 lg:p-6">
          {children}
        </main>
      </div>
    </div>
  );
}
