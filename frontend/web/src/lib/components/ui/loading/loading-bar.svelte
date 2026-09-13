<script lang="ts">
	import { navigating } from '$app/stores';
	import { cn } from '$lib/utils';

	let { class: className }: { class?: string } = $props();
</script>

{#if $navigating}
	<div
		class={cn(
			'fixed top-0 left-0 right-0 z-100 h-0.75 w-full overflow-hidden bg-primary/20',
			className
		)}
		role="progressbar"
		aria-label="Loading page..."
	>
		<div class="h-full w-full origin-left animate-loading-bar bg-primary shadow-xs shadow-primary/50"></div>
	</div>
{/if}

<style>
	@keyframes loading-bar {
		0% {
			transform: translateX(-100%) scaleX(0.2);
		}
		50% {
			transform: translateX(0%) scaleX(0.7);
		}
		100% {
			transform: translateX(100%) scaleX(0.3);
		}
	}
	.animate-loading-bar {
		animation: loading-bar 1.2s cubic-bezier(0.4, 0, 0.2, 1) infinite;
	}
</style>
