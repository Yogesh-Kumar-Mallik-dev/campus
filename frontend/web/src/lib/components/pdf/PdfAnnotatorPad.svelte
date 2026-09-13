<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import { toast } from '$lib/components/ui/sonner';
	import {
		IconPencil,
		IconTrash,
		IconCertificate,
		IconCheck,
		IconShieldCheck,
		IconSparkles
	} from '@tabler/icons-svelte';

	let {
		onApplyAnnotation
	}: {
		onApplyAnnotation: (options: {
			stampText?: string;
			watermarkText?: string;
			signatureImageDataUrl?: string;
			dateStamp?: boolean;
		}) => void;
	} = $props();

	let canvasRef = $state<HTMLCanvasElement | null>(null);
	let isDrawing = $state(false);
	let hasSignature = $state(false);
	let strokeColor = $state('#1e3a8a'); // Navy blue
	let selectedStamp = $state<'VERIFIED & APPROVED' | 'OFFICIAL INSTITUTIONAL SEAL' | 'CONFIDENTIAL' | 'PROVISIONAL DRAFT' | 'NONE'>('VERIFIED & APPROVED');
	let watermark = $state('');
	let includeDateStamp = $state(true);

	$effect(() => {
		if (canvasRef) {
			const ctx = canvasRef.getContext('2d');
			if (ctx) {
				ctx.lineWidth = 2.5;
				ctx.lineCap = 'round';
				ctx.lineJoin = 'round';
				ctx.strokeStyle = strokeColor;
			}
		}
	});

	function startDrawing(e: MouseEvent | TouchEvent) {
		isDrawing = true;
		const canvas = canvasRef;
		if (!canvas) return;
		const ctx = canvas.getContext('2d');
		if (!ctx) return;

		const rect = canvas.getBoundingClientRect();
		const clientX = 'touches' in e ? e.touches[0].clientX : e.clientX;
		const clientY = 'touches' in e ? e.touches[0].clientY : e.clientY;

		ctx.beginPath();
		ctx.moveTo(clientX - rect.left, clientY - rect.top);
	}

	function draw(e: MouseEvent | TouchEvent) {
		if (!isDrawing || !canvasRef) return;
		const ctx = canvasRef.getContext('2d');
		if (!ctx) return;

		const rect = canvasRef.getBoundingClientRect();
		const clientX = 'touches' in e ? e.touches[0].clientX : e.clientX;
		const clientY = 'touches' in e ? e.touches[0].clientY : e.clientY;

		ctx.strokeStyle = strokeColor;
		ctx.lineTo(clientX - rect.left, clientY - rect.top);
		ctx.stroke();
		hasSignature = true;
	}

	function stopDrawing() {
		isDrawing = false;
	}

	function clearCanvas() {
		if (!canvasRef) return;
		const ctx = canvasRef.getContext('2d');
		if (!ctx) return;
		ctx.clearRect(0, 0, canvasRef.width, canvasRef.height);
		hasSignature = false;
		toast.info('Signature canvas cleared');
	}

	function handleApply() {
		let sigUrl: string | undefined;
		if (hasSignature && canvasRef) {
			sigUrl = canvasRef.toDataURL('image/png');
		}

		onApplyAnnotation({
			stampText: selectedStamp !== 'NONE' ? selectedStamp : undefined,
			watermarkText: watermark.trim() ? watermark.trim() : undefined,
			signatureImageDataUrl: sigUrl,
			dateStamp: includeDateStamp
		});

		toast.success('Annotations & Seal Applied', {
			description: 'The document has been digitally stamped and updated.'
		});
	}
</script>

<div class="space-y-6">
	<!-- Digital Signature Pad -->
	<div class="space-y-3">
		<div class="flex items-center justify-between">
			<Label class="text-xs font-semibold flex items-center gap-1.5">
				<IconPencil class="size-4 text-primary" /> Draw Digital Signature
			</Label>
			<div class="flex items-center gap-2">
				<button
					type="button"
					onclick={() => (strokeColor = '#1e3a8a')}
					class="size-4 rounded-full bg-blue-900 border {strokeColor === '#1e3a8a' ? 'ring-2 ring-primary ring-offset-1' : ''}"
					title="Navy Blue"
				></button>
				<button
					type="button"
					onclick={() => (strokeColor = '#111827')}
					class="size-4 rounded-full bg-gray-900 border {strokeColor === '#111827' ? 'ring-2 ring-primary ring-offset-1' : ''}"
					title="Black"
				></button>
				<Button size="icon-xs" variant="ghost" onclick={clearCanvas} class="h-6 w-6 text-muted-foreground" title="Clear Pad">
					<IconTrash class="size-3.5" />
				</Button>
			</div>
		</div>

		<div class="rounded-lg border-2 border-dashed border-border/80 bg-background overflow-hidden relative">
			<canvas
				bind:this={canvasRef}
				width={380}
				height={120}
				onmousedown={startDrawing}
				onmousemove={draw}
				onmouseup={stopDrawing}
				onmouseleave={stopDrawing}
				ontouchstart={startDrawing}
				ontouchmove={draw}
				ontouchend={stopDrawing}
				class="w-full h-[120px] cursor-crosshair touch-none"
			></canvas>
			{#if !hasSignature}
				<div class="pointer-events-none absolute inset-0 flex items-center justify-center text-xs text-muted-foreground/60">
					Sign here with mouse or touch
				</div>
			{/if}
		</div>
	</div>

	<!-- Official Verification Seal -->
	<div class="space-y-2">
		<Label class="text-xs font-semibold flex items-center gap-1.5">
			<IconCertificate class="size-4 text-primary" /> Institutional Verification Seal
		</Label>
		<Select.Root type="single" bind:value={selectedStamp}>
			<Select.Trigger class="h-9 w-full text-xs">
				{#if selectedStamp === 'VERIFIED & APPROVED'}VERIFIED & APPROVED (Green Badge)
				{:else if selectedStamp === 'OFFICIAL INSTITUTIONAL SEAL'}OFFICIAL INSTITUTIONAL SEAL
				{:else if selectedStamp === 'CONFIDENTIAL'}CONFIDENTIAL DOCUMENT
				{:else if selectedStamp === 'PROVISIONAL DRAFT'}PROVISIONAL DRAFT COPY
				{:else if selectedStamp === 'NONE'}No Stamp
				{:else}{selectedStamp}{/if}
			</Select.Trigger>
			<Select.Content>
				<Select.Item value="VERIFIED & APPROVED">VERIFIED & APPROVED (Green Badge)</Select.Item>
				<Select.Item value="OFFICIAL INSTITUTIONAL SEAL">OFFICIAL INSTITUTIONAL SEAL</Select.Item>
				<Select.Item value="CONFIDENTIAL">CONFIDENTIAL DOCUMENT</Select.Item>
				<Select.Item value="PROVISIONAL DRAFT">PROVISIONAL DRAFT COPY</Select.Item>
				<Select.Item value="NONE">No Stamp</Select.Item>
			</Select.Content>
		</Select.Root>
	</div>

	<!-- Watermark Text -->
	<div class="space-y-2">
		<Label for="watermark-input" class="text-xs font-semibold">Background Watermark (Optional)</Label>
		<Input
			id="watermark-input"
			bind:value={watermark}
			placeholder="e.g. CAMPUS CONFIDENTIAL or COPY ONLY"
			class="h-9 text-xs"
		/>
	</div>

	<!-- Date Stamp Checkbox -->
	<div class="flex items-center gap-2 pt-1">
		<input
			type="checkbox"
			id="date-stamp-check"
			bind:checked={includeDateStamp}
			class="rounded border-input text-primary"
		/>
		<Label for="date-stamp-check" class="text-xs cursor-pointer select-none text-muted-foreground">
			Include live timestamp & date in stamp header
		</Label>
	</div>

	<Button
		type="button"
		onclick={handleApply}
		class="w-full h-9 gap-1.5 text-xs shadow-xs cursor-pointer"
	>
		<IconShieldCheck class="size-4" /> Apply Annotations & Stamp
	</Button>
</div>
