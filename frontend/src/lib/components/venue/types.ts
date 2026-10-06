/** Data wilayah administratif Indonesia (dari proxy /api/regions). */
export type Region = {
	code: string;
	name: string;
};

/** Pilihan alamat yang disusun dari dropdown wilayah. */
export type AddressSelection = {
	alamat: string;
	kota: string;
	kecamatan: string;
	kelurahan: string;
	provinsi: string;
	kodePos: string;
};

/** Status hasil pengecekan ketersediaan slug. */
export type SlugStatus =
	| 'idle'
	| 'checking'
	| 'available'
	| 'taken'
	| 'invalid'
	| 'error';

/** Satu slide pada carousel halaman auth. */
export type Slide = {
	src: string;
	title: string;
	subtitle: string;
};
