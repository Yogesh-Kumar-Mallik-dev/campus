<script lang="ts">
	/**
	 * BLOCK_WEB_ONBOARDING_FORM_001
	 * Subsystem: Rank 3 - Student & Staff Registration System (onboarding)
	 * Purpose:   Interactive multi-step KYC submission form with live validation and document attachments.
	 */
	import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Badge } from '$lib/components/ui/badge';
	import { Separator } from '$lib/components/ui/separator';
	import type { OnboardingType, Gender, DocumentType, OnboardingApplicant } from '$lib/types/onboarding';

	let {
		tenantId = 'tenant-demo',
		onSubmitSuccess = (_app: OnboardingApplicant) => {}
	}: {
		tenantId?: string;
		onSubmitSuccess?: (app: OnboardingApplicant) => void;
	} = $props();

	let currentStep = $state(1);
	let isSubmitting = $state(false);
	let errorMessage = $state<string | null>(null);
	let successApplicant = $state<OnboardingApplicant | null>(null);

	// Form State
	let applicantType = $state<OnboardingType>('STUDENT');
	let firstName = $state('');
	let lastName = $state('');
	let email = $state('');
	let phone = $state('');
	let dateOfBirth = $state('2006-01-01');
	let gender = $state<Gender>('MALE');
	let bloodGroup = $state('O+');
	let nationality = $state('Indian');

	let emergencyName = $state('');
	let emergencyPhone = $state('');
	let emergencyRelation = $state('Parent');

	let addressLine1 = $state('');
	let addressLine2 = $state('');
	let city = $state('');
	let stateName = $state('');
	let postalCode = $state('');

	let programCode = $state('BTECH_CSE');
	let departmentCode = $state('CSE');
	let academicYear = $state('2026-2027');
	let targetDesignation = $state('Assistant Professor');

	let guardianName = $state('');
	let guardianPhone = $state('');
	let guardianRelation = $state('Father');

	// Uploaded Documents state
	interface UploadItem {
		type: DocumentType;
		label: string;
		fileName?: string;
		fileSize?: number;
		fileKey?: string;
		isUploaded: boolean;
	}

	let documents = $state<UploadItem[]>([
		{ type: 'NATIONAL_ID', label: 'National ID / Aadhaar Card / Passport', isUploaded: false },
		{ type: 'ACADEMIC_TRANSCRIPT', label: '10th & 12th Grade / Previous Marksheets', isUploaded: false },
		{ type: 'PASSPORT_PHOTO', label: 'Passport Size Color Photograph', isUploaded: false },
	]);

	$effect(() => {
		if (applicantType === 'STUDENT') {
			documents = [
				{ type: 'NATIONAL_ID', label: 'National ID / Aadhaar Card / Passport', isUploaded: false },
				{ type: 'ACADEMIC_TRANSCRIPT', label: '10th & 12th Grade / Previous Marksheets', isUploaded: false },
				{ type: 'PASSPORT_PHOTO', label: 'Passport Size Color Photograph', isUploaded: false },
			];
		} else {
			documents = [
				{ type: 'NATIONAL_ID', label: 'National ID / Government ID Proof', isUploaded: false },
				{ type: 'DEGREE_CERTIFICATE', label: 'Highest Degree / PhD / Master Certificate', isUploaded: false },
				{ type: 'PASSPORT_PHOTO', label: 'Official Passport Size Photograph', isUploaded: false },
			];
		}
	});

	function handleSimulateUpload(index: number) {
		const doc = documents[index];
		doc.fileName = `${doc.type.toLowerCase()}_verified.pdf`;
		doc.fileSize = 180000;
		doc.fileKey = `uploads/kyc/${doc.type.toLowerCase()}_${Date.now()}.pdf`;
		doc.isUploaded = true;
	}

	async function handleSubmit() {
		isSubmitting = true;
		errorMessage = null;

		try {
			// Step 1: Create Draft
			const draftPayload = {
				tenant_id: tenantId,
				type: applicantType,
				first_name: firstName,
				last_name: lastName,
				email,
				phone,
				date_of_birth: dateOfBirth,
				gender,
				blood_group: bloodGroup,
				nationality,
				emergency_contact: {
					name: emergencyName,
					phone: emergencyPhone,
					relation: emergencyRelation,
				},
				address: {
					line1: addressLine1,
					line2: addressLine2,
					city,
					state: stateName,
					postal_code: postalCode,
					country: 'India',
				},
				program_id: applicantType === 'STUDENT' ? programCode : undefined,
				department_id: applicantType !== 'STUDENT' ? departmentCode : undefined,
				academic_year: academicYear,
				target_designation: applicantType !== 'STUDENT' ? targetDesignation : undefined,
				guardian: applicantType === 'STUDENT' ? {
					name: guardianName,
					phone: guardianPhone,
					relation: guardianRelation,
				} : undefined,
			};

			const draftRes = await fetch('/api/v1/onboarding/applicants', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(draftPayload),
			});

			if (!draftRes.ok) {
				const errData = await draftRes.json().catch(() => ({ detail: 'Failed to create application draft' }));
				throw new Error(errData.detail || 'Error creating application');
			}

			const draftData = await draftRes.json();
			const applicantId = draftData.applicant.id;

			// Step 2: Attach Documents
			for (const doc of documents) {
				if (doc.isUploaded) {
					await fetch(`/api/v1/onboarding/applicants/${applicantId}/documents`, {
						method: 'POST',
						headers: { 'Content-Type': 'application/json' },
						body: JSON.stringify({
							tenant_id: tenantId,
							document_type: doc.type,
							file_key: doc.fileKey,
							file_name: doc.fileName,
							file_size: doc.fileSize,
							mime_type: 'application/pdf',
						}),
					});
				}
			}

			// Step 3: Submit Application
			const submitRes = await fetch(`/api/v1/onboarding/applicants/${applicantId}/submit?tenant_id=${tenantId}`, {
				method: 'POST',
			});

			if (!submitRes.ok) {
				const errData = await submitRes.json().catch(() => ({ detail: 'Document validation or submission failed' }));
				throw new Error(errData.detail || 'Error submitting application');
			}

			const submitData = await submitRes.json();
			successApplicant = submitData.applicant;
			onSubmitSuccess(submitData.applicant);
		} catch (err: unknown) {
			if (err instanceof Error) {
				errorMessage = err.message;
			} else {
				errorMessage = 'An unexpected submission error occurred';
			}
		} finally {
			isSubmitting = false;
		}
	}
</script>

<div class="max-w-3xl mx-auto space-y-6">
	{#if successApplicant}
		<Card class="border-emerald-500/40 bg-emerald-500/5 shadow-sm">
			<CardHeader class="space-y-2">
				<div class="flex items-center gap-2">
					<Badge variant="outline" class="bg-emerald-500/10 text-emerald-700 dark:text-emerald-300 border-emerald-500/30">
						✓ APPLICATION SUBMITTED
					</Badge>
					<span class="text-xs font-mono text-muted-foreground">ID: {successApplicant.id}</span>
				</div>
				<CardTitle class="text-2xl font-bold tracking-tight text-foreground">
					Admissions KYC Successfully Registered
				</CardTitle>
				<CardDescription class="text-muted-foreground">
					Your onboarding application has been queued for registrar and departmental verification.
				</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4 text-sm">
				<div class="grid grid-cols-2 gap-4 p-4 rounded-lg bg-card border border-border">
					<div>
						<span class="text-xs text-muted-foreground">Applicant Name</span>
						<p class="font-semibold text-foreground">{successApplicant.first_name} {successApplicant.last_name}</p>
					</div>
					<div>
						<span class="text-xs text-muted-foreground">Category</span>
						<p class="font-semibold text-foreground">{successApplicant.type}</p>
					</div>
					<div>
						<span class="text-xs text-muted-foreground">Email Contact</span>
						<p class="font-semibold text-foreground">{successApplicant.email}</p>
					</div>
					<div>
						<span class="text-xs text-muted-foreground">Current Status</span>
						<p class="font-semibold text-emerald-600 dark:text-emerald-400">{successApplicant.status}</p>
					</div>
				</div>
				<Button onclick={() => { successApplicant = null; currentStep = 1; }} variant="outline" class="w-full">
					Submit Another Registration
				</Button>
			</CardContent>
		</Card>
	{:else}
		<!-- Stepper Header -->
		<Card class="shadow-sm">
			<CardHeader class="pb-4">
				<div class="flex items-center justify-between">
					<div>
						<CardTitle class="text-xl font-bold tracking-tight">Student & Staff Registration (KYC)</CardTitle>
						<CardDescription class="mt-0.5">
							Complete institutional registration, upload KYC documents, and secure your verified campus profile.
						</CardDescription>
					</div>
					<Badge variant="secondary" class="font-mono text-xs">Step {currentStep} of 4</Badge>
				</div>

				<!-- Step indicators -->
				<div class="grid grid-cols-4 gap-2 pt-4">
					{#each ['Identity', 'Contacts', 'Program', 'Documents'] as stepName, i}
						<div class="space-y-1">
							<div class="h-1.5 rounded-full {currentStep >= i + 1 ? 'bg-primary' : 'bg-muted'}"></div>
							<span class="text-[11px] font-medium {currentStep === i + 1 ? 'text-primary font-bold' : 'text-muted-foreground'}">
								{stepName}
							</span>
						</div>
					{/each}
				</div>
			</CardHeader>
		</Card>

		{#if errorMessage}
			<div class="p-4 rounded-lg border border-destructive/30 bg-destructive/10 text-destructive text-sm font-medium">
				⚠️ {errorMessage}
			</div>
		{/if}

		<!-- Step Content -->
		<Card class="shadow-sm">
			<CardContent class="p-6 space-y-6">
				{#if currentStep === 1}
					<!-- Step 1: Persona & Personal Details -->
					<div class="space-y-4">
						<div class="space-y-1.5">
							<Label class="text-xs font-semibold">Registration Category</Label>
							<div class="grid grid-cols-3 gap-3">
								{#each ['STUDENT', 'FACULTY', 'STAFF'] as type}
									<button
										type="button"
										onclick={() => applicantType = type as OnboardingType}
										class="p-3 rounded-lg border text-sm font-semibold transition text-center {applicantType === type ? 'border-primary bg-primary/10 text-primary' : 'border-border bg-card text-muted-foreground hover:border-primary/40'}"
									>
										{type === 'STUDENT' ? '🎓 Student' : type === 'FACULTY' ? '👨‍🏫 Faculty' : '🏢 Staff'}
									</button>
								{/each}
							</div>
						</div>

						<Separator />

						<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
							<div class="space-y-1.5">
								<Label for="firstName" class="text-xs font-semibold">First Name *</Label>
								<Input id="firstName" bind:value={firstName} placeholder="e.g. Aarav" required />
							</div>
							<div class="space-y-1.5">
								<Label for="lastName" class="text-xs font-semibold">Last Name *</Label>
								<Input id="lastName" bind:value={lastName} placeholder="e.g. Sharma" required />
							</div>
							<div class="space-y-1.5">
								<Label for="email" class="text-xs font-semibold">Email Address *</Label>
								<Input id="email" type="email" bind:value={email} placeholder="name@example.com" required />
							</div>
							<div class="space-y-1.5">
								<Label for="phone" class="text-xs font-semibold">Phone Number *</Label>
								<Input id="phone" type="tel" bind:value={phone} placeholder="+91 98765 43210" required />
							</div>
							<div class="space-y-1.5">
								<Label for="dob" class="text-xs font-semibold">Date of Birth *</Label>
								<Input id="dob" type="date" bind:value={dateOfBirth} required />
							</div>
							<div class="space-y-1.5">
								<Label for="gender" class="text-xs font-semibold">Gender</Label>
								<select
									id="gender"
									bind:value={gender}
									class="w-full h-9 px-3 rounded-md border border-input bg-background text-foreground text-sm focus:outline-none focus:ring-2 focus:ring-ring"
								>
									<option value="MALE">Male</option>
									<option value="FEMALE">Female</option>
									<option value="OTHER">Other</option>
								</select>
							</div>
						</div>
					</div>
				{:else if currentStep === 2}
					<!-- Step 2: Address & Emergency Contact -->
					<div class="space-y-4">
						<h3 class="text-sm font-bold tracking-tight text-foreground">Residential Address</h3>
						<div class="space-y-3">
							<div class="space-y-1.5">
								<Label for="addr1" class="text-xs font-semibold">Address Line 1 *</Label>
								<Input id="addr1" bind:value={addressLine1} placeholder="House / Flat No., Street, Locality" required />
							</div>
							<div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
								<div class="space-y-1.5">
									<Label for="city" class="text-xs font-semibold">City *</Label>
									<Input id="city" bind:value={city} placeholder="e.g. Bengaluru" required />
								</div>
								<div class="space-y-1.5">
									<Label for="state" class="text-xs font-semibold">State *</Label>
									<Input id="state" bind:value={stateName} placeholder="e.g. Karnataka" required />
								</div>
								<div class="space-y-1.5">
									<Label for="pincode" class="text-xs font-semibold">Postal Code *</Label>
									<Input id="pincode" bind:value={postalCode} placeholder="560001" required />
								</div>
							</div>
						</div>

						<Separator />

						<h3 class="text-sm font-bold tracking-tight text-foreground">Emergency Contact</h3>
						<div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
							<div class="space-y-1.5">
								<Label for="emName" class="text-xs font-semibold">Contact Name *</Label>
								<Input id="emName" bind:value={emergencyName} placeholder="Full Name" required />
							</div>
							<div class="space-y-1.5">
								<Label for="emPhone" class="text-xs font-semibold">Phone Number *</Label>
								<Input id="emPhone" bind:value={emergencyPhone} placeholder="+91..." required />
							</div>
							<div class="space-y-1.5">
								<Label for="emRel" class="text-xs font-semibold">Relationship *</Label>
								<Input id="emRel" bind:value={emergencyRelation} placeholder="e.g. Father, Mother" required />
							</div>
						</div>

						{#if applicantType === 'STUDENT'}
							<Separator />
							<h3 class="text-sm font-bold tracking-tight text-foreground">Parent / Guardian Details</h3>
							<div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
								<div class="space-y-1.5">
									<Label for="gdName" class="text-xs font-semibold">Guardian Name *</Label>
									<Input id="gdName" bind:value={guardianName} placeholder="Parent Name" required />
								</div>
								<div class="space-y-1.5">
									<Label for="gdPhone" class="text-xs font-semibold">Guardian Phone *</Label>
									<Input id="gdPhone" bind:value={guardianPhone} placeholder="+91..." required />
								</div>
								<div class="space-y-1.5">
									<Label for="gdRel" class="text-xs font-semibold">Relationship *</Label>
									<Input id="gdRel" bind:value={guardianRelation} placeholder="Father" required />
								</div>
							</div>
						{/if}
					</div>
				{:else if currentStep === 3}
					<!-- Step 3: Academic / Designation Target -->
					<div class="space-y-4">
						{#if applicantType === 'STUDENT'}
							<div class="space-y-3">
								<div class="space-y-1.5">
									<Label for="prog" class="text-xs font-semibold">Target Degree Program *</Label>
									<select
										id="prog"
										bind:value={programCode}
										class="w-full h-9 px-3 rounded-md border border-input bg-background text-foreground text-sm focus:outline-none focus:ring-2 focus:ring-ring"
									>
										<option value="BTECH_CSE">B.Tech - Computer Science & Engineering (4 Years)</option>
										<option value="BTECH_ECE">B.Tech - Electronics & Communication (4 Years)</option>
										<option value="BTECH_MECH">B.Tech - Mechanical Engineering (4 Years)</option>
										<option value="MTECH_AI">M.Tech - Artificial Intelligence (2 Years)</option>
									</select>
								</div>
								<div class="space-y-1.5">
									<Label for="ay" class="text-xs font-semibold">Academic Year</Label>
									<Input id="ay" bind:value={academicYear} placeholder="2026-2027" />
								</div>
							</div>
						{:else}
							<div class="space-y-3">
								<div class="space-y-1.5">
									<Label for="dept" class="text-xs font-semibold">Assigned Department *</Label>
									<select
										id="dept"
										bind:value={departmentCode}
										class="w-full h-9 px-3 rounded-md border border-input bg-background text-foreground text-sm focus:outline-none focus:ring-2 focus:ring-ring"
									>
										<option value="CSE">Department of Computer Science & Engineering</option>
										<option value="ECE">Department of Electronics & Communication</option>
										<option value="MECH">Department of Mechanical Engineering</option>
										<option value="ADMIN">Administrative & Operations Staff</option>
									</select>
								</div>
								<div class="space-y-1.5">
									<Label for="desig" class="text-xs font-semibold">Position / Designation *</Label>
									<Input id="desig" bind:value={targetDesignation} placeholder="e.g. Assistant Professor" required />
								</div>
							</div>
						{/if}
					</div>
				{:else if currentStep === 4}
					<!-- Step 4: Document Attachments -->
					<div class="space-y-4">
						<div class="space-y-1">
							<h3 class="text-sm font-bold tracking-tight text-foreground">Required KYC Document Checklist</h3>
							<p class="text-xs text-muted-foreground">
								All mandatory documents must be attached before submission to proceed with institutional verification.
							</p>
						</div>

						<div class="space-y-3">
							{#each documents as doc, idx}
								<div class="p-4 rounded-lg border border-border bg-card flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
									<div class="space-y-1">
										<div class="flex items-center gap-2">
											<span class="text-sm font-semibold text-foreground">{doc.label}</span>
											<Badge variant="outline" class="text-[10px] uppercase font-mono">{doc.type}</Badge>
										</div>
										{#if doc.isUploaded}
											<p class="text-xs text-emerald-600 dark:text-emerald-400 font-medium">
												✓ Attached: {doc.fileName} ({(doc.fileSize! / 1024).toFixed(0)} KB)
											</p>
										{:else}
											<p class="text-xs text-muted-foreground">PDF or JPEG, max 5MB</p>
										{/if}
									</div>

									<div>
										{#if doc.isUploaded}
											<Badge variant="outline" class="bg-emerald-500/10 text-emerald-700 dark:text-emerald-300 border-emerald-500/30">
												Ready
											</Badge>
										{:else}
											<Button size="sm" variant="secondary" onclick={() => handleSimulateUpload(idx)}>
												Upload File
											</Button>
										{/if}
									</div>
								</div>
							{/each}
						</div>
					</div>
				{/if}
			</CardContent>

			<CardFooter class="pt-4 border-t border-border flex items-center justify-between">
				<Button
					type="button"
					variant="outline"
					disabled={currentStep === 1 || isSubmitting}
					onclick={() => currentStep = Math.max(1, currentStep - 1)}
				>
					← Previous
				</Button>

				{#if currentStep < 4}
					<Button
						type="button"
						onclick={() => currentStep = Math.min(4, currentStep + 1)}
					>
						Next Step →
					</Button>
				{:else}
					<Button
						type="button"
						disabled={isSubmitting || documents.some(d => !d.isUploaded)}
						onclick={handleSubmit}
					>
						{isSubmitting ? 'Registering...' : 'Submit KYC Application'}
					</Button>
				{/if}
			</CardFooter>
		</Card>
	{/if}
</div>
