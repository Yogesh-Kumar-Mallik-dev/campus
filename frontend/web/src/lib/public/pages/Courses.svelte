<script lang="ts">
	import PublicShell from '$public/components/PublicShell.svelte';
	import PageHeader from '$public/components/PageHeader.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import * as Empty from '$lib/components/ui/empty';
	import IconBook from '@tabler/icons-svelte/icons/book';

	const courses = [
		{
			group: 'Technology',
			name: 'B.Tech',
			full: 'Bachelor of Technology',
			duration: '4 years',
			branches: [
				'Computer Science & Engineering (CSE)',
				'Artificial Intelligence & Machine Learning (AI-ML)',
				'Data Science & Engineering (DS)',
				'Electronics & Communication Engineering (ECE)',
				'Mechanical Engineering (ME)',
				'Civil Engineering (CE)',
				'Electrical & Electronics Engineering (EEE)'
			],
			sampleSubjects: [
				{ code: 'MATH-101', name: 'Engineering Mathematics-I', type: 'Required' },
				{ code: 'CS-101', name: 'Data Structures & Algorithms', type: 'Required' },
				{ code: 'CS-201', name: 'Operating Systems & Architecture', type: 'Required' },
				{ code: 'AI-301', name: 'Machine Learning & Pattern Recognition', type: 'Elective' },
				{ code: 'CYBER-301', name: 'Cyber Security & Network Defense', type: 'Elective' }
			]
		},
		{
			group: 'Technology',
			name: 'M.Tech',
			full: 'Master of Technology',
			duration: '2 years',
			branches: [
				'Computer Science & Engineering (CSE)',
				'VLSI Design & Embedded Systems'
			],
			sampleSubjects: [
				{ code: 'CS-501', name: 'Advanced Distributed Computing', type: 'Required' },
				{ code: 'VLSI-501', name: 'CMOS Analog Circuit Design', type: 'Required' }
			]
		},
		{
			group: 'Computing',
			name: 'BCA',
			full: 'Bachelor of Computer Applications',
			duration: '3 years',
			branches: [],
			sampleSubjects: [
				{ code: 'CS-101', name: 'Data Structures & Algorithms', type: 'Required' },
				{ code: 'CS-301', name: 'Database Management Systems', type: 'Required' },
				{ code: 'HU-101', name: 'Professional Communication', type: 'Required' }
			]
		},
		{
			group: 'Management',
			name: 'BBA',
			full: 'Bachelor of Business Administration',
			duration: '3 years',
			branches: [],
			sampleSubjects: [
				{ code: 'MGMT-101', name: 'Principles of Management', type: 'Required' },
				{ code: 'HU-101', name: 'Corporate Communication', type: 'Required' }
			]
		},
		{
			group: 'Management',
			name: 'MBA',
			full: 'Master of Business Administration',
			duration: '2 years',
			branches: [
				'Finance Management',
				'Marketing & Analytics',
				'Human Resource Management'
			],
			sampleSubjects: [
				{ code: 'MGMT-101', name: 'Organizational Leadership & Behaviour', type: 'Required' },
				{ code: 'CS-301', name: 'Enterprise Database & Management', type: 'Elective' }
			]
		},
		{
			group: 'Pharmacy',
			name: 'B.Pharm',
			full: 'Bachelor of Pharmacy',
			duration: '4 years',
			branches: [],
			sampleSubjects: [
				{ code: 'PH-101', name: 'Pharmaceutics & Formulation', type: 'Required' },
				{ code: 'ENV-101', name: 'Environmental Studies', type: 'Required' }
			]
		},
		{
			group: 'Pharmacy',
			name: 'D.Pharm',
			full: 'Diploma in Pharmacy',
			duration: '2 years',
			branches: [],
			sampleSubjects: [
				{ code: 'PH-101', name: 'Pharmaceutics & Drug Delivery', type: 'Required' }
			]
		},
		{
			group: 'Diploma',
			name: 'Polytechnic',
			full: 'Diploma in Engineering',
			duration: '3 years',
			branches: [
				'Mechanical Engineering',
				'Civil Engineering',
				'Electrical Engineering'
			],
			sampleSubjects: [
				{ code: 'ME-101', name: 'Applied Thermodynamics & Fluid Mechanics', type: 'Required' }
			]
		}
	];

	const departments = ['Technology', 'Computing', 'Management', 'Pharmacy', 'Education', 'Diploma'];
</script>

<svelte:head>
	<title>Courses · BBDIT</title>
</svelte:head>

<PublicShell>
	<main>
		<PageHeader
			eyebrow="Departments & courses"
			title="Find the programme that fits your plans"
			description="The complete programme list currently published by BBDIT, organised by academic department."
		/>

		<section class="mx-auto max-w-7xl space-y-16 px-4 py-16 sm:px-6 lg:px-8">
			{#each departments as department}
				{@const departmentCourses = courses.filter((course) => course.group === department)}
				<section aria-labelledby={`department-${department.toLowerCase()}`}>
					<div class="mb-6 flex items-center gap-4">
						<h2 id={`department-${department.toLowerCase()}`} class="font-heading text-3xl font-semibold">
							{department} Department
						</h2>
						<div class="h-px flex-1 bg-border"></div>
						<Badge variant="secondary">
							{departmentCourses.length} {departmentCourses.length === 1 ? 'course' : 'courses'}
						</Badge>
					</div>

					{#if departmentCourses.length === 0}
						<Empty.Root class="rounded-xl border border-dashed py-10 bg-card/30">
							<Empty.Header>
								<Empty.Media variant="icon">
									<IconBook class="size-5 text-muted-foreground" />
								</Empty.Media>
								<Empty.Title>No Courses Listed</Empty.Title>
								<Empty.Description>There are no courses currently listed under {department}.</Empty.Description>
							</Empty.Header>
						</Empty.Root>
					{:else}
						<div class="grid gap-5 md:grid-cols-2 lg:grid-cols-3">
							{#each departmentCourses as course}
								<Card class="flex flex-col shadow-none">
									<CardHeader>
										<div class="flex items-center justify-between gap-3">
											<Badge variant="outline">{department}</Badge>
											<span class="text-xs text-muted-foreground">{course.duration}</span>
										</div>
										<CardTitle class="pt-3 font-heading text-2xl">{course.name}</CardTitle>
										<p class="text-sm text-muted-foreground">{course.full}</p>
									</CardHeader>
									<CardContent class="flex flex-1 flex-col">
										{#if course.branches.length}
											<div class="mb-4">
												<p class="mb-2 text-xs font-semibold uppercase tracking-wider text-foreground">Branches & Specialisations</p>
												<ul class="space-y-1.5 text-xs leading-5 text-muted-foreground">
													{#each course.branches as branch}
														<li>• {branch}</li>
													{/each}
												</ul>
											</div>
										{/if}

										{#if course.sampleSubjects?.length}
											<div class="mb-4 space-y-2">
												<p class="text-xs font-semibold uppercase tracking-wider text-foreground">Key Curriculum Subjects</p>
												<div class="flex flex-wrap gap-1.5">
													{#each course.sampleSubjects as sub}
														<span class="inline-flex items-center gap-1 rounded-md px-2 py-0.5 text-[11px] font-medium border {sub.type === 'Required' ? 'bg-primary/5 text-primary border-primary/20' : 'bg-amber-500/10 text-amber-700 dark:text-amber-300 border-amber-500/20'}">
															<span class="font-mono font-bold">{sub.code}</span>: {sub.name}
															<span class="text-[9px] uppercase font-bold opacity-75">({sub.type})</span>
														</span>
													{/each}
												</div>
											</div>
										{/if}

										<Button
											href="/public/admissions"
											variant="link"
											class="mt-auto h-auto justify-start px-0 pt-4"
										>
											Enquire about {course.name}
										</Button>
									</CardContent>
								</Card>
							{/each}
						</div>
					{/if}
				</section>
			{/each}
		</section>
	</main>
</PublicShell>
