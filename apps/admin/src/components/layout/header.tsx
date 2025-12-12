import { useStore } from '@/store/useStore';
import { Button } from '@/components/ui/button';
import { Bell, Menu } from 'lucide-react';

interface HeaderProps {
  onMenuToggle?: () => void;
}

export function Header({ onMenuToggle }: HeaderProps) {
  const { currentProject } = useStore();

  return (
    <header className="sticky top-0 z-30 flex items-center justify-between w-full h-16 px-4 bg-white border-b lg:px-6">
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          onClick={onMenuToggle}
          className="hidden lg:flex"
        >
          <Menu className="h-5 w-5" />
        </Button>
        
        <div>
          <h1 className="text-lg font-semibold">
            {currentProject ? currentProject.name : 'Flagit'}
          </h1>
          {currentProject && (
            <p className="text-sm text-gray-500">
              {currentProject.description}
            </p>
          )}
        </div>
      </div>

      <div className="flex items-center gap-2">
        <Button variant="ghost" size="icon">
          <Bell className="h-5 w-5" />
        </Button>
      </div>
    </header>
  );
}
