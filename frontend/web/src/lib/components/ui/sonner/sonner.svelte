<script lang="ts">
	import { Toaster as Sonner, type ToasterProps as SonnerProps } from 'svelte-sonner';
	import { theme } from '$lib/theme.svelte';

	let { ...restProps }: SonnerProps = $props();
</script>

<Sonner
	theme={theme.isDark ? 'dark' : 'light'}
	position="bottom-left"
	class="toaster group"
	closeButton={false}
	offset="20px"
	gap={12}
	duration={4200}
	toastOptions={{
		classes: {
			toast:
				'group toast font-sans antialiased text-xs font-medium rounded-xl border border-border/80 bg-background/95 dark:bg-card/95 text-foreground shadow-[0_12px_36px_-6px_rgba(0,0,0,0.18),0_4px_12px_-2px_rgba(0,0,0,0.08)] dark:shadow-[0_16px_40px_-6px_rgba(0,0,0,0.55),0_4px_16px_-2px_rgba(0,0,0,0.35)] backdrop-blur-xl transition-all duration-200 py-3.5 px-4 flex items-start gap-3 select-none w-[calc(100vw-2rem)] sm:w-auto min-w-[min(280px,calc(100vw-2rem))] sm:min-w-[340px] max-w-[calc(100vw-2rem)] sm:max-w-[420px]',
			title: 'font-heading font-semibold text-xs text-foreground tracking-tight',
			description: 'text-[11px] text-muted-foreground leading-relaxed mt-0.5',
			actionButton:
				'bg-primary text-primary-foreground font-semibold text-[11px] h-7 px-3 rounded-lg shadow-xs hover:bg-primary/90 transition-all shrink-0 ml-auto',
			cancelButton:
				'bg-secondary text-secondary-foreground font-medium text-[11px] h-7 px-3 rounded-lg hover:bg-secondary/80 transition-all shrink-0',
			success:
				'!border-l-[3.5px] !border-l-emerald-500 !border-t-border/70 !border-r-border/70 !border-b-border/70 dark:!border-l-emerald-400 ring-1 ring-emerald-500/10 dark:ring-emerald-400/20',
			error:
				'!border-l-[3.5px] !border-l-rose-500 !border-t-border/70 !border-r-border/70 !border-b-border/70 dark:!border-l-rose-400 ring-1 ring-rose-500/10 dark:ring-rose-400/20',
			info:
				'!border-l-[3.5px] !border-l-sky-500 !border-t-border/70 !border-r-border/70 !border-b-border/70 dark:!border-l-sky-400 ring-1 ring-sky-500/10 dark:ring-sky-400/20',
			warning:
				'!border-l-[3.5px] !border-l-amber-500 !border-t-border/70 !border-r-border/70 !border-b-border/70 dark:!border-l-amber-400 ring-1 ring-amber-500/10 dark:ring-amber-400/20',
			loading:
				'!border-l-[3.5px] !border-l-primary !border-t-border/70 !border-r-border/70 !border-b-border/70 ring-1 ring-primary/10'
		}
	}}
	{...restProps}
/>
