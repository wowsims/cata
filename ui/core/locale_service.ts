const LOCALE_KEY = 'wowsims_locale';

export function getCurrentLang(): string {
  const record = localStorage.getItem(LOCALE_KEY);
  if (record) {
    try {
      return JSON.parse(record).lang || 'en';
    } catch {
      return 'en';
    }
  }
  return 'en';
}

export function setCurrentLang(lang: string): void {
  localStorage.setItem(LOCALE_KEY, JSON.stringify({ lang }));
}

export {};
