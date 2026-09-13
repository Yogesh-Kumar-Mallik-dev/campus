<script lang="ts">
	/**
	 * BLOCK_WEB_PAGE_HOME_001
	 * Purpose: The Hub Root Shell overview screen displaying accessible institutional workspaces.
	 */
	import { authState } from '$lib/auth-state.svelte';

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
	<div class="p-8 bg-white rounded-2xl border border-slate-200 shadow-sm flex flex-col md:flex-row items-start md:items-center justify-between gap-6">
		<div class="space-y-2">
			<div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-emerald-100 text-emerald-800">
				<span>●</span> Rank 1 IAM & Rank 2 Audit Active
			</div>
			<h1 class="text-3xl font-extrabold tracking-tight text-slate-900">
				The Hub — Universal Campus Shell
			</h1>
			<p class="text-sm text-slate-600 max-w-2xl">
				Unified access gateway for 17 campus subsystems. Single sign-on with scoped personas,
				tamper-evident audit logging, and role-specialized workspaces.
			</p>
		</div>

		<div class="flex flex-col sm:flex-row gap-3">
			<a
				href="/audit"
				class="h-11 px-5 rounded-xl bg-slate-900 hover:bg-slate-800 text-white text-xs font-bold uppercase tracking-wider flex items-center justify-center transition shadow-sm"
			>
				Explore Audit Ledger
			</a>
			<a
				href="/audit/compliance"
				class="h-11 px-5 rounded-xl border border-slate-300 hover:bg-slate-100 text-slate-800 text-xs font-bold uppercase tracking-wider flex items-center justify-center transition"
			>
				Accreditation Reports
			</a>
		</div>
	</div>

	<!-- Workspaces Grid -->
	<div class="space-y-4">
		<div class="flex items-center justify-between">
			<h2 class="text-lg font-bold text-slate-900">Institutional Workspaces</h2>
			<span class="text-xs text-slate-500 font-mono">30,000 Total Capacity · 5,000 Peak Active</span>
		</div>

		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
			{#each workspaces as ws}
				<div class="p-6 bg-white rounded-xl border border-slate-200 shadow-sm hover:border-slate-400 transition flex flex-col justify-between space-y-4">
					<div class="space-y-2">
						<div class="text-3xl">{ws.icon}</div>
						<h3 class="text-base font-bold text-slate-900">{ws.title}</h3>
						<p class="text-xs text-slate-500 leading-relaxed">{ws.desc}</p>
					</div>

					<div class="pt-2 border-t border-slate-100 flex items-center justify-between">
						<span class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Scoped Workspace</span>
						<a href={ws.path} class="text-xs font-bold text-blue-600 hover:text-blue-800 transition">
							Launch Workspace →
						</a>
					</div>
				</div>
			{/each}
		</div>
	</div>
</div>
