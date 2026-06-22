import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';

// =====================================================================
// 🌍 MEHRSPRACHIGKEIT (i18n)
// Ein einfacher Übersetzungs-Speicher. Komponenten nutzen $t('schluessel').
// Die Auswahl bleibt im localStorage erhalten.
//
// Hinweis: Das Admin-Menü wird bewusst NICHT übersetzt (bleibt deutsch).
// 'mc' = Minecraft-Verzauberungstisch (Standard Galactic Alphabet): nutzt die
// englischen Wörter, die per Spezial-Schriftart in Runen dargestellt werden.
// =====================================================================

export const SPRACHEN = [
  { code: 'de', label: 'Deutsch', flag: '🇩🇪' },
  { code: 'en', label: 'English', flag: '🇬🇧' },
  { code: 'es', label: 'Español', flag: '🇪🇸' },
  { code: 'ru', label: 'Русский', flag: '🇷🇺' },
  { code: 'ja', label: '日本語', flag: '🇯🇵' },
  { code: 'la', label: 'Latina', flag: '🏛️' },
  // Für Minecraft/Enchanting nutzen wir ein eigenes Bild-Icon (statt eines Emojis).
  { code: 'mc', label: 'Enchanting', flag: '🟪', icon: '/minecraft-icon-0.png' }
];

const woerterbuch = {
  de: {
    'nav.home': '🏠 Home',
    'nav.restaurants': '🍔 Restaurants',
    'nav.cart': '🛒 Warenkorb',
    'nav.orders': '🧾 Bestellungen',
    'nav.account': '👤 Account',
    'nav.login': '🔑 Login',
    'nav.admin': '🛠️ Admin',
    'home.hero_title': 'Willkommen bei Lieferino',
    'home.hero_sub': 'Dein Lieblingsessen, nur wenige Klicks entfernt 🍕',
    'home.cuisine_label': 'Essensart:',
    'home.min_order_label': 'Max. Mindestbestellwert:',
    'home.top10_title': '⭐ Unsere Top 10 Restaurants',
    'home.top10_sub': 'Live nach euren Bewertungen sortiert',
    'home.discover': '🔍 Alle Restaurants entdecken',
    'common.all_cuisines': '🌍 Alle Küchen',
    'common.search_placeholder': '🔍 Restaurant suchen…',
    'common.min': 'Min',
    'rest.title': '🍔 Alle Restaurants',
    'rest.subtitle': 'Stöbere durch unsere Auswahl und finde dein Lieblingsessen',
    'rest.found': '{n} Restaurant(s) gefunden',
    'rest.sort_standard': 'Sortierung: Standard',
    'rest.sort_rating': 'Beste Bewertung',
    'rest.sort_minorder': 'Niedrigster Mindestbestellwert',
    'rest.only_favs': '❤️ Nur Favoriten',
    'rest.all': '🤍 Alle',
    'rest.veg': '🌱 Vegetarisch',
    'rest.none': '😕 Keine Restaurants gefunden. Versuche einen anderen Filter.'
  },
  en: {
    'nav.home': '🏠 Home',
    'nav.restaurants': '🍔 Restaurants',
    'nav.cart': '🛒 Cart',
    'nav.orders': '🧾 Orders',
    'nav.account': '👤 Account',
    'nav.login': '🔑 Login',
    'nav.admin': '🛠️ Admin',
    'home.hero_title': 'Welcome to Lieferino',
    'home.hero_sub': 'Your favorite food, just a few clicks away 🍕',
    'home.cuisine_label': 'Cuisine:',
    'home.min_order_label': 'Max. minimum order:',
    'home.top10_title': '⭐ Our Top 10 Restaurants',
    'home.top10_sub': 'Sorted live by your ratings',
    'home.discover': '🔍 Discover all restaurants',
    'common.all_cuisines': '🌍 All cuisines',
    'common.search_placeholder': '🔍 Search restaurant…',
    'common.min': 'Min',
    'rest.title': '🍔 All restaurants',
    'rest.subtitle': 'Browse our selection and find your favorite meal',
    'rest.found': '{n} restaurant(s) found',
    'rest.sort_standard': 'Sort: Default',
    'rest.sort_rating': 'Best rating',
    'rest.sort_minorder': 'Lowest minimum order',
    'rest.only_favs': '❤️ Favorites only',
    'rest.all': '🤍 All',
    'rest.veg': '🌱 Vegetarian',
    'rest.none': '😕 No restaurants found. Try a different filter.'
  },
  es: {
    'nav.home': '🏠 Inicio',
    'nav.restaurants': '🍔 Restaurantes',
    'nav.cart': '🛒 Cesta',
    'nav.orders': '🧾 Pedidos',
    'nav.account': '👤 Cuenta',
    'nav.login': '🔑 Acceder',
    'nav.admin': '🛠️ Admin',
    'home.hero_title': 'Bienvenido a Lieferino',
    'home.hero_sub': 'Tu comida favorita, a solo unos clics 🍕',
    'home.cuisine_label': 'Tipo de cocina:',
    'home.min_order_label': 'Pedido mínimo máx.:',
    'home.top10_title': '⭐ Nuestros 10 mejores restaurantes',
    'home.top10_sub': 'Ordenados en vivo por tus valoraciones',
    'home.discover': '🔍 Descubre todos los restaurantes',
    'common.all_cuisines': '🌍 Todas las cocinas',
    'common.search_placeholder': '🔍 Buscar restaurante…',
    'common.min': 'Mín',
    'rest.title': '🍔 Todos los restaurantes',
    'rest.subtitle': 'Explora nuestra selección y encuentra tu comida favorita',
    'rest.found': '{n} restaurante(s) encontrado(s)',
    'rest.sort_standard': 'Orden: Estándar',
    'rest.sort_rating': 'Mejor valoración',
    'rest.sort_minorder': 'Pedido mínimo más bajo',
    'rest.only_favs': '❤️ Solo favoritos',
    'rest.all': '🤍 Todos',
    'rest.veg': '🌱 Vegetariano',
    'rest.none': '😕 No se encontraron restaurantes. Prueba otro filtro.'
  },
  ru: {
    'nav.home': '🏠 Главная',
    'nav.restaurants': '🍔 Рестораны',
    'nav.cart': '🛒 Корзина',
    'nav.orders': '🧾 Заказы',
    'nav.account': '👤 Аккаунт',
    'nav.login': '🔑 Вход',
    'nav.admin': '🛠️ Админ',
    'home.hero_title': 'Добро пожаловать в Lieferino',
    'home.hero_sub': 'Любимая еда всего в несколько кликов 🍕',
    'home.cuisine_label': 'Кухня:',
    'home.min_order_label': 'Макс. мин. заказ:',
    'home.top10_title': '⭐ Топ-10 наших ресторанов',
    'home.top10_sub': 'Сортировка по вашим оценкам в реальном времени',
    'home.discover': '🔍 Все рестораны',
    'common.all_cuisines': '🌍 Все кухни',
    'common.search_placeholder': '🔍 Поиск ресторана…',
    'common.min': 'Мин',
    'rest.title': '🍔 Все рестораны',
    'rest.subtitle': 'Просмотрите наш выбор и найдите любимое блюдо',
    'rest.found': 'Найдено ресторанов: {n}',
    'rest.sort_standard': 'Сортировка: по умолчанию',
    'rest.sort_rating': 'Лучший рейтинг',
    'rest.sort_minorder': 'Минимальный заказ',
    'rest.only_favs': '❤️ Только избранное',
    'rest.all': '🤍 Все',
    'rest.veg': '🌱 Вегетарианское',
    'rest.none': '😕 Рестораны не найдены. Попробуйте другой фильтр.'
  },
  ja: {
    'nav.home': '🏠 ホーム',
    'nav.restaurants': '🍔 レストラン',
    'nav.cart': '🛒 カート',
    'nav.orders': '🧾 注文',
    'nav.account': '👤 アカウント',
    'nav.login': '🔑 ログイン',
    'nav.admin': '🛠️ 管理',
    'home.hero_title': 'Lieferino へようこそ',
    'home.hero_sub': 'お気に入りの料理が数クリックで 🍕',
    'home.cuisine_label': '料理の種類:',
    'home.min_order_label': '最低注文額（最大）:',
    'home.top10_title': '⭐ 人気レストラン トップ10',
    'home.top10_sub': '評価をもとにリアルタイムで並び替え',
    'home.discover': '🔍 すべてのレストランを見る',
    'common.all_cuisines': '🌍 すべての料理',
    'common.search_placeholder': '🔍 レストランを検索…',
    'common.min': '最低',
    'rest.title': '🍔 すべてのレストラン',
    'rest.subtitle': '一覧からお気に入りの料理を見つけよう',
    'rest.found': '{n} 件のレストランが見つかりました',
    'rest.sort_standard': '並び替え: 標準',
    'rest.sort_rating': '評価が高い順',
    'rest.sort_minorder': '最低注文額が低い順',
    'rest.only_favs': '❤️ お気に入りのみ',
    'rest.all': '🤍 すべて',
    'rest.veg': '🌱 ベジタリアン',
    'rest.none': '😕 レストランが見つかりません。別のフィルターをお試しください。'
  },
  la: {
    'nav.home': '🏠 Domus',
    'nav.restaurants': '🍔 Popinae',
    'nav.cart': '🛒 Sporta',
    'nav.orders': '🧾 Iussa',
    'nav.account': '👤 Ratio',
    'nav.login': '🔑 Initus',
    'nav.admin': '🛠️ Praefectus',
    'home.hero_title': 'Salve apud Lieferino',
    'home.hero_sub': 'Cibus tuus dilectus, paucis digitis remotus 🍕',
    'home.cuisine_label': 'Genus cibi:',
    'home.min_order_label': 'Pretium minimum max.:',
    'home.top10_title': '⭐ Decem optimae popinae',
    'home.top10_sub': 'Secundum aestimationes vestras ordinatae',
    'home.discover': '🔍 Omnes popinas inveni',
    'common.all_cuisines': '🌍 Omnia genera',
    'common.search_placeholder': '🔍 Popinam quaere…',
    'common.min': 'Min',
    'rest.title': '🍔 Omnes popinae',
    'rest.subtitle': 'Perscrutare delectum nostrum et cibum dilectum inveni',
    'rest.found': '{n} popinae inventae',
    'rest.sort_standard': 'Ordo: Solitus',
    'rest.sort_rating': 'Optima aestimatio',
    'rest.sort_minorder': 'Minimum pretium minimum',
    'rest.only_favs': '❤️ Solae deliciae',
    'rest.all': '🤍 Omnes',
    'rest.veg': '🌱 Holeribus',
    'rest.none': '😕 Nullae popinae inventae. Aliud filtrum tempta.'
  }
};

function ladeStart() {
  if (!browser) return 'de';
  return localStorage.getItem('lieferino_sprache') || 'de';
}

export const sprache = writable(ladeStart());

if (browser) {
  sprache.subscribe((code) => {
    localStorage.setItem('lieferino_sprache', code);
    // <html lang> setzen (mc nutzt englische Wörter) + Flag für die Enchanting-Schrift.
    document.documentElement.lang = code === 'mc' ? 'en' : code;
    document.documentElement.dataset.sga = code === 'mc' ? 'true' : 'false';
  });
}

export function setzeSprache(code) {
  sprache.set(code);
}

// $t('schluessel') liefert den übersetzten Text (Fallback: Deutsch, dann der Schlüssel).
export const t = derived(sprache, ($s) => {
  const lang = $s === 'mc' ? 'en' : $s; // Enchanting nutzt die englischen Wörter
  return (schluessel) => woerterbuch[lang]?.[schluessel] ?? woerterbuch.de[schluessel] ?? schluessel;
});
