import { useState, useEffect } from 'react';
import { motion } from 'motion/react';
import { Car, CheckCircle, AlertTriangle, Settings } from 'lucide-react';

interface Vehicle {
  vehicle_id: string;
  brand: string;
  model: string;
  year: number;
  mileage_km: number;
  engine_type: string;
  errors: Array<{
    code: string;
    severity: string;
    description: string;
    recommended_action: string;
  }>;
  maintenance_alerts: string[];
  risk_level: string;
  location: string;
}

export function AdminPanel() {
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [selectedVehicle, setSelectedVehicle] = useState<Vehicle | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadVehicles();
    loadSelectedVehicle();
  }, []);

  const loadVehicles = async () => {
    try {
      const response = await fetch('/api/admin/vehicles');
      const data = await response.json();
      setVehicles(data);
    } catch (error) {
      console.error('Failed to load vehicles:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadSelectedVehicle = async () => {
    try {
      const response = await fetch('/api/admin/selected');
      const data = await response.json();
      setSelectedVehicle(data);
    } catch (error) {
      console.error('Failed to load selected vehicle:', error);
    }
  };

  const selectVehicle = async (vehicleId: string) => {
    try {
      const response = await fetch('/api/admin/selected', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ vehicle_id: vehicleId }),
      });
      const data = await response.json();
      setSelectedVehicle(data);
      alert('Vehicle selected! Dashboard will update.');
    } catch (error) {
      console.error('Failed to select vehicle:', error);
      alert('Failed to select vehicle');
    }
  };

  const getRiskColor = (risk: string) => {
    switch (risk) {
      case 'high':
        return 'text-red-400 bg-red-500/10 border-red-500/30';
      case 'medium':
        return 'text-yellow-400 bg-yellow-500/10 border-yellow-500/30';
      default:
        return 'text-emerald-400 bg-emerald-500/10 border-emerald-500/30';
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 pb-24 px-4 flex items-center justify-center transition-colors">
        <div className="text-zinc-900 dark:text-white">Loading...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 pb-24 px-4 transition-colors">
      {/* Header */}
      <div className="pt-6 pb-4">
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="flex items-center gap-3"
        >
          <div className="w-10 h-10 bg-cyan-500/10 rounded-full flex items-center justify-center border border-cyan-500/30">
            <Settings size={20} className="text-cyan-400" />
          </div>
          <div>
            <h1 className="text-2xl text-zinc-900 dark:text-white">Admin Panel</h1>
            <p className="text-sm text-zinc-600 dark:text-zinc-400 mt-1">Select vehicle to display</p>
          </div>
        </motion.div>
      </div>

      {/* Selected Vehicle Info */}
      {selectedVehicle && (
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          className="bg-gradient-to-br from-cyan-950/50 to-zinc-900 rounded-2xl p-5 border border-cyan-800/30 mb-5"
        >
          <div className="flex items-center gap-2 mb-2">
            <CheckCircle size={18} className="text-cyan-400" />
            <span className="text-sm text-cyan-400">Currently Selected</span>
          </div>
            <h2 className="text-xl text-zinc-900 dark:text-white mb-1">
              {selectedVehicle.brand} {selectedVehicle.model}
            </h2>
            <p className="text-sm text-zinc-600 dark:text-zinc-400">
              {selectedVehicle.year} • {selectedVehicle.mileage_km.toLocaleString()} km • {selectedVehicle.location}
            </p>
        </motion.div>
      )}

      {/* Vehicles List */}
      <div className="space-y-3">
        {vehicles.map((vehicle, index) => {
          const isSelected = selectedVehicle?.vehicle_id === vehicle.vehicle_id;
          const hasErrors = vehicle.errors.length > 0;

          return (
            <motion.div
              key={vehicle.vehicle_id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.05 }}
              className={`bg-zinc-100 dark:bg-zinc-900 rounded-xl border overflow-hidden transition-colors ${
                isSelected ? 'border-cyan-500' : 'border-zinc-300 dark:border-zinc-800'
              }`}
            >
              {/* Vehicle Header */}
              <div className="p-4 border-b border-zinc-300 dark:border-zinc-800 transition-colors">
                <div className="flex items-start justify-between mb-3">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      <Car size={18} className="text-cyan-400" />
                      <h3 className="text-zinc-900 dark:text-white text-lg">
                        {vehicle.brand} {vehicle.model}
                      </h3>
                      {isSelected && (
                        <span className="text-xs bg-cyan-500/10 text-cyan-400 px-2 py-1 rounded-full border border-cyan-500/30">
                          Selected
                        </span>
                      )}
                    </div>
                    <p className="text-sm text-zinc-600 dark:text-zinc-400">
                      {vehicle.year} • {vehicle.mileage_km.toLocaleString()} km • {vehicle.engine_type}
                    </p>
                    <p className="text-xs text-zinc-500 dark:text-zinc-500 mt-1">{vehicle.location}</p>
                  </div>
                  <div className={`px-3 py-1 rounded-full border text-xs ${getRiskColor(vehicle.risk_level)}`}>
                    {vehicle.risk_level.toUpperCase()}
                  </div>
                </div>

                {/* Errors */}
                {hasErrors && (
                  <div className="mt-3 space-y-2">
                    {vehicle.errors.map((error, idx) => (
                      <div
                        key={idx}
                        className={`p-2 rounded-lg border text-xs ${
                          error.severity === 'critical' || error.severity === 'high'
                            ? 'bg-red-500/10 border-red-500/30 text-red-400'
                            : error.severity === 'medium'
                            ? 'bg-yellow-500/10 border-yellow-500/30 text-yellow-400'
                            : 'bg-blue-500/10 border-blue-500/30 text-blue-400'
                        }`}
                      >
                        <div className="flex items-center gap-2 mb-1">
                          <AlertTriangle size={14} />
                          <span className="font-semibold">{error.code}</span>
                        </div>
                        <div className="text-zinc-700 dark:text-zinc-300">{error.description}</div>
                      </div>
                    ))}
                  </div>
                )}

                {/* Maintenance Alerts */}
                {vehicle.maintenance_alerts.length > 0 && (
                  <div className="mt-3">
                    <div className="text-xs text-zinc-500 mb-1">Maintenance:</div>
                    <div className="flex flex-wrap gap-2">
                      {vehicle.maintenance_alerts.map((alert, idx) => (
                        <span
                          key={idx}
                          className="text-xs bg-yellow-500/10 text-yellow-400 px-2 py-1 rounded-full border border-yellow-500/30"
                        >
                          {alert}
                        </span>
                      ))}
                    </div>
                  </div>
                )}
              </div>

              {/* Select Button */}
              <button
                onClick={() => selectVehicle(vehicle.vehicle_id)}
                disabled={isSelected}
                className={`w-full py-3 transition-colors ${
                  isSelected
                    ? 'bg-zinc-200 dark:bg-zinc-800 text-zinc-500 dark:text-zinc-500 cursor-not-allowed'
                    : 'bg-cyan-500 hover:bg-cyan-600 text-white'
                }`}
              >
                {isSelected ? 'Currently Selected' : 'Select This Vehicle'}
              </button>
            </motion.div>
          );
        })}
      </div>
    </div>
  );
}

