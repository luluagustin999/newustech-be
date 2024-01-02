// terbentuknya entity ini karena cycle yaitu di package payment kita manggil transaction dan di paket transaction kita memanggil payment jadi tidak boleh saling mengambil satu sama lain untuk package

package payment 

type Transaction struct {
	ID int
	Amount int
}

