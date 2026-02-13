import { getToken } from './authService';
import { DBPart } from './apiService';

async function authedFetch(url: string, init?: RequestInit) {
  const token = getToken();
  const headers: Record<string, string> = { ...(init?.headers as any) };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  return fetch(url, { ...init, headers });
}

export async function adminFetchParts(): Promise<DBPart[]> {
  const res = await authedFetch('/api/admin/parts');
  const data = await res.json();
  if (!res.ok) throw new Error(data?.error || 'Admin parts failed');
  return data.items || [];
}

export async function adminCreatePart(p: any) {
  const res = await authedFetch('/api/admin/parts', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(p),
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data?.error || 'Create failed');
  return data;
}

export async function adminUpdatePart(id: string, patch: any) {
  const res = await authedFetch(`/api/admin/parts/${id}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(patch),
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data?.error || 'Update failed');
  return data;
}

export async function adminDeletePart(id: string) {
  const res = await authedFetch(`/api/admin/parts/${id}`, { method: 'DELETE' });
  const data = await res.json();
  if (!res.ok) throw new Error(data?.error || 'Delete failed');
  return data;
}

export async function adminFetchOrders(): Promise<any[]> {
  const res = await authedFetch('/api/admin/orders');
  const data = await res.json();
  if (!res.ok) throw new Error(data?.error || 'Admin orders failed');
  return data.items || [];
}

export async function adminUpdateOrder(id: string, patch: any) {
  const res = await authedFetch(`/api/admin/orders/${id}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(patch),
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data?.error || 'Order update failed');
  return data;
}

export async function adminDeleteOrder(id: string) {
  const res = await authedFetch(`/api/admin/orders/${id}`, { method: 'DELETE' });
  const data = await res.json();
  if (!res.ok) throw new Error(data?.error || 'Order delete failed');
  return data;
}

// Backward-compatible aliases with part naming.
export const AdminFetchParts = adminFetchParts;
export const AdminCreatePart = adminCreatePart;
export const AdminUpdatePart = adminUpdatePart;
export const AdminDeletePart = adminDeletePart;

export const AdminFetchOrders = adminFetchOrders;
export const AdminUpdateOrder = adminUpdateOrder;
export const AdminDeleteOrder = adminDeleteOrder;