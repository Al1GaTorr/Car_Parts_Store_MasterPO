import { motion } from 'motion/react';
import { Battery, Droplet, Gauge, AlertCircle, CheckCircle, Clock } from 'lucide-react';

export function Dashboard() {
  const carInfo = {
    model: 'Toyota Camry XV70',
    year: '2022',
    plate: 'E 777 KZ',
    mileage: 45234,
  };

  const healthMetrics = [
    { label: 'Engine', value: 95, status: 'good', icon: Gauge, color: 'text-emerald-400' },
    { label: 'Battery', value: 87, status: 'good', icon: Battery, color: 'text-emerald-400' },
    { label: 'Oil Level', value: 62, status: 'warning', icon: Droplet, color: 'text-yellow-400' },
  ];

  const oilChangeData = {
    currentKm: 45234,
    nextChangeKm: 50000,
    daysRemaining: 28,
  };

  const kmRemaining = oilChangeData.nextChangeKm - oilChangeData.currentKm;
  const progressPercent = ((oilChangeData.currentKm % 5000) / 5000) * 100;

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 pb-24 px-4 transition-colors">
      {/* Header */}
      <div className="pt-6 pb-4">
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="flex items-center justify-between"
        >
          <div>
            <h1 className="text-2xl text-zinc-900 dark:text-white">Digital Service Book</h1>
            <p className="text-sm text-zinc-600 dark:text-zinc-400 mt-1">Welcome back, Owner</p>
          </div>
          <div className="w-12 h-12 bg-gradient-to-br from-cyan-500 to-blue-600 rounded-full flex items-center justify-center text-white text-lg">
            TC
          </div>
        </motion.div>
      </div>

      {/* Car Info Card */}
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ delay: 0.1 }}
        className="bg-gradient-to-br from-zinc-100 to-zinc-200 dark:from-zinc-900 dark:to-zinc-800 rounded-2xl p-5 border border-zinc-300 dark:border-zinc-700 mb-5 transition-colors"
      >
        <div className="flex items-start justify-between mb-4">
          <div>
            <h2 className="text-xl text-zinc-900 dark:text-white">{carInfo.model}</h2>
            <p className="text-sm text-zinc-600 dark:text-zinc-400">{carInfo.year} • {carInfo.plate}</p>
          </div>
          <div className="bg-emerald-500/10 px-3 py-1 rounded-full border border-emerald-500/30">
            <span className="text-xs text-emerald-400">Active</span>
          </div>
        </div>
        <div className="flex items-center gap-2 text-zinc-700 dark:text-zinc-300">
          <Gauge size={20} className="text-cyan-400" />
          <span className="text-2xl">{carInfo.mileage.toLocaleString()}</span>
          <span className="text-sm text-zinc-600 dark:text-zinc-500">km</span>
        </div>
      </motion.div>

      {/* Health Status */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2 }}
        className="mb-5"
      >
        <h3 className="text-lg text-zinc-900 dark:text-white mb-3 flex items-center gap-2">
          <CheckCircle size={20} className="text-emerald-400" />
          Vehicle Health
        </h3>
        <div className="grid grid-cols-3 gap-3">
          {healthMetrics.map((metric, index) => {
            const Icon = metric.icon;
            return (
              <motion.div
                key={metric.label}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.3 + index * 0.1 }}
                className="bg-zinc-100 dark:bg-zinc-900 rounded-xl p-4 border border-zinc-300 dark:border-zinc-800 transition-colors"
              >
                <Icon size={24} className={metric.color + ' mb-2'} />
                <div className="text-2xl text-zinc-900 dark:text-white mb-1">{metric.value}%</div>
                <div className="text-xs text-zinc-600 dark:text-zinc-400">{metric.label}</div>
              </motion.div>
            );
          })}
        </div>
      </motion.div>

      {/* Oil Change Countdown */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.4 }}
        className="bg-gradient-to-br from-cyan-950/50 to-zinc-900 rounded-2xl p-5 border border-cyan-800/30"
      >
        <div className="flex items-center gap-2 mb-4">
          <div className="w-10 h-10 bg-cyan-500/10 rounded-full flex items-center justify-center border border-cyan-500/30">
            <Droplet size={20} className="text-cyan-400" />
          </div>
          <div>
            <h3 className="text-white">Next Oil Change</h3>
            <p className="text-xs text-zinc-400">Recommended maintenance</p>
          </div>
        </div>
        
        <div className="space-y-3">
          <div className="flex items-end justify-between">
            <div>
              <div className="text-3xl text-cyan-400">{kmRemaining}</div>
              <div className="text-xs text-zinc-400">km remaining</div>
            </div>
            <div className="text-right">
              <div className="flex items-center gap-1 text-zinc-300">
                <Clock size={16} />
                <span className="text-lg">{oilChangeData.daysRemaining}</span>
              </div>
              <div className="text-xs text-zinc-400">days left</div>
            </div>
          </div>
          
          <div className="space-y-2">
            <div className="h-2 bg-zinc-800 rounded-full overflow-hidden">
              <div 
                className="h-full bg-gradient-to-r from-cyan-500 to-blue-500 transition-all"
                style={{ width: `${progressPercent}%` }}
              />
            </div>
            <div className="flex justify-between text-xs text-zinc-500">
              <span>{oilChangeData.currentKm.toLocaleString()} km</span>
              <span>{oilChangeData.nextChangeKm.toLocaleString()} km</span>
            </div>
          </div>
        </div>

        <motion.button
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
          className="w-full mt-4 bg-cyan-500 hover:bg-cyan-600 text-white py-3 rounded-xl transition-colors"
        >
          Schedule Service
        </motion.button>
      </motion.div>

      {/* Recent Alerts */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.5 }}
        className="mt-5"
      >
        <h3 className="text-lg text-white mb-3">Recent Alerts</h3>
        <div className="bg-yellow-950/30 border border-yellow-800/30 rounded-xl p-4 flex items-start gap-3">
          <AlertCircle size={20} className="text-yellow-400 mt-0.5 flex-shrink-0" />
          <div>
            <div className="text-sm text-yellow-200">Oil Level Low</div>
            <div className="text-xs text-yellow-300/70 mt-1">
              Oil level is at 62%. Consider topping up before next service.
            </div>
          </div>
        </div>
      </motion.div>
    </div>
  );
}
