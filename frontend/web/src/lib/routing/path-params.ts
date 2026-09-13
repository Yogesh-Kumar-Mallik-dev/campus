import type { RouteParameterValue } from './types';

export function fillRouteParameters(pattern: string, params: Readonly<Record<string, RouteParameterValue>>): string {
	const expected = [...pattern.matchAll(/:([A-Za-z_][A-Za-z0-9_]*)/g)].map((match) => match[1]);
	for (const name of expected) {
		if (!(name in params)) throw new Error(`Missing route parameter '${name}' for '${pattern}'`);
	}
	const extras = Object.keys(params).filter((name) => !expected.includes(name));
	if (extras.length > 0) {
		throw new Error(`Unexpected route parameter '${extras[0]}' for '${pattern}'`);
	}
	return pattern.replace(/:([A-Za-z_][A-Za-z0-9_]*)/g, (_, name: string) => encodeURIComponent(String(params[name])));
}
