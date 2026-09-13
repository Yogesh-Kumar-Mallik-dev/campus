<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent } from '$lib/components/ui/card';
	import { IconArrowLeft, IconBook2, IconMapPin, IconSchool } from '@tabler/icons-svelte';
</script>

<svelte:head><title>Page not found · Campus Management</title></svelte:head>

<section class="campus-scene relative isolate min-h-screen overflow-hidden bg-background px-5 py-10 sm:px-10 sm:py-14">
	<div class="pointer-events-none absolute inset-0 overflow-hidden" aria-hidden="true"><div class="glow glow-one"></div><div class="glow glow-two"></div><div class="cloud cloud-one"></div><div class="cloud cloud-two"></div><div class="absolute inset-x-0 bottom-0 h-1/3 bg-linear-to-t from-primary/10 to-transparent"></div></div>
	<div class="relative z-10 mx-auto flex max-w-3xl flex-col items-center text-center">
		<div class="mb-6 inline-flex items-center gap-2 rounded-full border bg-card/80 px-3 py-1.5 text-xs font-medium text-muted-foreground shadow-sm backdrop-blur-sm"><IconMapPin class="size-3.5 text-primary" /> Somewhere off the campus map</div>
		<div class="illustration mb-8" aria-hidden="true"><span class="digit digit-left">4</span><div class="campus-icon"><div class="icon-halo"></div><div class="icon-circle"><IconSchool class="size-16 stroke-[1.4] sm:size-20" /></div><IconBook2 class="floating-book size-7" /></div><span class="digit digit-right">4</span></div>
		<h1 class="font-heading text-4xl font-semibold tracking-tight text-balance sm:text-6xl">This hall isn’t on the map</h1>
		<p class="mt-4 max-w-xl text-base leading-7 text-muted-foreground text-balance sm:text-lg">The page may have graduated, moved to another building, or never enrolled. Let’s get you back to familiar grounds.</p>
		<div class="mt-8 flex flex-col gap-3 sm:flex-row"><Button href="/public" size="lg" class="h-11 rounded-xl px-5 shadow-lg shadow-primary/20"><IconArrowLeft data-icon="inline-start" /> Return to campus</Button><Button href="/" variant="outline" size="lg" class="h-11 rounded-xl bg-background/60 px-5 backdrop-blur-sm">Open portal</Button></div>
		<Card class="mt-10 w-full max-w-lg border-dashed bg-card/60 shadow-none backdrop-blur-sm"><CardContent class="flex items-start gap-3 p-4 text-left"><div class="mt-0.5 rounded-lg bg-primary/10 p-2 text-primary"><IconBook2 class="size-4" /></div><div><p class="text-sm font-medium">A quick study tip</p><p class="mt-0.5 text-sm leading-5 text-muted-foreground">Check the web address for a typo, or head home and start from the campus directory.</p></div></CardContent></Card>
	</div>
</section>

<style>
	.campus-scene { background-image: linear-gradient(to right, color-mix(in oklch, var(--border) 45%, transparent) 1px, transparent 1px), linear-gradient(to bottom, color-mix(in oklch, var(--border) 45%, transparent) 1px, transparent 1px); background-size: 32px 32px; }
	.glow { position: absolute; width: 24rem; height: 24rem; border-radius: 9999px; filter: blur(80px); opacity: .17; background: var(--primary); } .glow-one { top: -10rem; left: -8rem; } .glow-two { right: -10rem; bottom: -12rem; opacity: .1; }
	.cloud { position: absolute; width: 8rem; height: 1.5rem; border-radius: 9999px; background: color-mix(in oklch, var(--muted-foreground) 12%, transparent); filter: blur(1px); animation: drift 9s ease-in-out infinite alternate; } .cloud-one { top: 18%; left: 8%; } .cloud-two { top: 28%; right: 7%; scale: .65; animation-delay: -4s; }
	.illustration { display: flex; align-items: center; justify-content: center; gap: clamp(.5rem, 3vw, 1.5rem); } .digit { font-family: 'Outfit Variable', sans-serif; font-size: clamp(5rem, 16vw, 9rem); font-weight: 700; line-height: 1; color: color-mix(in oklch, var(--foreground) 88%, var(--primary)); text-shadow: 0 12px 35px color-mix(in oklch, var(--primary) 15%, transparent); animation: settle 700ms cubic-bezier(.22, 1, .36, 1) both; } .digit-right { animation-delay: 100ms; } .campus-icon { position: relative; animation: hover 4s ease-in-out infinite; }
	.icon-circle { position: relative; display: grid; width: clamp(6.5rem, 17vw, 9rem); height: clamp(6.5rem, 17vw, 9rem); place-items: center; border: 1px solid color-mix(in oklch, var(--primary) 28%, var(--border)); border-radius: 9999px; color: var(--primary); background: color-mix(in oklch, var(--card) 88%, transparent); box-shadow: 0 18px 45px color-mix(in oklch, var(--primary) 16%, transparent); backdrop-filter: blur(12px); } .icon-halo { position: absolute; inset: -.7rem; border: 1px dashed color-mix(in oklch, var(--primary) 35%, transparent); border-radius: 9999px; animation: spin 18s linear infinite; } .campus-icon :global(.floating-book) { position: absolute; right: -.7rem; bottom: .4rem; padding: .3rem; border-radius: .5rem; color: var(--primary-foreground); background: var(--primary); box-shadow: 0 8px 20px color-mix(in oklch, var(--primary) 28%, transparent); animation: bob 2.5s ease-in-out infinite; }
	@keyframes drift { to { transform: translateX(2rem); } } @keyframes hover { 50% { transform: translateY(-.5rem); } } @keyframes bob { 50% { transform: translateY(-.35rem) rotate(5deg); } } @keyframes spin { to { transform: rotate(360deg); } } @keyframes settle { from { opacity: 0; transform: translateY(1rem) scale(.92); } to { opacity: 1; transform: translateY(0) scale(1); } }
	@media (prefers-reduced-motion: reduce) { .cloud, .digit, .campus-icon, .icon-halo, .campus-icon :global(.floating-book) { animation: none; } }
</style>
