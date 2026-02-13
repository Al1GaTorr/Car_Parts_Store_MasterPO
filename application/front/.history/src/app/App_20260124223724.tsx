import { useState } from 'react';
import { motion, AnimatePresence } from 'motion/react';
import { Dashboard } from '@/app/components/dashboard';
import { DamageAssessment } from '@/app/components/damage-assessment';
import { ServiceHistory } from '@/app/components/service-history';
import { ServiceMap } from '@/app/components/service-map';
import { BottomNav } from '@/app/components/bottom-nav';
import { ThemeProvider } from '@/app/context/ThemeContext';

function AppContent() {
  const [activeScreen, setActiveScreen] = useState('dashboard');

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