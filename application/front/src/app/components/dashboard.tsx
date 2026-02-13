import { useState, useEffect } from 'react';
import { motion } from 'motion/react';
import { Battery, Droplet, Gauge, AlertCircle, CheckCircle, Clock, ShoppingCart } from 'lucide-react';

interface DashboardData {
  carInfo: {
    vin?: string;
    model: string;
    year: string;
    plate: string;
    mileage: number;
  };
  healthMetrics: Array<{
    label: string;
    value: number;
    status: string;
    icon: string;
    color: string;
  }>;
  oilChangeData: {
    currentKm: number;
    nextChangeKm: number;
    daysRemaining: number;
  };
  recentAlerts: Array<{
    type: string;
    message: string;
    description: string;
    issueCode?: string;
  }>;
}

export function Dashboard() {
  const [dashboardData, setDashboardData] = useState<DashboardData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadDashboardData();
    // Refresh every 5 seconds
    const interval = setInterval(loadDashboardData, 5000);
    return () => clearInterval(interval);
  }, []);

  const loadDashboardData = async () => {
    try {
      const response = await fetch('/api/dashboard');
      const data = await response.json();
      setDashboardData(data);
    } catch (error) {
      console.error('Failed to load dashboard data:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading || !dashboardData) {
    return (
      <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 pb-24 px-4 flex items-center justify-center">
        <div className="text-zinc-600 dark:text-zinc-400">Loading...</div>
      </div>
    );
  }

  const { carInfo, healthMetrics, oilChangeData, recentAlerts } = dashboardData;

  const kmRemaining = oilChangeData.nextChangeKm - oilChangeData.currentKm;
  const progressPercent = ((oilChangeData.currentKm % 5000) / 5000) * 100;

  const getIcon = (iconName: string) => {
    switch (iconName) {
      case 'Gauge':
        return Gauge;
      case 'Battery':
        return Battery;
      case 'Droplet':
        return Droplet;
      default:
        return Gauge;
    }
  };

  const mapIssueToShop = (code?: string) => {
    if (!code) return undefined;
    switch (code) {
      case 'oil_service_due':
        return 'oil_change';
      case 'brake_pads_worn':
        return 'brake_pads_worn';
      case 'battery_weak':
        return 'battery_dead';
      case 'wipers_bad':
        return 'wipers_bad';
      case 'headlight_out':
        return 'headlight_out';
      default:
        return undefined;
    }
  };

  const openShop = (issueCode?: string) => {
    const vin = (carInfo as any)?.vin;
    if (!vin) {
      alert('VIN not available for this vehicle');
      return;
    }
    const params = new URLSearchParams();
    params.set('vin', vin);
    const mapped = mapIssueToShop(issueCode);
    if (mapped) params.set('issue', mapped);
    const configuredShopUrl = (import.meta.env.VITE_SAIT_URL as string | undefined)?.trim();
    const fallbackShopUrl = `${window.location.protocol}//${window.location.hostname}:5174`;
    const baseShopUrl = configuredShopUrl || fallbackShopUrl;
    const url = new URL('/', baseShopUrl);
    url.search = params.toString();
    window.open(url.toString(), '_blank');
  };

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
            {carInfo.model.substring(0, 2).toUpperCase()}
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
            <p className="text-sm text-zinc-600 dark:text-zinc-400">
              {carInfo.year} • {carInfo.plate}
            </p>
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
            const Icon = getIcon(metric.icon);
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
        className="bg-gradient-to-br from-cyan-100 to-cyan-50 dark:from-cyan-950/50 dark:to-zinc-900 rounded-2xl p-5 border border-cyan-200 dark:border-cyan-800/30 transition-colors"
      >
        <div className="flex items-center gap-2 mb-4">
          <div className="w-10 h-10 bg-cyan-500/10 rounded-full flex items-center justify-center border border-cyan-500/30">
            <Droplet size={20} className="text-cyan-400" />
          </div>
          <div>
            <h3 className="text-zinc-900 dark:text-white">Next Oil Change</h3>
            <p className="text-xs text-zinc-600 dark:text-zinc-400">Recommended maintenance</p>
          </div>
        </div>
        
        <div className="space-y-3">
          <div className="flex items-end justify-between">
            <div>
              <div className="text-3xl text-cyan-400">{kmRemaining}</div>
              <div className="text-xs text-zinc-600 dark:text-zinc-400">km remaining</div>
            </div>
            <div className="text-right">
              <div className="flex items-center gap-1 text-zinc-700 dark:text-zinc-300">
                <Clock size={16} />
                <span className="text-lg">{oilChangeData.daysRemaining}</span>
              </div>
              <div className="text-xs text-zinc-600 dark:text-zinc-400">days left</div>
            </div>
          </div>
          
          <div className="space-y-2">
            <div className="h-2 bg-zinc-300 dark:bg-zinc-800 rounded-full overflow-hidden">
              <div 
                className="h-full bg-gradient-to-r from-cyan-500 to-blue-500 transition-all"
                style={{ width: `${progressPercent}%` }}
              />
            </div>
            <div className="flex justify-between text-xs text-zinc-600 dark:text-zinc-500">
              <span>{oilChangeData.currentKm.toLocaleString()} km</span>
              <span>{oilChangeData.nextChangeKm.toLocaleString()} km</span>
            </div>
          </div>
        </div>

        <motion.button
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
          onClick={() => {
            // Navigate to map page
            const event = new CustomEvent('navigate', { detail: 'map' });
            window.dispatchEvent(event);
          }}
          className="w-full mt-4 bg-cyan-500 hover:bg-cyan-600 text-white py-3 rounded-xl transition-colors"
        >
          Schedule Service
        </motion.button>
      </motion.div>

      {/* Recent Alerts */}
      {recentAlerts.length > 0 && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.5 }}
          className="mt-5"
        >
          <h3 className="text-lg text-zinc-900 dark:text-white mb-3">Recent Alerts</h3>
          <div className="space-y-2">
            {recentAlerts.map((alert, index) => {
              const isCritical = alert.type === 'critical';
              return (
                <div
                  key={index}
                  className={`rounded-xl p-4 border ${
                    isCritical
                      ? 'bg-red-50 dark:bg-red-950/30 border-red-200 dark:border-red-800/30'
                      : 'bg-yellow-50 dark:bg-yellow-950/30 border-yellow-200 dark:border-yellow-800/30'
                  } flex items-start gap-3 transition-colors`}
                >
                  <AlertCircle
                    size={20}
                    className={`${
                      isCritical
                        ? 'text-red-600 dark:text-red-400'
                        : 'text-yellow-600 dark:text-yellow-400'
                    } mt-0.5 flex-shrink-0`}
                  />
                  <div className="flex-1">
                    <div
                      className={`text-sm ${
                        isCritical
                          ? 'text-red-800 dark:text-red-200'
                          : 'text-yellow-800 dark:text-yellow-200'
                      }`}
                    >
                      {alert.message}
                    </div>
                    <div
                      className={`text-xs ${
                        isCritical
                          ? 'text-red-700 dark:text-red-300/70'
                          : 'text-yellow-700 dark:text-yellow-300/70'
                      } mt-1`}
                    >
                      {alert.description}
                    </div>
                    <div className="mt-3 flex flex-wrap gap-2">
                      {isCritical && (
                        <button
                          onClick={() => {
                            const event = new CustomEvent('navigate', { detail: 'map' });
                            window.dispatchEvent(event);
                          }}
                          className="text-xs bg-red-600 hover:bg-red-700 text-white px-3 py-1.5 rounded-lg transition-colors"
                        >
                          Найти сервис
                        </button>
                      )}

                      <button
                        onClick={() => openShop(alert.issueCode)}
                        className="text-xs bg-zinc-900 hover:bg-zinc-800 dark:bg-white dark:hover:bg-zinc-200 dark:text-zinc-900 text-white px-3 py-1.5 rounded-lg transition-colors inline-flex items-center gap-2"
                      >
                        <ShoppingCart size={14} />
                        Заказать запчасти
                      </button>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </motion.div>
      )}
    </div>
  );
}
