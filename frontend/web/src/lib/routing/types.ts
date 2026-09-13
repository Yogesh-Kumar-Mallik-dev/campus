export type RouteNameProps = {
	name: string;
	/** Stable ID that must be globally unique across all routable pages. */
	id: string;
	/** Parent route ID. Omitted values use the root sentinel `/`. */
	previousID?: string;
};

export type RouteParameterValue = string | number;
