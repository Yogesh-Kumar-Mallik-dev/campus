<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { downloadPdfBlob } from '$lib/pdf/pdf-service';
	import { toast } from '$lib/components/ui/sonner';
	import {
		IconDownload,
		IconPrinter,
		IconZoomIn,
		IconZoomOut,
		IconRotateClockwise,
		IconMaximize,
		IconFileText,
		IconExternalLink,
		IconEye,
		IconSparkles
	} from '@tabler/icons-svelte';

	let {
		pdfBytes = null,
		pdfUrl = null,
		title = 'Document.pdf',
		canDownload = true,
		canPrint = true
	}: {
		pdfBytes?: Uint8Array | null;
		pdfUrl?: string | null;
		title?: string;
		canDownload?: boolean;
		canPrint?: boolean;
	} = $props();

	let zoomLevel = $state(100);
	let rotation = $state(0);
	let objectUrl = $state<string | null>(null);
	let isFullscreen = $state(false);
	let containerRef = $state<HTMLDivElement | null>(null);

	$effect(() => {
		if (pdfBytes && pdfBytes.length > 0) {
			const blob = new Blob([pdfBytes as any], { type: 'application/pdf' });
			const url = URL.createObjectURL(blob);
			objectUrl = url;
			return () => {
				URL.revokeObjectURL(url);
			};
		} else if (pdfUrl) {
			objectUrl = pdfUrl;
		} else {
			objectUrl = null;
		}
	});

	function handleZoomIn() {
		if (zoomLevel < 200) zoomLevel = Math.min(200, zoomLevel + 15);
	}

	function handleZoomOut() {
		if (zoomLevel > 50) zoomLevel = Math.max(50, zoomLevel - 15);
	}

	function handleRotate() {
		rotation = (rotation + 90) % 360;
	}

	function handleDownload() {
		if (pdfBytes) {
			downloadPdfBlob(pdfBytes, title);
			toast.success('Downloaded PDF', { description: `${title} saved to downloads.` });
		} else if (objectUrl) {
			const a = document.createElement('a');
			a.href = objectUrl;
			a.download = title.endsWith('.pdf') ? title : `${title}.pdf`;
			document.body.appendChild(a);
			a.click();
			document.body.removeChild(a);
			toast.success('Downloaded PDF', { description: `${title} saved to downloads.` });
		}
	}

	function handlePrint() {
		if (objectUrl) {
			const printWindow = window.open(objectUrl, '_blank');
			if (printWindow) {
				printWindow.focus();
				toast.info('Opening Print View', { description: 'Use browser print dialog (Ctrl+P / Cmd+P).' });
			}
		}
	}

	function toggleFullscreen() {
		if (!containerRef) return;
		if (!document.fullscreenElement) {
			containerRef.requestFullscreen().catch(() => {});
			isFullscreen = true;
		} else {
			document.exitFullscreen().catch(() => {});
			isFullscreen = false;
		}
	}
</script>

<div
	bind:this={containerRef}
	class="flex flex-col h-full min-h-[540px] rounded-xl border border-border/60 bg-card overflow-hidden shadow-xs"
>
	<!-- Top Controls Toolbar -->
	<div class="flex flex-wrap items-center justify-between gap-3 px-4 py-2.5 bg-muted/40 border-b border-border/40 text-xs">
		<div class="flex items-center gap-2">
			<IconFileText class="size-4 text-primary shrink-0" />
			<span class="font-medium text-foreground truncate max-w-[200px] sm:max-w-xs">{title}</span>
			<Badge variant="outline" class="text-[10px] uppercase font-mono px-1.5 py-0">PDF</Badge>
		</div>

		<div class="flex flex-wrap items-center gap-1.5">
			<!-- Zoom controls -->
			<div class="flex items-center rounded-lg border border-border/50 bg-background/80 p-0.5 shadow-2xs">
				<Button
					size="icon-xs"
					variant="ghost"
					onclick={handleZoomOut}
					disabled={zoomLevel <= 50}
					title="Zoom Out"
					class="h-7 w-7 cursor-pointer"
				>
					<IconZoomOut class="size-3.5" />
				</Button>
				<span class="px-2 font-mono text-[11px] min-w-[3rem] text-center select-none">{zoomLevel}%</span>
				<Button
					size="icon-xs"
					variant="ghost"
					onclick={handleZoomIn}
					disabled={zoomLevel >= 200}
					title="Zoom In"
					class="h-7 w-7 cursor-pointer"
				>
					<IconZoomIn class="size-3.5" />
				</Button>
			</div>

			<!-- Rotation control -->
			<Button
				size="sm"
				variant="outline"
				onclick={handleRotate}
				title="Rotate 90 degrees"
				class="h-8 gap-1 text-xs bg-background/80 cursor-pointer"
			>
				<IconRotateClockwise class="size-3.5" />
				<span class="hidden sm:inline">{rotation}°</span>
			</Button>

			<!-- Fullscreen toggle -->
			<Button
				size="sm"
				variant="outline"
				onclick={toggleFullscreen}
				title="Fullscreen"
				class="h-8 text-xs bg-background/80 cursor-pointer"
			>
				<IconMaximize class="size-3.5" />
			</Button>

			{#if canPrint && objectUrl}
				<Button
					size="sm"
					variant="outline"
					onclick={handlePrint}
					class="h-8 gap-1 text-xs bg-background/80 cursor-pointer"
				>
					<IconPrinter class="size-3.5" />
					<span class="hidden sm:inline">Print</span>
				</Button>
			{/if}

			{#if canDownload && objectUrl}
				<Button
					size="sm"
					variant="default"
					onclick={handleDownload}
					class="h-8 gap-1 text-xs shadow-xs cursor-pointer"
				>
					<IconDownload class="size-3.5" />
					<span>Download PDF</span>
				</Button>
			{/if}
		</div>
	</div>

	<!-- Document Canvas / Frame Area -->
	<div class="relative flex-1 bg-muted/20 dark:bg-muted/10 overflow-auto flex items-center justify-center p-2 sm:p-4">
		{#if objectUrl}
			<div
				class="w-full h-full min-h-[460px] flex items-center justify-center transition-transform duration-200"
				style="transform: rotate({rotation}deg) scale({zoomLevel / 100}); transform-origin: center center;"
			>
				<object
					data={objectUrl}
					type="application/pdf"
					class="w-full h-full min-h-[460px] rounded-lg shadow-sm border border-border/30 bg-background"
					title={title}
				>
					<div class="flex flex-col items-center justify-center h-full p-8 text-center space-y-4">
						<IconFileText class="size-12 text-muted-foreground animate-pulse" />
						<p class="text-sm font-medium">PDF Preview</p>
						<p class="text-xs text-muted-foreground max-w-sm">
							Your browser is displaying the document. You can also download it directly or open in a new tab.
						</p>
						<div class="flex gap-2">
							<Button size="sm" onclick={handleDownload} class="gap-1.5 cursor-pointer">
								<IconDownload class="size-4" /> Download Document
							</Button>
							<a href={objectUrl} target="_blank" rel="noreferrer" class="inline-flex items-center gap-1 text-xs text-primary hover:underline pt-2">
								<IconExternalLink class="size-3.5" /> Open in New Tab
							</a>
						</div>
					</div>
				</object>
			</div>
		{:else}
			<div class="flex flex-col items-center justify-center h-full p-12 text-center space-y-3">
				<div class="rounded-full bg-primary/10 p-4">
					<IconFileText class="size-8 text-primary" />
				</div>
				<p class="text-sm font-semibold">No Document Loaded</p>
				<p class="text-xs text-muted-foreground max-w-sm">
					Select a template from the creator tab, fill an application form, or upload an external PDF file to view and edit.
				</p>
			</div>
		{/if}
	</div>
</div>
