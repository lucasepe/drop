package hostdir

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestLoadDict(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  map[string]string
	}{
		{
			name: "File valido con più utenti",
			lines: []string{
				"user1:$5$rounds=77777$short$JiO1O3ZpDAxGJeaDIuqCoEFysAe1mZNJRs3pw0KQRd/",
				"user2:$5$rounds=123456$asaltof16chars..$gP3VQ/6X7UUEW3HkBn2w1/Ptq2jxPyzV/cZKmF/wJvD",
				"user3:$5$rounds=1000$roundstoolow$yfvwcWrQ8l/K0DAWyuPMDNHpIVlTQebY9l/gL972bIC",
			},
			want: map[string]string{
				"user1": "$5$rounds=77777$short$JiO1O3ZpDAxGJeaDIuqCoEFysAe1mZNJRs3pw0KQRd/",
				"user2": "$5$rounds=123456$asaltof16chars..$gP3VQ/6X7UUEW3HkBn2w1/Ptq2jxPyzV/cZKmF/wJvD",
				"user3": "$5$rounds=1000$roundstoolow$yfvwcWrQ8l/K0DAWyuPMDNHpIVlTQebY9l/gL972bIC",
			},
		},
		{
			name: "File con riga malformata",
			lines: []string{
				"user1:$5$saltstring$5B8vYYiY.CVt1RlTTf8KbXBH3hsxY/GNooZaBBGWEc5",
				"malformed_line",
				"user2:$5$rounds=10000$saltstringsaltst$3xv.VbSHBb41AL9AvLeujZkZRBAwqFMz2.opqey6IcA",
			},
			want: map[string]string{
				"user1": "$5$saltstring$5B8vYYiY.CVt1RlTTf8KbXBH3hsxY/GNooZaBBGWEc5",
				"user2": "$5$rounds=10000$saltstringsaltst$3xv.VbSHBb41AL9AvLeujZkZRBAwqFMz2.opqey6IcA",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := createTempFile(t, tt.lines)
			defer os.Remove(filename) // Pulizia del file temporaneo

			users, err := loadDict(filename)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(users, tt.want) {
				t.Fatalf("got %v, want: %v", users, tt.want)
			}
		})
	}
}

// Crea un file temporaneo per il test
func createTempFile(t *testing.T, lines []string) string {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "test_users.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpfile.Close()

	for _, ln := range lines {
		_, err = fmt.Fprintln(tmpfile, ln)
		if err != nil {
			t.Fatalf("unable to write line %q: %v", ln, err)
		}
	}

	return tmpfile.Name()
}
