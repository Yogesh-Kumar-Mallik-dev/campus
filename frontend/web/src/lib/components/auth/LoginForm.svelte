<script lang="ts">
  /**
   * BLOCK_WEB_LOGIN_FORM_001
   * Purpose: Viewport-resilient authentication form with TOTP step and error handling.
   * Viewports: 280px to 4K responsive scaling.
   */
  import { authState } from '$lib/auth-state.svelte';

  let { onSuccess = () => {} }: { onSuccess?: () => void } = $props();

  let tenantId = $state('tenant_default');
  let identifier = $state('');
  let password = $state('');
  let mfaCode = $state('');
  let isMfaStep = $state(false);
  let localError = $state<string | null>(null);

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    localError = null;

    try {
      const result = await authState.login(tenantId, identifier, password, isMfaStep ? mfaCode : undefined);
      if (result.requiresMFA) {
        isMfaStep = true;
      } else {
        onSuccess();
      }
    } catch (err: any) {
      localError = authState.errorMessage || 'Authentication failed. Please check your credentials.';
    }
  }
</script>

<div class="w-full max-w-[min(420px,100%)] mx-auto p-6 sm:p-8 rounded-xl border border-slate-200 bg-white shadow-sm">
  <div class="mb-6 text-center">
    <h1 class="text-2xl font-bold tracking-tight text-slate-900">
      {isMfaStep ? 'Two-Factor Authentication' : 'Campus Portal Sign In'}
    </h1>
    <p class="text-sm text-slate-500 mt-1">
      {isMfaStep ? 'Enter the 6-digit code from your authenticator app' : 'Enter your institutional email or username'}
    </p>
  </div>

  {#if localError}
    <div class="mb-4 p-3 rounded-lg bg-red-50 border border-red-200 text-sm text-red-700 font-medium">
      {localError}
    </div>
  {/if}

  <form onsubmit={handleSubmit} class="space-y-4">
    {#if !isMfaStep}
      <div>
        <label for="identifier" class="block text-xs font-semibold uppercase tracking-wider text-slate-700 mb-1">
          Email or Username
        </label>
        <input
          id="identifier"
          type="text"
          bind:value={identifier}
          required
          placeholder="student@campus.edu or username"
          class="w-full h-11 px-3.5 rounded-lg border border-slate-300 text-slate-900 placeholder-slate-400 text-sm focus:outline-none focus:ring-2 focus:ring-slate-900 focus:border-transparent transition"
        />
      </div>

      <div>
        <div class="flex items-center justify-between mb-1">
          <label for="password" class="block text-xs font-semibold uppercase tracking-wider text-slate-700">
            Password
          </label>
          <a href="/forgot-password" class="text-xs text-blue-600 hover:underline">Forgot password?</a>
        </div>
        <input
          id="password"
          type="password"
          bind:value={password}
          required
          placeholder="••••••••••••"
          class="w-full h-11 px-3.5 rounded-lg border border-slate-300 text-slate-900 placeholder-slate-400 text-sm focus:outline-none focus:ring-2 focus:ring-slate-900 focus:border-transparent transition"
        />
      </div>
    {:else}
      <div>
        <label for="mfaCode" class="block text-xs font-semibold uppercase tracking-wider text-slate-700 mb-1">
          TOTP Passcode
        </label>
        <input
          id="mfaCode"
          type="text"
          inputmode="numeric"
          maxlength="6"
          bind:value={mfaCode}
          required
          placeholder="123456"
          class="w-full h-12 text-center tracking-[0.5em] font-mono text-xl rounded-lg border border-slate-300 text-slate-900 focus:outline-none focus:ring-2 focus:ring-slate-900 focus:border-transparent transition"
        />
      </div>
    {/if}

    <button
      type="submit"
      disabled={authState.isLoading}
      class="w-full h-11 bg-slate-900 hover:bg-slate-800 active:bg-slate-950 text-white font-medium text-sm rounded-lg transition duration-150 flex items-center justify-center disabled:opacity-50 min-h-[44px]"
    >
      {#if authState.isLoading}
        <span class="animate-spin mr-2">⟳</span> Authenticating...
      {:else}
        {isMfaStep ? 'Verify & Continue' : 'Sign In to Campus'}
      {/if}
    </button>
  </form>
</div>
