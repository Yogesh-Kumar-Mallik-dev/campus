<script lang="ts">
  import { onMount, type Snippet } from "svelte";
  import { Button } from "$lib/components/ui/button";
  import { Avatar, AvatarFallback } from "$lib/components/ui/avatar";
  import { Switch } from "$lib/components/ui/switch";
  import { page } from "$app/state";
  import logo from "@/assets/bbdit-logo-transparent.png";
  import IconMenu2 from "@tabler/icons-svelte/icons/menu-2";
  import IconX from "@tabler/icons-svelte/icons/x";
  import IconBrandFacebook from "@tabler/icons-svelte/icons/brand-facebook";
  import IconBrandInstagram from "@tabler/icons-svelte/icons/brand-instagram";
  import IconBrandX from "@tabler/icons-svelte/icons/brand-x";
  import IconMoon from "@tabler/icons-svelte/icons/moon";
  import IconSun from "@tabler/icons-svelte/icons/sun";
  import IconLayoutDashboard from "@tabler/icons-svelte/icons/layout-dashboard";
  import { theme } from "$lib/theme.svelte";
  import { authState } from "$lib/auth-state.svelte";

  let { children }: { children: Snippet } = $props();
  let menuOpen = $state(false);

  const currentUser = $derived(
    authState.user
      ? { displayName: authState.user.username, email: authState.user.email }
      : (page.data?.user as { displayName?: string; email?: string } | null)
  );
  const isAuthenticated = $derived(Boolean(currentUser) || authState.isAuthenticated);
  const initials = $derived(
    currentUser?.displayName
      ? currentUser.displayName
          .split(" ")
          .filter(Boolean)
          .map((n: string) => n[0])
          .slice(0, 2)
          .join("")
          .toUpperCase()
      : "CM"
  );

  const links = [
    ["/public", "Home"],
    ["/public/about", "About"],
    ["/public/courses", "Courses"],
    ["/public/academics", "Academics"],
    ["/public/placements", "Placements"],
    ["/public/campus", "Campus"],
    ["/public/news", "News"]
  ] as const;

  function isNavActive(path: (typeof links)[number][0]) {
    return page.url.pathname === path || (path === "/public/courses" && page.url.pathname === "/public/admissions");
  }
</script>

<!-- Shared shell for public-facing routes. -->
<div class="min-h-screen-dvh min-h-screen bg-background text-foreground flex flex-col">
  <header class="sticky top-0 z-40 border-b bg-background/95 text-foreground shadow-sm backdrop-blur supports-[backdrop-filter]:bg-background/85 pt-[env(safe-area-inset-top,0px)]">
    <div class="mx-auto flex h-16 sm:h-20 max-w-7xl items-center justify-between gap-2 px-[clamp(.75rem,4vw,2rem)]">
      <!-- Brand Logo -->
      <a href="/public" aria-label="BBDIT home" class="min-w-0 shrink rounded-lg border border-black/8 bg-white/95 px-2 py-1 shadow-xs ring-1 ring-black/3 dark:border-white/12 dark:bg-[#f2f5f7] dark:ring-white/5">
        <img src={logo} alt="Babu Banarsi Das Institute of Technology" class="h-auto max-h-8 sm:max-h-10 w-auto max-w-[min(180px,calc(100vw-9rem))] sm:max-w-[230px] object-contain" />
      </a>

      <!-- Desktop Navigation Links -->
      <nav class="hidden items-center gap-1 text-sm font-medium lg:flex" aria-label="Main navigation">
        {#each links as link}
          <a href={link[0]} class:active={isNavActive(link[0])} aria-current={isNavActive(link[0]) ? "page" : undefined}>{link[1]}</a>
        {/each}
      </nav>

      <!-- Right Header Actions (Side-by-side with Burger Menu on all screen sizes) -->
      <div class="flex items-center gap-1.5 sm:gap-2 shrink-0">
        {#if isAuthenticated}
          <!-- Authenticated header state -->
          <a
            href="/hub"
            class="flex items-center gap-2 rounded-full border border-border bg-card/60 p-1 sm:px-2.5 sm:py-1 shadow-2xs hover:bg-muted/40 transition-colors"
            title="Open Campus Dashboard"
          >
            <Avatar class="size-7 sm:size-7.5 border border-border">
              <AvatarFallback class="bg-primary/10 text-[11px] font-heading font-semibold text-primary">
                {initials}
              </AvatarFallback>
            </Avatar>
            <span class="hidden sm:inline text-xs font-semibold text-foreground max-w-[110px] md:max-w-[140px] truncate">
              {currentUser?.displayName}
            </span>
          </a>
          <Button href="/hub" size="sm" class="hidden min-[420px]:inline-flex h-8.5 sm:h-9 gap-1.5 shadow-xs">
            <IconLayoutDashboard class="size-4" /> The Hub
          </Button>
        {:else}
          <!-- Unauthenticated header state: Buttons remain visible until screen is < 480px -->
          <Button variant="ghost" size="sm" href="/auth/login" class="hidden min-[480px]:inline-flex h-8.5 sm:h-9 px-3">
            Sign in
          </Button>
          <Button size="sm" href="/onboarding" class="hidden min-[480px]:inline-flex h-8.5 sm:h-9 px-3.5 shadow-xs">
            Apply Now
          </Button>
        {/if}

        <!-- Mobile Hamburger Menu Button (Adjacent to action buttons) -->
        <Button
          variant="outline"
          size="icon"
          class="shrink-0 lg:hidden size-8.5 sm:size-9"
          onclick={() => menuOpen = !menuOpen}
          aria-label="Toggle navigation menu"
        >
          {#if menuOpen}<IconX class="size-4.5 sm:size-5" />{:else}<IconMenu2 class="size-4.5 sm:size-5" />{/if}
        </Button>
      </div>
    </div>

    <!-- Mobile Drawer -->
    {#if menuOpen}
      <nav class="grid max-h-[calc(100dvh-5rem)] gap-1 overflow-y-auto border-t px-[clamp(.75rem,4vw,2rem)] py-4 text-sm font-medium lg:hidden pb-[max(1rem,env(safe-area-inset-bottom))]">
        {#each links as link}
          <a href={link[0]} class:active={isNavActive(link[0])} aria-current={isNavActive(link[0]) ? "page" : undefined} onclick={() => menuOpen = false}>
            {link[1]}
          </a>
        {/each}

        <!-- Non-duplicated Auth State in Drawer -->
        {#if isAuthenticated}
          <div class="mt-3 flex items-center justify-between rounded-lg border border-border bg-muted/40 p-3 min-[420px]:hidden">
            <div class="flex items-center gap-2.5 min-w-0">
              <Avatar class="size-9 border border-border shrink-0">
                <AvatarFallback class="bg-primary/10 text-xs font-heading font-semibold text-primary">
                  {initials}
                </AvatarFallback>
              </Avatar>
              <div class="min-w-0 flex-1">
                <p class="truncate text-xs font-semibold text-foreground">{currentUser?.displayName}</p>
                <p class="truncate text-[10px] text-muted-foreground">{currentUser?.email}</p>
              </div>
            </div>
            <Button href="/hub" size="sm" class="h-8 text-xs shrink-0 ml-2">The Hub</Button>
          </div>
        {:else}
          <div class="mt-3 grid grid-cols-2 gap-2 min-[480px]:hidden">
            <Button variant="outline" class="min-w-0 h-9" href="/auth/login" onclick={() => menuOpen = false}>
              Sign in
            </Button>
            <Button class="min-w-0 h-9 shadow-xs" href="/onboarding" onclick={() => menuOpen = false}>
              Apply Now
            </Button>
          </div>
        {/if}
      </nav>
    {/if}
  </header>

  {#if !menuOpen}
    <div class="fixed left-3 top-[calc(4.75rem+env(safe-area-inset-top,0px))] sm:top-[calc(5.75rem+env(safe-area-inset-top,0px))] z-30 rounded-full border bg-background/95 p-1.5 shadow-lg backdrop-blur sm:left-4">
      <Switch size="icon" checked={theme.isDark} onCheckedChange={(val) => theme.set(val)} aria-label="Use dark mode" title={theme.isDark ? "Dark mode on" : "Dark mode off"}>
        {#snippet thumbContent()}
          {#if theme.isDark}
            <IconMoon class="size-3.5 text-primary" aria-hidden="true" />
          {:else}
            <IconSun class="size-3.5 text-foreground" aria-hidden="true" />
          {/if}
        {/snippet}
      </Switch>
    </div>
  {/if}

  <main class="flex-1 min-w-0">
    {@render children()}
  </main>

  <footer class="bg-[#071a32] py-12 text-white/70">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="grid gap-8 border-b border-white/10 pb-10 md:grid-cols-[1.5fr_1fr_1fr]">
        <div>
          <p class="font-heading text-xl font-semibold text-white">BBDIT Group of Institutions</p>
          <p class="mt-3 max-w-md text-sm leading-6">7th Km, Delhi–Meerut Road, near Duhai, Ghaziabad, Uttar Pradesh 201206</p>
        </div>
        <div>
          <h3>Follow BBDIT</h3>
          <div class="flex gap-2">
            <a class="social" href="https://www.facebook.com/Babu-Banarsi-Dass-Institute-of-Technology-378556229392036/" target="_blank" rel="noreferrer" aria-label="BBDIT on Facebook">
              <IconBrandFacebook />
            </a>
            <a class="social" href="https://www.twitter.com/bbdit_ghaziabad" target="_blank" rel="noreferrer" aria-label="BBDIT on X">
              <IconBrandX />
            </a>
            <a class="social" href="https://instagram.com/bbdit.ghaziabad" target="_blank" rel="noreferrer" aria-label="BBDIT on Instagram">
              <IconBrandInstagram />
            </a>
          </div>
        </div>
        <div>
          <h3>Contact</h3>
          <a href="tel:+911202675912">0120 2675912</a>
          <a href="mailto:contact@bbdit.edu.in">contact@bbdit.edu.in</a>
        </div>
      </div>
      <p class="pt-6 text-xs">© 2026 Babu Banarsi Das Institute of Technology</p>
    </div>
  </footer>
</div>

<style>
  header nav a {
    border-radius: 0.5rem;
    padding: 0.6rem 0.75rem;
    transition: background-color 0.15s, color 0.15s;
  }
  header nav a:hover {
    color: var(--primary);
    background: var(--muted);
  }
  header nav a.active {
    color: var(--primary);
    background: color-mix(in oklab, var(--primary) 11%, transparent);
    font-weight: 700;
  }
  footer h3 {
    margin-bottom: 0.75rem;
    font-size: 0.75rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.12em;
    color: white;
  }
  footer a, footer p {
    display: block;
    margin-top: 0.4rem;
    font-size: 0.875rem;
  }
  footer a:hover {
    color: white;
  }
  footer .social {
    display: grid;
    width: 2.5rem;
    height: 2.5rem;
    place-items: center;
    border: 1px solid rgb(255 255 255 / 0.2);
    border-radius: 999px;
    transition: 0.15s;
  }
  footer .social:hover {
    border-color: white;
    background: rgb(255 255 255 / 0.1);
  }
  footer .social :global(svg) {
    width: 1.1rem;
  }
</style>
