import { motion } from 'motion/react';
import { CheckCircle, FileText, Droplet, BadgeCheck, Calendar, MapPin } from 'lucide-react';
import { useEffect, useState } from 'react';

type RecordItem = {
  date: string;
  type?: string;
  description?: string;
  mileage?: number;
  cost?: number;
  shop?: string;
  location?: string;
  verified?: boolean;
};

export function ServiceHistory() {
  const [serviceRecords, setServiceRecords] = useState<RecordItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let es: EventSource | null = null;
    let cancelled = false;

    const fetchHistory = async (vehicleId: string | undefined) => {
      if (!vehicleId) return null;
      try {
        const r = await fetch(`/api/sto-panel/vehicles/${vehicleId}/changes`);
        if (!r.ok) return null;
        const ch = await r.json();
        const hist = (ch && ch.history) ? ch.history : [];
        return hist;
      } catch (e) {
        console.warn('fetchHistory failed', e);
        return null;
      }
    };

    const load = async () => {
      setLoading(true);
      try {
        const sel = await fetch('/api/admin/selected');
        if (sel.ok) {
          const vehicle = await sel.json();
          const vehicleId = vehicle && vehicle.vehicle_id;
          if (vehicleId) {
            const hist = await fetchHistory(vehicleId);
            if (hist && hist.length) {
              setServiceRecords(hist);
              setLoading(false);
            } else {
              setServiceRecords([]);
              setLoading(false);
            }

            // open SSE for this vehicle to receive live updates
            try {
              es = new EventSource(`/ws/cars/${vehicleId}`);
              es.onmessage = async (ev) => {
                if (cancelled) return;
                try {
                  const data = JSON.parse(ev.data);
                  if (data && data.type) {
                    if (data.type === 'SERVICE_RECORD_ADDED' || data.type === 'VEHICLE_STATE_UPDATED') {
                      // refresh persisted history
                      const updated = await fetchHistory(vehicleId);
                      if (updated) setServiceRecords(updated);
                    }
                  }
                } catch (e) {
                  console.warn('SSE parse error', e);
                }
              };
              es.onerror = (err) => {
                console.warn('SSE error', err);
              };
            } catch (e) {
              console.warn('SSE init failed', e);
            }
            return;
          }
        }
      } catch (e) {
        console.warn('Failed to load persisted history', e);
      }

      setServiceRecords([]);
      setLoading(false);
    };
    load();

    return () => {
      cancelled = true;
      if (es) es.close();
    };
  }, []);

  const getColorClasses = (color: string) => {
    const colors: Record<string, { bg: string; border: string; text: string; iconBg: string }> = {
      cyan: { bg: 'bg-cyan-500/10', border: 'border-cyan-500/30', text: 'text-cyan-400', iconBg: 'bg-cyan-500/20' },
      emerald: { bg: 'bg-emerald-500/10', border: 'border-emerald-500/30', text: 'text-emerald-400', iconBg: 'bg-emerald-500/20' },
      blue: { bg: 'bg-blue-500/10', border: 'border-blue-500/30', text: 'text-blue-400', iconBg: 'bg-blue-500/20' },
    };
    return colors['cyan'];
  };

  const totalServices = serviceRecords.length;
  const totalSpent = serviceRecords.reduce((sum, record) => sum + (record.cost || 0), 0);

  if (loading) {
    return (
      <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 pb-24 px-4 flex items-center justify-center">
        <div className="text-zinc-600 dark:text-zinc-400">Loading...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 pb-24 px-4 transition-colors">
      {/* Header */}
      <div className="pt-6 pb-4">
        <motion.div initial={{ opacity: 0, y: -20 }} animate={{ opacity: 1, y: 0 }}>
          <h1 className="text-2xl text-zinc-900 dark:text-white">Service History</h1>
          <p className="text-sm text-zinc-600 dark:text-zinc-400 mt-1">Complete maintenance records</p>
        </motion.div>
      </div>

      {/* Stats */}
      <motion.div initial={{ opacity: 0, scale: 0.95 }} animate={{ opacity: 1, scale: 1 }} transition={{ delay: 0.1 }} className="grid grid-cols-2 gap-3 mb-6">
        <div className="bg-gradient-to-br from-zinc-100 to-zinc-200 dark:from-zinc-900 dark:to-zinc-800 rounded-xl p-4 border border-zinc-300 dark:border-zinc-700 transition-colors">
          <FileText size={20} className="text-cyan-400 mb-2" />
          <div className="text-2xl text-zinc-900 dark:text-white mb-1">{totalServices}</div>
          <div className="text-xs text-zinc-600 dark:text-zinc-400">Total Services</div>
        </div>
        <div className="bg-gradient-to-br from-zinc-100 to-zinc-200 dark:from-zinc-900 dark:to-zinc-800 rounded-xl p-4 border border-zinc-300 dark:border-zinc-700 transition-colors">
          <BadgeCheck size={20} className="text-emerald-400 mb-2" />
          <div className="text-2xl text-zinc-900 dark:text-white mb-1">{totalSpent.toLocaleString()} ₸</div>
          <div className="text-xs text-zinc-600 dark:text-zinc-400">Total Spent</div>
        </div>
      </motion.div>

      {/* Timeline */}
      {serviceRecords.length === 0 && (
        <div className="bg-zinc-100 dark:bg-zinc-900 rounded-xl border border-zinc-300 dark:border-zinc-800 p-4 text-sm text-zinc-600 dark:text-zinc-400 mb-4">
          История обслуживания пока отсутствует.
        </div>
      )}
      <div className="relative">
        <div className="absolute left-[23px] top-0 bottom-0 w-0.5 bg-zinc-300 dark:bg-zinc-800 transition-colors" />

        <div className="space-y-6">
          {serviceRecords.map((record, index) => {
            const Icon = Droplet;
            const colors = getColorClasses('cyan');
            return (
              <motion.div key={index} initial={{ opacity: 0, x: -20 }} animate={{ opacity: 1, x: 0 }} transition={{ delay: 0.2 + index * 0.1 }} className="relative">
                <div className={`absolute left-0 w-12 h-12 ${colors.iconBg} rounded-full flex items-center justify-center border-2 ${colors.border} z-10`}>
                  <Icon size={20} className={colors.text} />
                </div>

                <div className="ml-16 bg-zinc-100 dark:bg-zinc-900 rounded-xl border border-zinc-300 dark:border-zinc-800 overflow-hidden transition-colors">
                  <div className={`${colors.bg} border-b ${colors.border} px-4 py-3`}>
                    <div className="flex items-start justify-between mb-2">
                      <div>
                        <h3 className="text-white mb-1">{record.type}</h3>
                        <p className="text-xs text-zinc-400">{record.description}</p>
                      </div>
                      {record.verified && (
                        <div className="flex items-center gap-1 bg-emerald-500/10 border border-emerald-500/30 px-2 py-1 rounded-full">
                          <CheckCircle size={12} className="text-emerald-400" />
                          <span className="text-xs text-emerald-400">Verified</span>
                        </div>
                      )}
                    </div>
                  </div>

                  <div className="px-4 py-3 space-y-2">
                    <div className="flex items-center gap-2 text-sm">
                      <Calendar size={14} className="text-zinc-500 dark:text-zinc-500" />
                      <span className="text-zinc-600 dark:text-zinc-400">{new Date(record.date).toLocaleDateString()}</span>
                    </div>
                    <div className="flex items-center gap-2 text-sm">
                      <MapPin size={14} className="text-zinc-500 dark:text-zinc-500" />
                      <div className="flex-1">
                        <div className="text-zinc-900 dark:text-white text-sm">{record.shop}</div>
                        <div className="text-xs text-zinc-600 dark:text-zinc-500">{record.location}</div>
                      </div>
                    </div>
                    <div className="flex items-center justify-between pt-2 border-t border-zinc-300 dark:border-zinc-800 transition-colors">
                      <div className="text-xs text-zinc-600 dark:text-zinc-500">Mileage: {(record.mileage||0).toLocaleString()} km</div>
                      <div className={`text-lg ${colors.text}`}>{(record.cost||0).toLocaleString()} ₸</div>
                    </div>
                  </div>
                </div>
              </motion.div>
            );
          })}
        </div>
      </div>

      <motion.button initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.6 }} whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }} className="w-full mt-6 bg-zinc-100 dark:bg-zinc-900 border border-zinc-300 dark:border-zinc-800 text-zinc-900 dark:text-white py-3 rounded-xl hover:bg-zinc-200 dark:hover:bg-zinc-800 transition-colors flex items-center justify-center gap-2">
        <FileText size={18} />
        Export Full Report
      </motion.button>
    </div>
  );
}
