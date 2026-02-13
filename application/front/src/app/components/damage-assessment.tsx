import { useState } from 'react';
import { motion, AnimatePresence } from 'motion/react';
import { Camera, ScanLine, AlertTriangle, X, Upload } from 'lucide-react';

export function DamageAssessment() {
  const [isScanning, setIsScanning] = useState(false);
  const [scanComplete, setScanComplete] = useState(false);
  const [damageDetected, setDamageDetected] = useState(false);

  const startScan = () => {
    setIsScanning(true);
    setScanComplete(false);
    setDamageDetected(false);

    // Simulate scanning process
    setTimeout(() => {
      setIsScanning(false);
      setScanComplete(true);
      setDamageDetected(true);
    }, 3000);
  };

  const resetScan = () => {
    setIsScanning(false);
    setScanComplete(false);
    setDamageDetected(false);
  };

  const damageResults = [
    { area: 'Front Bumper', severity: 'Minor', confidence: 94 },
    { area: 'Right Door', severity: 'Moderate', confidence: 87 },
  ];

  const estimatedCost = {
    min: 85000,
    max: 120000,
    currency: '₸',
  };

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 pb-24 px-4 transition-colors">
      {/* Header */}
      <div className="pt-6 pb-4">
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
        >
          <h1 className="text-2xl text-zinc-900 dark:text-white">AI Damage Assessment</h1>
          <p className="text-sm text-zinc-600 dark:text-zinc-400 mt-1">Scan your vehicle for damage</p>
        </motion.div>
      </div>

      {/* Camera View */}
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ delay: 0.1 }}
        className="relative bg-zinc-100 dark:bg-zinc-900 rounded-2xl overflow-hidden aspect-[3/4] border border-zinc-300 dark:border-zinc-800 mb-5 transition-colors"
      >
        {/* Simulated Camera Background */}
        <div className="absolute inset-0 bg-gradient-to-b from-zinc-200 to-zinc-300 dark:from-zinc-800 dark:to-zinc-900 flex items-center justify-center transition-colors">
          <Camera size={80} className="text-zinc-400 dark:text-zinc-700" />
        </div>

        {/* Corner Guides */}
        {!scanComplete && (
          <>
            <div className="absolute top-4 left-4 w-12 h-12 border-l-2 border-t-2 border-cyan-400" />
            <div className="absolute top-4 right-4 w-12 h-12 border-r-2 border-t-2 border-cyan-400" />
            <div className="absolute bottom-4 left-4 w-12 h-12 border-l-2 border-b-2 border-cyan-400" />
            <div className="absolute bottom-4 right-4 w-12 h-12 border-r-2 border-b-2 border-cyan-400" />
          </>
        )}

        {/* Scanning Effect */}
        <AnimatePresence>
          {isScanning && (
            <>
              {/* Scanning Lines */}
              <motion.div
                initial={{ top: 0 }}
                animate={{ top: '100%' }}
                exit={{ opacity: 0 }}
                transition={{ duration: 2, repeat: Infinity, ease: 'linear' }}
                className="absolute left-0 right-0 h-1 bg-gradient-to-r from-transparent via-cyan-400 to-transparent"
                style={{
                  boxShadow: '0 0 20px rgba(34, 211, 238, 0.8)',
                }}
              />
              
              {/* Grid Overlay */}
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 0.3 }}
                exit={{ opacity: 0 }}
                className="absolute inset-0"
                style={{
                  backgroundImage: `
                    linear-gradient(to right, rgba(34, 211, 238, 0.1) 1px, transparent 1px),
                    linear-gradient(to bottom, rgba(34, 211, 238, 0.1) 1px, transparent 1px)
                  `,
                  backgroundSize: '30px 30px',
                }}
              />
              
              {/* Scanning Text */}
              <div className="absolute inset-0 flex items-center justify-center">
                <motion.div
                  initial={{ opacity: 0, scale: 0.8 }}
                  animate={{ opacity: 1, scale: 1 }}
                  className="bg-zinc-950/80 backdrop-blur-sm px-6 py-4 rounded-xl border border-cyan-500/30"
                >
                  <div className="flex items-center gap-3">
                    <motion.div
                      animate={{ rotate: 360 }}
                      transition={{ duration: 2, repeat: Infinity, ease: 'linear' }}
                    >
                      <ScanLine size={24} className="text-cyan-400" />
                    </motion.div>
                    <div>
                      <div className="text-white">Analyzing...</div>
                      <div className="text-xs text-cyan-400">AI processing image</div>
                    </div>
                  </div>
                </motion.div>
              </div>
            </>
          )}
        </AnimatePresence>

        {/* Damage Detection Markers */}
        <AnimatePresence>
          {scanComplete && damageDetected && (
            <>
              <motion.div
                initial={{ opacity: 0, scale: 0 }}
                animate={{ opacity: 1, scale: 1 }}
                className="absolute top-1/4 left-1/4 w-20 h-20 border-2 border-red-500 rounded-lg"
              >
                <div className="absolute -top-6 left-0 bg-red-500 text-white text-xs px-2 py-1 rounded">
                  Bumper
                </div>
              </motion.div>
              <motion.div
                initial={{ opacity: 0, scale: 0 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ delay: 0.2 }}
                className="absolute top-1/2 right-1/4 w-24 h-16 border-2 border-orange-500 rounded-lg"
              >
                <div className="absolute -top-6 right-0 bg-orange-500 text-white text-xs px-2 py-1 rounded">
                  Door
                </div>
              </motion.div>
            </>
          )}
        </AnimatePresence>

        {/* Center Instruction */}
        {!isScanning && !scanComplete && (
          <div className="absolute inset-0 flex items-center justify-center">
              <div className="text-center">
              <Upload size={48} className="text-zinc-500 dark:text-zinc-600 mx-auto mb-2" />
              <div className="text-zinc-600 dark:text-zinc-500 text-sm">Position vehicle in frame</div>
            </div>
          </div>
        )}
      </motion.div>

      {/* Action Buttons */}
      {!scanComplete && (
        <motion.button
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
          onClick={startScan}
          disabled={isScanning}
          className="w-full bg-gradient-to-r from-cyan-500 to-blue-500 hover:from-cyan-600 hover:to-blue-600 text-white py-4 rounded-xl transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
        >
          {isScanning ? (
            <>
              <motion.div
                animate={{ rotate: 360 }}
                transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
              >
                <ScanLine size={20} />
              </motion.div>
              Scanning...
            </>
          ) : (
            <>
              <Camera size={20} />
              Start AI Scan
            </>
          )}
        </motion.button>
      )}

      {/* Results */}
      <AnimatePresence>
        {scanComplete && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: 20 }}
            className="space-y-4"
          >
            {/* Status */}
            <div className="bg-gradient-to-br from-red-50 to-zinc-100 dark:from-red-950/50 dark:to-zinc-900 border border-red-200 dark:border-red-800/30 rounded-2xl p-5 transition-colors">
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 bg-red-500/10 rounded-full flex items-center justify-center border border-red-500/30">
                    <AlertTriangle size={20} className="text-red-400" />
                  </div>
                  <div>
                    <h3 className="text-zinc-900 dark:text-white">Damage Detected</h3>
                    <p className="text-xs text-zinc-600 dark:text-zinc-400">2 areas identified</p>
                  </div>
                </div>
                <button
                  onClick={resetScan}
                  className="text-zinc-400 hover:text-white transition-colors"
                >
                  <X size={20} />
                </button>
              </div>

              {/* Damage List */}
              <div className="space-y-2 mb-4">
                {damageResults.map((damage, index) => (
                  <motion.div
                    key={index}
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: index * 0.1 }}
                    className="bg-zinc-900/50 rounded-lg p-3 border border-zinc-800"
                  >
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-zinc-900 dark:text-white text-sm">{damage.area}</span>
                      <span
                        className={`text-xs px-2 py-1 rounded ${
                          damage.severity === 'Minor'
                            ? 'bg-yellow-500/10 text-yellow-400 border border-yellow-500/30'
                            : 'bg-orange-500/10 text-orange-400 border border-orange-500/30'
                        }`}
                      >
                        {damage.severity}
                      </span>
                    </div>
                    <div className="flex items-center gap-2">
                      <div className="flex-1 h-1.5 bg-zinc-300 dark:bg-zinc-800 rounded-full overflow-hidden transition-colors">
                        <motion.div
                          initial={{ width: 0 }}
                          animate={{ width: `${damage.confidence}%` }}
                          transition={{ delay: 0.5 + index * 0.1, duration: 0.8 }}
                          className="h-full bg-gradient-to-r from-cyan-500 to-blue-500"
                        />
                      </div>
                      <span className="text-xs text-zinc-600 dark:text-zinc-400">{damage.confidence}%</span>
                    </div>
                  </motion.div>
                ))}
              </div>

              {/* Cost Estimate */}
              <div className="bg-zinc-200 dark:bg-zinc-950/50 rounded-xl p-4 border border-zinc-300 dark:border-zinc-800 transition-colors">
                <div className="text-xs text-zinc-600 dark:text-zinc-400 mb-2">Estimated Repair Cost</div>
                <div className="flex items-end gap-2">
                  <span className="text-2xl text-zinc-900 dark:text-white">
                    {estimatedCost.min.toLocaleString()}
                  </span>
                  <span className="text-zinc-600 dark:text-zinc-400 mb-1">-</span>
                  <span className="text-2xl text-zinc-900 dark:text-white">
                    {estimatedCost.max.toLocaleString()}
                  </span>
                  <span className="text-zinc-600 dark:text-zinc-400 mb-1">{estimatedCost.currency}</span>
                </div>
                <div className="text-xs text-zinc-500 dark:text-zinc-500 mt-1">Based on average market rates</div>
              </div>
            </div>

            {/* Actions */}
            <div className="grid grid-cols-2 gap-3">
              <motion.button
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
                className="bg-zinc-100 dark:bg-zinc-900 border border-zinc-300 dark:border-zinc-800 text-zinc-900 dark:text-white py-3 rounded-xl hover:bg-zinc-200 dark:hover:bg-zinc-800 transition-colors"
              >
                Save Report
              </motion.button>
              <motion.button
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
                onClick={() => {
                  const event = new CustomEvent('navigate', { detail: 'map' });
                  window.dispatchEvent(event);
                }}
                className="bg-cyan-500 hover:bg-cyan-600 text-white py-3 rounded-xl transition-colors"
              >
                Find Repair Shop
              </motion.button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
