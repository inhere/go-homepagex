import { writable } from 'svelte/store';

export const pageConfig = writable({});
export const currentRoute = writable('/');

function createPersistedStore(key, initialValue) {
  const storedValue = localStorage.getItem(key);
  const initial = storedValue ? JSON.parse(storedValue) : initialValue;
  const store = writable(initial);

  store.subscribe(value => {
    localStorage.setItem(key, JSON.stringify(value));
  });

  return store;
}

export const viewStyle = createPersistedStore('viewStyle', 'cards');

export const themes = [
  { id: 'default', name: '默认主题', colors: ['#667eea', '#764ba2'] },
  { id: 'ocean', name: '海洋蓝', colors: ['#2196F3', '#00BCD4'] },
  { id: 'forest', name: '森林绿', colors: ['#4CAF50', '#8BC34A'] },
  { id: 'sunset', name: '日落橙', colors: ['#FF6B6B', '#FF8E53'] },
  { id: 'midnight', name: '午夜紫', colors: ['#7C4DFF', '#E040FB'] },
  { id: 'dark', name: '深邃灰', colors: ['#37474F', '#546E7A'] },
  { id: 'minimal', name: '极简白', colors: ['#ffffff', '#f0f0f0'] },
  { id: 'github', name: 'GitHub 风', colors: ['#0969da', '#1f883d'] },
];

export const currentTheme = createPersistedStore('theme', 'default');

export function getThemeColors(themeId) {
  const theme = themes.find(t => t.id === themeId) || themes[0];
  return {
    primary: theme.colors[0],
    secondary: theme.colors[1],
  };
}
