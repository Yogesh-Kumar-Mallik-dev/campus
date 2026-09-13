<script lang="ts">
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { Card, CardContent, CardHeader } from '$lib/components/ui/card';
	import { cn } from '$lib/utils';

	let {
		title = true,
		statsCount = 4,
		cardsCount = 3,
		class: className
	}: {
		title?: boolean;
		statsCount?: number;
		cardsCount?: number;
		class?: string;
	} = $props();
</script>

<div class={cn('space-y-8 animate-in fade-in duration-300', className)}>
	<!-- Header Skeleton -->
	{#if title}
		<div class="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
			<div class="space-y-2">
				<Skeleton class="h-5 w-28 rounded-full" />
				<Skeleton class="h-8 w-64 sm:w-80" />
				<Skeleton class="h-4 w-96 max-w-full" />
			</div>
			<div class="flex gap-2">
				<Skeleton class="h-9 w-28 rounded-lg" />
				<Skeleton class="h-9 w-32 rounded-lg" />
			</div>
		</div>
	{/if}

	<!-- Stats Grid Skeleton -->
	{#if statsCount > 0}
		<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
			{#each Array(statsCount) as _}
				<Card class="border-border/60 shadow-xs">
					<CardHeader class="pb-2 space-y-2">
						<Skeleton class="h-3.5 w-20" />
						<Skeleton class="h-7 w-24" />
					</CardHeader>
					<CardContent>
						<Skeleton class="h-3.5 w-3/4" />
					</CardContent>
				</Card>
			{/each}
		</div>
	{/if}

	<!-- Main Content Cards Skeleton -->
	{#if cardsCount > 0}
		<div class="grid gap-5 md:grid-cols-2 lg:grid-cols-3">
			{#each Array(cardsCount) as _}
				<Card class="border-border/60 shadow-xs space-y-3 p-5">
					<div class="flex items-center justify-between">
						<Skeleton class="h-5 w-20 rounded-md" />
						<Skeleton class="size-6 rounded-full" />
					</div>
					<Skeleton class="h-6 w-3/4" />
					<Skeleton class="h-4 w-full" />
					<Skeleton class="h-4 w-2/3" />
					<div class="pt-2 flex justify-between items-center">
						<Skeleton class="h-8 w-24 rounded-md" />
						<Skeleton class="h-4 w-16" />
					</div>
				</Card>
			{/each}
		</div>
	{/if}
</div>
