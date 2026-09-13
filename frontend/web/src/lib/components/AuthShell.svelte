<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Switch } from '$lib/components/ui/switch';
	import logo from '$lib/../assets/bbdit-logo-transparent.png';
	import campus from '$lib/../assets/bbdit-campus-hero.webp';
	import { theme } from '$lib/theme.svelte';
	import { IconMoon, IconShieldCheck, IconSun } from '@tabler/icons-svelte';

	let { children, title, description }: { children: Snippet; title: string; description: string } = $props();
</script>

<div class="min-h-screen bg-muted/35 text-foreground lg:grid lg:grid-cols-[minmax(28rem,1fr)_minmax(32rem,1fr)]">
	<section class="relative hidden min-h-screen overflow-hidden bg-[#071a32] text-white lg:flex lg:min-h-[46rem] lg:flex-col lg:justify-between">
		<img src={campus} alt="Campus" class="absolute inset-0 size-full object-cover opacity-45" />
		<div class="absolute inset-0 bg-linear-to-b from-[#071a32]/45 via-[#071a32]/65 to-[#071a32]"></div>
		<div class="relative p-10 xl:p-14"><a href="/public" class="inline-flex rounded-lg bg-white/95 px-3 py-2 shadow-lg"><img src={logo} alt="BBDIT" class="h-11 w-auto" /></a></div>
		<div class="relative max-w-2xl p-10 xl:p-14"><div class="mb-5 inline-flex items-center gap-2 rounded-full border border-white/20 bg-white/10 px-3 py-1.5 text-sm backdrop-blur"><IconShieldCheck class="size-4" /> Secure campus gateway</div><h2 class="font-heading text-4xl font-semibold leading-tight xl:text-5xl">One campus.<br />One connected community.</h2><p class="mt-5 max-w-xl text-base leading-7 text-white/75">Access academic tools, campus services and communication for students, teachers and administrators.</p></div>
	</section>
	<main class="relative flex min-h-screen justify-center px-[clamp(1rem,5vw,3rem)] pb-12 pt-24 sm:pt-20 lg:min-h-[46rem] lg:items-start lg:pt-[clamp(6rem,14vh,9rem)]">
		<div class="absolute right-4 top-4 rounded-full border bg-background/90 p-1.5 shadow-sm sm:right-8 sm:top-7"><Switch size="icon" checked={theme.isDark} onCheckedChange={(val) => theme.set(val)} aria-label="Use dark mode">{#snippet thumbContent()}{#if theme.isDark}<IconMoon class="size-3.5 text-primary" />{:else}<IconSun class="size-3.5" />{/if}{/snippet}</Switch></div>
		<div class="w-full max-w-md"><a href="/public" class="mb-9 inline-flex max-w-full rounded-lg border bg-white/95 px-2 py-1.5 shadow-xs lg:hidden"><img src={logo} alt="BBDIT" class="h-auto max-h-10 w-auto max-w-full" /></a><p class="text-xs font-bold uppercase tracking-[.16em] text-primary">Campus portal</p><h1 class="mt-2 font-heading text-3xl font-semibold tracking-tight sm:text-4xl">{title}</h1><p class="mt-3 text-sm leading-6 text-muted-foreground">{description}</p><div class="mt-8 rounded-2xl border bg-card p-[clamp(1rem,5vw,2rem)] shadow-sm">{@render children()}</div><p class="mt-6 text-center text-xs leading-5 text-muted-foreground">By continuing, you agree to use this portal in accordance with BBDIT institutional policies.</p></div>
	</main>
</div>
