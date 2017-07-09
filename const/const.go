package konst

const DBGuru = "teachers"
const DBSiswa = "students"
const DBOrtu = "parents"
const DBUser = "users"
const DBName = "studenthack"

const HeaderToken = "Auth"
const HeaderSession = "Session"

const Letter = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

//GetRoleString berguna untuk mencari nama role dari kode role (1, 2, 3, dan 4)
func GetRoleString(role int) (bool, string) {
	/* Role
	   1. Siswa
	   2. Guru
	   3. Ortu
	   4. School Regulator
	*/
	if role == 1 {
		return true, "Siswa"
	} else if role == 2 {
		return true, "Guru"
	} else if role == 3 {
		return true, "Ortu"
	} else if role == 4 {
		return true, "SR"
	}

	return false, "Role Tidak Dikenal"
}
