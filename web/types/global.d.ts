import type Alpine from '@alpinejs/csp';

declare global {
  interface Window {
    Alpine: typeof Alpine;
  }
}

export { }; // Makes this file a module
