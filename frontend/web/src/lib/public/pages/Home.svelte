<script lang="ts">

  import PublicShell from "$public/components/PublicShell.svelte";
  import { Button } from "$lib/components/ui/button";
  import { Card, CardContent, CardHeader, CardTitle } from "$lib/components/ui/card";
  import hero from "@/assets/bbdit-campus-hero.webp";
  import library from "@/assets/bbdit-library.jpg";
  import convocation from "@/assets/bbdit-convocation.jpg";
  import IconArrowRight from "@tabler/icons-svelte/icons/arrow-right";
  import IconBook2 from "@tabler/icons-svelte/icons/book-2";
  import IconBuildingBank from "@tabler/icons-svelte/icons/building-bank";
  import IconCalendarEvent from "@tabler/icons-svelte/icons/calendar-event";
  import IconFlask2 from "@tabler/icons-svelte/icons/flask-2";
  import IconRosetteDiscountCheck from "@tabler/icons-svelte/icons/rosette-discount-check";
  import IconUsersGroup from "@tabler/icons-svelte/icons/users-group";

  const programmes = [
    ["Engineering", "B.Tech & M.Tech", "CSE, AI/ML, Data Science, ECE, Electrical, Mechanical and Civil", IconFlask2],
    ["Pharmacy", "B.Pharm & D.Pharm", "Pharmaceutical education supported by practical laboratory learning", IconBook2],
    ["Management", "BBA & MBA", "Business fundamentals, leadership and industry-oriented learning", IconBuildingBank],
  ] as const;

  const notices = [
    ["26 Apr 2026", "BBDIT students excel at Sports Zonal Fest 2025"],
    ["24 Apr 2026", "External university examination notice"],
    ["15 Jan 2026", "Students shortlisted for Motherson Sumi interviews"],
  ];

</script>


<svelte:head><title>BBDIT Group of Institutions, Ghaziabad</title><meta name="description" content="BBDIT Group of Institutions — Developing Nation Since 2000." /></svelte:head>

<PublicShell>
  <main id="top">
    <section class="relative min-h-[clamp(32rem,75svh,37rem)] overflow-hidden">
      <img src={hero} alt="The BBDIT academic campus in Ghaziabad" class="absolute inset-0 size-full object-cover object-center" />
      <div class="absolute inset-0 bg-linear-to-r from-[#071a32]/90 via-[#071a32]/50 to-transparent"></div>
      <div class="relative mx-auto flex min-h-[clamp(32rem,75svh,37rem)] max-w-7xl items-center px-[clamp(1rem,4vw,2rem)] py-[clamp(4rem,12vw,6rem)]">
        <div class="max-w-2xl text-white"><p class="mb-4 text-xs font-semibold uppercase tracking-[0.14em] text-white/80 min-[380px]:text-sm min-[380px]:tracking-[0.18em]">BBDIT Group of Institutions · Ghaziabad</p><h1 class="text-balance font-heading text-[clamp(2.25rem,10vw,3.75rem)] font-semibold leading-[1.08] tracking-tight">Contemporary education.<br />Responsible futures.</h1><p class="mt-5 max-w-xl text-[clamp(1rem,4vw,1.125rem)] leading-7 text-white/80">A multicultural learning community for engineering, pharmacy, management and education—built on academic excellence and ethical responsibility.</p><div class="mt-8 flex flex-col items-stretch gap-3 min-[380px]:flex-row min-[380px]:flex-wrap min-[380px]:items-center"><Button size="lg" href="/public/courses">Explore courses <IconArrowRight data-icon="inline-end" /></Button><Button size="lg" variant="outline" href="/public/about" class="border-white/60 bg-white/10 text-white hover:bg-white hover:text-[#0b2546]">Why BBDIT</Button></div></div>
      </div>
    </section>

    <section class="border-b bg-muted/50"><div class="mx-auto grid max-w-7xl divide-y sm:grid-cols-3 sm:divide-x sm:divide-y-0"><a href="/public/about#since-2000" class="fact"><strong>25 years</strong><span>of institution building · Read the claim</span></a><a href="/public/about#campus" class="fact"><strong>15-acre</strong><span>green campus in NCR · Read the claim</span></a><a href="/public/about#recognition" class="fact"><strong>AKTU · AICTE</strong><span>affiliated and recognized · Read the claim</span></a></div></section>

    <section id="programmes" class="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8">
      <div class="mb-10 max-w-2xl"><p class="label">Academic departments</p><h2 class="heading">Find your department</h2><p class="copy">Explore departments spanning technology, pharmacy, management, computer applications and education.</p></div>
      <div class="grid gap-5 md:grid-cols-3">{#each programmes as item}{@const ProgramIcon = item[3]}<Card class="group shadow-none transition hover:border-primary/50 hover:shadow-md"><CardHeader><div class="mb-3 grid size-11 place-items-center rounded-lg bg-primary/10 text-primary"><ProgramIcon class="size-5" /></div><p class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">{item[0]} department</p><CardTitle class="font-heading text-2xl">{item[1]}</CardTitle></CardHeader><CardContent><p class="min-h-16 text-sm leading-6 text-muted-foreground">{item[2]}</p><Button href="/public/courses" variant="link" class="mt-4 h-auto px-0">View department courses <IconArrowRight data-icon="inline-end" /></Button></CardContent></Card>{/each}</div>
    </section>

    <section id="about" class="bg-[#0b2546] py-20 text-white"><div class="mx-auto grid max-w-7xl gap-12 px-4 sm:px-6 lg:grid-cols-2 lg:items-center lg:px-8"><div><p class="label text-white/60">About the institute</p><h2 class="heading text-white">A centre of excellence with a difference.</h2><p class="mt-5 leading-7 text-white/75">Established in 2000 under Rajiv Educational Trust, BBDIT was founded to deliver contemporary, quality education. Its vision is a multicultural community whose students contribute productively, responsibly and ethically in India and beyond.</p><div class="mt-8 grid grid-cols-2 gap-4 text-sm"><div class="border-t border-white/20 pt-4">52,000+ library books</div><div class="border-t border-white/20 pt-4">400-bed boys’ hostel</div><div class="border-t border-white/20 pt-4">Smart classrooms</div><div class="border-t border-white/20 pt-4">Wi-Fi enabled campus</div></div></div><img src={library} alt="BBDIT students studying in the institute library" class="w-full rounded-xl object-cover shadow-2xl" /></div></section>

    <section id="placements" class="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8"><div class="grid gap-10 lg:grid-cols-[.8fr_1.2fr]"><div><p class="label">Recognitions</p><h2 class="heading">A record built over time</h2><p class="copy">Published recognitions and institutional milestones, shown with their original context rather than as unsupported headline claims.</p></div><div class="grid gap-4 sm:grid-cols-2"><Card class="shadow-none"><CardContent class="p-6"><IconRosetteDiscountCheck class="mb-4 size-7 text-primary" /><h3 class="font-heading text-xl font-semibold">Best Engineering College in U.P.</h3><p class="mt-2 text-sm leading-6 text-muted-foreground">Presented at the Indian Education and EdTech Summit 2022.</p></CardContent></Card><Card class="shadow-none"><CardContent class="p-6"><IconRosetteDiscountCheck class="mb-4 size-7 text-primary" /><h3 class="font-heading text-xl font-semibold">Academic Excellence</h3><p class="mt-2 text-sm leading-6 text-muted-foreground">Recognition for contribution to technical education by AITMC.</p></CardContent></Card><Card class="shadow-none sm:col-span-2"><CardContent class="p-6"><div class="flex gap-4"><IconRosetteDiscountCheck class="mt-1 size-7 shrink-0 text-primary" /><div><h3 class="font-heading text-xl font-semibold">Placement and career building</h3><p class="mt-2 text-sm leading-6 text-muted-foreground">Annual Summit sustainability award from ESDA and PHDCCI, alongside continuing industry workshops, visits and placement activity.</p></div></div></CardContent></Card></div></div></section>

    <section id="campus" class="border-y bg-muted/40 py-20"><div class="mx-auto grid max-w-7xl gap-10 px-4 sm:px-6 lg:grid-cols-2 lg:items-center lg:px-8"><img src={convocation} alt="BBDIT graduates celebrating at convocation" class="aspect-[4/3] w-full rounded-xl object-cover" /><div><p class="label">Student life</p><h2 class="heading">Learning beyond the timetable</h2><p class="copy">Technical society Prodhyogeeki, the student newsletter Sampark and social initiative Arpan sit alongside workshops, industry visits, cultural events and sport.</p><div class="mt-7 grid gap-3 sm:grid-cols-2"><div class="mini"><IconUsersGroup /> Student societies</div><div class="mini"><IconFlask2 /> Technical workshops</div><div class="mini"><IconBuildingBank /> Industry visits</div><div class="mini"><IconCalendarEvent /> Encore college fest</div></div></div></div></section>

    <section id="notices" class="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8"><div class="flex flex-wrap items-end justify-between gap-4"><div><p class="label">From campus</p><h2 class="heading">Latest announcements</h2></div><Button href="/public/news" variant="outline">View all announcements</Button></div><div class="mt-8 divide-y rounded-xl border bg-card">{#each notices as notice}<a href="/public/news" class="grid gap-2 p-5 transition hover:bg-muted/50 sm:grid-cols-[130px_1fr_auto] sm:items-center"><span class="text-sm text-muted-foreground">{notice[0]}</span><strong class="font-medium">{notice[1]}</strong><IconArrowRight class="hidden size-4 sm:block" /></a>{/each}</div></section>

    <section class="bg-primary py-14 text-primary-foreground"><div class="mx-auto flex max-w-7xl flex-col items-start justify-between gap-6 px-4 sm:px-6 md:flex-row md:items-center lg:px-8"><div><h2 class="font-heading text-3xl font-semibold">Planning to join BBDIT?</h2><p class="mt-2 text-primary-foreground/75">Talk to admissions about eligibility, programmes and campus visits.</p></div><Button href="/public/admissions" size="lg" variant="secondary">Admission enquiry <IconArrowRight data-icon="inline-end" /></Button></div></section>
  </main>

</PublicShell>

<style>
  :global(html){scroll-behavior:smooth}.label{margin-bottom:.75rem;font-size:.72rem;font-weight:700;letter-spacing:.16em;text-transform:uppercase;color:var(--primary)}.heading{font-family:'Outfit Variable',sans-serif;font-size:clamp(2rem,4vw,3.25rem);font-weight:600;line-height:1.12;letter-spacing:-.025em;text-wrap:balance}.copy{margin-top:1rem;max-width:38rem;line-height:1.75;color:var(--muted-foreground)}.fact{display:flex;flex-direction:column;padding:1.5rem 2rem}.fact strong{font-family:'Outfit Variable',sans-serif;font-size:1.5rem}.fact span{margin-top:.15rem;font-size:.8rem;color:var(--muted-foreground)}.mini{display:flex;align-items:center;gap:.75rem;border-radius:.6rem;border:1px solid var(--border);background:var(--card);padding:1rem;font-size:.875rem;font-weight:500}.mini :global(svg){width:1.1rem;color:var(--primary)}@media(prefers-reduced-motion:reduce){:global(html){scroll-behavior:auto}}
</style>
