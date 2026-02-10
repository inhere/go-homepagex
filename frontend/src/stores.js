import { writable } from 'svelte/store';

export const pageConfig = writable({});
export const currentRoute = writable('/');
export const viewStyle = writable('cards');
