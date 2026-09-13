<script lang="ts">
	/**
	 * BLOCK_WEB_LAYOUT_001
	 * Purpose: Root layout providing global header navigation, auth state awareness, and container framing.
	 */
	import { authState } from '$lib/auth-state.svelte';
	let { children } = $props();
</script>

<div class="min-h-screen flex flex-col bg-slate-50 text-slate-900 antialiased">
	<header class="h-14 border-b border-slate-200 bg-white px-4 lg:px-8 flex items-center justify-between sticky top-0 z-40">
		<div class="flex items-center gap-6">
			<a href="/" class="flex items-center gap-2 font-bold text-slate-900 text-base">
				<span class="w-8 h-8 rounded-lg bg-slate-900 text-white flex items-center justify-center text-sm">🏫</span>
				<span>Campus CMS</span>
			</a>
			<nav class="hidden md:flex items-center gap-4 text-sm font-medium text-slate-600">
				<a href="/" class="hover:text-slate-900 transition">The Hub</a>
				<a href="/audit" class="hover:text-slate-900 transition">Audit Logs</a>
				<a href="/audit/compliance" class="hover:text-slate-900 transition">Compliance</a>
			</nav>
		</div>
		<div class="flex items-center gap-3 text-xs">
			{#if authState.isAuthenticated && authState.user}
				<span class="font-medium text-slate-700">{authState.user.username} ({authState.user.tier})</span>
				<button
					onclick={() => authState.logout()}
					class="px-3 py-1.5 rounded-lg border border-slate-300 hover:bg-slate-100 font-semibold text-slate-700 transition"
				>
					Sign Out
				</button>
			{:else}
				<a
					href="/auth/login"
					class="px-3.5 py-1.5 rounded-lg bg-slate-900 text-white font-semibold hover:bg-slate-800 transition"
				>
					Sign In
				</a>
			{/if}
		</div>
	</header>

	<main class="flex-1 max-w-7xl w-full mx-auto p-4 sm:p-6 lg:p-8">
		{@render children()}
	</main>

	<footer class="border-t border-slate-200 bg-white py-4 px-6 text-center text-xs text-slate-500">
		Campus Management System · Enterprise Modular Monolith · Proprietary
	</footer>
</div>
