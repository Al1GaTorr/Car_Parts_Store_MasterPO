import { User } from '../types';
import { apiUrl } from './http';

const TOKEN_KEY = 'bazarpo_token';

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}
export function setToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token);
}
export function clearToken() {
  localStorage.removeItem(TOKEN_KEY);
}

export async function registerUser(payload: {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
}): Promise<{ role: 'user' | 'admin' }> {
  const res = await fetch(apiUrl('/api/auth/register'), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data?.error || 'Registration failed');
  setToken(data.token);
  return { role: data.role };
}

export async function loginUser(payload: { email: string; password: string }): Promise<{ role: 'user' | 'admin' }> {
  const res = await fetch(apiUrl('/api/auth/login'), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data?.error || 'Login failed');
  setToken(data.token);
  return { role: data.role };
}

export async function fetchMe(): Promise<User> {
  const token = getToken();
  if (!token) throw new Error('No token');
  const res = await fetch(apiUrl('/api/auth/me'), {
    headers: { Authorization: `Bearer ${token}` },
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data?.error || 'Unauthorized');
  return {
    id: data.id,
    email: data.email,
    name: `${data.firstName || ''} ${data.lastName || ''}`.trim() || data.email,
    role: data.role,
  };
}
