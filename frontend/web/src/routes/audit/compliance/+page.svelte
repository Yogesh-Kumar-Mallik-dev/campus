<script lang="ts">
	/**
	 * BLOCK_WEB_PAGE_COMPLIANCE_001
	 * Purpose: Dedicated Institutional Accreditation & Compliance Reporting page.
	 */
	import { onMount } from 'svelte';
	import { apiClient } from '@campus/api-client';
	import ComplianceReportView, { type ComplianceReport } from '$lib/components/audit/ComplianceReportView.svelte';

	let report = $state<ComplianceReport | null>(null);
	let isLoading = $state(false);

	async function fetchReport(framework: string = 'NAAC') {
		isLoading = true;
		try {
			const res = await apiClient.get<any>(`/audit/compliance-reports?tenant_id=tenant_default&framework=${framework}`);
			report = res.data || null;
		} catch (err) {
			console.error('Failed to compile compliance report:', err);
		} finally {
			isLoading = false;
		}
	}

	onMount(() => {
		fetchReport('NAAC');
	});
</script>

<div class="space-y-6">
	<ComplianceReportView
		{report}
		{isLoading}
		onGenerate={fetchReport}
	/>
</div>
