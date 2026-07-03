// src/lib/data/restaurants.js

// 🍕 Demo-Restaurants (ersetze mit deinen echten Daten)
export const restaurants = [
  {
    id: '1',
    name: 'Pizza Paradiso',
    slug: 'pizza-paradiso',
    image: '/images/pizza.jpg',
    rating: 4.5,
    deliveryTime: '20-30 min',
    categories: ['Pizza', 'Italienisch'],
    menu: [
      { id: 'p1', name: 'Margherita', preis: 8.99 },
      { id: 'p2', name: 'Salami', preis: 9.99 },
    ]
  },
  {
    id: '2',
    name: 'Burger Boss',
    slug: 'burger-boss',
    image: '/images/burger.jpg',
    rating: 4.3,
    deliveryTime: '15-25 min',
    categories: ['Burger', 'Amerikanisch'],
    menu: [
      { id: 'b1', name: 'Classic Burger', preis: 7.99 },
      { id: 'b2', name: 'Cheese Burger', preis: 8.99 },
    ]
  },
  {
    id: '3',
    name: 'Sushi Master',
    slug: 'sushi-master',
    image: '/images/sushi.jpg',
    rating: 4.7,
    deliveryTime: '25-35 min',
    categories: ['Sushi', 'Japanisch'],
    menu: [
      { id: 's1', name: 'Maki Set', preis: 12.99 },
      { id: 's2', name: 'Nigiri Mix', preis: 14.99 },
    ]
  }
];

export function getRestaurant(slugOderName) {
  return restaurants.find(
    (r) => r.slug === slugOderName || r.name === slugOderName
  );
}

export function getAllRestaurants() {
  return restaurants;
}