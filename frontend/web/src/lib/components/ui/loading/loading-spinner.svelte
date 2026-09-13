<script lang="ts">
	import { cn } from '$lib/utils';

	let {
		size = 'md',
		label,
		class: className
	}: {
		size?: 'sm' | 'md' | 'lg' | 'xl';
		label?: string;
		class?: string;
	} = $props();

	const sizeClasses = {
		sm: 'size-4 border-2',
		md: 'size-6 border-2.5',
		lg: 'size-9 border-3',
		xl: 'size-12 border-4'
	};
</script>

<div class={cn('flex flex-col items-center justify-center gap-3', className)}>
	<div
		class={cn(
			'animate-spin rounded-full border-primary border-t-transparent',
			sizeClasses[size]
		)}
		role="status"
		aria-label={label ?? 'Loading'}
	></div>
	{#if label}
		<p class="text-xs font-medium text-muted-foreground animate-pulse">{label}</p>
	{/if}
</div>
