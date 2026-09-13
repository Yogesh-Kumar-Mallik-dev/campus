<script lang="ts">
	/**
	 * BLOCK_WEB_PAGE_AUDIT_001
	 * Purpose: Dedicated Audit Ledger Explorer page with live API query & chain verification.
	 */
	import { onMount } from 'svelte';
	import { apiClient } from '@campus/api-client';
	import AuditLogTable, { type AuditLogItem } from '$lib/components/audit/AuditLogTable.svelte';

	let logs = $state<AuditLogItem[]>([]);
	let total = $state(0);
	let isLoading = $state(false);
	let verificationStatus = $state<{ isIntact: boolean; count: number } | null>(null);

	async function loadLogs(filter: { action?: string; status?: string } = {}) {
		isLoading = true;
		try {
			const queryParams = new URLSearchParams({
				tenant_id: 'tenant_default',
				limit: '50',
				offset: '0'
			});
			if (filter.action) queryParams.set('action', filter.action);
			if (filter.status) queryParams.set('status', filter.status);

			const res = await apiClient.get<any>(`/audit/logs?${queryParams.toString()}`);
			logs = res.data || [];
			total = res.total || logs.length;
		} catch (err) {
			console.error('Failed to load audit logs:', err);
		} finally {
			isLoading = false;
		}
	}

	async function handleVerifyChain() {
		try {
			const res = await apiClient.post<any>('/audit/verify', {
				tenant_id: 'tenant_default'
			});
			if (res.data) {
				verificationStatus = {
					isIntact: res.data.is_chain_intact,
					count: res.data.total_verified
				};
			}
		} catch (err) {
			alert('Chain verification failed to execute.');
		}
	}

	onMount(() => {
		loadLogs();
	});
</script>

<div class="space-y-6">
	<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Immutable Audit Ledger</h1>
			<p class="text-sm text-slate-500 mt-1">
				Real-time tamper-evident event stream with SHA-256 Merkle chain verification.
			</p>
		</div>

		{#if verificationStatus}
			<div class="px-4 py-2 rounded-xl text-xs font-bold {verificationStatus.isIntact ? 'bg-emerald-100 text-emerald-800' : 'bg-red-100 text-red-800'} flex items-center gap-2 shadow-sm">
				<span>{verificationStatus.isIntact ? '✓' : '⚠'}</span>
				<span>{verificationStatus.isIntact ? `Hash Chain Intact (${verificationStatus.count} verified)` : 'Tampering Detected in Ledger!'}</span>
			</div>
		{/if}
	</div>

	<AuditLogTable
		{logs}
		{total}
		{isLoading}
		onVerifyChain={handleVerifyChain}
		onFilterChange={loadLogs}
	/>
</div>
