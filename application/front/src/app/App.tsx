import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'motion/react';
import { Dashboard } from '@/app/components/dashboard';
import { DamageAssessment } from '@/app/components/damage-assessment';
import { ServiceHistory } from '@/app/components/service-history';
import { ServiceMap } from '@/app/components/service-map';
import { AdminPanel } from '@/app/components/admin-panel';
import { BottomNav } from '@/app/components/bottom-nav';
import { ThemeProvider } from '@/app/context/ThemeContext';

function AppContent() {
  const [activeScreen, setActiveScreen] = useState('dashboard');

  useEffect(() => {
    const handler = (e: any) => {
      const d = e?.detail;
      if (typeof d === 'string') setActiveScreen(d);
    };
    window.addEventListener('app-navigate', handler as EventListener);
    return () => window.removeEventListener('app-navigate', handler as EventListener);
  }, []);

  const renderScreen = () => {
    switch (activeScreen) {
      case 'dashboard':
        return <Dashboard />;
      case 'damage':
        return <DamageAssessment />;
      case 'history':
        return <ServiceHistory />;
      case 'map':
        return <ServiceMap />;
      case 'admin':
        return <AdminPanel />;
      default:
        return <Dashboard />;
    }
  };

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 transition-colors">
      {/* Mobile Container */}
      <div className="max-w-md mx-auto bg-zinc-50 dark:bg-zinc-950 min-h-screen relative">
        {/* Screen Content */}
        <AnimatePresence mode="wait">
          <motion.div
            key={activeScreen}
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -20 }}
            transition={{ duration: 0.2 }}
            className="relative"
          >
            {renderScreen()}
          </motion.div>
        </AnimatePresence>

        {/* Bottom Navigation */}
        <BottomNav activeScreen={activeScreen} onNavigate={setActiveScreen} />
      </div>
    </div>
  );
}

export default function App() {
  return (
    <ThemeProvider>
      <AppContent />
    </ThemeProvider>
  );
}

// Listen for global navigate events to switch mobile screens
if (typeof window !== 'undefined') {
  window.addEventListener('navigate', (e: any) => {
    const detail = e?.detail;
    // Dispatch a custom event to the app content via window to be handled there if needed
    // The AppContent uses a BottomNav setter; we simulate by finding mounted React root or rely on CustomEvent handling in components.
    // For simplicity, trigger a window-level event others can listen to.
    const nav = new CustomEvent('app-navigate', { detail });
    window.dispatchEvent(nav);
  });
}