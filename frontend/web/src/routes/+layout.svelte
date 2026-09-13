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

	let isSubsystemsMenuOpen = $state<boolean>(false);

	const allSubsystems = [
		{ name: '1. Auth & Permissions', path: '/auth', icon: '🔑' },
		{ name: '2. Audit Ledger', path: '/audit', icon: '🔒' },
		{ name: '3. Admissions Onboarding', path: '/onboarding', icon: '🎓' },
		{ name: '4. Attendance', path: '/attendance', icon: '📍' },
		{ name: '5. Billing & Invoices', path: '/billing', icon: '💳' },
		{ name: '6. Campus Notices', path: '/notices', icon: '📢' },
		{ name: '7. Hostel & Gate Pass', path: '/hostel', icon: '🏢' },
		{ name: '8. Mess & QR Dining', path: '/mess', icon: '🍽️' },
		{ name: '9. E-Library Catalog', path: '/library', icon: '📚' },
		{ name: '10. Study Hub & LMS', path: '/studyhub', icon: '📝' },
		{ name: '11. Mentorship Tracker', path: '/mentorship', icon: '🌱' },
		{ name: '12. Campus Events', path: '/events', icon: '🎉' },
		{ name: '13. Helpdesk & SLA', path: '/helpdesk', icon: '🎫' },
		{ name: '14. SOS Emergency Radar', path: '/sos', icon: '🚨' },
		{ name: '15. Whistleblower Vault', path: '/whistleblower', icon: '🛡️' },
		{ name: '16. Public Web Portal', path: '/portal', icon: '🌐' },
		{ name: '17. The Hub Super-App', path: '/hub', icon: '⚡' },
	];
</script>

<div class="min-h-screen flex flex-col bg-background text-foreground antialiased transition-colors duration-200">
	<header class="h-14 border-b border-border bg-card/80 backdrop-blur-md px-4 lg:px-8 flex items-center justify-between sticky top-0 z-40">
		<div class="flex items-center gap-6">
			<a href="/" class="flex items-center gap-2.5 font-bold text-foreground text-base">
				<span class="w-8 h-8 rounded-lg bg-primary text-primary-foreground flex items-center justify-center text-sm shadow-sm font-black">🏫</span>
				<span class="tracking-tight">Campus CMS</span>
			</a>
			<nav class="hidden md:flex items-center gap-4 text-sm font-medium text-muted-foreground">
				<a href="/" class="hover:text-foreground transition">Overview</a>
				<a href="/hub" class="hover:text-foreground transition font-bold text-foreground">The Hub</a>
				<a href="/audit" class="hover:text-foreground transition">Audit Ledger</a>
				<a href="/audit/compliance" class="hover:text-foreground transition">Accreditation</a>

				<!-- Subsystems dropdown toggle -->
				<div class="relative">
					<button
						onclick={() => (isSubsystemsMenuOpen = !isSubsystemsMenuOpen)}
						class="flex items-center gap-1 hover:text-foreground transition px-2 py-1 rounded-lg hover:bg-muted text-xs font-bold"
					>
						<span>All 17 Subsystems</span>
						<span class="text-[10px]">▼</span>
					</button>

					{#if isSubsystemsMenuOpen}
						<!-- Dropdown Panel -->
						<div class="absolute left-0 top-full mt-2 w-72 p-2 bg-card rounded-2xl border border-border shadow-2xl z-50 grid grid-cols-1 gap-1 max-h-96 overflow-y-auto">
							{#each allSubsystems as s}
								<a
									href={s.path}
									onclick={() => (isSubsystemsMenuOpen = false)}
									class="p-2 rounded-xl text-xs font-medium text-foreground hover:bg-muted flex items-center gap-2 transition"
								>
									<span>{s.icon}</span>
									<span>{s.name}</span>
								</a>
							{/each}
						</div>
					{/if}
				</div>
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
