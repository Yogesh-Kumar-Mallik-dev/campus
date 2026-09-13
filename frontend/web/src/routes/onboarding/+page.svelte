<script lang="ts">
	/**
	 * BLOCK_WEB_PAGE_ONBOARDING_PUBLIC_001
	 * Purpose: Public Admissions & Staff Registration KYC Portal.
	 */
	import KYCApplicationForm from '$lib/components/onboarding/KYCApplicationForm.svelte';
	import { authState } from '$lib/auth-state.svelte';

	let tenantId = $derived(authState.user?.tenant_id || 'tenant-demo');
</script>

<div class="max-w-4xl mx-auto space-y-8">
	<div class="text-center space-y-2">
		<h1 class="text-3xl font-extrabold tracking-tight text-foreground">
			Campus Admissions & Registration Portal
		</h1>
		<p class="text-sm text-muted-foreground max-w-xl mx-auto">
			Submit your digital KYC profile, attach mandatory credentials, and initiate fast-track institutional enrollment.
		</p>
	</div>

	<KYCApplicationForm {tenantId} />
</div>
