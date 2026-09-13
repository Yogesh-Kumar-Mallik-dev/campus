<script lang="ts">
	import { dev } from '$app/environment';
	import { enhance } from '$app/forms';
	import * as Tabs from '$lib/components/ui/tabs';
	import * as Select from '$lib/components/ui/select';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { toast } from '$lib/components/ui/sonner';
	import {
		IconArrowRight,
		IconId,
		IconLock,
		IconMail,
		IconSparkles,
		IconTestPipe,
		IconLoader2
	} from '@tabler/icons-svelte';

	type Mode = 'sign-in' | 'sign-up';
	let { initialMode, form }: { initialMode: Mode; form?: { error?: string } | null } = $props();
	let mode = $state<Mode>('sign-in');
	let role = $state('STUDENT');
	let emailInput = $state('');
	let passwordInput = $state('');
	let isSubmitting = $state(false);
	let submitType = $state<'login' | 'dev-login' | 'signup' | null>(null);

	function fillDevCredentials() {
		emailInput = 'dev@campus.internal';
		passwordInput = 'DevPassword123!';
		toast.info('Development credentials filled.');
	}

	$effect(() => {
		mode = initialMode;
	});
</script>

<Tabs.Root bind:value={mode}>
	<Tabs.List class="grid h-11 w-full grid-cols-2 p-1">
		<Tabs.Trigger value="sign-in">Log in</Tabs.Trigger>
		<Tabs.Trigger value="sign-up">Sign up</Tabs.Trigger>
	</Tabs.List>

	<Tabs.Content value="sign-in" class="mt-6">
		<div class="border-b pb-4">
			<h2 class="font-heading text-xl font-semibold">Welcome back</h2>
			<p class="mt-1 text-sm text-muted-foreground">Continue with your institutional account.</p>
		</div>

		<form
			method="POST"
			action="/public/sign-in"
			class="mt-6 grid gap-5"
			use:enhance={() => {
				isSubmitting = true;
				submitType = 'login';
				toast.loading('Authenticating institutional credentials...', { id: 'auth-toast' });
				return async ({ result, update }) => {
					isSubmitting = false;
					submitType = null;
					if (result.type === 'failure') {
						const errorMsg = (result.data as any)?.error || 'Authentication failed. Check your credentials.';
						toast.error(errorMsg, { id: 'auth-toast' });
					} else if (result.type === 'redirect') {
						toast.success('Signed in successfully! Redirecting...', { id: 'auth-toast' });
					} else {
						toast.dismiss('auth-toast');
					}
					await update();
				};
			}}
		>
			<div class="grid gap-2">
				<Label for="login-email">Email</Label>
				<div class="relative">
					<IconId class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
					<Input
						id="login-email"
						name="email"
						type="email"
						autocomplete="email"
						class="h-10 pl-10"
						required
						placeholder="you@institution.edu"
						bind:value={emailInput}
						disabled={isSubmitting}
					/>
				</div>
			</div>

			<div class="grid gap-2">
				<div class="flex items-center justify-between">
					<Label for="login-password">Password</Label>
					<span class="text-xs text-muted-foreground">Case sensitive</span>
				</div>
				<div class="relative">
					<IconLock class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
					<Input
						id="login-password"
						name="password"
						type="password"
						autocomplete="current-password"
						class="h-10 pl-10"
						required
						placeholder="Enter your password"
						bind:value={passwordInput}
						disabled={isSubmitting}
					/>
				</div>
			</div>

			<div class="min-h-6" aria-live="polite">
				{#if initialMode === 'sign-in' && form?.error}
					<p class="rounded-lg border border-destructive/25 bg-destructive/10 px-3 py-2.5 text-sm text-destructive" role="alert">
						{form.error}
					</p>
				{/if}
			</div>

			<Button type="submit" size="lg" class="h-11 w-full gap-2 cursor-pointer" disabled={isSubmitting}>
				{#if isSubmitting && submitType === 'login'}
					<IconLoader2 class="size-4.5 animate-spin" />
					<span>Authenticating...</span>
				{:else}
					<span>Sign in</span>
					<IconArrowRight class="size-4" />
				{/if}
			</Button>
		</form>

		{#if dev}
			<div class="mt-6 rounded-xl border border-primary/30 bg-primary/5 p-4 text-xs">
				<div class="flex items-center justify-between">
					<div class="flex items-center gap-1.5 font-semibold text-primary">
						<IconTestPipe class="size-4" /> Development Test Mode
					</div>
					<span class="rounded bg-primary/10 px-1.5 py-0.5 text-[10px] font-bold text-primary">DEV ONLY</span>
				</div>
				<p class="mt-1 text-muted-foreground">
					Pre-configured test account equipped with Admin, Teacher, and Student permissions to test the entire Web UI.
				</p>
				<div class="mt-2 space-y-0.5 text-[11px] font-mono text-muted-foreground">
					<p><span class="font-bold text-foreground">Email:</span> dev@campus.internal</p>
					<p><span class="font-bold text-foreground">Password:</span> DevPassword123!</p>
					<p><span class="font-bold text-foreground">MFA TOTP Key:</span> HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ</p>
					<p><span class="font-bold text-foreground">Backup Code:</span> DEV-RECV-0001</p>
				</div>
				<div class="mt-3 flex gap-2">
					<form
						method="POST"
						action="/public/sign-in"
						class="flex-1"
						use:enhance={() => {
							isSubmitting = true;
							submitType = 'dev-login';
							toast.loading('Logging into dev environment...', { id: 'auth-toast' });
							return async ({ result, update }) => {
								isSubmitting = false;
								submitType = null;
								if (result.type === 'failure') {
									const errorMsg = (result.data as any)?.error || 'Dev login failed.';
									toast.error(errorMsg, { id: 'auth-toast' });
								} else if (result.type === 'redirect') {
									toast.success('Signed in as dev user! Redirecting...', { id: 'auth-toast' });
								} else {
									toast.dismiss('auth-toast');
								}
								await update();
							};
						}}
					>
						<input type="hidden" name="devLogin" value="true" />
						<Button type="submit" size="sm" class="h-8 w-full gap-1 cursor-pointer" disabled={isSubmitting}>
							{#if isSubmitting && submitType === 'dev-login'}
								<IconLoader2 class="size-3.5 animate-spin" />
								<span>Logging in...</span>
							{:else}
								<IconSparkles class="size-3.5" />
								<span>1-Click Dev Sign In</span>
							{/if}
						</Button>
					</form>
					<Button type="button" variant="outline" size="sm" class="h-8 cursor-pointer" onclick={fillDevCredentials} disabled={isSubmitting}>
						Fill Form
					</Button>
				</div>
			</div>
		{/if}
	</Tabs.Content>

	<Tabs.Content value="sign-up" class="mt-6">
		<div class="border-b pb-4">
			<h2 class="font-heading text-xl font-semibold">Create your account</h2>
			<p class="mt-1 text-sm text-muted-foreground">Use details that match college records.</p>
		</div>

		<form
			method="POST"
			action="/public/sign-up"
			class="mt-6 grid gap-4"
			use:enhance={() => {
				isSubmitting = true;
				submitType = 'signup';
				toast.loading('Creating institutional account...', { id: 'auth-toast' });
				return async ({ result, update }) => {
					isSubmitting = false;
					submitType = null;
					if (result.type === 'failure') {
						const errorMsg = (result.data as any)?.error || 'Account creation failed.';
						toast.error(errorMsg, { id: 'auth-toast' });
					} else if (result.type === 'redirect') {
						toast.success('Account created! Redirecting...', { id: 'auth-toast' });
					} else {
						toast.dismiss('auth-toast');
					}
					await update();
				};
			}}
		>
			<div class="grid gap-2">
				<Label for="displayName">Full name</Label>
				<Input id="displayName" name="displayName" autocomplete="name" class="h-9" required placeholder="As shown in college records" disabled={isSubmitting} />
			</div>
			<div class="grid gap-2">
				<Label for="institutionalId">Institutional ID</Label>
				<Input id="institutionalId" name="institutionalId" autocomplete="username" class="h-9" required placeholder="Enrollment or employee ID" disabled={isSubmitting} />
			</div>
			<div class="grid gap-2">
				<Label for="signup-email">Email</Label>
				<div class="relative">
					<IconMail class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
					<Input id="signup-email" name="email" type="email" autocomplete="email" class="h-9 pl-10" required placeholder="you@institution.edu" disabled={isSubmitting} />
				</div>
			</div>
			<div class="grid gap-2">
				<Label for="signup-password">Create password</Label>
				<Input id="signup-password" name="password" type="password" autocomplete="new-password" minlength={12} class="h-9" required placeholder="At least 12 characters" disabled={isSubmitting} />
			</div>
			<div class="grid gap-2">
				<Label for="role">Requested role</Label>
				<input type="hidden" name="role" value={role} />
				<Select.Root type="single" bind:value={role} disabled={isSubmitting}>
					<Select.Trigger id="role" class="h-9 w-full">
						{role === 'GROUND_STAFF' ? 'Ground staff' : role[0] + role.slice(1).toLowerCase().replace('_', ' ')}
					</Select.Trigger>
					<Select.Content>
						<Select.Item value="STUDENT">Student</Select.Item>
						<Select.Item value="TEACHER">Teacher</Select.Item>
						<Select.Item value="SECURITY">Security</Select.Item>
						<Select.Item value="HELP">Help</Select.Item>
						<Select.Item value="GROUND_STAFF">Ground staff</Select.Item>
						<Select.Item value="MODERATOR">Moderator</Select.Item>
						<Select.Item value="ADMIN">Admin</Select.Item>
					</Select.Content>
				</Select.Root>
			</div>
			<div class="min-h-6" aria-live="polite">
				{#if initialMode === 'sign-up' && form?.error}
					<p class="rounded-lg border border-destructive/25 bg-destructive/10 px-3 py-2.5 text-sm text-destructive" role="alert">
						{form.error}
					</p>
				{/if}
			</div>
			<Button type="submit" size="lg" class="h-11 w-full gap-2 cursor-pointer" disabled={isSubmitting}>
				{#if isSubmitting && submitType === 'signup'}
					<IconLoader2 class="size-4.5 animate-spin" />
					<span>Creating account...</span>
				{:else}
					<span>Create account</span>
					<IconArrowRight class="size-4" />
				{/if}
			</Button>
		</form>
	</Tabs.Content>
</Tabs.Root>
