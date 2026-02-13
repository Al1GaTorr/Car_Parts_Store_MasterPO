import { Home, Camera, FileText, MapPin, Moon, Sun } from 'lucide-react';
import { useTheme } from '@/app/context/ThemeContext';

interface BottomNavProps {
  activeScreen: string;
  onNavigate: (screen: string) => void;
}

export function BottomNav({ activeScreen, onNavigate }: BottomNavProps) {
  const { theme, toggleTheme } = useTheme();

  const navItems = [
    { id: 'dashboard', icon: Home, label: 'Dashboard' },
    { id: 'damage', icon: Camera, label: 'Scan' },
    { id: 'history', icon: FileText, label: 'History' },
    { id: 'map', icon: MapPin, label: 'Map' },
  ];

  return (
    <div className="fixed bottom-0 left-0 right-0 bg-zinc-100 dark:bg-zinc-900 border-t border-zinc-200 dark:border-zinc-800 transition-colors">
      <div className="max-w-md mx-auto flex justify-around items-center px-4 py-3">
        {navItems.map((item) => {
          const Icon = item.icon;
          const isActive = activeScreen === item.id;
          return (
            <button
              key={item.id}
              onClick={() => onNavigate(item.id)}
              className="flex flex-col items-center gap-1 transition-colors"
            >
              <Icon
                size={24}
                className={isActive ? 'text-cyan-400' : 'text-zinc-500 dark:text-zinc-400'}
                strokeWidth={isActive ? 2.5 : 2}
              />
              <span
                className={`text-xs ${
                  isActive ? 'text-cyan-400' : 'text-zinc-500 dark:text-zinc-400'
                }`}
              >
                {item.label}
              </span>
            </button>
          );
        })}
        
        {/* Theme Toggle */}
        <button
          onClick={toggleTheme}
          className="flex flex-col items-center gap-1 transition-colors"
          title={`Switch to ${theme === 'dark' ? 'light' : 'dark'} mode`}
        >
          {theme === 'dark' ? (
            <Sun size={24} className="text-zinc-500 dark:text-zinc-400" />
          ) : (
            <Moon size={24} className="text-zinc-500 dark:text-zinc-400" />
          )}
          <span className="text-xs text-zinc-500 dark:text-zinc-400">
            {theme === 'dark' ? 'Light' : 'Dark'}
          </span>
        </button>
      </div>
    </div>
  );
}
