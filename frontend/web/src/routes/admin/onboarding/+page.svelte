<script lang="ts">
	/**
	 * BLOCK_WEB_PAGE_ONBOARDING_ADMIN_001
	 * Purpose: Registrar and Admissions Admin Desk for KYC Verification and Cohort Enrollment.
	 */
	import VerificationQueue from '$lib/components/onboarding/VerificationQueue.svelte';
	import { authState } from '$lib/auth-state.svelte';
	import type { OnboardingApplicant } from '$lib/types/onboarding';

	let tenantId = $derived(authState.user?.tenant_id || 'tenant-demo');
	let adminId = $derived(authState.user?.id || 'usr_admissions_admin');

	let applicants = $state<OnboardingApplicant[]>([
		{
			id: 'app_001_demo',
			tenant_id: 'tenant-demo',
			type: 'STUDENT',
			status: 'SUBMITTED',
			first_name: 'Aarav',
			last_name: 'Sharma',
			email: 'aarav.sharma@example.com',
			phone: '+91 98765 43210',
			date_of_birth: '2006-03-15',
			gender: 'MALE',
			nationality: 'Indian',
			program_id: 'BTECH_CSE',
			academic_year: '2026-2027',
			emergency_contact: { name: 'Rajesh Sharma', phone: '+91 98765 43211', relation: 'Father' },
			address: { line1: '42 Tech Enclave', city: 'Bengaluru', state: 'Karnataka', postal_code: '560001', country: 'India' },
			documents: [
				{ id: 'doc_1', tenant_id: 'tenant-demo', applicant_id: 'app_001_demo', document_type: 'NATIONAL_ID', file_key: 'uploads/id.pdf', file_name: 'aadhaar_card.pdf', file_size: 150000, mime_type: 'application/pdf', status: 'PENDING_REVIEW', created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
				{ id: 'doc_2', tenant_id: 'tenant-demo', applicant_id: 'app_001_demo', document_type: 'ACADEMIC_TRANSCRIPT', file_key: 'uploads/tr.pdf', file_name: 'marksheet_12th.pdf', file_size: 450000, mime_type: 'application/pdf', status: 'PENDING_REVIEW', created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
				{ id: 'doc_3', tenant_id: 'tenant-demo', applicant_id: 'app_001_demo', document_type: 'PASSPORT_PHOTO', file_key: 'uploads/ph.jpg', file_name: 'aarav_photo.jpg', file_size: 80000, mime_type: 'image/jpeg', status: 'PENDING_REVIEW', created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
			],
			created_at: new Date(Date.now() - 3600000 * 2).toISOString(),
			updated_at: new Date().toISOString(),
		},
		{
			id: 'app_002_demo',
			tenant_id: 'tenant-demo',
			type: 'FACULTY',
			status: 'UNDER_REVIEW',
			first_name: 'Dr. Ananya',
			last_name: 'Iyer',
			email: 'ananya.iyer@example.com',
			phone: '+91 98765 43333',
			date_of_birth: '1985-11-22',
			gender: 'FEMALE',
			nationality: 'Indian',
			department_id: 'CSE',
			academic_year: '2026-2027',
			target_designation: 'Associate Professor',
			emergency_contact: { name: 'Karthik Iyer', phone: '+91 98765 43334', relation: 'Spouse' },
			address: { line1: '88 Faculty Avenue', city: 'Chennai', state: 'Tamil Nadu', postal_code: '600025', country: 'India' },
			documents: [
				{ id: 'doc_4', tenant_id: 'tenant-demo', applicant_id: 'app_002_demo', document_type: 'NATIONAL_ID', file_key: 'uploads/id.pdf', file_name: 'passport.pdf', file_size: 120000, mime_type: 'application/pdf', status: 'VERIFIED', created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
				{ id: 'doc_5', tenant_id: 'tenant-demo', applicant_id: 'app_002_demo', document_type: 'DEGREE_CERTIFICATE', file_key: 'uploads/phd.pdf', file_name: 'phd_doctorate.pdf', file_size: 600000, mime_type: 'application/pdf', status: 'VERIFIED', created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
				{ id: 'doc_6', tenant_id: 'tenant-demo', applicant_id: 'app_002_demo', document_type: 'PASSPORT_PHOTO', file_key: 'uploads/ph2.jpg', file_name: 'photo_iyer.jpg', file_size: 70000, mime_type: 'image/jpeg', status: 'VERIFIED', created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
			],
			created_at: new Date(Date.now() - 3600000 * 24).toISOString(),
			updated_at: new Date().toISOString(),
		}
	]);

	let isLoading = $state(false);

	async function fetchApplicants() {
		isLoading = true;
		try {
			const res = await fetch(`/api/v1/onboarding/applicants?tenant_id=${tenantId}&limit=50`);
			if (res.ok) {
				const data = await res.json();
				if (data.applicants && data.applicants.length > 0) {
					applicants = data.applicants;
				}
			}
		} catch (err) {
			console.error('Failed to load queue', err);
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="space-y-8">
	<div class="space-y-1">
		<h1 class="text-2xl sm:text-3xl font-extrabold tracking-tight text-foreground">
			Admissions & Staff Enrollment Desk
		</h1>
		<p class="text-sm text-muted-foreground">
			Accredited verification portal for applicant KYC approvals, credential audits, and deterministic roll number allocation.
		</p>
	</div>

	<VerificationQueue
		{tenantId}
		{adminId}
		{applicants}
		{isLoading}
		onRefresh={fetchApplicants}
	/>
</div>
