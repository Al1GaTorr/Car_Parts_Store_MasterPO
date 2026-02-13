export enum PartCategory {
  CHEMISTRY = "Автохимия",
  ACCESSORIES = "Автоаксессуары",
  OILS = "Масла и жидкости",
  TOOLS = "Инструменты",
  WHEELS = "Диски и шины",
  LAMPS = "Автолампы",
  
  // Group 1: Comfort & Power (From Screenshot)
  SEATS = "Накидки на сидения",
  HEATED_SEATS = "Накидки с обогревом",
  CHILD_SEATS = "Автокресла и бустеры",
  PREHEATERS = "Предпусковые обогреватели",
  CHARGERS = "Пуско зарядные устройства",
  JUMP_CABLES = "Провода пусковые",

  // Group 2: Emergency & Safety
  CANISTERS = "Канистры",
  JACKS = "Домкраты",
  COMPRESSORS = "Компрессоры",
  FIRE_EXTINGUISHERS = "Огнетушители",
  EMERGENCY_KITS = "Наборы автомобилиста",
  FIRST_AID = "Аптечки",
  TOW_ROPES = "Тросы буксировочные",
  WARNING_SIGNS = "Знаки аварийной остановки",

  // Group 3: Maintenance
  BATTERIES = "Аккумуляторы",
  WIPERS = "Щетки дворников"
}

export interface Part {
  id: string;
  name: string;
  category: PartCategory;
  price: number;
  description: string;
  image: string; // URL
  make: string; // "Toyota", "BMW", or "Universal"
  model: string; // "Camry", "X5", or "Universal"
  compatibleVins: string[]; // List of specific VINs this fits (for high precision match)
  stock: number;
  isVisible: boolean;
  partNumber: string;
}

export interface CartItem extends Part {
  quantity: number;
}

export interface User {
  id: string;
  email: string;
  phone?: string;
  name?: string;
  role: 'user' | 'admin';
  address?: string;
}

export interface Order {
  id: string;
  userId: string;
  items: CartItem[];
  totalAmount: number;
  status: 'pending' | 'processing' | 'shipped' | 'completed' | 'cancelled';
  date: string;
  shippingAddress: string;
  contactInfo: string;
}

// Navigation state to simulate routing without react-router
export enum PageView {
  HOME = 'HOME',
  CATALOG = 'CATALOG',
  CART = 'CART',
  CHECKOUT = 'CHECKOUT',
  ADMIN_DASHBOARD = 'ADMIN_DASHBOARD',
  LOGIN = 'LOGIN' // Modal usually, but can be a view state
}
