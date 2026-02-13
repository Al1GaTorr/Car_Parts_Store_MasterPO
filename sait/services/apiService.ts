import { Part, PartCategory } from '../types';

export type VehicleFit = { make: string; model: string; year_from: number; year_to: number };

export type DBPart = {
  _id?: string;
  sku: string;
  name: string;
  category: string; // Russian category label (PartCategory)
  type?: string;
  brand?: string;
  price_kzt: number;
  currency?: string;
  stock_qty: number;
  is_visible: boolean;
  images?: string[];
  compatibility?:
    | { type: 'vin'; vins: string[] }
    | { type: 'vehicle'; vehicles: VehicleFit[] }
    | { type: 'universal' };
  recommend_for_issue_codes?: string[];
  created_at?: string;
};

function toPartCategory(cat: string): PartCategory {
  const values = Object.values(PartCategory) as string[];
  if (values.includes(cat)) return cat as PartCategory;
  return PartCategory.ACCESSORIES;
}

function mapDbPart(p: DBPart): Part {
  const comp = p.compatibility;
  let make = 'Universal';
  let model = 'Universal';
  let compatibleVins: string[] = [];
  if (comp?.type === 'vin') compatibleVins = comp.vins || [];
  if (comp?.type === 'vehicle' && comp.vehicles && comp.vehicles.length > 0) {
    make = comp.vehicles[0].make;
    model = comp.vehicles[0].model;
  }
  const image = p.images && p.images.length > 0 ? p.images[0] : '';
  const category = toPartCategory(p.category || '');

  const descBits = [p.brand ? `Бренд: ${p.brand}` : '', p.type ? `Тип: ${p.type}` : '']
    .filter(Boolean)
    .join(' • ');

  return {
    id: p._id || p.sku,
    name: p.name,
    category,
    price: p.price_kzt,
    description: descBits || 'Запчасть для автомобиля',
    image,
    make,
    model,
    compatibleVins,
    stock: p.stock_qty,
    isVisible: p.is_visible,
    partNumber: p.sku,
  };
}

export async function fetchParts(params: {
  vin?: string;
  issue?: string;
  q?: string;
  make?: string;
  model?: string;
  category?: string | null;
}): Promise<Part[]> {
  const usp = new URLSearchParams();
  if (params.vin) usp.set('vin', params.vin);
  if (params.issue) usp.set('issue', params.issue);
  if (params.q) usp.set('q', params.q);
  if (params.make && params.make !== 'Universal') usp.set('make', params.make);
  if (params.model && params.model !== 'All') usp.set('model', params.model);
  if (params.category) usp.set('category', params.category);

  const res = await fetch(`/api/parts?${usp.toString()}`);
  if (!res.ok) throw new Error(`API error: ${res.status}`);
  const data = (await res.json()) as { items: DBPart[] };
  return (data.items || []).map(mapDbPart);
}

export async function fetchVehicleOptions(): Promise<{
  makes: string[];
  modelsByMake: Record<string, string[]>;
}> {
  const res = await fetch('/api/parts');
  if (!res.ok) throw new Error(`API error: ${res.status}`);
  const data = (await res.json()) as { items: DBPart[] };
  const items = data.items || [];

  const modelsByMake = new Map<string, Set<string>>();
  for (const p of items) {
    const comp = p.compatibility;
    if (comp?.type !== 'vehicle' || !Array.isArray(comp.vehicles)) continue;
    for (const v of comp.vehicles) {
      const make = (v.make || '').trim();
      const model = (v.model || '').trim();
      if (!make || !model) continue;
      if (!modelsByMake.has(make)) modelsByMake.set(make, new Set<string>());
      modelsByMake.get(make)?.add(model);
    }
  }

  const makes = Array.from(modelsByMake.keys()).sort((a, b) => a.localeCompare(b));
  const out: Record<string, string[]> = {};
  for (const m of makes) {
    out[m] = Array.from(modelsByMake.get(m) || []).sort((a, b) => a.localeCompare(b));
  }
  return { makes, modelsByMake: out };
}
