export type BlogBlock =
	| { type: 'paragraph'; text: string }
	| { type: 'heading'; text: string }
	| { type: 'list'; items: string[] }
	| { type: 'orderedList'; items: string[] }
	| { type: 'quote'; text: string; cite?: string }
	| { type: 'image'; src: string; alt?: string }
	| { type: 'divider' };

export type BlogPost = {
	id: string;
	slug: string;
	title: string;
	excerpt: string;
	category: string;
	image: string;
	author: string;
	authorImg: string;
	authorRole?: string;
	readTime: string;
	publishedAt: string;
	views?: number;
	content: BlogBlock[];
};

export type BlogVenue = {
	name: string;
	logoUrl?: string;
	city?: string;
	province?: string;
};

export const blogPosts: BlogPost[] = [
	{
		id: 'blog-1',
		slug: 'cara-mengurangi-booking-bentrok',
		title: 'Cara Mengurangi Booking Bentrok di Lapangan Olahraga',
		excerpt: 'Langkah praktis menyatukan jadwal admin dan booking pelanggan agar setiap slot selalu akurat.',
		category: 'Operasional',
		image: '/images/dashboard-jadwal.png',
		author: 'Raka Pratama',
		authorImg: '/avatars/raka.png',
		authorRole: 'Owner Futsal Arena',
		readTime: '6 menit',
		publishedAt: '28 September 2026',
		views: 284,
		content: [
			{ type: 'paragraph', text: 'Booking bentrok biasanya terjadi ketika jadwal dicatat di beberapa tempat. Admin menerima pesan WhatsApp, pelanggan lain menelepon, lalu perubahan tidak langsung masuk ke kalender utama.' },
			{ type: 'heading', text: 'Gunakan satu sumber jadwal' },
			{ type: 'paragraph', text: 'Pastikan booking online dan input manual admin masuk ke kalender yang sama. Setiap perubahan status harus langsung memperbarui ketersediaan slot.' },
			{ type: 'list', items: ['Tentukan durasi slot yang konsisten.', 'Blokir waktu perawatan dan kegiatan internal.', 'Catat DP dan batas waktu pembayaran.', 'Batasi akses perubahan jadwal sesuai peran staff.'] },
			{ type: 'quote', text: 'Jadwal yang rapi bukan hanya memudahkan admin, tetapi juga menjaga kepercayaan pelanggan.', cite: 'Tim LapanganKu.id' },
			{ type: 'heading', text: 'Evaluasi setiap minggu' },
			{ type: 'paragraph', text: 'Periksa penyebab pembatalan, perubahan mendadak, dan slot yang sering bentrok. Data ini membantu Anda memperbaiki aturan booking.' }
		]
	},
	{
		id: 'blog-2',
		slug: 'strategi-meningkatkan-booking-jam-sepi',
		title: 'Strategi Meningkatkan Booking pada Jam Sepi',
		excerpt: 'Gunakan data jadwal untuk menyusun promo yang tepat tanpa menurunkan harga di jam ramai.',
		category: 'Pemasaran',
		image: '/images/dashboard-admin.png',
		author: 'Dina Lestari',
		authorImg: '/avatars/dina.png',
		authorRole: 'Manajer Badminton Center',
		readTime: '5 menit',
		publishedAt: '24 September 2026',
		views: 197,
		content: [
			{ type: 'paragraph', text: 'Promo yang efektif dimulai dari data. Cari jam dengan tingkat keterisian rendah selama beberapa minggu, lalu tentukan pelanggan yang paling mungkin memakai slot tersebut.' },
			{ type: 'heading', text: 'Buat penawaran yang spesifik' },
			{ type: 'orderedList', items: ['Pilih jam sepi yang konsisten.', 'Buat paket untuk komunitas atau pelajar.', 'Kirim penawaran ke pelanggan yang relevan.', 'Ukur perubahan okupansi dan pendapatan.'] },
			{ type: 'paragraph', text: 'Hindari diskon menyeluruh. Fokus pada waktu dan segmen yang memang membutuhkan dorongan.' }
		]
	},
	{
		id: 'blog-3',
		slug: 'laporan-harian-venue-yang-wajib-dipantau',
		title: 'Laporan Harian Venue yang Wajib Dipantau',
		excerpt: 'Ringkasan angka penting untuk membantu owner memahami kondisi operasional setiap hari.',
		category: 'Keuangan',
		image: '/images/dashboard-admin.png',
		author: 'Maya Anggraini',
		authorImg: '/avatars/maya.png',
		authorRole: 'Finance Sport Venue',
		readTime: '7 menit',
		publishedAt: '20 September 2026',
		views: 331,
		content: [
			{ type: 'paragraph', text: 'Laporan harian sebaiknya singkat tetapi cukup untuk menunjukkan pendapatan, piutang, refund, dan perbedaan kas.' },
			{ type: 'heading', text: 'Angka utama' },
			{ type: 'list', items: ['Total booking dan okupansi.', 'Pendapatan tunai dan non-tunai.', 'DP yang belum dilunasi.', 'Pembatalan dan refund.', 'Selisih kas akhir shift.'] },
			{ type: 'paragraph', text: 'Bandingkan angka tersebut dengan hari dan minggu sebelumnya agar perubahan mudah terlihat.' }
		]
	},
	{
		id: 'blog-4',
		slug: 'sop-admin-lapangan-saat-jam-ramai',
		title: 'SOP Admin Lapangan Saat Jam Ramai',
		excerpt: 'Alur sederhana untuk menangani konfirmasi, pembayaran, dan pergantian slot dengan cepat.',
		category: 'Operasional',
		image: '/images/dashboard-jadwal.png',
		author: 'Bima Santoso',
		authorImg: '/avatars/bima.png',
		readTime: '4 menit',
		publishedAt: '16 September 2026',
		views: 155,
		content: [
			{ type: 'paragraph', text: 'Jam ramai membutuhkan pembagian tugas yang jelas. Satu admin memantau kedatangan, sementara staff lain menangani pembayaran dan kebutuhan lapangan.' },
			{ type: 'list', items: ['Konfirmasi pelanggan sebelum slot dimulai.', 'Tandai pembayaran segera setelah diterima.', 'Catat perubahan durasi atau lapangan.', 'Siapkan prosedur keterlambatan yang konsisten.'] }
		]
	},
	{
		id: 'blog-5',
		slug: 'membangun-loyalitas-pelanggan-venue',
		title: 'Membangun Loyalitas Pelanggan Venue Olahraga',
		excerpt: 'Pengalaman booking yang konsisten membantu pelanggan kembali tanpa promo berlebihan.',
		category: 'Pelanggan',
		image: '/images/dashboard-admin.png',
		author: 'Nadia Putri',
		authorImg: '/avatars/nadia.png',
		readTime: '5 menit',
		publishedAt: '12 September 2026',
		views: 208,
		content: [
			{ type: 'paragraph', text: 'Loyalitas tumbuh dari pengalaman yang mudah diprediksi: jadwal akurat, respons cepat, lapangan siap, dan pembayaran jelas.' },
			{ type: 'heading', text: 'Kenali kebiasaan pelanggan' },
			{ type: 'paragraph', text: 'Riwayat booking membantu tim menawarkan slot dan layanan yang relevan tanpa mengirim pesan massal yang mengganggu.' }
		]
	},
	{
		id: 'blog-6',
		slug: 'memilih-software-manajemen-lapangan',
		title: 'Panduan Memilih Software Manajemen Lapangan',
		excerpt: 'Daftar kebutuhan penting sebelum venue memindahkan operasional dari catatan manual.',
		category: 'Teknologi',
		image: '/images/dashboard-jadwal.png',
		author: 'Yoga Firmansyah',
		authorImg: '/avatars/yoga.png',
		readTime: '8 menit',
		publishedAt: '8 September 2026',
		views: 412,
		content: [
			{ type: 'paragraph', text: 'Software yang tepat harus menyelesaikan masalah harian, bukan menambah langkah kerja.' },
			{ type: 'orderedList', items: ['Petakan alur booking saat ini.', 'Tentukan data yang wajib dipindahkan.', 'Uji kemudahan penggunaan untuk admin.', 'Periksa laporan, akses staff, dan dukungan.', 'Mulai dari satu venue sebelum ekspansi.'] }
		]
	},
	{
		id: 'blog-7',
		slug: 'mengukur-okupansi-lapangan',
		title: 'Cara Mengukur Okupansi Lapangan dengan Benar',
		excerpt: 'Pahami perbedaan slot tersedia, terjual, diblokir, dan dibatalkan sebelum membaca performa venue.',
		category: 'Analitik',
		image: '/images/dashboard-admin.png',
		author: 'Salsa Maharani',
		authorImg: '/avatars/salsa.png',
		authorRole: 'Founder CourtHub',
		readTime: '6 menit',
		publishedAt: '3 September 2026',
		views: 267,
		content: [
			{ type: 'paragraph', text: 'Okupansi adalah persentase slot yang terjual dari slot yang benar-benar tersedia untuk dijual.' },
			{ type: 'quote', text: 'Jangan memasukkan waktu perawatan sebagai slot kosong. Angka okupansi akan terlihat lebih rendah dari kondisi sebenarnya.' },
			{ type: 'paragraph', text: 'Pisahkan laporan berdasarkan hari, jam, jenis lapangan, dan sumber booking untuk melihat peluang perbaikan.' }
		]
	}
];

export const blogCategories = ['Semua', ...new Set(blogPosts.map((post) => post.category))];

export function getBlogPost(slug: string) {
	return blogPosts.find((post) => post.slug === slug);
}

export function getRelatedPosts(post: BlogPost, limit = 3) {
	return blogPosts
		.filter((candidate) => candidate.id !== post.id)
		.sort((a, b) => Number(b.category === post.category) - Number(a.category === post.category))
		.slice(0, limit);
}
