<script lang="ts">
	/**
	 * BLOCK_WEB_PAGE_HOME_001
	 * Purpose: The Hub Root Shell overview screen displaying accessible institutional workspaces.
	 */
	import { authState } from '$lib/auth-state.svelte';
	import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { buttonVariants } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';

	const workspaces = [
		{ id: 'student_hub', title: 'Student Hub', desc: 'Courses, timetable, attendance, assignments, mess dining', icon: '🎓', path: '/student' },
		{ id: 'teacher_hub', title: 'Teacher Hub', desc: 'Lecture rosters, gradebook, syllabus planner, leave requests', icon: '👨‍🏫', path: '/teacher' },
		{ id: 'admin_hub', title: 'Admin Hub', desc: 'User provisioning, fee structures, admissions KYC, accreditation', icon: '⚙️', path: '/admin' },
		{ id: 'librarian_hub', title: 'Librarian Hub', desc: 'Book cataloging, barcode scanning, issue/return checkout', icon: '📚', path: '/library' },
		{ id: 'warden_hub', title: 'Warden Hub', desc: 'Room allocations, gate pass approvals, curfew inspection logs', icon: '🏢', path: '/warden' },
		{ id: 'parent_hub', title: 'Parent Hub', desc: 'Fee payments, academic progress reports, attendance tracker', icon: '👨‍👩‍👧', path: '/parent' },
	];
</script>

<div class="space-y-8">
	<!-- Hero Section -->
	<Card class="p-4 sm:p-6 shadow-sm border-border">
		<div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-6">
			<div class="space-y-3">
				<Badge variant="outline" class="gap-1.5 bg-emerald-500/10 text-emerald-700 dark:text-emerald-300 border-emerald-500/30">
					<span class="inline-block h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
					Rank 1 IAM & Rank 2 Audit Active
				</Badge>
				<h1 class="text-3xl sm:text-4xl font-extrabold tracking-tight text-foreground">
					The Hub — Universal Campus Shell
				</h1>
				<p class="text-sm text-muted-foreground max-w-2xl leading-relaxed">
					Unified access gateway for 17 campus subsystems. Single sign-on with scoped personas,
					tamper-evident audit logging, and role-specialized workspaces.
				</p>
			</div>

			<div class="flex flex-wrap gap-2.5 shrink-0">
				<a
					href="/onboarding"
					class={cn(buttonVariants({ variant: 'default' }), 'font-semibold')}
				>
					Admissions KYC
				</a>
				<a
					href="/admin/onboarding"
					class={cn(buttonVariants({ variant: 'secondary' }), 'font-semibold')}
				>
					Admissions Desk
				</a>
				<a
					href="/audit"
					class={cn(buttonVariants({ variant: 'outline' }), 'font-semibold')}
				>
					Audit Ledger
				</a>
			</div>
		</div>
	</Card>

	<!-- Workspaces Grid -->
	<div class="space-y-4">
		<div class="flex items-center justify-between">
			<h2 class="text-lg font-bold tracking-tight text-foreground">Institutional Workspaces</h2>
			<Badge variant="secondary" class="font-mono text-xs">
				30,000 Capacity · 5,000 Peak Active
			</Badge>
		</div>

		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
			{#each workspaces as ws}
				<Card class="hover:border-primary/50 transition duration-200 flex flex-col justify-between shadow-sm">
					<CardHeader class="pb-3">
						<div class="text-3xl mb-2">{ws.icon}</div>
						<CardTitle class="text-base font-bold">{ws.title}</CardTitle>
						<CardDescription class="text-xs leading-relaxed">{ws.desc}</CardDescription>
					</CardHeader>

					<CardFooter class="pt-3 border-t border-border flex items-center justify-between">
						<span class="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Scoped Workspace</span>
						<a href={ws.path} class="text-xs font-bold text-primary hover:underline transition">
							Launch Workspace →
						</a>
					</CardFooter>
				</Card>
			{/each}
		</div>
	</div>
</div>
