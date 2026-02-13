import { getToken } from './authService';
import { DBPart } from './apiService';

async function authedFetch(url: string, init?: RequestInit) {
  const token = getToken();
  const headers: Record<string, string> = { ...(init?.headers as any) };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  return fetch(url, { ...init, headers });
}

export async function adminFetchProducts(): Promise<DBPart[]> {
  const res = await authedFetch('/api/admin/products');
  const data = await res.json();
  if (!res.ok) throw new Error(data?.error || 'Admin products failed');
  return data.items || [];
}

export async function adminCreateProduct(p: any) {
  const res = await authedFetch('/api/admin/products', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(p),
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data?.error || 'Create failed');
  return data;
}

export async function adminUpdateProduct(id: string, patch: any) {
  const res = await authedFetch(`/api/admin/products/${id}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(patch),
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data?.error || 'Update failed');
  return data;
}

export async function adminDeleteProduct(id: string) {
  const res = await authedFetch(`/api/admin/products/${id}`, { method: 'DELETE' });
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
export const adminFetchParts = adminFetchProducts;
export const adminCreatePart = adminCreateProduct;
export const adminUpdatePart = adminUpdateProduct;
export const adminDeletePart = adminDeleteProduct;
