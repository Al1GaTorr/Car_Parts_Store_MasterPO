import { useState, useEffect, useRef } from 'react';
import { Wrench, Search, Bell } from 'lucide-react';

interface Car {
  id: string;
  vin: string;
  brand: string;
  model: string;
  licensePlate: string;
  mileage: number;
  lastOilChange: number;
}

export function STOPanel() {
  const [query, setQuery] = useState('');
  const [cars, setCars] = useState<Car[]>([]);
  const [found, setFound] = useState<Car | null>(null);
  const [events, setEvents] = useState<any[]>([]);
  const esRef = useRef<EventSource | null>(null);

  useEffect(() => {
    fetch('/api/cars')
      .then((r) => r.json())
      .then((data) => setCars(data))
      .catch(() => setCars([]));

    return () => {
      if (esRef.current) {
        esRef.current.close();
      }
    };
  }, []);

  const search = () => {
    const q = query.trim().toLowerCase();
    const car = cars.find((c) => c.licensePlate?.toLowerCase() === q || c.vin?.toLowerCase() === q || c.id?.toLowerCase() === q);
    if (car) {
      setFound(car);
      attachEvents(car.id);
    } else {
      alert('Car not found');
      setFound(null);
      if (esRef.current) {
        esRef.current.close();
        esRef.current = null;
      }
    }
  };

  const attachEvents = (carId: string) => {
    if (esRef.current) {
      esRef.current.close();
    }
    const url = `/ws/cars/${carId}`;
    const es = new EventSource(url);
    es.onmessage = (ev) => {
      try {
        const data = JSON.parse(ev.data);
        setEvents((s) => [{ receivedAt: new Date().toISOString(), ...data }, ...s].slice(0, 50));
      } catch (e) {
        // ignore
      }
    };
    es.onerror = () => {
      // keep it simple for demo
    };
    esRef.current = es;
  };

  const submitRecord = async (e: any) => {
    e.preventDefault();
    if (!found) return alert('Select a car first');
    const form = new FormData(e.target);
    const payload = {
      date: form.get('date'),
      mileage: Number(form.get('mileage')) || found.mileage || 0,
      description: String(form.get('description') || ''),
      serviceName: String(form.get('serviceName') || ''),
    };

    try {
      await fetch(`/api/cars/${found.id}/records`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });
      alert('Record submitted — clients will receive an update');
      // Optionally show locally
      setEvents((s) => [{ receivedAt: new Date().toISOString(), type: 'SERVICE_RECORD_ADDED', payload }, ...s].slice(0,50));
      e.target.reset();
    } catch (err) {
      alert('Failed to submit');
    }
  };

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 pb-24 px-4 transition-colors">
      <div className="pt-6 pb-4 flex items-center gap-3">
        <div className="w-10 h-10 bg-cyan-500/10 rounded-full flex items-center justify-center border border-cyan-500/30">
          <Wrench size={20} className="text-cyan-400" />
        </div>
        <div>
          <h1 className="text-2xl text-zinc-900 dark:text-white">STO Attendant</h1>
          <p className="text-sm text-zinc-600 dark:text-zinc-400 mt-1">Enter car number (license plate / VIN / id)</p>
        </div>
      </div>

      <div className="space-y-4">
        <div className="flex gap-2">
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="License plate, VIN or car id"
            className="flex-1 p-3 rounded-xl border bg-white dark:bg-zinc-900"
          />
          <button onClick={search} className="px-4 py-3 bg-cyan-500 text-white rounded-xl flex items-center gap-2">
            <Search size={16} /> Search
          </button>
        </div>

        {found && (
          <div className="bg-zinc-100 dark:bg-zinc-900 rounded-xl p-4 border">
            <div className="flex items-center justify-between mb-2">
              <div>
                <div className="text-sm text-zinc-500">{found.brand} {found.model} • {found.year}</div>
                <div className="text-lg text-zinc-900 dark:text-white">{found.licensePlate} • {found.id}</div>
              </div>
              <div className="text-sm text-zinc-500">Mileage: {found.mileage?.toLocaleString()}</div>
            </div>

            <form onSubmit={submitRecord} className="space-y-2">
              <div className="flex gap-2">
                <input name="date" type="date" className="flex-1 p-2 rounded-md border bg-white dark:bg-zinc-900" required />
                <input name="mileage" type="number" placeholder="mileage" className="w-32 p-2 rounded-md border bg-white dark:bg-zinc-900" />
              </div>
              <input name="serviceName" placeholder="Service shop/name" className="w-full p-2 rounded-md border bg-white dark:bg-zinc-900" />
              <textarea name="description" placeholder="Description" className="w-full p-2 rounded-md border bg-white dark:bg-zinc-900" />
              <div className="flex justify-end">
                <button type="submit" className="px-4 py-2 bg-cyan-500 text-white rounded-md flex items-center gap-2">
                  <Wrench size={14} /> Submit
                </button>
              </div>
            </form>
          </div>
        )}

        <div className="mt-2">
          <h3 className="flex items-center gap-2 text-sm text-zinc-600 dark:text-zinc-400 mb-2"><Bell size={16} /> Recent events</h3>
          <div className="space-y-2">
            {events.length === 0 && <div className="text-xs text-zinc-500">No events yet</div>}
            {events.map((ev, idx) => (
              <div key={idx} className="p-3 rounded-lg border bg-white dark:bg-zinc-900 text-sm">
                <div className="text-xs text-zinc-500">{ev.receivedAt}</div>
                <div className="font-medium">{ev.type}</div>
                <pre className="text-xs mt-1 whitespace-pre-wrap">{JSON.stringify(ev.payload || ev, null, 2)}</pre>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
