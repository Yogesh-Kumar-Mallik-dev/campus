<script lang="ts">
	import { onMount, type Snippet } from 'svelte';
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Avatar, AvatarFallback, AvatarImage } from '$lib/components/ui/avatar';
	import { Switch } from '$lib/components/ui/switch';
	import { Separator } from '$lib/components/ui/separator';
	import logo from '@/assets/bbdit-logo-transparent.png';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import {
		IconLayoutDashboard,
		IconBell,
		IconBook,
		IconCalendar,
		IconCertificate,
		IconClipboardCheck,
		IconFileDescription,
		IconFileText,
		IconLogout,
		IconShieldCheck,
		IconUserCheck,
		IconMenu2,
		IconSun,
		IconMoon,
		IconBuildingCommunity,
		IconUsersGroup,
		IconShieldLock,
		IconFolder,
		IconExternalLink,
		IconChevronRight,
		IconLoader2,
		IconUser
	} from '@tabler/icons-svelte';
	import { getCachedAvatar, setCachedAvatar, fileToDataUrl } from '$lib/avatar';

	import { theme } from '$lib/theme.svelte';

	type User = {
		id?: string;
		displayName?: string;
		email?: string;
		status?: string;
		mfaRequired?: boolean;
		mfaSetupRequired?: boolean;
		assignments?: Array<{ role: string; scopeId?: string }>;
	} | null;

	let { children, user, unreadCount = 0 }: { children: Snippet; user: User; unreadCount?: number } = $props();

	let mobileNavOpen = $state(false);
	let isLogoutModalOpen = $state(false);
	let isLoggingOut = $state(false);
	let cachedAvatar = $state<string | null>(null);

	onMount(() => {
		if (user?.id) {
			cachedAvatar = getCachedAvatar(user.id);
			if (!cachedAvatar) {
				fetch(`/api/v1/users/${user.id}/avatar`, { credentials: 'include' })
					.then((res) => {
						if (res.ok) return res.blob();
						throw new Error();
					})
					.then((blob) => fileToDataUrl(blob))
					.then((dataUrl) => {
						cachedAvatar = dataUrl;
						if (user?.id) setCachedAvatar(user.id, dataUrl);
					})
					.catch(() => {});
			}
		}

		const handleAvatarUpdate = (e: Event) => {
			const detail = (e as CustomEvent).detail;
			if (detail && (!user?.id || detail.userId === user.id)) {
				cachedAvatar = detail.avatar;
			}
		};
		window.addEventListener('campus:avatar-updated', handleAvatarUpdate);
		return () => {
			window.removeEventListener('campus:avatar-updated', handleAvatarUpdate);
		};
	});

	const roles = $derived(user?.assignments?.map((assignment) => assignment.role) ?? []);
	const primaryRole = $derived(roles[0] ?? 'STUDENT');
	const canApprove = $derived(roles.some((role) => ['SUPER_ADMIN', 'ADMIN', 'MODERATOR'].includes(role)));
	const canPublish = $derived(roles.some((role) => ['SUPER_ADMIN', 'ADMIN', 'MODERATOR', 'TEACHER'].includes(role)));
	const isMFAEnabled = $derived(Boolean(user?.mfaRequired && !user?.mfaSetupRequired));

	const initials = $derived(
		user?.displayName
			? user.displayName
					.split(' ')
					.filter(Boolean)
					.map((n) => n[0])
					.slice(0, 2)
					.join('')
					.toUpperCase()
			: 'CM'
	);

	function isNavActive(path: string) {
		if (path === '/') return page.url.pathname === '/';
		return page.url.pathname === path || page.url.pathname.startsWith(path + '/');
	}

	async function handleSignOut() {
		try {
			await fetch('/api/v1/auth/sign-out', {
				method: 'POST',
				credentials: 'include'
			});
		} catch {
			/* ignore */
		}
		window.location.href = '/public/sign-in';
	}
</script>

<svelte:window
	onkeydown={(e) => {
		if (e.key === 'Escape' && mobileNavOpen) {
			mobileNavOpen = false;
		}
	}}
/>

<div class="min-h-screen-dvh min-h-screen bg-background text-foreground flex flex-col">
	<!-- Mobile Drawer Backdrop -->
	{#if mobileNavOpen}
		<button
			class="fixed inset-0 z-40 bg-black/60 backdrop-blur-xs lg:hidden"
			onclick={() => (mobileNavOpen = false)}
			aria-label="Close mobile sidebar"
		></button>
	{/if}

	<!-- Left Sidebar (Detached, Fixed & Independently Scrollable) -->
	<aside
		class="fixed inset-y-0 left-0 z-50 flex h-dvh h-screen w-72 max-w-[min(18rem,calc(100vw-3rem))] flex-col border-r border-sidebar-border bg-sidebar transition-transform duration-200 ease-in-out lg:translate-x-0 {mobileNavOpen
			? 'translate-x-0 shadow-2xl'
			: '-translate-x-full'}"
	>
		<!-- 1. Header with Logo (Fixed / Non-scrolling & Centered) -->
		<div class="flex h-20 shrink-0 items-center justify-center border-b border-sidebar-border px-4 py-3 bg-sidebar pt-[max(0.75rem,env(safe-area-inset-top))]">
			<a
				href="/public"
				aria-label="BBDIT Public Home"
				class="flex items-center justify-center rounded-lg border border-black/8 bg-white/95 px-3 py-1.5 shadow-xs ring-1 ring-black/3 dark:border-white/12 dark:bg-[#f2f5f7] dark:ring-white/5 transition-opacity hover:opacity-95"
				onclick={() => (mobileNavOpen = false)}
			>
				<img
					src={logo}
					alt="Babu Banarsi Das Institute of Technology"
					class="h-auto max-h-10 w-auto max-w-[200px] object-contain mx-auto"
				/>
			</a>
		</div>

		<!-- 2. Navigation Menu Links (Independently Scrollable with hidden scrollbar) -->
		<div class="flex-1 overflow-y-auto overscroll-contain px-3 py-4 space-y-6 hide-scrollbar">
			<!-- Section: Main -->
			<div class="space-y-1">
				<p class="px-3 text-[11px] font-bold uppercase tracking-wider text-muted-foreground">Overview</p>
				<a
					href="/"
					class="nav-link {isNavActive('/') ? 'nav-link-active' : ''}"
					onclick={() => (mobileNavOpen = false)}
				>
					<IconLayoutDashboard class="size-4.5 shrink-0" />
					<span class="flex-1">Dashboard</span>
				</a>
				<a
					href="/profile"
					class="nav-link {isNavActive('/profile') ? 'nav-link-active' : ''}"
					onclick={() => (mobileNavOpen = false)}
				>
					<IconUser class="size-4.5 shrink-0" />
					<span class="flex-1">My Profile</span>
				</a>
				<a
					href="/notifications"
					class="nav-link {isNavActive('/notifications') ? 'nav-link-active' : ''}"
					onclick={() => (mobileNavOpen = false)}
				>
					<IconBell class="size-4.5 shrink-0" />
					<span class="flex-1">Notifications</span>
					{#if unreadCount > 0}
						<Badge variant="default" class="text-[10px] font-bold px-1.5 py-0 h-4.5 bg-primary text-primary-foreground">
							{unreadCount}
						</Badge>
					{/if}
				</a>
			</div>

			<!-- Section: Academics & Operations -->
			<div class="space-y-1">
				<p class="px-3 text-[11px] font-bold uppercase tracking-wider text-muted-foreground">Academics & Hub</p>
				<a
					href="/academics"
					class="nav-link {isNavActive('/academics') ? 'nav-link-active' : ''}"
					onclick={() => (mobileNavOpen = false)}
				>
					<IconBook class="size-4.5 shrink-0" />
					<span class="flex-1">Academics</span>
				</a>
				<a
					href="/attendance"
					class="nav-link {isNavActive('/attendance') ? 'nav-link-active' : ''}"
					onclick={() => (mobileNavOpen = false)}
				>
					<IconUserCheck class="size-4.5 shrink-0" />
					<span class="flex-1">Attendance</span>
				</a>
				<a
					href="/calendar"
					class="nav-link {isNavActive('/calendar') ? 'nav-link-active' : ''}"
					onclick={() => (mobileNavOpen = false)}
				>
					<IconCalendar class="size-4.5 shrink-0" />
					<span class="flex-1">Calendar & Events</span>
				</a>
				<a
					href="/results"
					class="nav-link {isNavActive('/results') ? 'nav-link-active' : ''}"
					onclick={() => (mobileNavOpen = false)}
				>
					<IconCertificate class="size-4.5 shrink-0" />
					<span class="flex-1">Examination Results</span>
				</a>
			</div>

			<!-- Section: Student & Administrative Services -->
			<div class="space-y-1">
				<p class="px-3 text-[11px] font-bold uppercase tracking-wider text-muted-foreground">Services & Workflow</p>
				<a
					href="/applications"
					class="nav-link {isNavActive('/applications') ? 'nav-link-active' : ''}"
					onclick={() => (mobileNavOpen = false)}
				>
					<IconFileDescription class="size-4.5 shrink-0" />
					<span class="flex-1">Applications Portal</span>
				</a>
				<a
					href="/pdf-studio"
					class="nav-link {isNavActive('/pdf-studio') ? 'nav-link-active' : ''}"
					onclick={() => (mobileNavOpen = false)}
				>
					<IconCertificate class="size-4.5 shrink-0 text-primary" />
					<span class="flex-1">PDF Studio & Suite</span>
					<Badge variant="outline" class="text-[9px] px-1.5 py-0 h-4 bg-primary/10 text-primary border-primary/30 font-semibold">Pro</Badge>
				</a>
				{#if canPublish}
					<a
						href="/publishing"
						class="nav-link {isNavActive('/publishing') ? 'nav-link-active' : ''}"
						onclick={() => (mobileNavOpen = false)}
					>
						<IconFileText class="size-4.5 shrink-0" />
						<span class="flex-1">Publishing Desk</span>
					</a>
				{/if}
				{#if canApprove}
					<a
						href="/approvals"
						class="nav-link {isNavActive('/approvals') ? 'nav-link-active' : ''}"
						onclick={() => (mobileNavOpen = false)}
					>
						<IconClipboardCheck class="size-4.5 shrink-0" />
						<span class="flex-1">Role Approvals</span>
					</a>
				{/if}
				<a
					href="/mfa/setup"
					class="nav-link {isNavActive('/mfa/setup') ? 'nav-link-active' : ''}"
					onclick={() => (mobileNavOpen = false)}
				>
					<IconShieldLock class="size-4.5 shrink-0" />
					<span class="flex-1">Identity & MFA</span>
					{#if isMFAEnabled}
						<Badge variant="outline" class="text-[9px] px-1.5 py-0 h-4 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/30 font-semibold">Enabled</Badge>
					{:else}
						<Badge variant="outline" class="text-[9px] px-1.5 py-0 h-4 bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/30 font-semibold">Not Enabled</Badge>
					{/if}
				</a>
			</div>

			<!-- Section: Residential & Campus Life -->
			<div class="space-y-1">
				<p class="px-3 text-[11px] font-bold uppercase tracking-wider text-muted-foreground">Campus Life</p>
				<a
					href="/resources"
					class="nav-link {isNavActive('/resources') ? 'nav-link-active' : ''}"
					onclick={() => (mobileNavOpen = false)}
				>
					<IconFolder class="size-4.5 shrink-0" />
					<span class="flex-1">Course Resources & LMS</span>
				</a>
				<a
					href="/library"
					class="nav-link {isNavActive('/library') ? 'nav-link-active' : ''}"
					onclick={() => (mobileNavOpen = false)}
				>
					<IconBook class="size-4.5 shrink-0" />
					<span class="flex-1">E-Library & Resources</span>
				</a>
				<a
					href="/clubs"
					class="nav-link {isNavActive('/clubs') ? 'nav-link-active' : ''}"
					onclick={() => (mobileNavOpen = false)}
				>
					<IconUsersGroup class="size-4.5 shrink-0" />
					<span class="flex-1">Clubs & Societies</span>
				</a>
				<a
					href="/hostel"
					class="nav-link {isNavActive('/hostel') ? 'nav-link-active' : ''}"
					onclick={() => (mobileNavOpen = false)}
				>
					<IconBuildingCommunity class="size-4.5 shrink-0" />
					<span class="flex-1">Hostel & Housing</span>
				</a>
			</div>
		</div>

		<!-- 3. Footer (Fixed / Non-scrolling) -->
		<div class="shrink-0 border-t border-sidebar-border bg-sidebar-accent/30 p-3.5 space-y-3 pb-[max(0.875rem,env(safe-area-inset-bottom))]">
			<a href="/profile" class="flex items-center gap-2.5 min-w-0 p-1 -m-1 rounded-lg hover:bg-sidebar-accent transition-colors">
				<Avatar class="size-9 sm:size-10 border border-sidebar-border shadow-xs shrink-0 overflow-hidden">
					{#if cachedAvatar}
						<AvatarImage src={cachedAvatar} alt={user?.displayName ?? 'Profile'} class="size-full object-cover" />
					{/if}
					<AvatarFallback class="bg-primary/10 font-heading font-semibold text-primary text-xs sm:text-sm">
						{initials}
					</AvatarFallback>
				</Avatar>
				<div class="flex-1 min-w-0">
					<p class="truncate text-xs font-semibold leading-tight text-sidebar-foreground">
						{user?.displayName ?? 'Campus User'}
					</p>
					<p class="truncate text-[10px] sm:text-[11px] text-muted-foreground">{user?.email}</p>
				</div>
				<Badge variant="secondary" class="text-[9px] sm:text-[10px] font-bold shrink-0 px-1.5 py-0.5">
					{primaryRole}
				</Badge>
			</a>

			<Separator class="bg-sidebar-border" />

			<div class="flex items-center justify-between gap-2">
				<div class="flex items-center gap-1.5">
					<!-- Compact Theme Toggle Reusing /public/ design scaled down -->
					<div class="rounded-full border border-sidebar-border bg-background/90 p-0.5 shadow-2xs backdrop-blur-xs scale-[0.72] origin-left -mr-4.5">
						<Switch
							size="icon"
							checked={theme.isDark}
							onCheckedChange={(val) => theme.set(val)}
							aria-label="Use dark mode"
							title={theme.isDark ? 'Dark mode on' : 'Dark mode off'}
						>
							{#snippet thumbContent()}
								{#if theme.isDark}
									<IconMoon class="size-3 text-primary" aria-hidden="true" />
								{:else}
									<IconSun class="size-3 text-foreground" aria-hidden="true" />
								{/if}
							{/snippet}
						</Switch>
					</div>
					<span class="text-xs font-medium text-muted-foreground">{theme.isDark ? 'Dark' : 'Light'}</span>
				</div>

				<div class="flex items-center gap-1">
					<Button href="/public" variant="ghost" size="icon" class="size-8 text-muted-foreground hover:text-foreground" title="Public Campus Website">
						<IconExternalLink class="size-4" />
					</Button>
					<Button
						type="button"
						variant="ghost"
						size="icon"
						class="size-8 text-destructive hover:bg-destructive/10 cursor-pointer"
						title="Sign out"
						aria-label="Sign out"
						onclick={() => (isLogoutModalOpen = true)}
					>
						<IconLogout class="size-4" />
					</Button>
				</div>
			</div>
		</div>
	</aside>

	<!-- Main App Content Area (Detached from Sidebar via lg:pl-72) -->
	<div class="flex-1 flex flex-col min-w-0 lg:pl-72">
		<!-- Top App Header -->
		<header class="sticky top-0 z-30 flex h-14 sm:h-16 items-center justify-between border-b border-border bg-background/95 px-3 sm:px-6 lg:px-8 backdrop-blur-md shadow-2xs pt-[env(safe-area-inset-top,0px)]">
			<div class="flex items-center gap-2 sm:gap-3 min-w-0">
				<Button
					variant="outline"
					size="icon"
					class="lg:hidden size-8.5 sm:size-9 shrink-0"
					onclick={() => (mobileNavOpen = true)}
					aria-label="Open sidebar navigation"
				>
					<IconMenu2 class="size-4.5 sm:size-5" />
				</Button>

				<div class="flex items-center gap-1.5 sm:gap-2 text-xs text-muted-foreground min-w-0">
					<a href="/" class="hover:text-foreground shrink-0">Campus App</a>
					<IconChevronRight class="size-3.5 opacity-50 shrink-0" />
					<span class="font-medium text-foreground capitalize truncate max-w-[120px] min-[380px]:max-w-[180px] sm:max-w-none">
						{page.url.pathname === '/' ? 'Dashboard' : page.url.pathname.split('/')[1] || 'Dashboard'}
					</span>
				</div>
			</div>

			<div class="flex items-center gap-2 sm:gap-3 shrink-0">
				<Button href="/notifications" variant="ghost" size="icon" class="relative size-8.5 sm:size-9" aria-label="Notifications">
					<IconBell class="size-4 sm:size-4.5" />
					{#if unreadCount > 0}
						<span class="absolute right-2 top-2 size-2 rounded-full bg-primary ring-2 ring-background"></span>
					{/if}
				</Button>

				<a
					href="/profile"
					class="flex items-center gap-2 p-0.5 rounded-full hover:ring-2 hover:ring-primary/40 transition-all cursor-pointer"
					title="My Profile"
					aria-label="My Profile"
				>
					<Avatar class="size-8 sm:size-8.5 border border-border shadow-2xs overflow-hidden">
						{#if cachedAvatar}
							<AvatarImage src={cachedAvatar} alt={user?.displayName ?? 'Profile'} class="size-full object-cover" />
						{/if}
						<AvatarFallback class="bg-primary/10 font-heading font-semibold text-primary text-xs">
							{initials}
						</AvatarFallback>
					</Avatar>
				</a>
			</div>
		</header>

		<!-- Main Workspace View -->
		<main class="flex-1 px-3 py-4 min-[480px]:px-4 min-[480px]:py-6 sm:px-6 sm:py-8 lg:px-8 max-w-7xl w-full mx-auto min-w-0">
			{@render children()}
		</main>

		<!-- Footer -->
		<footer class="border-t border-border bg-card/40 px-3 py-4 sm:px-6 sm:py-6 lg:px-8 text-xs text-muted-foreground pb-[max(1rem,env(safe-area-inset-bottom))]">
			<div class="mx-auto flex max-w-7xl flex-col items-center justify-between gap-2 sm:gap-3 text-center sm:text-left sm:flex-row">
				<p class="flex items-center gap-1.5">
					<IconShieldCheck class="size-4 text-primary shrink-0" /> Scoped access enforced by BBDIT Institutional Policy.
				</p>
				<p>© 2026 Babu Banarsi Das Institute of Technology</p>
			</div>
		</footer>
	</div>
</div>

<!-- Logout Confirmation Modal -->
<AlertDialog.Root bind:open={isLogoutModalOpen}>
	<AlertDialog.Content class="w-[94vw] sm:max-w-md p-4 sm:p-6">
		<AlertDialog.Header class="space-y-1.5">
			<AlertDialog.Title class="text-base sm:text-lg font-bold text-destructive flex items-center gap-2">
				<IconLogout class="size-5 shrink-0" />
				<span>Log Out</span>
			</AlertDialog.Title>
			<AlertDialog.Description class="text-xs sm:text-sm text-muted-foreground leading-relaxed">
				Would you like to log out of your campus account? Any unsaved changes in active forms will be lost.
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer class="flex flex-col-reverse sm:flex-row gap-2 sm:gap-0 pt-3">
			<AlertDialog.Cancel class="h-9 text-xs cursor-pointer">No, Go back</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={handleSignOut}
				disabled={isLoggingOut}
				class="h-9 text-xs bg-destructive text-destructive-foreground hover:bg-destructive/90 cursor-pointer"
			>
				{#if isLoggingOut}
					<IconLoader2 class="size-3.5 animate-spin mr-1" />
				{/if}
				Yes, Log out
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

<style>
	.nav-link {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		border-radius: var(--radius);
		padding: 0.55rem 0.75rem;
		font-size: 0.84rem;
		font-weight: 500;
		color: var(--sidebar-foreground);
		transition: all 0.15s ease-in-out;
	}
	.nav-link:hover {
		background-color: var(--sidebar-accent);
		color: var(--sidebar-accent-foreground);
	}
	.nav-link-active {
		background-color: color-mix(in oklab, var(--primary) 12%, transparent);
		color: var(--primary);
		font-weight: 600;
	}
	.hide-scrollbar {
		scrollbar-width: none;
		-ms-overflow-style: none;
	}
	.hide-scrollbar::-webkit-scrollbar {
		display: none;
	}
</style>
