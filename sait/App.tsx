import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { CATEGORY_IMAGES, SUBCATEGORIES } from './constants';
import { Part, User, CartItem, Order, PageView, PartCategory } from './types';
import { fetchParts, fetchVehicleMakes, fetchModelsByMake, fetchYearsByMakeModel } from './services/apiService';
import { loginUser, registerUser, fetchMe, clearToken, getToken } from './services/authService';
import { AdminFetchParts, AdminCreatePart, AdminUpdatePart, AdminDeletePart, AdminFetchOrders, AdminUpdateOrder, AdminDeleteOrder } from './services/adminService';
import { createOrder } from './services/orderService';
import { 
  ShoppingCart, 
  Search, 
  Menu, 
  User as UserIcon, 
  X, 
  Trash2, 
  Plus, 
  Minus,
  Settings,
  Package,
  LogOut,
  ChevronRight,
  Filter,
  ArrowRight
} from 'lucide-react';

// --- Components ---

// 1. Navbar - Matches the "bazarPO" design
const Navbar: React.FC<{
  cartCount: number;
  user: User | null;
  onNavigate: (page: PageView) => void;
  onLogout: () => void;
  onLoginClick: () => void;
  searchQuery: string;
  setSearchQuery: (q: string) => void;
}> = ({ cartCount, user, onNavigate, onLogout, onLoginClick, searchQuery, setSearchQuery }) => {
  return (
    <nav className="bg-white text-slate-900 sticky top-0 z-50 shadow-sm border-b border-slate-200 py-3">
      <div className="container mx-auto px-4 flex flex-wrap items-center justify-between gap-4">
        
        {/* Logo Section */}
        <div 
          className="flex items-center gap-3 cursor-pointer select-none" 
          onClick={() => onNavigate(PageView.HOME)}
        >
           <div className="relative">
             {/* Logo Icon simulation */}
             <div className="w-10 h-10 bg-slate-800 transform -skew-x-12 flex items-center justify-center">
                <div className="w-1 h-6 bg-cyan-400 absolute left-2"></div>
                <div className="w-1 h-6 bg-white absolute left-4"></div>
             </div>
           </div>
           <div className="flex flex-col">
             <span className="font-black text-2xl tracking-tighter leading-none text-slate-900">bazarPO</span>
             <span className="text-[10px] font-bold text-cyan-600 bg-cyan-50 px-1 uppercase tracking-widest w-fit">Внимание к деталям</span>
           </div>
        </div>

        {/* Catalog Button */}
        <button 
          onClick={() => onNavigate(PageView.CATALOG)}
          className="bg-slate-800 text-white px-5 py-2.5 rounded hover:bg-slate-700 transition-colors flex items-center gap-2 font-bold text-sm uppercase"
        >
          <Menu size={20} />
          Каталоги
        </button>

        {/* Search Bar - Wide */}
        <div className="flex-1 min-w-[300px] relative">
           <input 
             type="text" 
             placeholder="Введите код запчасти или VIN номер автомобиля"
             className="w-full bg-slate-100 border border-slate-200 rounded px-4 py-2.5 pl-4 pr-12 focus:outline-none focus:ring-2 focus:ring-slate-300 text-sm"
             value={searchQuery}
             onChange={(e) => setSearchQuery(e.target.value)}
           />
           <button 
             onClick={() => onNavigate(PageView.CATALOG)}
             className="absolute right-2 top-1/2 -translate-y-1/2 text-slate-500 hover:text-slate-800"
           >
             <Search size={20} />
           </button>
        </div>

        {/* Right Actions */}
        <div className="flex items-center gap-4">
          {user?.role === 'admin' && (
            <button
              onClick={() => onNavigate(PageView.ADMIN_DASHBOARD)}
              className="hidden md:flex items-center gap-2 px-3 py-2 rounded border border-slate-200 text-slate-800 font-bold text-xs uppercase hover:bg-slate-50"
            >
              <Settings size={16} />
              Админ
            </button>
          )}
          <div 
            onClick={user ? undefined : onLoginClick}
            className="flex items-center gap-2 cursor-pointer hover:bg-slate-50 p-2 rounded transition-colors"
          >
             <div className="bg-slate-800 text-white p-2 rounded-full">
               <UserIcon size={18} />
             </div>
             <div className="hidden lg:flex flex-col text-sm">
                {user ? (
                  <>
                    <span className="font-bold">{user.name || 'User'}</span>
                    <span className="text-xs text-slate-500" onClick={(e) => { e.stopPropagation(); onLogout(); }}>Выйти</span>
                  </>
                ) : (
                  <>
                    <span className="font-bold">Войти</span>
                    <span className="text-xs text-slate-500">Регистрация</span>
                  </>
                )}
             </div>
          </div>

          <div 
            onClick={() => onNavigate(PageView.CART)}
            className="flex items-center gap-2 cursor-pointer hover:bg-slate-50 p-2 rounded transition-colors"
          >
             <div className="bg-slate-800 text-white p-2 rounded-full relative">
               <ShoppingCart size={18} />
               {cartCount > 0 && (
                 <span className="absolute -top-1 -right-1 bg-cyan-500 border-2 border-white text-white text-[10px] font-bold h-4 w-4 flex items-center justify-center rounded-full">
                   {cartCount}
                 </span>
               )}
             </div>
          </div>
        </div>

      </div>
    </nav>
  );
};

// 2. Hero Section - Exact Design
const HeroSection: React.FC<{
  vin: string;
  setVin: (v: string) => void;
  onSearch: () => void;
  make: string;
  setMake: (m: string) => void;
  model: string;
  setModel: (m: string) => void;
  year: string;
  setYear: (y: string) => void;
  makes: string[];
  modelsByMake: Record<string, string[]>;
  yearsByMakeModel: Record<string, number[]>;
}> = ({ vin, setVin, onSearch, make, setMake, model, setModel, year, setYear, makes, modelsByMake, yearsByMakeModel }) => {
  const availableModels = modelsByMake[make] || [];
  const availableYears = yearsByMakeModel[`${make}::${model}`] || [];

  return (
    <div className="bg-white pt-6 pb-2 border-b border-slate-200">
      <div className="container mx-auto px-4">
        <div className="grid lg:grid-cols-2 gap-6">
          
          {/* Block 1: VIN Search */}
          <div className="bg-slate-100 p-8 rounded-md relative overflow-hidden">
            <h2 className="text-2xl font-bold text-slate-800 mb-1">Поиск запчастей по VIN — номеру</h2>
            <p className="text-slate-500 text-sm mb-6">VIN автомобиля является самым надежным и точным идентификатором</p>
            
            <div className="flex gap-0">
               <input 
                 type="text" 
                 placeholder="Введите VIN" 
                 className="flex-1 border border-slate-300 rounded-l px-4 py-3 focus:outline-none focus:ring-1 focus:ring-slate-400 uppercase text-sm"
                 value={vin}
                 onChange={(e) => setVin(e.target.value.toUpperCase())}
               />
               <button 
                 onClick={onSearch}
                 className="bg-slate-800 text-white px-8 py-3 rounded-r font-bold text-sm hover:bg-slate-700 transition-colors uppercase tracking-wider"
               >
                 Поиск
               </button>
            </div>
            <div className="flex justify-between mt-3 text-xs">
              <span className="text-slate-400">Например: VIN XW8AN2NE3JH035743</span>
              <a href="#" className="underline font-bold text-slate-800">Где взять VIN-код?</a>
            </div>
          </div>

          {/* Block 2: Parameters */}
          <div className="bg-slate-100 p-8 rounded-md relative overflow-hidden">
             <h2 className="text-2xl font-bold text-slate-800 mb-1">Выберите модель по параметрам</h2>
             <p className="text-slate-500 text-sm mb-6">Если не помните VIN-номер, то заполните параметры ниже</p>
             
             <div className="space-y-4">
               <div className="grid grid-cols-4 items-center gap-4">
                  <label className="font-bold text-sm text-slate-700 col-span-1">Марка</label>
                  <select 
                    className="col-span-3 border border-slate-300 bg-white rounded px-3 py-2 w-full text-sm outline-none focus:border-slate-500"
                    value={make}
                    onChange={(e) => { setMake(e.target.value); setModel('All'); setYear('All'); }}
                  >
                    <option value="Universal">Не выбрано</option>
                    {makes.map(m => (
                      <option key={m} value={m}>{m}</option>
                    ))}
                  </select>
               </div>
               <div className="grid grid-cols-4 items-center gap-4">
                  <label className="font-bold text-sm text-slate-700 col-span-1">Модель</label>
                  <select 
                    className="col-span-3 border border-slate-300 bg-white rounded px-3 py-2 w-full text-sm outline-none focus:border-slate-500"
                    value={model}
                    onChange={(e) => { setModel(e.target.value); setYear('All'); }}
                    disabled={make === 'Universal'}
                  >
                    <option value="All">Не выбрано</option>
                    {availableModels.map((m: string) => (
                      <option key={m} value={m}>{m}</option>
                    ))}
                  </select>
               </div>
               <div className="grid grid-cols-4 items-center gap-4">
                  <label className="font-bold text-sm text-slate-700 col-span-1">Год</label>
                  <select
                    className="col-span-3 border border-slate-300 bg-white rounded px-3 py-2 w-full text-sm outline-none focus:border-slate-500"
                    value={year}
                    onChange={(e) => setYear(e.target.value)}
                    disabled={make === 'Universal' || model === 'All'}
                  >
                    <option value="All">Не выбрано</option>
                    {availableYears.map((y) => (
                      <option key={y} value={String(y)}>{y}</option>
                    ))}
                  </select>
               </div>
               
               <div className="flex justify-end mt-2">
                 <button onClick={onSearch} className="text-slate-800 font-bold text-sm hover:underline flex items-center gap-1">
                   Перейти в каталог <ArrowRight size={14}/>
                 </button>
               </div>
             </div>
          </div>

        </div>

      </div>
    </div>
  );
};

// 3. Category Grid - Updated with bottom popups and new SmallCards
const CategoryGrid: React.FC<{
  onSelectCategory: (cat: string) => void;
}> = ({ onSelectCategory }) => {
  
  // Big Card with Popup
  const Card: React.FC<{ 
    category: PartCategory; 
    className?: string; 
    onClick: () => void; 
    bg?: string;
    customImage?: string;
    popupDirection?: 'left' | 'right' | 'bottom';
    children?: React.ReactNode;
  }> = ({ category, className, onClick, bg = "bg-slate-100", customImage, popupDirection = 'right', children }) => {
    const img = customImage || CATEGORY_IMAGES[category];
    const subs = SUBCATEGORIES[category] || [];

    // Calculate Popup Classes based on direction
    let popupClasses = "";
    let arrowClasses = "";
    
    if (popupDirection === 'right') {
        popupClasses = "left-[95%] top-0 min-h-[100%] group-hover:left-[102%]";
        arrowClasses = "top-8 -left-2 transform rotate-45";
    } else if (popupDirection === 'left') {
        popupClasses = "right-[95%] top-0 min-h-[100%] group-hover:right-[102%]";
        arrowClasses = "top-8 -right-2 transform rotate-45";
    } else if (popupDirection === 'bottom') {
        popupClasses = "top-[95%] left-0 w-full min-h-0 group-hover:top-[102%]";
        arrowClasses = "-top-2 left-1/2 transform -translate-x-1/2 rotate-45";
    }

    return (
      <div 
        className={`relative group h-full hover:z-50 ${className}`} 
      >
        {/* Visual Card Content */}
        <div 
          onClick={onClick}
          className={`${bg} rounded-md p-5 relative overflow-hidden cursor-pointer shadow-sm group-hover:shadow-xl transition-all h-full z-10 min-h-[240px]`}
        >
          {/* Animated Dark Circle Background - Smaller circle (scale 1.8) */}
          <div className="absolute -top-12 -left-12 w-48 h-48 bg-slate-900 rounded-full transition-transform duration-500 scale-0 group-hover:scale-[1.8] z-0 origin-center pointer-events-none"></div>

          {/* Title */}
          <h3 className="text-lg font-bold text-slate-800 relative z-10 transition-colors duration-300 group-hover:text-white">{category}</h3>
          
          {/* Image */}
          <img 
            src={img} 
            alt={category} 
            className="absolute right-0 bottom-0 w-32 h-32 object-contain transition-transform duration-500 group-hover:scale-110 z-0" 
          />
          
          {children}
        </div>

        {/* Subcategories List Overlay */}
        {subs.length > 0 && (
          <div 
            className={`
              absolute bg-slate-900 text-white p-5 rounded-lg shadow-2xl 
              opacity-0 invisible group-hover:opacity-100 group-hover:visible 
              transition-all duration-300 z-50 flex flex-col justify-center
              ${popupDirection === 'bottom' ? 'w-full' : 'w-64'}
              ${popupClasses}
            `}
          >
             {/* Decorative Arrow */}
             <div className={`absolute w-4 h-4 bg-slate-900 ${arrowClasses}`}></div>

             <ul className={`text-sm font-medium ${popupDirection === 'bottom' ? 'grid grid-cols-2 gap-2' : 'space-y-2'}`}>
               {subs.slice(0, 10).map(s => (
                 <li key={s} className="hover:text-cyan-400 transition-colors cursor-pointer flex items-center gap-2">
                   <span className="w-1.5 h-1.5 bg-cyan-500 rounded-full"></span> {s}
                 </li>
               ))}
               {subs.length > 10 && <li className="text-cyan-500 italic text-[10px] mt-1 col-span-full">и еще {subs.length - 10}...</li>}
             </ul>
          </div>
        )}
      </div>
    );
  };

  // New SmallCard for the flat categories - Compact Rectangle
  const SmallCard: React.FC<{ category: PartCategory }> = ({ category }) => (
    <div 
      onClick={() => onSelectCategory(category)}
      className="bg-slate-100 rounded-lg h-28 relative overflow-hidden cursor-pointer transition-all duration-300 hover:shadow-md hover:scale-[1.03] group"
    >
       <h3 className="font-bold text-xs sm:text-sm text-slate-800 absolute top-3 left-3 z-10 leading-tight max-w-[65%]">{category}</h3>
       <img 
         src={CATEGORY_IMAGES[category]} 
         alt={category}
         className="absolute bottom-0 right-0 w-20 h-20 object-contain transition-transform duration-300 group-hover:scale-105"
       />
    </div>
  );

  // List Block matching the screenshot - Distributed items to match height
  const ListBlock: React.FC = () => {
      const items = [
        PartCategory.SEATS,
        PartCategory.HEATED_SEATS,
        PartCategory.CHILD_SEATS,
        PartCategory.PREHEATERS,
        PartCategory.CHARGERS,
        PartCategory.JUMP_CABLES
      ];

      return (
        <div className="bg-[#1f2937] rounded-lg py-2 px-1 h-full flex flex-col shadow-md">
            <ul className="flex-1 flex flex-col justify-between">
                {items.map(item => (
                    <li 
                        key={item}
                        onClick={() => onSelectCategory(item)}
                        className="text-white font-bold text-sm px-4 py-2 rounded cursor-pointer transition-colors hover:bg-slate-700 hover:text-cyan-400 flex items-center h-full"
                    >
                        {item}
                    </li>
                ))}
            </ul>
        </div>
      );
  };

  // Remaining categories to render as Small Cards
  const blockCategories = [
    PartCategory.CANISTERS, PartCategory.JACKS, PartCategory.COMPRESSORS, PartCategory.FIRE_EXTINGUISHERS,
    PartCategory.EMERGENCY_KITS, PartCategory.FIRST_AID, PartCategory.TOW_ROPES, PartCategory.WARNING_SIGNS,
    PartCategory.BATTERIES, PartCategory.WIPERS
  ];

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Top Grid */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
        <Card 
          category={PartCategory.CHEMISTRY}
          onClick={() => onSelectCategory(PartCategory.CHEMISTRY)}
        />
        <Card 
          category={PartCategory.ACCESSORIES}
          onClick={() => onSelectCategory(PartCategory.ACCESSORIES)}
        />
        
        {/* Oils Card */}
        <Card 
          category={PartCategory.OILS}
          onClick={() => onSelectCategory(PartCategory.OILS)}
          bg="bg-slate-200"
        >
           {/* Volume Tags */}
           <div className="absolute bottom-4 left-4 right-4 flex justify-center gap-2 z-10 group-hover:opacity-0 transition-opacity">
             {['1л', '4л', '20л', '208л'].map(v => (
               <span key={v} className="bg-white/80 border border-slate-300 text-[10px] px-2 py-1 rounded-full">{v}</span>
             ))}
          </div>
        </Card>

        {/* Tools - Popup to Bottom */}
        <Card 
          category={PartCategory.TOOLS}
          onClick={() => onSelectCategory(PartCategory.TOOLS)}
          popupDirection="bottom"
        />
      </div>

      {/* Second Row */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
         {/* Wheels - Large Animated Card */}
         <div 
            onClick={() => onSelectCategory(PartCategory.WHEELS)}
            className="md:col-span-2 relative group h-full min-h-[240px] cursor-pointer overflow-hidden rounded-md bg-slate-900 transition-colors duration-500 hover:bg-slate-200"
         >
             {/* Dark Curve Bottom Left (Visible on Hover) */}
             <div className="absolute -bottom-32 -left-32 w-80 h-80 bg-slate-900 rounded-full transition-transform duration-500 scale-0 group-hover:scale-100 z-0"></div>

             {/* Title */}
             <h3 className="absolute left-8 top-1/2 -translate-y-1/2 text-3xl font-bold text-white z-20 transition-all duration-500 group-hover:text-slate-900 group-hover:top-8 group-hover:translate-y-0 group-hover:left-6">
                Диски и шины
             </h3>

             {/* Wheel Image */}
             <img 
               src="  \assets\img\koleso.png" 
               alt="Wheel" 
               className="absolute top-1/2 -translate-y-1/2 right-12 w-48 h-48 object-contain z-10 transition-all duration-700 ease-in-out group-hover:right-1/2 group-hover:translate-x-1/2 group-hover:rotate-[360deg] group-hover:scale-110" 
             />

             {/* Right Menu Buttons */}
             <div className="absolute right-8 top-0 h-full flex flex-col justify-center gap-3 z-20 opacity-0 translate-x-10 transition-all duration-500 delay-100 group-hover:opacity-100 group-hover:translate-x-0">
                 {["Диски", "Шины", "Мотошины", "Грузовые шины"].map(item => (
                     <button key={item} className="px-6 py-2 border-2 border-slate-900 text-slate-900 font-bold rounded-full hover:bg-slate-900 hover:text-white transition-colors bg-transparent uppercase text-xs tracking-wider min-w-[140px]">
                         {item}
                     </button>
                 ))}
             </div>
         </div>

         {/* Lamps - Popup to Bottom */}
         <Card 
            category={PartCategory.LAMPS}
            onClick={() => onSelectCategory(PartCategory.LAMPS)}
            popupDirection="bottom"
         />
      </div>

      {/* Third Section - Flex Layout for ListBlock and 5-col Grid */}
      <h3 className="font-bold text-xl mb-4 mt-8 text-slate-800">Популярные категории</h3>
      
      <div className="flex flex-col lg:flex-row gap-4 items-stretch">
         
         {/* List Block - Seats etc. */}
         <div className="lg:w-72 flex-shrink-0">
            <ListBlock />
         </div>

         {/* Small Cards Grid - Canisters etc. (5 in a row on large screens) */}
         <div className="flex-1">
             <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4 h-full">
                 {blockCategories.map(cat => (
                     <SmallCard key={cat} category={cat} />
                 ))}
             </div>
         </div>
      </div>

    </div>
  );
};

// 4. Catalog View (Refined)
const Catalog: React.FC<{
  parts: Part[];
  onAddToCart: (part: Part) => void;
  selectedCategory: string | null;
  onSelectCategory: (cat: string | null) => void;
  searchQuery: string;
}> = ({ parts, onAddToCart, selectedCategory, onSelectCategory, searchQuery }) => {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex flex-col md:flex-row gap-8">
        {/* Sidebar */}
        <div className="w-full md:w-64 flex-shrink-0">
          <div className="bg-white border rounded shadow-sm p-4 sticky top-24">
            <div className="flex items-center gap-2 mb-4 text-lg font-bold border-b pb-2">
              <Filter size={20} /> Фильтры
            </div>
            
            <h3 className="font-semibold mb-2 text-sm text-slate-500 uppercase">Категории</h3>
            <ul className="space-y-1">
              <li 
                className={`cursor-pointer px-3 py-2 rounded text-sm transition-colors ${selectedCategory === null ? 'bg-slate-800 text-white font-bold' : 'hover:bg-slate-100'}`}
                onClick={() => onSelectCategory(null)}
              >
                Все категории
              </li>
              {Object.values(PartCategory).map(cat => (
                <li 
                  key={cat}
                  className={`cursor-pointer px-3 py-2 rounded text-sm transition-colors flex justify-between items-center ${selectedCategory === cat ? 'bg-slate-800 text-white font-bold' : 'hover:bg-slate-100'}`}
                  onClick={() => onSelectCategory(cat)}
                >
                  {cat}
                </li>
              ))}
            </ul>
          </div>
        </div>

        {/* Grid */}
        <div className="flex-1">
          <h2 className="text-2xl font-bold mb-6">
            {selectedCategory || (searchQuery ? `Результаты поиска: "${searchQuery}"` : 'Каталог товаров')} 
            <span className="text-sm font-normal text-slate-500 ml-2">({parts.length} товаров)</span>
          </h2>

          {parts.length === 0 ? (
             <div className="text-center py-20 bg-white border rounded">
               <Package size={48} className="mx-auto text-slate-300 mb-4" />
               <p className="text-slate-500 text-lg">Товары не найдены.</p>
               <button onClick={() => onSelectCategory(null)} className="text-cyan-600 font-bold mt-2 hover:underline">Сбросить фильтры</button>
             </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {parts.map(part => (
                <div key={part.id} className="bg-white border rounded hover:shadow-xl transition-all duration-300 flex flex-col">
                  <div className="h-48 overflow-hidden relative bg-white p-4 flex items-center justify-center">
                    <img src={part.image} alt={part.name} className="max-h-full max-w-full object-contain" />
                    {part.stock === 0 && (
                      <span className="absolute inset-0 bg-white/80 flex items-center justify-center text-slate-900 font-bold border-2 border-slate-900 m-4">Нет в наличии</span>
                    )}
                  </div>
                  <div className="p-4 flex-1 flex flex-col">
                    <div className="text-xs text-slate-500 mb-1">{part.make} {part.model !== 'Universal' ? part.model : ''}</div>
                    <h3 className="font-bold text-slate-800 mb-2 leading-tight line-clamp-2 min-h-[2.5rem]">{part.name}</h3>
                    <div className="text-xs text-slate-400 mb-4">{part.partNumber}</div>
                    <div className="mt-auto flex justify-between items-center">
                      <span className="text-xl font-black text-slate-900">{part.price} ₸</span>
                      <button 
                        disabled={part.stock === 0}
                        onClick={() => onAddToCart(part)}
                        className={`px-4 py-2 rounded text-sm font-bold transition-colors ${part.stock > 0 ? 'bg-cyan-500 text-white hover:bg-cyan-600' : 'bg-slate-200 text-slate-400 cursor-not-allowed'}`}
                      >
                        В корзину
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

// 5. Cart View
const CartView: React.FC<{
  cart: CartItem[];
  onUpdateQty: (id: string, delta: number) => void;
  onRemove: (id: string) => void;
  onCheckout: () => void;
}> = ({ cart, onUpdateQty, onRemove, onCheckout }) => {
  const total = cart.reduce((sum, item) => sum + (item.price * item.quantity), 0);

  if (cart.length === 0) {
    return (
      <div className="container mx-auto px-4 py-20 text-center">
        <ShoppingCart size={64} className="mx-auto text-slate-300 mb-6" />
        <h2 className="text-2xl font-bold text-slate-700 mb-2">Корзина пуста</h2>
        <p className="text-slate-500 mb-8">Вы еще ничего не добавили.</p>
        <button onClick={onCheckout} className="text-cyan-600 font-bold hover:underline">Перейти в каталог</button>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Корзина</h1>
      <div className="flex flex-col lg:flex-row gap-8">
        <div className="flex-1 bg-white rounded shadow-sm border overflow-hidden">
          {cart.map(item => (
            <div key={item.id} className="flex items-center gap-4 p-4 border-b last:border-b-0">
              <img src={item.image} alt={item.name} className="w-20 h-20 object-contain p-2 border rounded bg-slate-50" />
              <div className="flex-1">
                <h3 className="font-bold text-slate-800">{item.name}</h3>
                <p className="text-sm text-slate-500">{item.partNumber}</p>
              </div>
              <div className="flex items-center gap-3">
                <button onClick={() => onUpdateQty(item.id, -1)} className="p-1 hover:bg-slate-100 rounded text-slate-500"><Minus size={16} /></button>
                <span className="font-bold w-6 text-center">{item.quantity}</span>
                <button onClick={() => onUpdateQty(item.id, 1)} className="p-1 hover:bg-slate-100 rounded text-slate-500"><Plus size={16} /></button>
              </div>
              <div className="text-right min-w-[100px]">
                <div className="font-bold">{item.price * item.quantity} ₸</div>
                <div className="text-xs text-slate-400">{item.price} ₸ / шт</div>
              </div>
              <button onClick={() => onRemove(item.id)} className="text-red-400 hover:text-red-600 p-2"><Trash2 size={18} /></button>
            </div>
          ))}
        </div>
        
        <div className="w-full lg:w-80">
          <div className="bg-white rounded shadow-sm border p-6 sticky top-24">
            <h3 className="text-lg font-bold mb-4">Сумма заказа</h3>
            <div className="flex justify-between mb-2 text-slate-600">
              <span>Товары</span>
              <span>{total} ₸</span>
            </div>
            <div className="border-t pt-4 flex justify-between font-bold text-xl mb-6">
              <span>Итого</span>
              <span>{total} ₸</span>
            </div>
            <button 
              onClick={onCheckout}
              className="w-full bg-slate-800 text-white font-bold py-3 rounded shadow hover:bg-slate-700 transition-colors flex items-center justify-center gap-2"
            >
              Оформить заказ <ChevronRight size={18} />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

// 6. Checkout View
const CheckoutView: React.FC<{
  user: User;
  cart: CartItem[];
  onPlaceOrder: (address: string, contact: string) => void;
  onCancel: () => void;
}> = ({ user, cart, onPlaceOrder, onCancel }) => {
  const [address, setAddress] = useState(user.address || "");
  const [contact, setContact] = useState(user.phone || user.email);
  const total = cart.reduce((sum, item) => sum + (item.price * item.quantity), 0);

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      <div className="bg-white rounded shadow-sm border p-8">
        <h2 className="text-2xl font-bold mb-6 flex items-center gap-2">
           Оформление заказа
        </h2>
        
        <div className="mb-6 bg-slate-50 p-4 rounded border">
          <h3 className="font-bold text-sm text-slate-500 uppercase mb-2">Ваш заказ</h3>
          <ul className="text-sm space-y-1">
            {cart.map(item => (
              <li key={item.id} className="flex justify-between">
                <span>{item.quantity}x {item.name}</span>
                <span className="font-medium">{item.price * item.quantity} ₸</span>
              </li>
            ))}
          </ul>
          <div className="border-t mt-3 pt-2 flex justify-between font-bold">
            <span>К оплате</span>
            <span>{total} ₸</span>
          </div>
        </div>

        <div className="space-y-4">
          <div>
            <label className="block text-sm font-bold text-slate-700 mb-1">Контактные данные (Телефон)</label>
            <input 
              type="text" 
              value={contact}
              onChange={(e) => setContact(e.target.value)}
              className="w-full border border-slate-300 rounded px-3 py-2 focus:ring-2 focus:ring-slate-500 outline-none"
              placeholder="+7 (999) 000-00-00"
            />
          </div>
          <div>
            <label className="block text-sm font-bold text-slate-700 mb-1">Адрес доставки</label>
            <textarea 
              value={address}
              onChange={(e) => setAddress(e.target.value)}
              className="w-full border border-slate-300 rounded px-3 py-2 focus:ring-2 focus:ring-slate-500 outline-none h-24"
              placeholder="Город, улица, дом, квартира..."
            ></textarea>
          </div>
        </div>

        <div className="flex gap-4 mt-8">
          <button onClick={onCancel} className="flex-1 py-3 border border-slate-300 rounded font-bold text-slate-600 hover:bg-slate-50">Назад</button>
          <button 
            onClick={() => onPlaceOrder(address, contact)}
            disabled={!address || !contact}
            className={`flex-1 py-3 rounded font-bold text-white transition-colors ${!address || !contact ? 'bg-slate-300 cursor-not-allowed' : 'bg-green-600 hover:bg-green-700'}`}
          >
            Подтвердить
          </button>
        </div>
      </div>
    </div>
  );
};

// 8. Admin Panel (Briefly adapted for language)


// 7. Auth View (Sign Up / Log In)
const AuthView: React.FC<{
  onDone: (user: User) => void;
}> = ({ onDone }) => {
  const [tab, setTab] = useState<'signup'|'login'>('signup');
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [err, setErr] = useState<string>('');

  const submit = async () => {
    setErr('');
    try {
      if (tab === 'signup') {
        await registerUser({ firstName, lastName, email, password });
      } else {
        await loginUser({ email, password });
      }
      const me = await fetchMe();
      onDone(me);
    } catch (e:any) {
      setErr(e?.message || 'Ошибка');
    }
  };

  return (
    <div className="min-h-[calc(100vh-80px)] flex items-center justify-center px-4 py-10 bg-slate-50">
      <div className="w-full max-w-2xl bg-slate-800 text-white rounded shadow overflow-hidden">
        <div className="grid grid-cols-2">
          <button onClick={()=>setTab('signup')} className={`${tab==='signup'?'bg-emerald-500':'bg-slate-700'} py-4 font-bold`}>Sign Up</button>
          <button onClick={()=>setTab('login')} className={`${tab==='login'?'bg-emerald-500':'bg-slate-700'} py-4 font-bold`}>Log In</button>
        </div>

        <div className="p-10">
          <h2 className="text-4xl font-light text-center mb-8">{tab==='signup'?'Sign Up for Free':'Log In'}</h2>

          {tab==='signup' && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <input className="bg-slate-900/30 border border-slate-500 rounded px-4 py-3 outline-none focus:border-emerald-400" placeholder="First Name *" value={firstName} onChange={e=>setFirstName(e.target.value)} />
              <input className="bg-slate-900/30 border border-slate-500 rounded px-4 py-3 outline-none focus:border-emerald-400" placeholder="Last Name *" value={lastName} onChange={e=>setLastName(e.target.value)} />
            </div>
          )}

          <div className="space-y-4">
            <input className="w-full bg-slate-900/30 border border-slate-500 rounded px-4 py-3 outline-none focus:border-emerald-400" placeholder="Email Address *" value={email} onChange={e=>setEmail(e.target.value)} />
            <input type="password" className="w-full bg-slate-900/30 border border-slate-500 rounded px-4 py-3 outline-none focus:border-emerald-400" placeholder="Set A Password *" value={password} onChange={e=>setPassword(e.target.value)} />
          </div>

          {err && <div className="text-red-300 text-sm mt-4">{err}</div>}

          <button onClick={submit} className="w-full mt-8 bg-emerald-500 hover:bg-emerald-600 transition-colors text-2xl font-black py-4 rounded">
            GET STARTED
          </button>

          <div className="text-center text-slate-300 text-xs mt-4">
            Админ входит по данным из .env (ADMIN_EMAIL / ADMIN_PASSWORD)
          </div>
        </div>
      </div>
    </div>
  );
};


const AdminPanel: React.FC<{
  onExit: () => void;
}> = ({ onExit }) => {
  const [activeTab, setActiveTab] = useState<'analytics' | 'parts' | 'orders'>('analytics');
  const [dbParts, setDbParts] = useState<any[]>([]);
  const [dbOrders, setDbOrders] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [err, setErr] = useState<string>('');

  const [newSku, setNewSku] = useState('');
  const [newName, setNewName] = useState('');
  const [newCategory, setNewCategory] = useState<string>(PartCategory.ACCESSORIES);
  const [newBrand, setNewBrand] = useState('');
  const [newType, setNewType] = useState('');
  const [newPrice, setNewPrice] = useState<number>(1000);
  const [newStock, setNewStock] = useState<number>(10);
  const [newVisible, setNewVisible] = useState<boolean>(true);
  const [newImages, setNewImages] = useState('');

  const load = useCallback(async () => {
    setLoading(true);
    setErr('');
    try {
      const items = await AdminFetchParts();
      setDbParts(items);
      const orders = await AdminFetchOrders();
      setDbOrders(orders);
    } catch (e: any) {
      setErr(e?.message || 'Ошибка загрузки');
    } finally {
      setLoading(false);
    }
  }, [])

  useEffect(() => { load(); }, [load]);

  const analytics = useMemo(() => {
    const orderTotal = (o: any) => {
      if (typeof o.total_kzt === 'number') return o.total_kzt;
      return (o.items || []).reduce((sum: number, it: any) => sum + (it.price_kzt || 0) * (it.qty || 0), 0);
    };

    const totalOrders = dbOrders.length;
    const totalRevenue = dbOrders.reduce((sum, o) => sum + orderTotal(o), 0);
    const avgOrder = totalOrders > 0 ? Math.round(totalRevenue / totalOrders) : 0;
    const pendingCount = dbOrders.filter((o) => (o.status || 'pending') === 'pending').length;
    const lowStock = dbParts.filter((p) => Number(p.stock_qty || 0) <= 5);

    const topMap = new Map<string, { sku: string; name: string; qty: number; revenue: number }>();
    dbOrders.forEach((o) => {
      (o.items || []).forEach((it: any) => {
        const sku = it.sku || it.partNumber || 'unknown';
        const prev = topMap.get(sku) || { sku, name: it.name || sku, qty: 0, revenue: 0 };
        prev.qty += Number(it.qty || 0);
        prev.revenue += Number(it.price_kzt || 0) * Number(it.qty || 0);
        topMap.set(sku, prev);
      });
    });
    const topParts = Array.from(topMap.values()).sort((a, b) => b.revenue - a.revenue).slice(0, 5);

    return { totalOrders, totalRevenue, avgOrder, pendingCount, lowStock, topParts };
  }, [dbOrders, dbParts]);

  const updateOrderStatus = async (id: string, status: string) => {
    try {
      await AdminUpdateOrder(id, { status });
      await load();
    } catch (e: any) {
      alert(e?.message || 'Не удалось обновить заказ');
    }
  };

  const deleteOrder = async (id: string) => {
    if (!confirm('Удалить заказ?')) return;
    try {
      await AdminDeleteOrder(id);
      await load();
    } catch (e: any) {
      alert(e?.message || 'Не удалось удалить заказ');
    }
  };

  const updatePart = async (id: string, patch: any) => {
    try {
      await AdminUpdatePart(id, patch);
      await load();
    } catch (e: any) {
      alert(e?.message || 'Не удалось обновить');
    }
  };

  const delPart = async (id: string) => {
    if (!confirm('Удалить товар?')) return;
    try {
      await adminDeletePart(id);
      await load();
    } catch (e: any) {
      alert(e?.message || 'Не удалось удалить');
    }
  };

  const addPart = async () => {
    if (!newSku.trim() || !newName.trim()) {
      alert('SKU и название обязательны');
      return;
    }
    try {
      const images = newImages
        .split(',')
        .map((s) => s.trim())
        .filter(Boolean);
      await adminCreatePart({
        sku: newSku.trim(),
        name: newName.trim(),
        category: newCategory,
        brand: newBrand.trim() || undefined,
        type: newType.trim() || undefined,
        price_kzt: Number(newPrice) || 0,
        currency: "KZT",
        stock_qty: Number(newStock) || 0,
        is_visible: !!newVisible,
        images,
        compatibility: { type: "universal" },
        recommend_for_issue_codes: []
      });
      setNewSku('');
      setNewName('');
      setNewBrand('');
      setNewType('');
      setNewImages('');
      await load();
    } catch (e: any) {
      alert(e?.message || 'Не удалось добавить');
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-black tracking-tight">Админ панель</h2>
        <button onClick={onExit} className="px-4 py-2 rounded bg-slate-800 text-white font-bold hover:bg-slate-700">Выйти</button>
      </div>

      <div className="flex gap-2 mb-6">
        <button onClick={() => setActiveTab('analytics')} className={`px-4 py-2 rounded font-bold ${activeTab==='analytics'?'bg-cyan-600 text-white':'bg-slate-100 text-slate-700'}`}>Аналитика</button>
        <button onClick={() => setActiveTab('parts')} className={`px-4 py-2 rounded font-bold ${activeTab==='parts'?'bg-cyan-600 text-white':'bg-slate-100 text-slate-700'}`}>Товары</button>
        <button onClick={() => setActiveTab('orders')} className={`px-4 py-2 rounded font-bold ${activeTab==='orders'?'bg-cyan-600 text-white':'bg-slate-100 text-slate-700'}`}>Заказы</button>
      </div>

      {err && <div className="mb-4 text-red-600 text-sm">{err}</div>}
      {loading && <div className="mb-4 text-slate-600 text-sm">Загрузка...</div>}

      {activeTab === 'analytics' && (
        <div className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="bg-white border rounded p-4">
              <div className="text-xs uppercase text-slate-500 font-bold">Заказы</div>
              <div className="text-2xl font-black">{analytics.totalOrders}</div>
              <div className="text-xs text-slate-500 mt-1">В ожидании: {analytics.pendingCount}</div>
            </div>
            <div className="bg-white border rounded p-4">
              <div className="text-xs uppercase text-slate-500 font-bold">Выручка</div>
              <div className="text-2xl font-black">{analytics.totalRevenue} ₸</div>
              <div className="text-xs text-slate-500 mt-1">Средний чек: {analytics.avgOrder} ₸</div>
            </div>
            <div className="bg-white border rounded p-4">
              <div className="text-xs uppercase text-slate-500 font-bold">Товары</div>
              <div className="text-2xl font-black">{dbParts.length}</div>
              <div className="text-xs text-slate-500 mt-1">Низкий остаток: {analytics.lowStock.length}</div>
            </div>
            <div className="bg-white border rounded p-4">
              <div className="text-xs uppercase text-slate-500 font-bold">Обновление</div>
              <div className="text-2xl font-black">Сейчас</div>
              <div className="text-xs text-slate-500 mt-1">Данные получены из базы</div>
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <div className="bg-white border rounded p-4">
              <div className="font-bold mb-3">Топ товаров по выручке</div>
              {analytics.topParts.length === 0 && <div className="text-sm text-slate-600">Нет данных</div>}
              <div className="space-y-2">
                {analytics.topParts.map((p) => (
                  <div key={p.sku} className="flex items-center justify-between text-sm">
                    <div className="font-medium">{p.name}</div>
                    <div className="text-slate-600">{p.revenue} ₸</div>
                  </div>
                ))}
              </div>
            </div>
            <div className="bg-white border rounded p-4">
              <div className="font-bold mb-3">Низкий остаток</div>
              {analytics.lowStock.length === 0 && <div className="text-sm text-slate-600">Все в наличии</div>}
              <div className="space-y-2">
                {analytics.lowStock.slice(0, 8).map((p: any) => (
                  <div key={p._id} className="flex items-center justify-between text-sm">
                    <div className="font-medium">{p.name}</div>
                    <div className="text-slate-600">Остаток: {p.stock_qty}</div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}

      {activeTab === 'parts' && (
        <div className="space-y-6">
          <div className="bg-white border rounded p-4">
            <div className="font-bold mb-3">Добавить товар</div>
            <div className="grid grid-cols-1 md:grid-cols-8 gap-3">
              <input className="border rounded px-3 py-2 text-sm" placeholder="SKU" value={newSku} onChange={e=>setNewSku(e.target.value)} />
              <input className="border rounded px-3 py-2 text-sm md:col-span-2" placeholder="Название" value={newName} onChange={e=>setNewName(e.target.value)} />
              <select className="border rounded px-3 py-2 text-sm" value={newCategory} onChange={e=>setNewCategory(e.target.value)}>
                {Object.values(PartCategory).map(c => <option key={c} value={c}>{c}</option>)}
              </select>
              <input className="border rounded px-3 py-2 text-sm" placeholder="Бренд" value={newBrand} onChange={e=>setNewBrand(e.target.value)} />
              <input className="border rounded px-3 py-2 text-sm" placeholder="Тип" value={newType} onChange={e=>setNewType(e.target.value)} />
              <input type="number" className="border rounded px-3 py-2 text-sm" placeholder="Цена ₸" value={newPrice} onChange={e=>setNewPrice(Number(e.target.value))} />
              <input type="number" className="border rounded px-3 py-2 text-sm" placeholder="Склад" value={newStock} onChange={e=>setNewStock(Number(e.target.value))} />
            </div>
            <div className="flex items-center gap-3 mt-3">
              <input className="flex-1 border rounded px-3 py-2 text-sm" placeholder="Изображения (URL, через запятую)" value={newImages} onChange={e=>setNewImages(e.target.value)} />
              <label className="flex items-center gap-2 text-sm">
                <input type="checkbox" checked={newVisible} onChange={e=>setNewVisible(e.target.checked)} />
                Доступен
              </label>
              <button onClick={addPart} className="ml-auto px-4 py-2 rounded bg-slate-800 text-white font-bold hover:bg-slate-700">Добавить</button>
            </div>
          </div>

          <div className="bg-white border rounded overflow-x-auto">
            <table className="min-w-full text-sm">
              <thead className="bg-slate-50">
                <tr>
                  <th className="text-left p-3">SKU</th>
                  <th className="text-left p-3">Название</th>
                  <th className="text-left p-3">Категория</th>
                  <th className="text-left p-3">Бренд</th>
                  <th className="text-left p-3">Тип</th>
                  <th className="text-left p-3">Цена ₸</th>
                  <th className="text-left p-3">Склад</th>
                  <th className="text-left p-3">Доступен</th>
                  <th className="text-left p-3">Действия</th>
                </tr>
              </thead>
              <tbody>
                {dbParts.map((p:any) => (
                  <tr key={p._id} className="border-t">
                    <td className="p-3 font-mono text-xs">{p.sku}</td>
                    <td className="p-3">
                      <input className="border rounded px-2 py-1 w-full" value={p.name} onChange={e=>setDbParts(prev=>prev.map(x=>x._id===p._id?{...x,name:e.target.value}:x))} onBlur={()=>updatePart(p._id,{name:p.name})} />
                    </td>
                    <td className="p-3">
                      <select className="border rounded px-2 py-1 w-full" value={p.category} onChange={e=>updatePart(p._id,{category:e.target.value})}>
                        {Object.values(PartCategory).map(c => <option key={c} value={c}>{c}</option>)}
                      </select>
                    </td>
                    <td className="p-3">
                      <input className="border rounded px-2 py-1 w-full" value={p.brand || ''} onChange={e=>setDbParts(prev=>prev.map(x=>x._id===p._id?{...x,brand:e.target.value}:x))} onBlur={()=>updatePart(p._id,{brand:p.brand || ''})} />
                    </td>
                    <td className="p-3">
                      <input className="border rounded px-2 py-1 w-full" value={p.type || ''} onChange={e=>setDbParts(prev=>prev.map(x=>x._id===p._id?{...x,type:e.target.value}:x))} onBlur={()=>updatePart(p._id,{type:p.type || ''})} />
                    </td>
                    <td className="p-3">
                      <input type="number" className="border rounded px-2 py-1 w-28" value={p.price_kzt} onChange={e=>setDbParts(prev=>prev.map(x=>x._id===p._id?{...x,price_kzt:Number(e.target.value)}:x))} onBlur={()=>updatePart(p._id,{price_kzt:p.price_kzt})} />
                    </td>
                    <td className="p-3">
                      <input type="number" className="border rounded px-2 py-1 w-20" value={p.stock_qty} onChange={e=>setDbParts(prev=>prev.map(x=>x._id===p._id?{...x,stock_qty:Number(e.target.value)}:x))} onBlur={()=>updatePart(p._id,{stock_qty:p.stock_qty})} />
                    </td>
                    <td className="p-3">
                      <input type="checkbox" checked={!!p.is_visible} onChange={e=>updatePart(p._id,{is_visible:e.target.checked})} />
                    </td>
                    <td className="p-3">
                      <button onClick={()=>delPart(p._id)} className="px-3 py-1 rounded bg-red-600 text-white font-bold">Удалить</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {activeTab === 'orders' && (
        <div className="space-y-4">
          {dbOrders.map((o:any) => (
            <div key={o.id || o._id} className="bg-white border rounded p-4">
              <div className="flex flex-wrap items-center justify-between gap-3">
                <div className="font-bold">Заказ #{(o.id||o._id||'').toString().slice(-6)} — {o.total_kzt || 0} ₸</div>
                <div className="flex items-center gap-2">
                  <select
                    className="border rounded px-2 py-1 text-sm"
                    value={o.status || 'pending'}
                    onChange={(e)=>updateOrderStatus(o.id||o._id, e.target.value)}
                  >
                    {['pending','processing','shipped','completed','cancelled'].map(s => (
                      <option key={s} value={s}>{s}</option>
                    ))}
                  </select>
                  <button onClick={()=>deleteOrder(o.id||o._id)} className="px-3 py-1 rounded bg-red-600 text-white font-bold text-sm">Удалить</button>
                </div>
              </div>
              <div className="text-sm text-slate-600">{o.contactInfo}</div>
              <div className="text-sm text-slate-600">{o.shippingAddress}</div>
              <div className="mt-2 text-sm">
                {(o.items||[]).map((it:any, idx:number)=>(
                  <div key={idx} className="flex justify-between">
                    <span>{it.qty}x {it.name}</span>
                    <span className="font-medium">{it.price_kzt * it.qty} ₸</span>
                  </div>
                ))}
              </div>
            </div>
          ))}
          {dbOrders.length===0 && !loading && <div className="text-sm text-slate-600">Заказов пока нет</div>}
        </div>
      )}
    </div>
  );
};


// --- Main App ---

const App: React.FC = () => {
  const [parts, setParts] = useState<Part[]>([]);
  const [cart, setCart] = useState<CartItem[]>([]);
  const [user, setUser] = useState<User | null>(null);
  const [orders, setOrders] = useState<Order[]>([]);
  const [currentView, setCurrentView] = useState<PageView>(PageView.HOME);
  
  const [vin, setVin] = useState('');
  const [make, setMake] = useState('Universal');
  const [model, setModel] = useState('All');
  const [year, setYear] = useState('All');
  const [vehicleMakes, setVehicleMakes] = useState<string[]>([]);
  const [modelsByMake, setModelsByMake] = useState<Record<string, string[]>>({});
  const [yearsByMakeModel, setYearsByMakeModel] = useState<Record<string, number[]>>({});
  const [category, setCategory] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [issue, setIssue] = useState<string>('');

  const [isLoadingParts, setIsLoadingParts] = useState(false);
  const [partsError, setPartsError] = useState<string>('');
  const [authLoading, setAuthLoading] = useState(true);

  // Parts are already filtered/sorted on the backend
  const displayedParts = useMemo(() => parts, [parts]);

  const cartCount = cart.reduce((acc, item) => acc + item.quantity, 0);

  const [toastMsg, setToastMsg] = useState<string>('');
  useEffect(() => {
    if (!toastMsg) return;
    const t = setTimeout(() => setToastMsg(''), 2000);
    return () => clearTimeout(t);
  }, [toastMsg]);

  useEffect(() => {
    const token = getToken();
    if (!token) {
      setAuthLoading(false);
      return;
    }
    fetchMe()
      .then((me) => {
        setUser(me);
      })
      .catch(() => {
        clearToken();
      })
      .finally(() => setAuthLoading(false));
  }, []);

  const loadParts = useCallback(async () => {
    setIsLoadingParts(true);
    setPartsError('');
    try {
      const items = await fetchParts({
        vin: vin || undefined,
        issue: issue || undefined,
        q: searchQuery || undefined,
        make,
        model,
        year,
        category,
      });
      setParts(items);
    } catch (e:any) {
      setParts([]);
      setPartsError(e?.message || 'Не удалось загрузить товары');
    } finally {
      setIsLoadingParts(false);
    }
  }, [vin, issue, searchQuery, make, model, year, category]);

  const handleAddToCart = (part: Part) => {
    if (!user) {
      setCurrentView(PageView.LOGIN);
      return;
    }
    if (part.stock <= 0) {
      setToastMsg('Товар закончился на складе');
      return;
    }
    setCart(prev => {
      const existing = prev.find(item => item.id === part.id);
      if (existing) {
        if (existing.quantity >= part.stock) {
          setToastMsg(`Доступно только ${part.stock} шт.`);
          return prev;
        }
        setToastMsg('Товар добавлен в корзину');
        return prev.map(item => item.id === part.id ? { ...item, quantity: item.quantity + 1 } : item);
      }
      setToastMsg('Товар добавлен в корзину');
      return [...prev, { ...part, quantity: 1 }];
    });
  };

  const handleUpdateCartQty = (id: string, delta: number) => {
    let exceeded: CartItem | null = null;
    setCart(prev => prev
      .map(item => {
        if (item.id !== id) return item;
        if (delta > 0 && item.quantity >= item.stock) {
          exceeded = item;
          return item;
        }
        return { ...item, quantity: Math.max(0, item.quantity + delta) };
      })
      .filter(item => item.quantity > 0));
    if (exceeded) {
      setToastMsg(`Доступно только ${exceeded.stock} шт.`);
    }
  };

  const handleSearch = () => {
    setCurrentView(PageView.CATALOG);
    // Logic handled by displayedParts dependency on state
  };

  // 1) Read URL params (?vin=...&issue=...) to support redirect from startup app
  useEffect(() => {
    const usp = new URLSearchParams(window.location.search);
    const vinParam = (usp.get('vin') || '').trim().toUpperCase();
    const issueParam = (usp.get('issue') || '').trim();
    if (vinParam) {
      setVin(vinParam);
      setCurrentView(PageView.CATALOG);
    }
    if (issueParam) {
      setIssue(issueParam);
      setCurrentView(PageView.CATALOG);
    }
  }, []);

  // Load selectable vehicle makes from backend once.
  useEffect(() => {
    let cancelled = false;
    fetchVehicleMakes()
      .then((makes) => {
        if (cancelled) return;
        setVehicleMakes(makes || []);
      })
      .catch(() => {
        if (cancelled) return;
        setVehicleMakes([]);
      });
    return () => {
      cancelled = true;
    };
  }, []);

  // Load models lazily: only when a make is selected.
  useEffect(() => {
    if (!make || make === 'Universal') return;
    if (modelsByMake[make]) return;

    let cancelled = false;
    fetchModelsByMake(make)
      .then((models) => {
        if (cancelled) return;
        setModelsByMake((prev) => ({ ...prev, [make]: models || [] }));
      })
      .catch(() => {
        if (cancelled) return;
        setModelsByMake((prev) => ({ ...prev, [make]: [] }));
      });

    return () => {
      cancelled = true;
    };
  }, [make, modelsByMake]);

  // Load years lazily: only when both make and model are selected.
  useEffect(() => {
    if (!make || make === 'Universal') return;
    if (!model || model === 'All') return;
    const key = `${make}::${model}`;
    if (yearsByMakeModel[key]) return;

    let cancelled = false;
    fetchYearsByMakeModel(make, model)
      .then((years) => {
        if (cancelled) return;
        setYearsByMakeModel((prev) => ({ ...prev, [key]: years || [] }));
      })
      .catch(() => {
        if (cancelled) return;
        setYearsByMakeModel((prev) => ({ ...prev, [key]: [] }));
      });

    return () => {
      cancelled = true;
    };
  }, [make, model, yearsByMakeModel]);

  // Load parts only in Catalog view.
  useEffect(() => {
    if (currentView !== PageView.CATALOG) return;
    loadParts();
  }, [currentView, loadParts]);

  
  return (
    <div className="min-h-screen flex flex-col bg-white">
      {toastMsg && (
        <div className="fixed top-4 right-4 bg-slate-900 text-white px-4 py-3 rounded shadow-lg text-sm z-[9999]">
          {toastMsg}
        </div>
      )}
      <Navbar 
        cartCount={cartCount} 
        user={user} 
        onNavigate={setCurrentView} 
        onLogout={() => { clearToken(); setUser(null); setCurrentView(PageView.HOME); }}
        onLoginClick={() => setCurrentView(PageView.LOGIN)}
        searchQuery={searchQuery}
        setSearchQuery={setSearchQuery}
      />

      <main className="flex-1">
        {authLoading && (
          <div className="container mx-auto px-4 py-3 text-sm text-slate-600">Проверка сессии...</div>
        )}
        {currentView === PageView.LOGIN && (
          <AuthView onDone={(me) => { setUser(me); if (me.role==='admin') setCurrentView(PageView.ADMIN_DASHBOARD); else setCurrentView(PageView.HOME); }} />
        )}

        {currentView === PageView.HOME && (
          <>
            <HeroSection 
              vin={vin} setVin={setVin}
              onSearch={handleSearch}
              make={make} setMake={setMake}
              model={model} setModel={setModel}
              year={year} setYear={setYear}
              makes={vehicleMakes}
              modelsByMake={modelsByMake}
              yearsByMakeModel={yearsByMakeModel}
            />
            <CategoryGrid 
              onSelectCategory={(cat) => { setCategory(cat); setCurrentView(PageView.CATALOG); }} 
            />
          </>
        )}

        {currentView === PageView.CATALOG && (
          <>
            {isLoadingParts && (
              <div className="container mx-auto px-4 py-3 text-sm text-slate-600">Загрузка товаров...</div>
            )}
            {partsError && !isLoadingParts && (
              <div className="container mx-auto px-4 py-3 text-sm text-red-600">{partsError}</div>
            )}
            <Catalog 
              parts={displayedParts}
              onAddToCart={handleAddToCart}
              selectedCategory={category}
              onSelectCategory={setCategory}
              searchQuery={searchQuery}
            />
          </>
        )}

        {currentView === PageView.CART && (
          <CartView cart={cart} onUpdateQty={handleUpdateCartQty} onRemove={(id) => setCart(c => c.filter(i => i.id !== id))} onCheckout={() => setCurrentView(user ? PageView.CHECKOUT : PageView.HOME)} />
        )}

        {currentView === PageView.CHECKOUT && user && (
           <CheckoutView
             user={user}
             cart={cart}
             onPlaceOrder={async (address, contact) => {
               try {
                 await createOrder({ cart, shippingAddress: address, contactInfo: contact });
                 setCart([]);
                 await loadParts();
                 alert('Заказ оформлен!');
                 setCurrentView(PageView.HOME);
               } catch (e:any) {
                 alert(e?.message || 'Не удалось оформить заказ');
               }
             }}
             onCancel={() => setCurrentView(PageView.CART)}
           />
        )}

        {currentView === PageView.ADMIN_DASHBOARD && user?.role === 'admin' && (
           <AdminPanel onExit={() => setCurrentView(PageView.HOME)} />
        )}
      </main>

      <footer className="bg-slate-900 text-slate-400 py-8 border-t border-slate-800">
        <div className="container mx-auto px-4 text-center">
           <p className="font-bold text-white mb-2">bazarPO</p>
           <p className="text-sm">Внимание к деталям. Интернет-магазин автозапчастей.</p>
        </div>
      </footer>
    </div>
  );
};

export default App;
