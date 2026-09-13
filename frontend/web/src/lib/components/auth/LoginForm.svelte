<script lang="ts">
  /**
   * BLOCK_WEB_LOGIN_FORM_001
   * Purpose: Viewport-resilient authentication form with TOTP step and error handling.
   * Viewports: 280px to 4K responsive scaling with shadcn-svelte primitives.
   */
  import { authState } from '$lib/auth-state.svelte';
  import * as Card from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';

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

<div class="w-full max-w-[min(440px,100%)] mx-auto">
  <Card.Root class="shadow-sm border-border bg-card">
    <Card.Header class="text-center pb-4">
      <Card.Title class="text-2xl font-bold tracking-tight text-foreground">
        {isMfaStep ? 'Two-Factor Verification' : 'Institutional Sign In'}
      </Card.Title>
      <Card.Description class="text-sm text-muted-foreground mt-1">
        {isMfaStep ? 'Enter the 6-digit TOTP code from your authenticator app' : 'Enter your institutional email or username'}
      </Card.Description>
    </Card.Header>

    <Card.Content class="space-y-4">
      {#if localError}
        <div class="p-3 rounded-lg bg-destructive/10 border border-destructive/20 text-xs font-semibold text-destructive">
          {localError}
        </div>
      {/if}

      <form onsubmit={handleSubmit} class="space-y-4">
        {#if !isMfaStep}
          <div class="space-y-1.5">
            <Label for="identifier" class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Institutional Email / Username
            </Label>
            <Input
              id="identifier"
              type="text"
              bind:value={identifier}
              required
              placeholder="student@campus.edu or username"
              class="h-11 text-sm bg-background"
            />
          </div>

          <div class="space-y-1.5">
            <div class="flex items-center justify-between">
              <Label for="password" class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                Password
              </Label>
              <a href="/forgot-password" class="text-xs text-primary hover:underline">Forgot password?</a>
            </div>
            <Input
              id="password"
              type="password"
              bind:value={password}
              required
              placeholder="••••••••••••"
              class="h-11 text-sm bg-background"
            />
          </div>
        {:else}
          <div class="space-y-1.5">
            <Label for="mfaCode" class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              6-Digit Passcode
            </Label>
            <Input
              id="mfaCode"
              type="text"
              inputmode="numeric"
              maxlength={6}
              bind:value={mfaCode}
              required
              placeholder="123456"
              class="h-14 text-center tracking-[0.4em] font-mono text-xl bg-background"
            />
          </div>
        {/if}

        <Button
          type="submit"
          disabled={authState.isLoading}
          class="w-full h-11 text-sm font-semibold mt-2"
        >
          {#if authState.isLoading}
            <span class="animate-spin mr-2">⟳</span> Authenticating...
          {:else}
            {isMfaStep ? 'Verify & Continue' : 'Sign In to Campus'}
          {/if}
        </Button>
      </form>
    </Card.Content>
  </Card.Root>
</div>
