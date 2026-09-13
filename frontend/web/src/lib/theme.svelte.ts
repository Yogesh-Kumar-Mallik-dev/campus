import { browser } from '$app/environment';

const THEME_KEY = 'bbdit-theme';
const LEGACY_THEME_KEY = 'campus-theme';

class ThemeState {
	isDark = $state(false);

	constructor() {
		if (browser) {
			// Read current class that was set synchronously by <head> script
			this.isDark = document.documentElement.classList.contains('dark');

			// Listen to system preference changes if no manual setting
			const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
			mediaQuery.addEventListener('change', (e) => {
				const saved = localStorage.getItem(THEME_KEY) || localStorage.getItem(LEGACY_THEME_KEY);
				if (!saved) {
					this.set(e.matches);
				}
			});
		}
	}

	set(dark: boolean) {
		this.isDark = dark;
		if (browser) {
			document.documentElement.classList.toggle('dark', dark);
			document.documentElement.style.colorScheme = dark ? 'dark' : 'light';
			const mode = dark ? 'dark' : 'light';
			localStorage.setItem(THEME_KEY, mode);
			localStorage.setItem(LEGACY_THEME_KEY, mode);
		}
	}

	toggle() {
		this.set(!this.isDark);
	}
}

export const theme = new ThemeState();
