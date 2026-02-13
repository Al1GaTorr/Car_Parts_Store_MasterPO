import { useState, useEffect } from 'react';
import { motion } from 'motion/react';
import { MapPin, Star, Phone, Clock, Navigation, Filter as FilterIcon, AlertTriangle } from 'lucide-react';
import { Badge } from '@/app/components/ui/badge';

interface RepairShop {
  id: number;
  name: string;
  rating: number;
  reviews: number;
  distance: number;
  address: string;
  phone: string;
  hours: string;
  services: string[];
  verified: boolean;
  priceLevel: number;
  latitude?: number;
  longitude?: number;
}

interface Vehicle {
  vehicle_id: string;
  brand: string;
  model: string;
  errors: Array<{
    code: string;
    severity: string;
  }>;
  location: string;
}

export function ServiceMap() {
  const [repairShops, setRepairShops] = useState<RepairShop[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedVehicle, setSelectedVehicle] = useState<Vehicle | null>(null);
  const [hasCriticalErrors, setHasCriticalErrors] = useState(false);

  useEffect(() => {
    loadSelectedVehicle();
  }, []);

  useEffect(() => {
    if (selectedVehicle) {
      loadRepairShops();
    }
  }, [selectedVehicle]);

  const loadSelectedVehicle = async () => {
    try {
      const response = await fetch('/api/admin/selected');
      const data = await response.json();
      setSelectedVehicle(data);
      
      // Check for critical errors
      const critical = data.errors?.some((e: any) => 
        e.severity === 'critical' || e.severity === 'high'
      );
      setHasCriticalErrors(critical);
    } catch (error) {
      console.error('Failed to load selected vehicle:', error);
    }
  };

  const loadRepairShops = async () => {
    setLoading(true);
    try {
      // Default location (Almaty, Kazakhstan) - can be improved with GPS
      const lat = 43.2220;
      const lon = 76.8512;
      const radius = 10;
      
      // Filter by vehicle brand if available
      const brand = selectedVehicle?.brand || '';
      const url = `/api/repair-shops?lat=${lat}&lon=${lon}&radius=${radius}&brand=${encodeURIComponent(brand)}`;
      
      const response = await fetch(url);
      const data = await response.json();
      setRepairShops(data);
    } catch (error) {
      console.error('Failed to load repair shops:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCall = (phone: string) => {
    window.location.href = `tel:${phone}`;
  };

  const handleNavigate = (shop: RepairShop) => {
    if (shop.latitude && shop.longitude) {
      // Open in 2GIS or Google Maps
      const url = `https://2gis.ru/routeSearch/rsType/car/to/${shop.longitude},${shop.latitude}`;
      window.open(url, '_blank');
    } else {
      alert('Координаты недоступны для навигации');
    }
  };

  // Location mapping for Kazakhstan cities
  const getLocationCoords = (location: string) => {
    const locations: Record<string, { lat: number; lon: number }> = {
      'Almaty': { lat: 43.2220, lon: 76.8512 },
      'Astana': { lat: 51.1694, lon: 71.4491 },
      'Shymkent': { lat: 42.3419, lon: 69.5901 },
      'Karaganda': { lat: 49.8014, lon: 73.1025 },
      'Aktobe': { lat: 50.2833, lon: 57.1667 },
      'Taraz': { lat: 42.9000, lon: 71.3667 },
      'Kostanay': { lat: 53.2144, lon: 63.6246 },
    };
    return locations[location] || locations['Almaty'];
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
            <h1 className="text-2xl text-zinc-900 dark:text-white">Nearby Services</h1>
            <p className="text-sm text-zinc-600 dark:text-zinc-400 mt-1">
              {selectedVehicle 
                ? `Service stations for ${selectedVehicle.brand} ${selectedVehicle.model}`
                : 'Find verified repair shops'}
            </p>
          </div>
          <button 
            onClick={loadRepairShops}
            className="w-10 h-10 bg-zinc-100 dark:bg-zinc-900 rounded-full flex items-center justify-center border border-zinc-300 dark:border-zinc-800 transition-colors"
          >
            <FilterIcon size={18} className="text-zinc-600 dark:text-zinc-400" />
          </button>
        </motion.div>
      </div>

      {/* Critical Error Alert */}
      {hasCriticalErrors && (
        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          className="bg-red-50 dark:bg-red-950/30 border border-red-200 dark:border-red-800/30 rounded-xl p-4 mb-5 flex items-start gap-3"
        >
          <AlertTriangle size={20} className="text-red-600 dark:text-red-400 mt-0.5 flex-shrink-0" />
          <div>
            <div className="text-sm font-semibold text-red-800 dark:text-red-200">
              Критическая проблема обнаружена!
            </div>
            <div className="text-xs text-red-700 dark:text-red-300/70 mt-1">
              Рекомендуется срочно обратиться в сервис. Ниже показаны ближайшие станции обслуживания.
            </div>
          </div>
        </motion.div>
      )}

      {/* Map Placeholder */}
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ delay: 0.1 }}
        className="relative bg-zinc-100 dark:bg-zinc-900 rounded-2xl overflow-hidden h-48 border border-zinc-300 dark:border-zinc-800 mb-5 transition-colors"
      >
        {/* Map Background with Grid */}
        <div className="absolute inset-0 bg-gradient-to-br from-zinc-200 to-zinc-300 dark:from-zinc-800 dark:to-zinc-900">
          <div
            className="absolute inset-0 opacity-10"
            style={{
              backgroundImage: `
                linear-gradient(to right, rgba(255, 255, 255, 0.1) 1px, transparent 1px),
                linear-gradient(to bottom, rgba(255, 255, 255, 0.1) 1px, transparent 1px)
              `,
              backgroundSize: '40px 40px',
            }}
          />
        </div>

        {/* Map Markers */}
        <motion.div
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{ delay: 0.3 }}
          className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2"
        >
          <div className="relative">
            {/* Current Location */}
            <div className="w-4 h-4 bg-cyan-500 rounded-full border-2 border-white shadow-lg" />
            <div className="absolute top-1/2 left-1/2 w-16 h-16 -translate-x-1/2 -translate-y-1/2 bg-cyan-500/20 rounded-full animate-pulse" />
          </div>
        </motion.div>

        {/* Shop Markers */}
        {repairShops.slice(0, 4).map((shop, index) => (
          <motion.div
            key={shop.id}
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ delay: 0.4 + index * 0.1 }}
            className="absolute"
            style={{ 
              top: `${20 + index * 20}%`, 
              left: `${50 + (index % 2) * 20}%` 
            }}
          >
            <MapPin size={24} className="text-emerald-400 drop-shadow-lg" fill="currentColor" />
          </motion.div>
        ))}

        {/* 2GIS Badge */}
        <div className="absolute bottom-3 left-3 bg-white dark:bg-zinc-900 px-3 py-1.5 rounded-lg flex items-center gap-2 border border-zinc-300 dark:border-zinc-800 transition-colors">
          <div className="w-5 h-5 bg-gradient-to-br from-orange-500 to-red-500 rounded" />
          <span className="text-xs font-medium text-zinc-900 dark:text-white">2GIS</span>
        </div>
      </motion.div>

      {/* Loading State */}
      {loading && (
        <div className="text-center py-8 text-zinc-600 dark:text-zinc-400">
          Загрузка сервисных станций...
        </div>
      )}

      {/* Shop List */}
      {!loading && repairShops.length === 0 && (
        <div className="text-center py-8 text-zinc-600 dark:text-zinc-400">
          Сервисные станции не найдены
        </div>
      )}

      {!loading && repairShops.length > 0 && (
      <div className="space-y-4">
        {repairShops.map((shop, index) => (
          <motion.div
            key={shop.id}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 + index * 0.1 }}
              className="bg-zinc-100 dark:bg-zinc-900 rounded-xl border border-zinc-300 dark:border-zinc-800 overflow-hidden transition-colors"
          >
            {/* Shop Header */}
              <div className="p-4 border-b border-zinc-300 dark:border-zinc-800">
              <div className="flex items-start justify-between mb-2">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-1">
                      <h3 className="text-white dark:text-white">{shop.name}</h3>
                    {shop.verified && (
                      <Badge className="bg-emerald-500/10 text-emerald-400 border-emerald-500/30 text-xs px-1.5 py-0">
                        Verified
                      </Badge>
                    )}
                  </div>
                  <div className="flex items-center gap-3 text-sm">
                    <div className="flex items-center gap-1">
                      <Star size={14} className="text-yellow-400" fill="currentColor" />
                        <span className="text-zinc-900 dark:text-white">{shop.rating}</span>
                        <span className="text-zinc-600 dark:text-zinc-500">({shop.reviews})</span>
                    </div>
                      <div className="flex items-center gap-1 text-zinc-600 dark:text-zinc-400">
                      <MapPin size={14} />
                        <span>{shop.distance.toFixed(1)} km</span>
                      </div>
                  </div>
                </div>
                <div className="flex gap-0.5">
                  {[...Array(3)].map((_, i) => (
                    <div
                      key={i}
                      className={`w-1.5 h-3 rounded-full ${
                        i < shop.priceLevel ? 'bg-emerald-400' : 'bg-zinc-700'
                      }`}
                    />
                  ))}
                </div>
              </div>

              {/* Services */}
              <div className="flex flex-wrap gap-2 mb-3">
                {shop.services.map((service, idx) => (
                  <span
                    key={idx}
                      className="text-xs bg-zinc-200 dark:bg-zinc-800 text-zinc-700 dark:text-zinc-300 px-2 py-1 rounded-full transition-colors"
                  >
                    {service}
                  </span>
                ))}
              </div>

              {/* Contact Info */}
              <div className="space-y-1.5">
                  <div className="flex items-center gap-2 text-sm text-zinc-600 dark:text-zinc-400">
                  <MapPin size={14} className="flex-shrink-0" />
                  <span>{shop.address}</span>
                </div>
                  {shop.phone && (
                    <div className="flex items-center gap-2 text-sm text-zinc-600 dark:text-zinc-400">
                  <Phone size={14} className="flex-shrink-0" />
                  <span>{shop.phone}</span>
                </div>
                  )}
                  {shop.hours && (
                <div className="flex items-center gap-2 text-sm">
                      <Clock size={14} className="flex-shrink-0 text-zinc-600 dark:text-zinc-400" />
                  <span className="text-emerald-400">{shop.hours}</span>
                    </div>
                  )}
              </div>
            </div>

            {/* Actions */}
            <div className="grid grid-cols-2 gap-0">
                {shop.phone && (
              <motion.button
                whileHover={{ backgroundColor: 'rgba(39, 39, 42, 1)' }}
                whileTap={{ scale: 0.98 }}
                    onClick={() => handleCall(shop.phone)}
                    className="py-3 text-zinc-700 dark:text-zinc-300 hover:text-white transition-all border-r border-zinc-300 dark:border-zinc-800 flex items-center justify-center gap-2"
              >
                <Phone size={16} />
                <span className="text-sm">Call</span>
              </motion.button>
                )}
              <motion.button
                whileHover={{ backgroundColor: 'rgba(6, 182, 212, 0.1)' }}
                whileTap={{ scale: 0.98 }}
                  onClick={() => handleNavigate(shop)}
                className="py-3 text-cyan-400 hover:text-cyan-300 transition-all flex items-center justify-center gap-2"
              >
                <Navigation size={16} />
                  <span className="text-sm">Route</span>
              </motion.button>
            </div>
          </motion.div>
        ))}
      </div>
      )}
    </div>
  );
}
