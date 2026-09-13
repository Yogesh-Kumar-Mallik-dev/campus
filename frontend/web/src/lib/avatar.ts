/**
 * Profile avatar client-side local caching and state synchronization utilities.
 */

export function getAvatarStorageKey(userId: string): string {
	return `campus_avatar_${userId}`;
}

/**
 * Retrieve cached avatar Data URL or base64 string from browser/desktop localStorage.
 */
export function getCachedAvatar(userId: string): string | null {
	if (typeof window === 'undefined' || !window.localStorage) {
		return null;
	}
	try {
		return window.localStorage.getItem(getAvatarStorageKey(userId));
	} catch {
		return null;
	}
}

/**
 * Store user avatar in localStorage per instance.
 */
export function setCachedAvatar(userId: string, dataUrlOrBase64: string): void {
	if (typeof window === 'undefined' || !window.localStorage) {
		return;
	}
	try {
		window.localStorage.setItem(getAvatarStorageKey(userId), dataUrlOrBase64);
		window.dispatchEvent(
			new CustomEvent('campus:avatar-updated', {
				detail: { userId, avatar: dataUrlOrBase64 }
			})
		);
	} catch {
		// Ignore storage quota exceeded errors gracefully
	}
}

/**
 * Clear cached avatar from localStorage.
 */
export function clearCachedAvatar(userId: string): void {
	if (typeof window === 'undefined' || !window.localStorage) {
		return;
	}
	try {
		window.localStorage.removeItem(getAvatarStorageKey(userId));
		window.dispatchEvent(
			new CustomEvent('campus:avatar-updated', {
				detail: { userId, avatar: null }
			})
		);
	} catch {
		// Ignore
	}
}

/**
 * Converts a Blob or File into a base64 Data URL string.
 */
export function fileToDataUrl(file: Blob | File): Promise<string> {
	return new Promise((resolve, reject) => {
		const reader = new FileReader();
		reader.onload = () => {
			if (typeof reader.result === 'string') {
				resolve(reader.result);
			} else {
				reject(new Error('Failed to convert file to data URL'));
			}
		};
		reader.onerror = () => reject(reader.error || new Error('File read error'));
		reader.readAsDataURL(file);
	});
}
