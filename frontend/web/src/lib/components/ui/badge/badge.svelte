<script lang="ts" module>
	import { type VariantProps, tv } from "tailwind-variants";

	export const badgeVariants = tv({
		base: "relative h-5.5 gap-1.5 border border-transparent px-2.5 py-0.5 text-xs font-semibold tracking-tight transition-all has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2 [&>svg]:size-3! group/badge inline-flex w-fit shrink-0 items-center justify-center whitespace-nowrap transition-colors focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&>svg]:pointer-events-none select-none drop-shadow-xs",
		variants: {
			variant: {
				default: "bg-primary text-primary-foreground shadow-xs [a]:hover:bg-primary/85",
				secondary: "bg-secondary text-secondary-foreground shadow-xs [a]:hover:bg-secondary/80",
				destructive: "bg-destructive/15 text-destructive border-destructive/30 dark:bg-destructive/25 dark:text-red-300 dark:border-destructive/40 focus-visible:ring-destructive/20 [a]:hover:bg-destructive/25",
				outline: "border-border/80 bg-background/80 text-foreground shadow-2xs [a]:hover:bg-muted [a]:hover:text-muted-foreground",
				ghost: "hover:bg-muted hover:text-muted-foreground dark:hover:bg-muted/50",
				link: "text-primary underline-offset-4 hover:underline",
				bookmark: "bg-primary text-primary-foreground font-bold shadow-xs",
			},
			shape: {
				bookmark: "[clip-path:polygon(0%_0%,100%_0%,calc(100%-6px)_50%,100%_100%,0%_100%)] rounded-l-xs pl-2.5 pr-4",
				ribbon: "[clip-path:polygon(6px_0%,100%_0%,calc(100%-6px)_50%,100%_100%,6px_100%,0%_50%)] px-3.5",
				tag: "[clip-path:polygon(0%_0%,calc(100%-6px)_0%,100%_50%,calc(100%-6px)_100%,0%_100%)] rounded-l-xs pl-2.5 pr-4",
				banner: "[clip-path:polygon(0%_0%,100%_0%,100%_100%,50%_calc(100%-4px),0%_100%)] pb-1 px-3",
				pill: "rounded-full px-2.5 py-0.5",
			},
		},
		defaultVariants: {
			variant: "default",
			shape: "bookmark",
		},
	});

	export type BadgeVariant = VariantProps<typeof badgeVariants>["variant"];
	export type BadgeShape = VariantProps<typeof badgeVariants>["shape"];
</script>

<script lang="ts">
	import { cn, type WithElementRef } from "$lib/utils.js";
	import type { HTMLAnchorAttributes } from "svelte/elements";

	let {
		ref = $bindable(null),
		href,
		class: className,
		variant = "default",
		shape = "bookmark",
		children,
		...restProps
	}: WithElementRef<HTMLAnchorAttributes> & {
		variant?: BadgeVariant;
		shape?: BadgeShape;
	} = $props();
</script>

<svelte:element
	this={href ? "a" : "span"}
	bind:this={ref}
	data-slot="badge"
	{href}
	class={cn(badgeVariants({ variant, shape }), className)}
	{...restProps}
>
	{@render children?.()}
</svelte:element>
