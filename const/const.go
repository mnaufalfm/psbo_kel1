package konst

//Nama koleksi Guru pada database
const DBGuru = "teachers"

//Nama koleksi Siswa pada database
const DBSiswa = "students"

//Nama koleksi Orang Tua (Ortu) pada database
const DBOrtu = "parents"

//Nama koleksi User pada database. Digunakan untuk Login dan Logout
const DBUser = "users"

//Nama koleksi Post pada database
const DBPost = "posts"

//Nama dari database
const DBName = "studenthack"

//Nama dari Header untuk menyimpan token/jwt/authorization
const HeaderToken = "Auth"

//Nama dari Header untuk menyimpan session
const HeaderSession = "Session"

//List dari gabungan huruf kecil, huruf kapital, dan angka. Dimulai huruf kecil (a...z), huruf besar (A...Z), dan diakhiri angka (0...9)
const Letter = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

//GetRoleString berguna untuk mencari nama role dari kode role (1, 2, dan 3)
func GetRoleString(role int) (bool, string) {
	/* Role
	   1. Siswa
	   2. Guru
	   3. Ortu
	*/
	if role == 1 {
		return true, "siswa"
	} else if role == 2 {
		return true, "guru"
	} else if role == 3 {
		return true, "ortu"
	}
	// } else if role == 4 {
	// 	return true, "SR"
	// }

	return false, "Role Tidak Dikenal"
}
