<script lang="ts">
	/**
	 * BLOCK_WEB_LAYOUT_001
	 * Purpose: Root layout providing global header navigation, theme toggling, auth state awareness, and container framing.
	 */
	import './+layout.css';
	import { authState } from '$lib/auth-state.svelte';
	import { theme } from '$lib/theme.svelte';
	import { Button } from '$lib/components/ui/button';

	let { children } = $props();
</script>

<div class="min-h-screen flex flex-col bg-background text-foreground antialiased transition-colors duration-200">
	<header class="h-14 border-b border-border bg-card/80 backdrop-blur-md px-4 lg:px-8 flex items-center justify-between sticky top-0 z-40">
		<div class="flex items-center gap-6">
			<a href="/" class="flex items-center gap-2.5 font-bold text-foreground text-base">
				<span class="w-8 h-8 rounded-lg bg-primary text-primary-foreground flex items-center justify-center text-sm shadow-sm font-black">🏫</span>
				<span class="tracking-tight">Campus CMS</span>
			</a>
			<nav class="hidden md:flex items-center gap-5 text-sm font-medium text-muted-foreground">
				<a href="/" class="hover:text-foreground transition">The Hub</a>
				<a href="/audit" class="hover:text-foreground transition">Audit Ledger</a>
				<a href="/audit/compliance" class="hover:text-foreground transition">Accreditation</a>
			</nav>
		</div>
		<div class="flex items-center gap-3 text-xs">
			<button
				onclick={() => theme.toggle()}
				class="h-8 w-8 rounded-lg border border-border bg-card hover:bg-muted text-foreground flex items-center justify-center transition"
				title="Toggle Theme"
			>
				{theme.isDark ? '☀️' : '🌙'}
			</button>

			{#if authState.isAuthenticated && authState.user}
				<span class="font-medium text-muted-foreground hidden sm:inline">{authState.user.username} ({authState.user.tier})</span>
				<Button
					variant="outline"
					size="sm"
					onclick={() => authState.logout()}
				>
					Sign Out
				</Button>
			{:else}
				<Button
					href="/auth/login"
					size="sm"
				>
					Sign In
				</Button>
			{/if}
		</div>
	</header>

	<main class="flex-1 max-w-7xl w-full mx-auto p-4 sm:p-6 lg:p-8">
		{@render children()}
	</main>

	<footer class="border-t border-border bg-card py-4 px-6 text-center text-xs text-muted-foreground">
		Campus Management System · Enterprise Modular Monolith · Proprietary
	</footer>
</div>
