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
  year?: string;
  category?: string | null;
}): Promise<Part[]> {
  const usp = new URLSearchParams();
  if (params.vin) usp.set('vin', params.vin);
  if (params.issue) usp.set('issue', params.issue);
  if (params.q) usp.set('q', params.q);
  if (params.make && params.make !== 'Universal') usp.set('make', params.make);
  if (params.model && params.model !== 'All') usp.set('model', params.model);
  if (params.year && params.year !== 'All') usp.set('year', params.year);
  if (params.category) usp.set('category', params.category);

  const res = await fetch(`/api/parts?${usp.toString()}`);
  if (!res.ok) throw new Error(`API error: ${res.status}`);
  const data = (await res.json()) as { items: DBPart[] };
  return (data.items || []).map(mapDbPart);
}

let makesCache: string[] | null = null;
let makesInFlight: Promise<string[]> | null = null;
const modelsCacheByMake = new Map<string, string[]>();
const modelsInFlightByMake = new Map<string, Promise<string[]>>();
const yearsCacheByMakeModel = new Map<string, number[]>();
const yearsInFlightByMakeModel = new Map<string, Promise<number[]>>();

function normalizeStringItems(items?: unknown[]): string[] {
  return (items || [])
    .map((v) => String(v || '').trim())
    .filter(Boolean)
    .sort((a, b) => a.localeCompare(b));
}

export async function fetchVehicleMakes(): Promise<string[]> {
  if (makesCache) return makesCache;
  if (makesInFlight) return makesInFlight;

  makesInFlight = (async () => {
    const res = await fetch('/api/cars/makes');
    if (!res.ok) throw new Error(`API error: ${res.status}`);
    const data = (await res.json()) as { items?: unknown[] };
    const makes = normalizeStringItems(data.items);
    makesCache = makes;
    return makes;
  })();

  try {
    return await makesInFlight;
  } finally {
    makesInFlight = null;
  }
}

export async function fetchModelsByMake(make: string): Promise<string[]> {
  const makeName = String(make || '').trim();
  if (!makeName) return [];

  if (modelsCacheByMake.has(makeName)) {
    return modelsCacheByMake.get(makeName) || [];
  }
  if (modelsInFlightByMake.has(makeName)) {
    return modelsInFlightByMake.get(makeName) || [];
  }

  const p = (async () => {
    const qs = new URLSearchParams({ make: makeName });
    const res = await fetch(`/api/cars/models?${qs.toString()}`);
    if (!res.ok) throw new Error(`API error: ${res.status}`);
    const data = (await res.json()) as { items?: unknown[] };
    const models = normalizeStringItems(data.items);
    modelsCacheByMake.set(makeName, models);
    return models;
  })();

  modelsInFlightByMake.set(makeName, p);
  try {
    return await p;
  } finally {
    modelsInFlightByMake.delete(makeName);
  }
}

function normalizeNumberItems(items?: unknown[]): number[] {
  return (items || [])
    .map((v) => Number(v))
    .filter((v) => Number.isFinite(v))
    .sort((a, b) => a - b);
}

export async function fetchYearsByMakeModel(make: string, model: string): Promise<number[]> {
  const makeName = String(make || '').trim();
  const modelName = String(model || '').trim();
  if (!makeName || !modelName) return [];

  const key = `${makeName}::${modelName}`;
  if (yearsCacheByMakeModel.has(key)) {
    return yearsCacheByMakeModel.get(key) || [];
  }
  if (yearsInFlightByMakeModel.has(key)) {
    return yearsInFlightByMakeModel.get(key) || [];
  }

  const p = (async () => {
    const qs = new URLSearchParams({ make: makeName, model: modelName });
    const res = await fetch(`/api/cars/years?${qs.toString()}`);
    if (!res.ok) throw new Error(`API error: ${res.status}`);
    const data = (await res.json()) as { items?: unknown[] };
    const years = normalizeNumberItems(data.items);
    yearsCacheByMakeModel.set(key, years);
    return years;
  })();

  yearsInFlightByMakeModel.set(key, p);
  try {
    return await p;
  } finally {
    yearsInFlightByMakeModel.delete(key);
  }
}
