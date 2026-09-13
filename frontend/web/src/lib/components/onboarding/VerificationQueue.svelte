<script lang="ts">
	/**
	 * BLOCK_WEB_ONBOARDING_QUEUE_001
	 * Subsystem: Rank 3 - Student & Staff Registration System (onboarding)
	 * Purpose:   Registrar & HR admissions verification desk with document review and enrollment triggers.
	 */
	import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '$lib/components/ui/card';
	import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '$lib/components/ui/table';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import type { OnboardingApplicant, OnboardingStatus, OnboardingType } from '$lib/types/onboarding';

	let {
		tenantId = 'tenant-demo',
		adminId = 'usr_registrar_admin',
		applicants = [],
		isLoading = false,
		onRefresh = () => {}
	}: {
		tenantId?: string;
		adminId?: string;
		applicants: OnboardingApplicant[];
		isLoading?: boolean;
		onRefresh?: () => void;
	} = $props();

	let selectedType = $state<string>('ALL');
	let selectedStatus = $state<string>('ALL');
	let searchQuery = $state('');
	let selectedApplicant = $state<OnboardingApplicant | null>(null);
	let actionMessage = $state<string | null>(null);
	let isActionLoading = $state(false);

	let filteredApplicants = $derived(
		applicants.filter(app => {
			if (selectedType !== 'ALL' && app.type !== selectedType) return false;
			if (selectedStatus !== 'ALL' && app.status !== selectedStatus) return false;
			if (searchQuery.trim() !== '') {
				const q = searchQuery.toLowerCase();
				const matchesName = `${app.first_name} ${app.last_name}`.toLowerCase().includes(q);
				const matchesEmail = app.email.toLowerCase().includes(q);
				if (!matchesName && !matchesEmail) return false;
			}
			return true;
		})
	);

	let totalSubmitted = $derived(applicants.filter(a => a.status === 'SUBMITTED').length);
	let totalUnderReview = $derived(applicants.filter(a => a.status === 'UNDER_REVIEW').length);
	let totalVerified = $derived(applicants.filter(a => a.status === 'VERIFIED').length);
	let totalEnrolled = $derived(applicants.filter(a => a.status === 'ENROLLED').length);

	async function handleAssignReview(applicantId: string) {
		isActionLoading = true;
		actionMessage = null;
		try {
			const res = await fetch(`/api/v1/onboarding/applicants/${applicantId}/review`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					tenant_id: tenantId,
					reviewer_id: adminId,
				}),
			});
			if (!res.ok) throw new Error('Failed to start review');
			actionMessage = `Review started for applicant ${applicantId}`;
			onRefresh();
		} catch (err: unknown) {
			actionMessage = err instanceof Error ? err.message : 'Error updating status';
		} finally {
			isActionLoading = false;
		}
	}

	async function handleVerifyDoc(applicantId: string, docId: string, approved: boolean) {
		isActionLoading = true;
		try {
			const res = await fetch(`/api/v1/onboarding/applicants/${applicantId}/documents/${docId}`, {
				method: 'PATCH',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					tenant_id: tenantId,
					approved,
					verifier_id: adminId,
					rejection_reason: approved ? undefined : 'Unreadable or invalid documentation',
				}),
			});
			if (!res.ok) throw new Error('Failed to update document status');
			// Update local modal state if open
			if (selectedApplicant) {
				const doc = selectedApplicant.documents.find(d => d.id === docId);
				if (doc) doc.status = approved ? 'VERIFIED' : 'REJECTED';
			}
			onRefresh();
		} catch (err: unknown) {
			actionMessage = err instanceof Error ? err.message : 'Error verifying document';
		} finally {
			isActionLoading = false;
		}
	}

	async function handleVerifyApplication(applicantId: string) {
		isActionLoading = true;
		actionMessage = null;
		try {
			const res = await fetch(`/api/v1/onboarding/applicants/${applicantId}/verify?tenant_id=${tenantId}`, {
				method: 'POST',
			});
			if (!res.ok) {
				const errData = await res.json().catch(() => ({ detail: 'Verification failed' }));
				throw new Error(errData.detail || 'Verification failed');
			}
			actionMessage = 'KYC successfully verified! Ready for enrollment.';
			selectedApplicant = null;
			onRefresh();
		} catch (err: unknown) {
			actionMessage = err instanceof Error ? err.message : 'Error verifying application';
		} finally {
			isActionLoading = false;
		}
	}

	async function handleEnroll(applicant: OnboardingApplicant) {
		isActionLoading = true;
		actionMessage = null;
		try {
			if (applicant.type === 'STUDENT') {
				const res = await fetch(`/api/v1/onboarding/applicants/${applicant.id}/enroll`, {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({
						tenant_id: tenantId,
						user_id: `usr_student_${applicant.first_name.toLowerCase()}`,
						cohort_id: 'cohort-2026-cse-a',
						admin_id: adminId,
					}),
				});
				if (!res.ok) throw new Error('Enrollment failed');
				const data = await res.json();
				actionMessage = `Student Enrolled! Generated Roll Number: ${data.roll_number}`;
			} else {
				const res = await fetch(`/api/v1/onboarding/applicants/${applicant.id}/provision`, {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({
						tenant_id: tenantId,
						user_id: `usr_staff_${applicant.first_name.toLowerCase()}`,
						role_key: applicant.type.toLowerCase(),
						admin_id: adminId,
					}),
				});
				if (!res.ok) throw new Error('Staff provisioning failed');
				const data = await res.json();
				actionMessage = `Staff Provisioned! Generated Employee ID: ${data.employee_id}`;
			}
			onRefresh();
		} catch (err: unknown) {
			actionMessage = err instanceof Error ? err.message : 'Error during enrollment';
		} finally {
			isActionLoading = false;
		}
	}

	function getStatusBadgeVariant(status: OnboardingStatus): 'default' | 'secondary' | 'outline' | 'destructive' {
		switch (status) {
			case 'ENROLLED': return 'default';
			case 'VERIFIED': return 'outline';
			case 'UNDER_REVIEW': return 'secondary';
			case 'REJECTED': return 'destructive';
			default: return 'outline';
		}
	}
</script>

<div class="space-y-6">
	<!-- KPI Summary Metrics -->
	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
		<Card class="shadow-sm">
			<CardContent class="p-5">
				<span class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">New Submissions</span>
				<div class="mt-2 text-3xl font-bold tracking-tight text-foreground">{totalSubmitted}</div>
				<div class="mt-1 text-xs text-muted-foreground">Pending review assignment</div>
			</CardContent>
		</Card>

		<Card class="shadow-sm">
			<CardContent class="p-5">
				<span class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Under Review</span>
				<div class="mt-2 text-3xl font-bold tracking-tight text-amber-600 dark:text-amber-400">{totalUnderReview}</div>
				<div class="mt-1 text-xs text-muted-foreground">Document verification in progress</div>
			</CardContent>
		</Card>

		<Card class="shadow-sm">
			<CardContent class="p-5">
				<span class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Verified Ready</span>
				<div class="mt-2 text-3xl font-bold tracking-tight text-blue-600 dark:text-blue-400">{totalVerified}</div>
				<div class="mt-1 text-xs text-muted-foreground">Ready for Roll No / ID assignment</div>
			</CardContent>
		</Card>

		<Card class="shadow-sm">
			<CardContent class="p-5">
				<span class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Total Enrolled</span>
				<div class="mt-2 text-3xl font-bold tracking-tight text-emerald-600 dark:text-emerald-400">{totalEnrolled}</div>
				<div class="mt-1 text-xs text-muted-foreground">Active profiles provisioned</div>
			</CardContent>
		</Card>
	</div>

	{#if actionMessage}
		<div class="p-4 rounded-lg bg-primary/10 border border-primary/20 text-primary text-sm font-medium">
			ℹ️ {actionMessage}
		</div>
	{/if}

	<!-- Filter Controls & Search -->
	<Card class="shadow-sm">
		<CardHeader class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 pb-4">
			<div>
				<CardTitle class="text-lg font-bold tracking-tight">Admissions Verification Queue</CardTitle>
				<CardDescription>Review applicant KYC submissions, verify documents, and assign roll numbers.</CardDescription>
			</div>

			<div class="flex flex-wrap items-center gap-2.5 w-full md:w-auto">
				<Input
					bind:value={searchQuery}
					placeholder="Search applicant name or email..."
					class="w-full md:w-64 h-9"
				/>

				<select
					bind:value={selectedType}
					class="h-9 px-3 rounded-md border border-input bg-background text-foreground text-xs font-medium focus:outline-none focus:ring-2 focus:ring-ring"
				>
					<option value="ALL">All Categories</option>
					<option value="STUDENT">Students</option>
					<option value="FACULTY">Faculty</option>
					<option value="STAFF">Staff</option>
				</select>

				<select
					bind:value={selectedStatus}
					class="h-9 px-3 rounded-md border border-input bg-background text-foreground text-xs font-medium focus:outline-none focus:ring-2 focus:ring-ring"
				>
					<option value="ALL">All Statuses</option>
					<option value="SUBMITTED">Submitted</option>
					<option value="UNDER_REVIEW">Under Review</option>
					<option value="VERIFIED">Verified</option>
					<option value="ENROLLED">Enrolled</option>
					<option value="REJECTED">Rejected</option>
				</select>

				<Button size="sm" variant="outline" onclick={onRefresh} disabled={isLoading}>
					{isLoading ? 'Refreshing...' : 'Refresh'}
				</Button>
			</div>
		</CardHeader>

		<CardContent class="p-0">
			<Table>
				<TableHeader>
					<TableRow>
						<TableHead class="w-[180px]">Applicant</TableHead>
						<TableHead>Category</TableHead>
						<TableHead>Program / Dept</TableHead>
						<TableHead>Submission Date</TableHead>
						<TableHead>Status</TableHead>
						<TableHead class="text-right">Actions</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					{#if filteredApplicants.length === 0}
						<TableRow>
							<TableCell colspan={6} class="text-center py-8 text-muted-foreground text-sm">
								No applicants match the selected criteria.
							</TableCell>
						</TableRow>
					{:else}
						{#each filteredApplicants as app}
							<TableRow>
								<TableCell>
									<div class="font-medium text-foreground">{app.first_name} {app.last_name}</div>
									<div class="text-xs text-muted-foreground">{app.email}</div>
								</TableCell>
								<TableCell>
									<Badge variant="outline" class="text-[10px] uppercase font-mono">{app.type}</Badge>
								</TableCell>
								<TableCell class="text-xs text-foreground">
									{app.type === 'STUDENT' ? (app.program_id || 'BTECH_CSE') : (app.department_id || 'CSE')}
								</TableCell>
								<TableCell class="text-xs text-muted-foreground">
									{new Date(app.created_at).toLocaleDateString()}
								</TableCell>
								<TableCell>
									<Badge variant={getStatusBadgeVariant(app.status)} class="text-xs">
										{app.status}
									</Badge>
								</TableCell>
								<TableCell class="text-right space-x-1.5">
									{#if app.status === 'SUBMITTED'}
										<Button size="sm" onclick={() => handleAssignReview(app.id)} disabled={isActionLoading}>
											Start Review
										</Button>
									{:else if app.status === 'UNDER_REVIEW'}
										<Button size="sm" variant="secondary" onclick={() => selectedApplicant = app}>
											Verify KYC ({app.documents?.length || 0} Docs)
										</Button>
									{:else if app.status === 'VERIFIED'}
										<Button size="sm" class="bg-emerald-600 hover:bg-emerald-700 text-white" onclick={() => handleEnroll(app)} disabled={isActionLoading}>
											Enroll Profile →
										</Button>
									{:else if app.status === 'ENROLLED'}
										<Badge variant="outline" class="bg-emerald-500/10 text-emerald-700 dark:text-emerald-300 border-emerald-500/30">
											✓ Active ID
										</Badge>
									{/if}
								</TableCell>
							</TableRow>
						{/each}
					{/if}
				</TableBody>
			</Table>
		</CardContent>
	</Card>

	<!-- Document Verification Modal / Drawer -->
	{#if selectedApplicant}
		<Card class="border-primary/50 shadow-md">
			<CardHeader class="flex flex-row items-center justify-between pb-3">
				<div>
					<CardTitle class="text-lg font-bold">
						KYC Verification: {selectedApplicant.first_name} {selectedApplicant.last_name}
					</CardTitle>
					<CardDescription class="text-xs">
						Review uploaded proofs and approve/reject each mandatory document.
					</CardDescription>
				</div>
				<Button size="sm" variant="outline" onclick={() => selectedApplicant = null}>Close</Button>
			</CardHeader>

			<CardContent class="space-y-4">
				<div class="grid grid-cols-1 sm:grid-cols-3 gap-3 p-3 rounded-lg bg-muted/50 text-xs">
					<div><strong>Category:</strong> {selectedApplicant.type}</div>
					<div><strong>Academic Year:</strong> {selectedApplicant.academic_year}</div>
					<div><strong>Phone:</strong> {selectedApplicant.phone}</div>
				</div>

				<div class="space-y-3">
					<h4 class="text-xs font-bold uppercase tracking-wider text-muted-foreground">Document Verification Checklist</h4>
					{#each selectedApplicant.documents as doc}
						<div class="p-3 rounded-lg border border-border bg-card flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
							<div class="space-y-0.5">
								<div class="flex items-center gap-2">
									<span class="text-sm font-semibold text-foreground">{doc.file_name}</span>
									<Badge variant="outline" class="text-[10px] font-mono">{doc.document_type}</Badge>
								</div>
								<span class="text-xs text-muted-foreground">Status: <strong>{doc.status}</strong></span>
							</div>

							<div class="flex items-center gap-2">
								{#if doc.status === 'VERIFIED'}
									<Badge variant="outline" class="bg-emerald-500/10 text-emerald-700 dark:text-emerald-300 border-emerald-500/30">
										✓ Verified
									</Badge>
								{:else}
									<Button size="sm" variant="outline" class="text-emerald-600 hover:bg-emerald-50" onclick={() => handleVerifyDoc(selectedApplicant!.id, doc.id, true)}>
										Approve
									</Button>
									<Button size="sm" variant="outline" class="text-destructive hover:bg-destructive/10" onclick={() => handleVerifyDoc(selectedApplicant!.id, doc.id, false)}>
										Reject
									</Button>
								{/if}
							</div>
						</div>
					{/each}
				</div>

				<div class="pt-3 border-t border-border flex justify-end gap-2">
					<Button
						disabled={isActionLoading || selectedApplicant.documents.some(d => d.status !== 'VERIFIED')}
						onclick={() => handleVerifyApplication(selectedApplicant!.id)}
					>
						Complete KYC Verification ✓
					</Button>
				</div>
			</CardContent>
		</Card>
	{/if}
</div>
