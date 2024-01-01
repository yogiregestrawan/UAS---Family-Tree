package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Struct untuk menyimpan data keluarga
type Family struct {
	Name         string   `json:"name"`
	Gender       string   `json:"gender"`
	Ancestors      map[string]string `json:"Ancestors"`
	Children     []string `json:"children"`
}

// Slice untuk menyimpan semua data keluarga
var families []Family

const filename = "family_data.json"

// Fungsi untuk menyimpan data keluarga ke file JSON
func simpanDataKeFile() error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(families)
	if err != nil {
		return err
	}

	return nil
}

// Fungsi untuk membaca data keluarga dari file JSON
func ambilDataDariFile() error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// File tidak ditemukan
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&families)
	if err != nil {
		return err
	}

	return nil
}

// Fungsi untuk menambahkan keluarga baru
func tambahKeluarga() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Nama Keluarga: ")
	scanner.Scan()
	name := scanner.Text()

	fmt.Print("Jenis Kelamin (l/p): ")
	scanner.Scan()
	gender := scanner.Text()

	newFamily := Family{
		Name:   name,
		Gender: gender,
		Ancestors:  make(map[string]string),
	}

	families = append(families, newFamily)

	err := simpanDataKeFile()
	if err != nil {
		fmt.Println("Gagal menyimpan data ke file:", err)
		return
	}

	fmt.Println("Keluarga berhasil ditambahkan!")
}

// Fungsi untuk mencari keluarga berdasarkan nama
func cariKeluarga() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Masukkan nama keluarga: ")
	scanner.Scan()
	name := scanner.Text()

	found := false
	for _, family := range families {
		if strings.EqualFold(family.Name, name) {
			fmt.Println("Keluarga ditemukan!")
			fmt.Printf("Nama Keluarga: %s\nJenis Kelamin: %s\n", family.Name, family.Gender)
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Keluarga tidak ditemukan.")
	}
}

// Fungsi untuk menghapus keluarga berdasarkan nama
func hapusKeluarga() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Daftar Keluarga:")
	for i, family := range families {
		fmt.Printf("%d. %s\n", i+1, family.Name)
	}

	fmt.Print("Pilih Keluarga yang akan dihapus [1-" + fmt.Sprint(len(families)) + "]: ")
	scanner.Scan()
	index := scanner.Text()

	i := 0
	fmt.Sscan(index, &i)
	i--

	if i >= 0 && i < len(families) {
		// Hapus keluarga dan silsilahnya
		hapusSilSilah(families[i].Name)

		families = append(families[:i], families[i+1:]...)

		err := simpanDataKeFile()
		if err != nil {
			fmt.Println("Gagal menyimpan data ke file:", err)
			return
		}

		fmt.Println("Keluarga berhasil dihapus!")
	} else {
		fmt.Println("Pilihan tidak valid.")
	}
}

// Fungsi untuk menghapus silsilah keluarga
func hapusSilSilah(nama string) {
	for i, family := range families {
		// Hapus anak dari keluarga ini
		for j, child := range family.Children {
			if child == nama {
				families[i].Children = append(families[i].Children[:j], families[i].Children[j+1:]...)				
				break
			}
		}

		// Hapus keluarga ini dari orang tua lainnya
		for parentName := range family.Ancestors {
			if parentName == nama {
				delete(families[i].Ancestors, parentName) // hapus orang tua				
				break
			}
		}
	}
}

// Fungsi untuk mengaitkan keluarga
func kaitkanKeluarga() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Daftar Keluarga:")
	for i, family := range families {
		fmt.Printf("%d. %s\n", i+1, family.Name)
	}

	fmt.Print("Pilih Keluarga yang akan dikaitkan [1-" + fmt.Sprint(len(families)) + "]: ")
	scanner.Scan()
	index1 := scanner.Text()

	i1 := 0
	fmt.Sscan(index1, &i1)
	i1--

	if i1 >= 0 && i1 < len(families) {
		fmt.Println("Daftar Keluarga yang akan dikaitkan:")
		for i, family := range families {
			fmt.Printf("%d. %s\n", i+1, family.Name)
		}

		fmt.Print("Pilih Keluarga yang akan dikaitkan [1-" + fmt.Sprint(len(families)) + "]: ")
		scanner.Scan()
		index2 := scanner.Text()

		i2 := 0
		fmt.Sscan(index2, &i2)
		i2--

		if i2 >= 0 && i2 < len(families) {
			fmt.Print(families[i1].Name, " dengan merupakan [kakek buyut/nenek buyut/kakek/nenek/ayah/ibu/anak] dari ", families[i2].Name, ": ")
			
			scanner.Scan()
			relation := scanner.Text()

			switch relation {
				case "ayah", "ibu", "kakek", "nenek", "kakek buyut", "nenek buyut":
					families[i2].Ancestors[families[i1].Name] = relation
				case "anak":
					families[i1].Children = append(families[i1].Children, families[i2].Name)
				default:
					fmt.Println("Hubungan tidak valid.")
					return
			}

			err := simpanDataKeFile()
			if err != nil {
				fmt.Println("Gagal menyimpan data ke file:", err)
				return
			}

			fmt.Println("Keluarga berhasil dikaitkan!")
		} else {
			fmt.Println("Pilihan tidak valid.")
		}
	} else {
		fmt.Println("Pilihan tidak valid.")
	}
}

// Fungsi untuk mencari silsilah keluarga
func cariSilsilah() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Daftar Keluarga:")
	for i, family := range families {
		fmt.Printf("%d. %s\n", i+1, family.Name)
	}

	fmt.Print("Pilih Keluarga yang akan dicari silsilahnya [1-" + fmt.Sprint(len(families)) + "]: ")
	scanner.Scan()
	index := scanner.Text()

	i := 0
	fmt.Sscan(index, &i)
	i--

	if i >= 0 && i < len(families) {
		// Menampilkan data leluhur
		fmt.Println("Data leluhur dari", families[i].Name, "adalah:")
		tampilkanLeluhur(families[i].Ancestors)
		fmt.Println()

		// Menampilkan data keturunan
		fmt.Println("Data keturunan dari", families[i].Name, "adalah:")
		tampilkanKeturunan(families[i].Children)

	} else {
		fmt.Println("Pilihan tidak valid.")
	}
}

// Fungsi untuk mencari keluarga berdasarkan nama
func cariKeluargaBerdasarkanNama(name string) Family {
	for _, family := range families {
		if strings.EqualFold(family.Name, name) {
			return family
		}
	}
	return Family{}
}

// Fungsi untuk menampilkan leluhur keluarga
func tampilkanLeluhur(ancestors map[string]string) {
	for parentName, relation := range ancestors {
		parent := cariKeluargaBerdasarkanNama(parentName)
		switch relation {
			case "ibu":
				fmt.Println("Ibu:", parent.Name)
			case "ayah":
				fmt.Println("Ayah:", parent.Name)
			case "kakek":
				fmt.Println("Kakek:", parent.Name)
			case "nenek":
				fmt.Println("Nenek:", parent.Name)
			case "kakek buyut":
				fmt.Println("Kakek Buyut:", parent.Name)
			case "nenek buyut":
				fmt.Println("Nenek Buyut:", parent.Name)
			default:
				fmt.Println("Hubungan tidak valid.")
		}		
	}
}

// Fungsi untuk menampilkan keturunan keluarga
func tampilkanKeturunan(children []string) {
	fmt.Print("Anak: ")
	for i, childName := range children {
		child := cariKeluargaBerdasarkanNama(childName)

		fmt.Print(child.Name)		

		// Tambahkan koma jika bukan elemen terakhir
		if i < len(children)-1 {
				fmt.Print(", ")
		}
	}
}

// Fungsi utama program
func main() {
	// Membaca data dari file saat program dimulai
	err := ambilDataDariFile()
	if err != nil {
		fmt.Println("Gagal membaca data dari file:", err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nProgram Pencatatan Silsilah Keluarga")
		fmt.Println("====================================")
		fmt.Println("1. Tambah Keluarga")
		fmt.Println("2. Cari Keluarga")
		fmt.Println("3. Hapus Keluarga")
		fmt.Println("4. Kaitkan Keluarga")
		fmt.Println("5. Cari Silsilah")
		fmt.Println("6. Keluar")

		fmt.Print("Pilih Menu [1-6]: ")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			tambahKeluarga()
		case "2":
			cariKeluarga()
		case "3":
			hapusKeluarga()
		case "4":
			kaitkanKeluarga()
		case "5":
			cariSilsilah()
		case "6":
			// Menyimpan data ke file sebelum keluar dari program
			err := simpanDataKeFile()
			if err != nil {
				fmt.Println("Gagal menyimpan data ke file:", err)
			}
			fmt.Println("Terima kasih telah menggunakan program ini. Sampai jumpa!")
			os.Exit(0)
		default:
			fmt.Println("Pilihan tidak valid. Silakan pilih menu 1-6.")
		}
	}
}
