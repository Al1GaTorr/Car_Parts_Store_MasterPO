import { PartCategory } from './types';

// Subcategories for Hover Menu
export const SUBCATEGORIES: Record<string, string[]> = {
  [PartCategory.CHEMISTRY]: ["Клеи", "Смазки", "Герметики", "Антисептики", "Электролиты", "Уход автомобиля", "Присадки в масла", "Присадки топливные"],
  [PartCategory.ACCESSORIES]: ["Ароматизаторы", "Коврики", "Чехлы", "Оплетки", "Держатели", "Органайзеры", "Шторки", "Зеркала"],
  [PartCategory.OILS]: ["Моторные масла", "Трансмиссионные", "Тормозные жидкости", "Антифризы", "Стеклоомыватели", "Промывочные"],
  [PartCategory.TOOLS]: ["Наборы инструментов", "Ключи", "Отвертки", "Домкраты", "Электроинструмент", "Съемники", "Ящики"],
  [PartCategory.WHEELS]: ["Литые диски", "Штампованные", "Шины летние", "Шины зимние", "Колпаки", "Крепеж"],
  [PartCategory.LAMPS]: ["Галогеновые", "Ксеноновые", "Светодиодные", "Накаливания", "Противотуманные", "ДХО"],
};

// Realistic images for categories
export const CATEGORY_IMAGES: Record<PartCategory, string> = {
  [PartCategory.CHEMISTRY]: "https://f.nodacdn.net/326121", 
  [PartCategory.ACCESSORIES]: "https://f.nodacdn.net/630131", 
  [PartCategory.OILS]: "https://f.nodacdn.net/628448", 
  [PartCategory.TOOLS]: "./assets/img/tols.png", 
  [PartCategory.WHEELS]: "https://images.unsplash.com/photo-1580274455191-1c62238fa333?auto=format&fit=crop&q=80&w=800", 
  [PartCategory.LAMPS]: "https://f.nodacdn.net/628948", 
  
  // Placeholders for list items (images used if they appear in search results)
  [PartCategory.SEATS]: "https://images.unsplash.com/photo-1605218427306-6354db69e563?auto=format&fit=crop&q=80&w=500", 
  [PartCategory.HEATED_SEATS]: "https://images.unsplash.com/photo-1605218427306-6354db69e563?auto=format&fit=crop&q=80&w=500",
  [PartCategory.CHILD_SEATS]: "https://images.unsplash.com/photo-1518115502127-142f15321f82?auto=format&fit=crop&q=80&w=500", // Child seat
  [PartCategory.PREHEATERS]: "https://images.unsplash.com/photo-1619642751034-765dfdf7c58e?auto=format&fit=crop&q=80&w=500", // Engine bay
  [PartCategory.CHARGERS]: "https://images.unsplash.com/photo-1620882367151-540291df1843?auto=format&fit=crop&q=80&w=500", // Charger
  [PartCategory.JUMP_CABLES]: "https://images.unsplash.com/photo-1594818379496-da1e345b0ded?auto=format&fit=crop&q=80&w=500", // Cables

  [PartCategory.CANISTERS]: "	https://f.nodacdn.net/630805", 
  [PartCategory.JACKS]: "	https://f.nodacdn.net/631716?auto=format&fit=crop&q=80&w=500", 
  [PartCategory.COMPRESSORS]: "	https://f.nodacdn.net/631740", 
  [PartCategory.FIRE_EXTINGUISHERS]: "	https://f.nodacdn.net/631747", 
  [PartCategory.EMERGENCY_KITS]: "	https://f.nodacdn.net/631767", 
  [PartCategory.BATTERIES]: "https://f.nodacdn.net/630803", 
  [PartCategory.WIPERS]: "	https://f.nodacdn.net/631811",
  [PartCategory.FIRST_AID]: "https://f.nodacdn.net/631750",
  [PartCategory.TOW_ROPES]: "https://f.nodacdn.net/631765",
  [PartCategory.WARNING_SIGNS]: "	https://f.nodacdn.net/631771"
};
