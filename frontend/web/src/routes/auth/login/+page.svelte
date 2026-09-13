<script lang="ts">
	/**
	 * BLOCK_WEB_PAGE_LOGIN_001
	 * Purpose: Dedicated Sign-In page mounting LoginForm.svelte with redirection.
	 */
	import LoginForm from '$lib/components/auth/LoginForm.svelte';

	function handleSuccess() {
		window.location.href = '/';
	}
</script>

<div class="py-12 flex items-center justify-center">
	<LoginForm onSuccess={handleSuccess} />
</div>
