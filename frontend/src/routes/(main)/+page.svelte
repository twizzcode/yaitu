<script>
    import Button from "#lib/components/ui/button/button.svelte";
    import CardContent from "#lib/components/ui/card/card-content.svelte";
    import Card from "#lib/components/ui/card/card.svelte";
    import BadgeMain from "#lib/components/main/badge-main.svelte"
    import Badge from "#lib/components/ui/badge/badge.svelte";
    
    import ChevronRight from "@lucide/svelte/icons/chevron-right";
    import FileCheck from "@lucide/svelte/icons/file-check";
    import TrendingUp from "@lucide/svelte/icons/trending-up";
    import CheckCheck from "@lucide/svelte/icons/check-check";
    import Star from "@lucide/svelte/icons/star";
    import PanelsTopLeft from "@lucide/svelte/icons/panels-top-left";
    import Workflow from "@lucide/svelte/icons/workflow";
    import ChartNoAxesCombined from "@lucide/svelte/icons/chart-no-axes-combined";
    import CalendarCheck from "@lucide/svelte/icons/calendar-check";
    import CreditCard from "@lucide/svelte/icons/credit-card";
    import Users from "@lucide/svelte/icons/users";
    import WalletCards from "@lucide/svelte/icons/wallet-cards";
    import LineChart from "@lucide/svelte/icons/line-chart";
    import Bot from "@lucide/svelte/icons/bot";
    import Dumbbell from "@lucide/svelte/icons/dumbbell";
    import Volleyball from "@lucide/svelte/icons/volleyball";
    import Baby from "@lucide/svelte/icons/baby";
    import PersonStanding from "@lucide/svelte/icons/person-standing";
    import Trees from "@lucide/svelte/icons/trees";
    import Armchair from "@lucide/svelte/icons/armchair";
    import Footprints from "@lucide/svelte/icons/footprints";
    import Flower2 from "@lucide/svelte/icons/flower-2";
    import Bike from "@lucide/svelte/icons/bike";
    import Fence from "@lucide/svelte/icons/fence";
    import PawPrint from "@lucide/svelte/icons/paw-print";
    import Waves from "@lucide/svelte/icons/waves";
    import Mail from "@lucide/svelte/icons/mail";
    
    import AccordionTrigger from "#lib/components/ui/accordion/accordion-trigger.svelte";
    import Accordion from "#lib/components/ui/accordion/accordion.svelte";
    import AccordionItem from "#lib/components/ui/accordion/accordion-item.svelte";
    import AccordionContent from "#lib/components/ui/accordion/accordion-content.svelte";
    import { fly } from "svelte/transition";
    
    

    const coreValues = [
      {
        icon: CheckCheck,
        title: "Booking lebih rapi",
        content:
          "Kelola jadwal lapangan real-time supaya pelanggan bisa memilih slot tanpa antre chat admin.",
      },
      {
        icon: TrendingUp,
        title: "Operasional lebih cepat",
        content:
          "Pantau pelanggan, pembayaran, dan riwayat booking dari satu dashboard yang mudah dipakai tim.",
      },
      {
        icon: FileCheck,
        title: "Laporan lebih jelas",
        content:
          "Lihat pemasukan, performa lapangan, dan tren booking untuk mengambil keputusan bisnis lebih tepat.",
      },
    ]

    const features = [
      {
        icon: PanelsTopLeft,
        title: "Unified workspace",
        content:
          "Kelola jadwal, pelanggan, transaksi, dan laporan lapangan dari satu dashboard yang mudah dipakai.",
      },
      {
        icon: Workflow,
        title: "Workflow automation",
        content:
          "Kurangi kerja manual untuk booking, konfirmasi pembayaran, dan update jadwal harian.",
      },
      {
        icon: ChartNoAxesCombined,
        title: "Actionable analytics",
        content:
          "Lihat performa lapangan, jam ramai, dan pemasukan supaya keputusan bisnis lebih cepat.",
      },
    ]

    const amenitys = [
      { icon: Dumbbell, label: "Booking Online" },
      { icon: Volleyball, label: "Jadwal Real-time" },
      { icon: Baby, label: "Manajemen Lapangan" },
      { icon: PersonStanding, label: "Data Pelanggan" },
      { icon: Trees, label: "Booking Manual" },
      { icon: Armchair, label: "Riwayat Booking" },
      { icon: Footprints, label: "Transaksi" },
      { icon: Flower2, label: "Pemasukan" },
      { icon: Bike, label: "Laporan Keuangan" },
      { icon: Fence, label: "Dashboard Analytics" },
      { icon: PawPrint, label: "Manajemen Staff" },
      { icon: Waves, label: "Website Bisnis" },
    ]

    const testimonials = [
      {
        name: "Raka Pratama",
        role: "Owner Futsal Arena",
        quote: "Jadwal booking jauh lebih rapi.",
        img: "/avatars/raka.png",
      },
      {
        name: "Dina Lestari",
        role: "Manajer Badminton Center",
        quote:
          "Pembayaran, pelanggan, dan laporan harian sekarang bisa dipantau dari satu dashboard tanpa rekap manual tiap malam.",
        img: "/avatars/dina.png",
      },
      {
        name: "Bima Santoso",
        role: "Pengelola Mini Soccer",
        quote:
          "Pelanggan suka karena bisa lihat slot kosong sendiri. Admin kami jadi lebih hemat waktu, terutama pas jam ramai dan akhir pekan.",
        img: "/avatars/bima.png",
      },
      {
        name: "Nadia Putri",
        role: "Owner Tennis Court",
        quote:
          "Kami jadi tahu jam ramai, pemasukan harian, dan lapangan mana yang paling sering dipakai.",
        img: "/avatars/nadia.png",
      },
      {
        name: "Arif Hidayat",
        role: "Operator Sport Center",
        quote: "Sistemnya simpel.",
        img: "/avatars/arif.png",
      },
      {
        name: "Salsa Maharani",
        role: "Founder CourtHub",
        quote:
          "Dulu banyak booking bentrok. Setelah pakai LapanganKu.id, semua slot lebih jelas, mudah dicek, dan gampang dibagikan ke tim operasional.",
        img: "/avatars/salsa.png",
      },
      {
        name: "Yoga Firmansyah",
        role: "Owner Basket Court",
        quote:
          "Notifikasi dan riwayat transaksi bikin follow-up pelanggan lebih cepat.",
        img: "/avatars/yoga.png",
      },
      {
        name: "Maya Anggraini",
        role: "Finance Sport Venue",
        quote:
          "Laporan pemasukan tidak perlu direkap manual lagi. Data harian langsung kelihatan, jadi pengecekan kas tutup shift jauh lebih cepat dan minim salah hitung.",
        img: "/avatars/maya.png",
      },
    ]

    const faqGroup = [
      {
        label: "Basics",
        items: [
          {
            question: "Apa itu LapanganKu.id?",
            answer:
              "LapanganKu.id adalah platform untuk mengelola booking, jadwal, pelanggan, pembayaran, dan laporan bisnis lapangan olahraga dari satu dashboard.",
          },
          {
            question: "Siapa yang cocok memakai LapanganKu.id?",
            answer:
              "LapanganKu.id cocok untuk pengelola futsal, badminton, mini soccer, basket, tenis, dan sport center dengan sistem booking harian.",
          },
          {
            question: "Apakah pelanggan bisa booking sendiri?",
            answer:
              "Bisa. Pelanggan dapat melihat slot tersedia dan mengirim booking tanpa harus menunggu balasan admin.",
          },
        ],
      },
      {
        label: "Project Approach",
        items: [
          {
            question: "Apakah website bisa disesuaikan dengan brand kami?",
            answer:
              "Bisa. Tampilan, informasi venue, daftar lapangan, dan alur booking dapat disesuaikan dengan kebutuhan operasional Anda.",
          },
          {
            question: "Apakah data lama bisa dimasukkan?",
            answer:
              "Bisa dibantu untuk data penting seperti daftar lapangan, pelanggan, harga, dan jadwal dasar.",
          },
          {
            question: "Berapa lama setup awal?",
            answer:
              "Setup dasar biasanya cepat setelah data venue, harga, jam operasional, dan aturan booking sudah lengkap.",
          },
        ],
      },
      {
        label: "Industries",
        items: [
          {
            question: "Apakah hanya untuk lapangan futsal?",
            answer:
              "Tidak. Sistem dapat digunakan untuk berbagai venue olahraga yang memakai jadwal slot dan pembayaran booking.",
          },
          {
            question: "Bisa untuk banyak cabang?",
            answer:
              "Bisa. Setiap cabang dapat dikelola dengan data lapangan, jadwal, dan laporan masing-masing.",
          },
          {
            question: "Bisa untuk event atau turnamen?",
            answer:
              "Bisa dibuatkan pengaturan khusus jika venue sering menerima event, block schedule, atau booking rombongan.",
          },
        ],
      },
      {
        label: "Process",
        items: [
          {
            question: "Bagaimana proses mulai pakai?",
            answer:
              "Mulai dari konsultasi kebutuhan, pengisian data venue, setup sistem, uji alur booking, lalu go live.",
          },
          {
            question: "Apakah ada bantuan setelah launch?",
            answer:
              "Ada. Tim dapat membantu penyesuaian, perbaikan data, dan pertanyaan operasional setelah sistem berjalan.",
          },
          {
            question: "Apakah bisa request fitur tambahan?",
            answer:
              "Bisa. Request akan dicek berdasarkan kebutuhan bisnis, dampak operasional, dan scope pengerjaan.",
          },
        ],
      },
    ]

    let activeGroup = $state(faqGroup[0]);
    let activeQuestion = $state(faqGroup[0].items[0].question);

    /** @param {(typeof faqGroup)[number]} group */
    function selectFaqGroup(group) {
      activeGroup = group;
      activeQuestion = group.items[0].question;
    }
</script>

<main class="container max-w-7xl mx-auto">
    
    <section>
        <div class="relative z-10 mx-auto flex max-w-6xl flex-col items-center px-4 pt-12 text-center md:pt-20">
            <BadgeMain>
                <span>Platform Manajemen Lapangan Olahraga</span>
                <ChevronRight size={16} />
            </BadgeMain>
            
            <h1 class="mt-7 text-4xl font-bold tracking-tight md:text-5xl lg:text-6xl">
                Semua yang Anda Butuhkan <br class="hidden md:block" />
                untuk Mengelola Bisnis Lapangan.
            </h1>
            <p class="mt-6 max-w-4xl text-base text-muted-foreground md:text-lg">
                Dari website booking hingga pengelolaan jadwal, pelanggan, transaksi,
                dan keuangan. Satu platform untuk membuat operasional lebih efisien dan
                memberikan layanan yang lebih baik kepada pelanggan.
            </p>
    
            <div class="mt-8 flex items-center gap-3 sm:flex-row">
                <Button
                    variant="default"
                    href="/auth/register"
                    class="px-6 py-5 shadow-lg transition hover:-translate-y-0.5 font-semibold"
                >
                    Lihat Demo
                </Button>
                <Button
                    variant="outline"
                    href="/auth/register"
                    class="px-6 py-5 shadow-lg transition hover:-translate-y-0.5 font-semibold"
                >
                    Mulai Sekarang
                </Button>
            </div>
            
            <p class="mt-4 text-xs text-muted-foreground">
                ✓ Mudah digunakan    ✓ Tanpa instalasi    ✓ Bisa diakses dari mana saja
            </p>
        </div>
    
        <div class="relative z-10 mx-auto mt-12 w-full max-w-7xl pb-4 md:mt-16 px-4">
          <div class="rounded-2xl border border-white/30 bg-white/20 p-2 shadow-2xl backdrop-blur-md md:rounded-[2rem] md:p-3">
            <a href="/demo" aria-label="Buka demo dashboard LapanganKu.id">
              <img
                src="/images/dashboard-admin.png"
                alt="Tampilan dashboard admin LapanganKu.id"
                width={1440}
                height={1000}
                sizes="(max-width: 1280px) 100vw, 1280px"
                class="w-full rounded-xl border border-black/5 shadow-lg transition duration-300 hover:brightness-[0.98] md:rounded-3xl"
              />
            </a>
          </div>
        </div>
    </section>

    <section class="py-16 md:py-24">
        <section>
          <div class="lg:py-20 sm:py-16 py-8">
            <div class="mx-auto max-w-7xl px-4 sm:px-8">
              <div class="flex flex-col gap-8 md:gap-16">
                <div
                  class="flex flex-col items-center justify-center gap-4 max-w-2xl mx-auto"
                >
                    
                  <BadgeMain>
                    <span>Core Values</span>
                  </BadgeMain>
                  
                  <h2 class="text-3xl md:text-4xl font-semibold text-center tracking-[-1px]">
                    Dirancang untuk Membuat Operasional
                    Bisnis Lapangan Lebih Sederhana.
                  </h2>
                </div>
                <div
                  class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6"
                >
                  {#each coreValues as coreValue}
                      <div>
                        <Card class="py-10 h-full border-t-4 border-t-transparent transition-all duration-300 hover:border-t-primary hover:shadow-lg">
                          <CardContent class="px-8 flex flex-col gap-6">
                            <coreValue.icon
                              class="w-8 h-8 text-primary"
                              strokeWidth={1.2}
                            />
                            <div class="flex flex-col gap-3">
                              <h3 class="text-xl font-semibold">
                                {coreValue.title}
                              </h3>
                              <p class="text-base font-normal text-muted-foreground">
                                {coreValue.content}
                              </p>
                            </div>
                          </CardContent>
                        </Card>
                      </div>
                  {/each}
                </div>
                <div
                  class="flex flex-col items-center justify-center gap-5"
                >
                  <div class="flex items-center gap-2 text-muted-foreground">
                    <Star size={16} />
                    <p class="font-normal text-sm">
                      Satu platform untuk membantu operasional lapangan berjalan lebih sederhana.
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>
    </section>

    <section class="py-16 md:py-24">
      <div class="mx-auto max-w-7xl px-4 sm:px-8">
        <div class="grid gap-10 lg:grid-cols-[1.1fr_0.9fr] lg:items-center">
          <div
              class="aspect-[4/3] overflow-hidden rounded-3xl border bg-muted/40 p-3 shadow-sm"
            >
              <a
                href="/demo/operasional/jadwal"
                aria-label="Buka demo pengelolaan jadwal LapanganKu.id"
                class="block h-full"
              >
                <img
                  src="/images/dashboard-jadwal.png"
                  alt="Tampilan pengelolaan jadwal LapanganKu.id"
                  width={1440}
                  height={1000}
                  sizes="(max-width: 1024px) 100vw, 704px"
                  class="h-full w-full rounded-2xl object-cover object-left-top transition duration-300 hover:scale-[1.01]"
                />
              </a>
          </div>

          <div class="flex flex-col justify-center gap-7">
            <div class="max-w-xl">
              
                <div
                  class="w-fit"
                >
                  <BadgeMain>
                    <span>Features</span>
                  </BadgeMain>
                </div>
              
                <h2
                  class="mt-4 text-3xl font-semibold tracking-tight md:text-4xl"
                >
                  Cara lebih baik mengelola semuanya
                </h2>
                <p
                  class="mt-4 text-muted-foreground"
                >
                  Satukan booking, pembayaran, dan operasional tim dalam satu tempat yang rapi.
                </p>
              </div>

              {#each features as feature}
                  <div
                    class="flex gap-4"
                  >
                    <div class="flex size-11 shrink-0 items-center justify-center rounded-xl text-primary border">
                      <feature.icon class="size-5" strokeWidth={1.8} />
                    </div>
                    <div>
                      <p class="font-semibold">{feature.title}</p>
                      <p class="mt-2 text-sm leading-6 text-muted-foreground">
                        {feature.content}
                      </p>
                    </div>
                  </div>
              {/each}

          </div>
        </div>
      </div>
    </section>

    <section class="py-16 md:py-24">
      <div class="mx-auto max-w-7xl px-4 sm:px-8">
          <div class="grid gap-8 md:grid-cols-[1fr_0.85fr] md:items-end">
          <div>
            <BadgeMain>
              <span>Cara Kerja</span>
            </BadgeMain>
            
            <h2 class="mt-5 max-w-2xl text-3xl font-semibold tracking-tight md:text-5xl">
              Semua operasional lapangan, dalam satu workspace
            </h2>
          </div>
          <p class="text-muted-foreground md:text-base md:leading-7">
            Pantau booking, pembayaran, pelanggan, dan laporan harian secara real-time dari dashboard yang familiar untuk tim Anda.
          </p>
        </div>

        <div class="mt-12 overflow-hidden rounded-3xl border bg-card shadow-sm">
          <div class="grid divide-y md:grid-cols-3 md:divide-x md:divide-y-0">
              <div class="flex min-h-[360px] flex-col justify-between p-6 md:p-8">
                  <div class="flex min-h-40 items-center justify-center rounded-2xl bg-muted/40 p-4">
                      <div class="relative w-full max-w-xs">
                          <div class="rounded-2xl border bg-background p-4 shadow-lg">
                              <div class="mb-4 flex items-center justify-between">
                          <div>
                            <p class="text-xs text-muted-foreground">Hari ini</p>
                            <p class="text-sm font-semibold">Lapangan A</p>
                          </div>
                          <CalendarCheck class="size-5 text-primary" />
                        </div>
                        <div class="space-y-2">
                            <div class="flex items-center justify-between rounded-xl border bg-muted/30 px-3 py-2 text-xs">
                              <span class="font-medium">08:00</span>
                              <span class="text-muted-foreground">Booked</span>
                            </div>
                            <div class="flex items-center justify-between rounded-xl border bg-muted/30 px-3 py-2 text-xs">
                              <span class="font-medium">10:00</span>
                              <span class="text-muted-foreground">Kosong</span>
                            </div>
                            <div class="flex items-center justify-between rounded-xl border bg-muted/30 px-3 py-2 text-xs">
                              <span class="font-medium">12:00</span>
                              <span class="text-muted-foreground">Istirahat</span>
                            </div>
                        </div>
                      </div>
                      <div class="absolute -right-3 -top-3 rounded-xl border bg-background px-3 py-2 text-xs shadow-md">
                        Slot update real-time
                      </div>
                    </div>
                </div>
                <div class="mt-8 text-center">
                  <div class="mx-auto mb-4 flex size-11 items-center justify-center rounded-xl border bg-background text-primary shadow-xs">
                      <CalendarCheck class="size-5" strokeWidth={1.8} />
                  </div>
                  <h3 class="text-lg font-semibold">Jadwal selalu sinkron</h3>
                  <p class="mx-auto mt-3 max-w-xs text-sm leading-6 text-muted-foreground">
                    Booking online dan input admin masuk ke kalender yang sama, jadi slot kosong dan terisi terlihat jelas.
                  </p>
                </div>
              </div>

              <div class="flex min-h-[360px] flex-col justify-between p-6 md:p-8">
                  <div class="flex min-h-40 items-center justify-center rounded-2xl bg-muted/40 p-4">
                      <div class="relative w-full max-w-xs">
                          <div class="rounded-2xl border bg-background p-4 shadow-lg">
                              <div class="grid grid-cols-2 gap-3">
                                  <div class="rounded-xl border bg-muted/30 p-3">
                                      <CreditCard class="mb-5 size-5 text-primary" />
                                      <p class="text-xs text-muted-foreground">Pembayaran</p>
                                      <p class="text-sm font-semibold">Rp 1,8jt</p>
                          </div>
                          <div class="rounded-xl border bg-muted/30 p-3">
                            <Users class="mb-5 size-5 text-primary" />
                            <p class="text-xs text-muted-foreground">Pelanggan</p>
                            <p class="text-sm font-semibold">36 aktif</p>
                          </div>
                        </div>
                        <div class="mt-3 rounded-xl border bg-primary px-3 py-2 text-xs font-medium text-primary-foreground">
                          12 transaksi sudah cocok
                        </div>
                      </div>
                      <div class="absolute -bottom-3 left-5 rounded-full border bg-background px-3 py-1 text-xs shadow-md">
                        Rekap otomatis
                      </div>
                    </div>
                </div>
                <div class="mt-8 text-center">
                  <div class="mx-auto mb-4 flex size-11 items-center justify-center rounded-xl border bg-background text-primary shadow-xs">
                      <WalletCards class="size-5" strokeWidth={1.8} />
                  </div>
                  <h3 class="text-lg font-semibold">Pembayaran tercatat rapi</h3>
                  <p class="mx-auto mt-3 max-w-xs text-sm leading-6 text-muted-foreground">
                    DP, pelunasan, refund, dan status pembayaran tersimpan di satu tempat untuk mengurangi rekap manual.
                  </p>
                </div>
              </div>

              <div class="flex min-h-[360px] flex-col justify-between p-6 md:p-8">
                  <div class="flex min-h-40 items-center justify-center rounded-2xl bg-muted/40 p-4">
                      <div class="w-full max-w-xs rounded-2xl border bg-background p-4 shadow-lg">
                          <div class="mb-4 flex items-center gap-2">
                              <div class="flex size-8 items-center justify-center rounded-full bg-primary text-primary-foreground">
                                  <Bot class="size-4" />
                        </div>
                        <div>
                          <p class="text-sm font-semibold">Draft siap kirim</p>
                          <p class="text-xs text-muted-foreground">Pengingat booking</p>
                        </div>
                      </div>
                      <div class="rounded-xl border bg-muted/30 p-3 text-xs leading-5 text-muted-foreground">
                        Halo Kak, booking futsal jam 19:00 sudah dikonfirmasi. Mohon datang 10 menit lebih awal.
                      </div>
                      <div class="mt-3 flex items-center gap-2 rounded-xl border bg-muted/30 px-3 py-2 text-xs">
                        <LineChart class="size-4 text-primary" />
                        <span>Follow-up pelanggan repeat</span>
                      </div>
                    </div>
                </div>
                <div class="mt-8 text-center">
                  <div class="mx-auto mb-4 flex size-11 items-center justify-center rounded-xl border bg-background text-primary shadow-xs">
                      <Bot class="size-5" strokeWidth={1.8} />
                  </div>
                  <h3 class="text-lg font-semibold">Follow-up lebih cepat</h3>
                  <p class="mx-auto mt-3 max-w-xs text-sm leading-6 text-muted-foreground">
                    Kirim pengingat booking, konfirmasi pembayaran, dan promo pelanggan tanpa bolak-balik cek catatan.
                  </p>
                </div>
              </div>
          </div>
        </div>
        
      </div>
    </section>

    <section class="py-16 md:py-24">
      <div class="mx-auto max-w-7xl px-4 sm:px-8">
        <div class="mx-auto flex max-w-3xl flex-col items-center gap-4 text-center">
          <div>
            
            <BadgeMain>
              Lebih Banyak Fitur
            </BadgeMain>
            
          </div>
          <h2
            class="text-3xl font-semibold tracking-tight md:text-5xl"
          >
            Fitur yang Membantu Setiap Bagian dari Operasional Anda.
          </h2>
          <p
            class="text-muted-foreground md:text-lg"
          >
            Dari pelanggan melakukan booking hingga Anda melihat laporan bisnis, setiap fitur dirancang untuk membuat pekerjaan sehari-hari menjadi lebih sederhana.
          </p>
        </div>

        <div
          class="mx-auto mt-12 grid max-w-5xl grid-cols-2 overflow-hidden rounded-2xl border sm:grid-cols-3 lg:grid-cols-6"
        >
            {#each amenitys as amenity}
                <div
                class="group relative flex min-h-36 flex-col items-center justify-center gap-4 border-b border-r p-6 text-center transition hover:bg-muted/70"
                >
                <amenity.icon
                    class="size-5 text-muted-foreground transition group-hover:text-foreground"
                    strokeWidth={1.8}
                />
                <p class="text-sm font-medium text-muted-foreground transition group-hover:text-foreground">
                    {amenity.label}
                </p>
                <span class="absolute inset-x-0 bottom-0 h-0.5 origin-left scale-x-0 bg-foreground transition-transform duration-300 group-hover:scale-x-100"></span>
                </div>
            {/each}
        </div>
      </div>
    </section>

    <section class="mx-auto max-w-7xl px-4 py-16">
      <div class="mx-auto max-w-3xl text-center">
        <BadgeMain>
          <span>Testimonials</span>
        </BadgeMain>
        <h2 class="mt-3 text-3xl font-bold tracking-tight md:text-5xl">
          Dipakai Owner Lapangan Setiap Hari
        </h2>
        <p class="mt-4 text-muted-foreground">
          Cerita singkat dari pengelola venue yang sudah merapikan operasional dengan LapanganKu.id.
        </p>
      </div>

      <div class="mt-14 columns-1 gap-4 sm:columns-2 lg:columns-3 xl:columns-4">
        {#each testimonials as testimonial}
          <figure
            class="mb-4 break-inside-avoid rounded-xl border bg-background p-5 shadow-xs"
          >
            <blockquote class="text-sm leading-6 text-foreground">
              “{testimonial.quote}”
            </blockquote>
            <figcaption class="mt-6 flex items-center gap-3">
              <img
                src={testimonial.img}
                alt=""
                width={100}
                height={100}
                class="size-8 rounded-full border object-cover"
              />
              <div>
                <div class="text-sm font-semibold">{testimonial.name}</div>
                <div class="text-xs text-muted-foreground">
                  {testimonial.role}
                </div>
              </div>
            </figcaption>
          </figure>
	    {/each}
      </div>
    </section>

    <section class="py-16 md:py-24">
      <div class="mx-auto max-w-7xl px-4 sm:px-8">
        <div class="grid gap-10 lg:grid-cols-[0.9fr_1.1fr]">
            <div class="flex flex-col gap-10">
            <div
              class="flex flex-col items-start gap-4"
            >
                <Badge variant="outline" class="px-3 py-1 text-sm">
                  FAQs
                </Badge>
                <h2 class="text-3xl font-semibold tracking-tight md:text-5xl">
                  Semua yang Perlu Anda Ketahui.
                </h2>
                <p class="max-w-md text-muted-foreground">
                  Masih punya pertanyaan tentang LapanganKu.id? Berikut beberapa hal yang paling sering ditanyakan oleh pengelola lapangan.
                </p>
            </div>

            <div
              class="flex flex-col items-start gap-4"
            >
              <p class="font-medium">Masih punya pertanyaan? Tim kami siap membantu.</p>
              <Button size="lg">
                <Mail data-icon="inline-start" />
                Hubungi Kami
              </Button>
            </div>
          </div>

          <div
            class="flex flex-col gap-6"
          >
            <div class="flex flex-wrap gap-2 rounded-2xl border bg-background p-1">
              {#each faqGroup as group (group.label)}
                <button
                  type="button"
                  onclick={() => selectFaqGroup(group)}
                  class="rounded-xl px-4 py-2 text-sm font-medium text-muted-foreground transition hover:-translate-y-px hover:bg-muted/60 hover:text-foreground active:scale-[0.98] data-[active=true]:bg-muted data-[active=true]:text-foreground data-[active=true]:shadow-xs"
                  data-active={activeGroup.label === group.label}
                >
                  {group.label}
                </button>
              {/each}
            </div>

            {#key activeGroup.label}
              <Accordion type="single" bind:value={activeQuestion} class="gap-2 border-none">
                {#each activeGroup.items as item, index (item.question)}
                  <div in:fly={{ x: 24, duration: 400, delay: index * 60 }}>
                    <AccordionItem
                      value={item.question}
                      class="rounded-2xl border px-5 not-last:mb-3 not-last:border-b"
                    >
                      <AccordionTrigger class="py-5 text-base hover:no-underline">
                        {item.question}
                      </AccordionTrigger>
                      <AccordionContent class="pb-5 text-muted-foreground">
                        {item.answer}
                      </AccordionContent>
                    </AccordionItem>
                  </div>
                {/each}
              </Accordion>
            {/key}
          </div>
        </div>
      </div>
    </section>
</main>
