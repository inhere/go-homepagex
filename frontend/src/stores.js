import { writable } from 'svelte/store';

export const pageConfig = writable({});
export const currentRoute = writable('/');
// 当前登录用户信息，null 表示游客
export const userInfo = writable(null);

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
  { 
    id: 'ocean-depths',
    name: '海洋深处', 
    colors: ['#1a2332', '#2d8b8b', '#a8dadc', '#f1faee'],
    description: '专业宁静的海洋主题'
  },
  { 
    id: 'tech-innovation',
    name: '科技创新', 
    colors: ['#1e1e1e', '#0066ff', '#00ffff', '#ffffff'],
    description: '科技感十足的现代主题'
  },
  { 
    id: 'modern-minimalist',
    name: '现代极简', 
    colors: ['#36454f', '#708090', '#d3d3d3', '#ffffff'],
    description: '简洁专业的极简主义'
  },
  { 
    id: 'midnight-galaxy',
    name: '午夜星河', 
    colors: ['#2b1e3e', '#4a4e8f', '#a490c2', '#e6e6fa'],
    description: '深邃神秘的宇宙主题'
  },
  { 
    id: 'forest-canopy',
    name: '森林树冠', 
    colors: ['#2d4a2b', '#7d8471', '#a4ac86', '#faf9f6'],
    description: '自然清新的森林主题'
  },
  { 
    id: 'arctic-frost',
    name: '北极冰霜', 
    colors: ['#4a6fa5', '#4a6fa5', '#c0c0c0', '#fafafa'],
    description: '清爽现代的冰雪主题'
  },
];

export const currentTheme = createPersistedStore('theme', 'ocean-depths');

export function getThemeColors(themeId) {
  const theme = themes.find(t => t.id === themeId) || themes[0];
  return {
    // 主色调（深色背景或主色）
    primary: theme.colors[0],
    // 次要色（强调色或辅助色）
    secondary: theme.colors[1],
    // 强调色（高亮或装饰）
    accent: theme.colors[2],
    // 背景色（浅色或文字色）
    background: theme.colors[3],
  };
}
