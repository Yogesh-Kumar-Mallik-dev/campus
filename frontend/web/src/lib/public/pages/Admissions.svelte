<script lang="ts">
	import PublicShell from "$public/components/PublicShell.svelte";
	import PageHeader from "$public/components/PageHeader.svelte";
	import { Button } from "$lib/components/ui/button";
	import { Card, CardContent, CardHeader, CardTitle } from "$lib/components/ui/card";
	import { Input } from "$lib/components/ui/input";
	import { Label } from "$lib/components/ui/label";
	import { Textarea } from "$lib/components/ui/textarea";
	import * as Select from "$lib/components/ui/select";
	import * as Accordion from "$lib/components/ui/accordion";
	import { toast } from "$lib/components/ui/sonner";
	import { IconCheck, IconLoader2, IconSend } from "@tabler/icons-svelte";

	let submitted = $state(false);
	let isSubmitting = $state(false);
	let programme = $state("");
	let applicantName = $state("");
	let applicantPhone = $state("");
	let applicantEmail = $state("");
	let applicantQuestion = $state("");

	const steps = [
		["1", "Choose a programme", "Review course duration, eligibility and the relevant university requirements."],
		["2", "Send an enquiry", "Share your contact details so the admissions team can guide you."],
		["3", "Complete verification", "Submit required academic records and complete the institute process."]
	];

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		isSubmitting = true;
		const toastId = toast.loading("Submitting admission enquiry...");
		setTimeout(() => {
			isSubmitting = false;
			submitted = true;
			toast.success("Enquiry Submitted", {
				id: toastId,
				description: "Our admissions counselor will contact you within 24 business hours."
			});
		}, 600);
	}
</script>

<svelte:head><title>Admissions & Enquiry · BBDIT</title></svelte:head>

<PublicShell>
	<main>
		<PageHeader
			eyebrow="Admissions"
			title="Begin your application with clear guidance"
			description="Explore the process and send an enquiry for engineering, pharmacy, management, computer applications, education or diploma programmes."
		/>

		<section class="mx-auto grid max-w-7xl gap-12 px-4 py-16 sm:px-6 lg:grid-cols-[.8fr_1.2fr] lg:px-8">
			<div>
				<h2 class="font-heading text-3xl font-semibold">How admission works</h2>
				<div class="mt-7 space-y-7">
					{#each steps as step}
						<div class="grid grid-cols-[2.5rem_1fr] gap-4">
							<span class="grid size-10 place-items-center rounded-full bg-primary text-sm font-bold text-primary-foreground">
								{step[0]}
							</span>
							<div>
								<h3 class="font-semibold">{step[1]}</h3>
								<p class="mt-1 text-sm leading-6 text-muted-foreground">{step[2]}</p>
							</div>
						</div>
					{/each}
				</div>

				<Card class="mt-9 bg-muted/40 shadow-none">
					<CardContent class="p-5 text-sm leading-6">
						<strong>College codes</strong>
						<p class="mt-1 text-muted-foreground">B.Tech: 035 · B.Pharm: 877</p>
					</CardContent>
				</Card>
			</div>

			<Card class="border-border/80 shadow-xs">
				<CardHeader>
					<CardTitle class="font-heading text-2xl">Admission enquiry</CardTitle>
				</CardHeader>
				<CardContent>
					{#if submitted}
						<div class="rounded-xl border border-emerald-500/30 bg-emerald-500/5 p-8 text-center space-y-3">
							<div class="inline-flex size-12 items-center justify-center rounded-full bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
								<IconCheck class="size-6" />
							</div>
							<h3 class="font-heading text-2xl font-semibold text-foreground">Enquiry Received</h3>
							<p class="text-sm text-muted-foreground max-w-md mx-auto">
								Thank you, <strong>{applicantName}</strong>. Your enquiry for {programme.toUpperCase() || 'the selected program'} has been recorded. Our admissions desk will reach out shortly.
							</p>
							<Button variant="outline" class="mt-4 cursor-pointer" onclick={() => {
								submitted = false;
								applicantName = "";
								applicantPhone = "";
								applicantEmail = "";
								applicantQuestion = "";
								programme = "";
							}}>
								Send Another Enquiry
							</Button>
						</div>
					{:else}
						<form class="grid gap-5" onsubmit={handleSubmit}>
							<div class="grid gap-2">
								<Label for="name">Full name</Label>
								<Input id="name" bind:value={applicantName} required placeholder="Your full name" disabled={isSubmitting} />
							</div>

							<div class="grid gap-4 sm:grid-cols-2">
								<div class="grid gap-2">
									<Label for="phone">Phone</Label>
									<Input id="phone" bind:value={applicantPhone} required type="tel" placeholder="+91 98765 43210" disabled={isSubmitting} />
								</div>
								<div class="grid gap-2">
									<Label for="email">Email</Label>
									<Input id="email" bind:value={applicantEmail} required type="email" placeholder="you@example.com" disabled={isSubmitting} />
								</div>
							</div>

							<div class="grid gap-2">
								<Label for="programme">Programme</Label>
								<Select.Root type="single" bind:value={programme} disabled={isSubmitting}>
									<Select.Trigger id="programme" class="w-full">
										{programme ? programme.toUpperCase() : "Select a programme"}
									</Select.Trigger>
									<Select.Content>
										<Select.Item value="btech">B.Tech (Computer Science, AI/ML, ECE, ME)</Select.Item>
										<Select.Item value="bpharm">B.Pharm (Pharmacy)</Select.Item>
										<Select.Item value="mba">MBA (Master of Business Administration)</Select.Item>
										<Select.Item value="bba">BBA (Bachelor of Business Administration)</Select.Item>
										<Select.Item value="bca">BCA (Bachelor of Computer Applications)</Select.Item>
										<Select.Item value="diploma">Diploma / Polytechnic</Select.Item>
									</Select.Content>
								</Select.Root>
							</div>

							<div class="grid gap-2">
								<Label for="question">Your question or comments</Label>
								<Textarea id="question" bind:value={applicantQuestion} placeholder="Tell us what you would like to know..." rows={3} disabled={isSubmitting} />
							</div>

							<Button type="submit" size="lg" class="gap-1.5 cursor-pointer" disabled={isSubmitting}>
								{#if isSubmitting}
									<IconLoader2 class="size-4 animate-spin" />
									<span>Submitting Enquiry...</span>
								{:else}
									<IconSend class="size-4" />
									<span>Submit Enquiry</span>
								{/if}
							</Button>
						</form>
					{/if}
				</CardContent>
			</Card>
		</section>

		<section class="border-t bg-muted/40">
			<div class="mx-auto max-w-4xl px-4 py-16 sm:px-6">
				<h2 class="font-heading text-3xl font-semibold">Common questions</h2>
				<Accordion.Root type="single" class="mt-6">
					<Accordion.Item value="eligibility">
						<Accordion.Trigger>Where can I confirm eligibility?</Accordion.Trigger>
						<Accordion.Content>
							Eligibility varies by programme and university rules. Send an enquiry with your prior qualification for current guidance.
						</Accordion.Content>
					</Accordion.Item>
					<Accordion.Item value="visit">
						<Accordion.Trigger>Can I visit the campus?</Accordion.Trigger>
						<Accordion.Content>
							Yes. Contact admissions to coordinate a suitable date and the department you want to visit.
						</Accordion.Content>
					</Accordion.Item>
					<Accordion.Item value="documents">
						<Accordion.Trigger>Which documents are usually required?</Accordion.Trigger>
						<Accordion.Content>
							Academic records, identity documents and photographs are generally requested. The admissions team will provide the current programme-specific checklist.
						</Accordion.Content>
					</Accordion.Item>
				</Accordion.Root>
			</div>
		</section>
	</main>
</PublicShell>
