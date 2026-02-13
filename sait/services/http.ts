const rawApiBase = (import.meta as any)?.env?.VITE_API_BASE_URL || '';
const API_BASE = String(rawApiBase).trim().replace(/\/+$/, '');

export function apiUrl(path: string): string {
  const normalized = path.startsWith('/') ? path : `/${path}`;
  if (!API_BASE) return normalized;
  return `${API_BASE}${normalized}`;
}

