<script lang="ts">

  import PublicShell from "$public/components/PublicShell.svelte";
  import PageHeader from "$public/components/PageHeader.svelte";
  import * as Empty from "$lib/components/ui/empty";
  import IconPhoto from "@tabler/icons-svelte/icons/photo";
  import building from "@/assets/bbdit-building-real.webp";
  import hostel from "@/assets/bbdit-hostel.jpg";
  import labs from "@/assets/bbdit-lab.jpg";
  import ground from "@/assets/bbdit-ground-real.jpg";
  import canteen from "@/assets/bbdit-mess.jpg";
  import mess from "@/assets/bbdit-mess.jpg";

  const sections = [
    { eyebrow: "Academic environment", title: "College building", image: building, alt: "BBDIT academic building and green campus", copy: "The academic complex is the centre of teaching at BBDIT. Large, well-lit lecture rooms, tutorial rooms, seminar halls and conference spaces bring daily classes, presentations and academic discussion together. The building sits within the institute’s green campus on the Delhi–Meerut corridor, with classrooms and shared spaces connected to the wider student environment." },
    { eyebrow: "Residential life", title: "Hostel", image: hostel, alt: "BBDIT hostel building", copy: "BBDIT publishes an on-campus, 400-bed hostel for boys, keeping resident students close to classrooms, laboratories and common facilities. The hostel includes Wi-Fi connectivity, continuous power backup and the everyday amenities needed for campus living. The institute also states that accommodation for girls is provided off campus in a peaceful and safe environment." },
    { eyebrow: "Dining", title: "Mess", image: mess, alt: "Dining area used by BBDIT students", copy: "The hostel mess is part of the daily residential routine, giving students a shared place for regular meals. More than a dining facility, it is an everyday community space where hostel residents meet outside class and maintain a dependable schedule around academic work." },
    { eyebrow: "Breaks and conversation", title: "Canteen", image: canteen, alt: "BBDIT college canteen", copy: "The campus canteen offers refreshments and an informal meeting place between lectures. Its open dining arrangement gives students space to pause, talk with classmates and spend time together without leaving the academic campus." },
    { eyebrow: "Sport and recreation", title: "College ground", image: ground, alt: "Students taking part in outdoor activity at BBDIT", copy: "Open grounds support outdoor games, practice sessions and college sporting events. BBDIT’s campus culture includes cricket, badminton and inter-college sporting activity, helping students build fitness, teamwork and a healthy break from classroom schedules." },
    { eyebrow: "Practical learning", title: "Laboratories", image: labs, alt: "Computer laboratory at BBDIT", copy: "Laboratory learning connects theory with supervised practice. Alongside its central, air-conditioned computer centre, BBDIT states that individual departments maintain computing and programme-specific laboratories. These spaces support software work, engineering tools, experiments and the practical requirements of engineering and pharmacy education." },
  ] as const;
</script>


<svelte:head><title>Campus Life · BBDIT</title></svelte:head>

<PublicShell>
  <main>
    <PageHeader compact eyebrow="Campus life" title="Life on campus" description="A closer look at the academic, residential, dining and recreational spaces shared by the BBDIT community." />

    <div class="divide-y">
      {#each sections as section, index}
        <section class={index % 2 === 1 ? "bg-muted/40" : ""}>
          <div class={`mx-auto grid max-w-6xl gap-8 px-4 py-12 sm:px-6 sm:py-14 lg:items-center lg:gap-14 lg:px-8 ${index % 2 === 1 ? "lg:grid-cols-[7fr_5fr]" : "lg:grid-cols-[5fr_7fr]"}`}>
            <img src={section.image} alt={section.alt} class={`h-56 w-full rounded-xl object-cover shadow-sm sm:h-64 lg:h-72 ${index % 2 === 1 ? "lg:order-2" : ""}`} />
            <div class={index % 2 === 1 ? "lg:order-1" : ""}>
              <p class="text-xs font-bold uppercase tracking-[.16em] text-primary">{section.eyebrow}</p>
              <h2 class="mt-2 font-heading text-2xl font-semibold tracking-tight sm:text-3xl">{section.title}</h2>
              <p class="mt-4 max-w-2xl text-[.95rem] leading-7 text-muted-foreground">{section.copy}</p>
            </div>
          </div>
        </section>
      {/each}
    </div>

    <section class="border-t bg-muted/40">
      <div class="mx-auto max-w-6xl px-4 py-14 sm:px-6 lg:px-8">
        <p class="text-xs font-bold uppercase tracking-widest text-primary">Gallery</p>
        <h2 class="mt-2 font-heading text-2xl font-semibold sm:text-3xl">Life At BBDIT</h2>
        <div class="mt-7 grid gap-5 md:grid-cols-2 lg:grid-cols-3">
          {#each ["Campus activities", "Student events", "Everyday moments"] as category}
            <Empty.Root class="min-h-52 border bg-background">
              <Empty.Header>
                <Empty.Media variant="icon"><IconPhoto /></Empty.Media>
                <Empty.Title>{category}</Empty.Title>
                <Empty.Description>Photos will be added here.</Empty.Description>
              </Empty.Header>
            </Empty.Root>
          {/each}
        </div>
      </div>
    </section>
  </main>
</PublicShell>
