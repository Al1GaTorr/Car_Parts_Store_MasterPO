import { CartItem } from '../types';
import { getToken } from './authService';
import { apiUrl } from './http';

export async function createOrder(payload: { cart: CartItem[]; shippingAddress: string; contactInfo: string }) {
  const token = getToken();
  if (!token) throw new Error('Not logged in');
  const items = payload.cart.map((i) => ({
    sku: i.partNumber,
    name: i.name,
    price_kzt: i.price,
    qty: i.quantity,
    image: i.image,
  }));
  const res = await fetch(apiUrl('/api/orders'), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
    body: JSON.stringify({ items, shippingAddress: payload.shippingAddress, contactInfo: payload.contactInfo }),
  });
  const data = await res.json();
  if (!res.ok) {
    if (data?.error === 'insufficient stock' && Array.isArray(data?.issues)) {
      const details = data.issues
        .map((i: { sku?: string; requested?: number; available?: number }) => {
          const sku = i?.sku || 'unknown';
          const requested = Number(i?.requested || 0);
          const available = Number(i?.available || 0);
          return `${sku}: запрошено ${requested}, доступно ${available}`;
        })
        .join('; ');
      throw new Error(`Недостаточно товара на складе: ${details}`);
    }
    throw new Error(data?.error || 'Order failed');
  }
  return data;
}
